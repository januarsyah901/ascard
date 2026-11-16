package authapi

import (
	"context"

	"github.com/januarsyah901/ascard/api/internal/db"
)

// Admin is a row of "Admin" as the handlers need it.
type Admin struct {
	ID           string
	Username     string
	PasswordHash string
}

// Repo loads admins. Both methods return db.ErrNotFound when no row matches.
type Repo interface {
	AdminByUsername(ctx context.Context, username string) (Admin, error)
	AdminByID(ctx context.Context, id string) (Admin, error)
}

type pgRepo struct{ db db.DB }

// NewRepo returns a Repo backed by PostgreSQL.
func NewRepo(d db.DB) Repo { return pgRepo{db: d} }

func (p pgRepo) AdminByUsername(ctx context.Context, username string) (Admin, error) {
	var a Admin
	err := p.db.QueryRow(ctx,
		`SELECT "id", "username", "password" FROM "Admin" WHERE "username" = $1`, username).
		Scan(&a.ID, &a.Username, &a.PasswordHash)
	return a, db.NotFound(err)
}

func (p pgRepo) AdminByID(ctx context.Context, id string) (Admin, error) {
	var a Admin
	err := p.db.QueryRow(ctx,
		`SELECT "id", "username", "password" FROM "Admin" WHERE "id" = $1`, id).
		Scan(&a.ID, &a.Username, &a.PasswordHash)
	return a, db.NotFound(err)
}
