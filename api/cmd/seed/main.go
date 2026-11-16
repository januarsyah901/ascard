// Command seed creates the admin account and the sample articles (replaces prisma/seed.js).
//
//	DATABASE_URL=... ADMIN_USERNAME=admin ADMIN_PASSWORD='at-least-12-chars' seed
//
// It is idempotent: an existing admin is left untouched (its password is never
// overwritten) and sample articles are only inserted when their title is not taken.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/bcrypt"

	"github.com/januarsyah901/ascard/api/internal/db"
)

// MinPasswordLength matches prisma/seed.js.
const MinPasswordLength = 12

// bcryptCost matches bcrypt.hash(password, 10) in prisma/seed.js.
const bcryptCost = 10

type article struct {
	title, content, category string
	published                bool
}

// sampleArticles are the articles of prisma/seed.js.
var sampleArticles = []article{
	{
		title:     "Anatomi Jebakan Dopamin",
		content:   "Dopamin bukan hormon kebahagiaan, melainkan hormon antisipasi. Di balik setiap putaran mesin slot, ada lonjakan dopamin yang memicu rasa penasaran berlebih. Secara klinis, sistem reward otak dibajak oleh ketidakpastian yang dirancang secara matematis.",
		category:  "Psikologi",
		published: true,
	},
	{
		title:     "Mitos Keberuntungan dalam RNG",
		content:   "Random Number Generator (RNG) tidak mengenal pola. Setiap putaran bersifat independen. Mempercayai \"pola\" adalah awal dari kekalahan finansial. Secara statistik, probabilitas tidak memiliki memori.",
		category:  "Matematika",
		published: true,
	},
	{
		title:     "Biaya Tersembunyi: Sunk Cost Fallacy",
		content:   "Semakin banyak uang yang Anda keluarkan, semakin sulit Anda berhenti. Ini adalah jebakan psikologis yang sering digunakan dalam desain kasino. Otak kita membenci kehilangan, sehingga kita terus \"berinvestasi\" pada kekalahan.",
		category:  "Finansial",
		published: true,
	},
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Getenv, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "seed:", err)
		os.Exit(1)
	}
}

type settings struct {
	databaseURL, username, password string
}

func loadSettings(getenv func(string) string) (settings, error) {
	s := settings{
		databaseURL: getenv("DATABASE_URL"),
		username:    getenv("ADMIN_USERNAME"),
		password:    getenv("ADMIN_PASSWORD"),
	}
	if s.username == "" {
		s.username = "admin"
	}
	// Length in characters like JavaScript's password.length (runes are close enough).
	if len([]rune(s.password)) < MinPasswordLength {
		return s, fmt.Errorf("set ADMIN_PASSWORD (min. %d characters) before seeding", MinPasswordLength)
	}
	if s.databaseURL == "" {
		return s, errors.New("DATABASE_URL is required")
	}
	return s, nil
}

// hashPassword hashes like bcryptjs, which silently uses only the first 72 bytes.
func hashPassword(password string) (string, error) {
	b := []byte(password)
	if len(b) > 72 {
		b = b[:72]
	}
	h, err := bcrypt.GenerateFromPassword(b, bcryptCost)
	return string(h), err
}

func run(ctx context.Context, getenv func(string) string, out io.Writer) error {
	s, err := loadSettings(getenv)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Starting seeding...")

	pool, err := db.Open(ctx, s.databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := db.EnsureSchema(ctx, pool); err != nil {
		return err
	}

	hash, err := hashPassword(s.password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	tag, err := pool.Exec(ctx,
		`INSERT INTO "Admin" ("id", "username", "password") VALUES (gen_random_uuid()::text, $1, $2)
		 ON CONFLICT ("username") DO NOTHING`, s.username, hash)
	if err != nil {
		return fmt.Errorf("upsert admin: %w", err)
	}
	if tag.RowsAffected() == 1 {
		fmt.Fprintln(out, "Admin created:", s.username)
	} else {
		fmt.Fprintln(out, "Admin exists (password unchanged):", s.username)
	}

	created := 0
	for _, a := range sampleArticles {
		tag, err := pool.Exec(ctx,
			`INSERT INTO "Article" ("id", "title", "content", "category", "published", "updatedAt")
			 VALUES (gen_random_uuid()::text, $1, $2, $3, $4, now())
			 ON CONFLICT ("title") DO NOTHING`, a.title, a.content, a.category, a.published)
		if err != nil {
			return fmt.Errorf("insert article %q: %w", a.title, err)
		}
		created += int(tag.RowsAffected())
	}
	fmt.Fprintf(out, "Sample articles: %d created, %d already present.\n", created, len(sampleArticles)-created)
	fmt.Fprintln(out, "Seeding finished successfully.")
	return nil
}
