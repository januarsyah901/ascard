package articles

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/januarsyah901/ascard/api/internal/db"
)

// Repo is the article storage. Get, Update and Delete return db.ErrNotFound for a
// missing row; Create and Update return an error matching db.IsUniqueViolation for a
// duplicate title.
type Repo interface {
	List(ctx context.Context, publishedOnly bool) ([]Article, error)
	GetPublished(ctx context.Context, id string) (Article, error)
	Create(ctx context.Context, in Input) (Article, error)
	Update(ctx context.Context, id string, in Input) (Article, error)
	Delete(ctx context.Context, id string) error
}

type pgRepo struct{ db db.DB }

// NewRepo returns the PostgreSQL implementation of Repo.
func NewRepo(d db.DB) Repo { return pgRepo{db: d} }

const columns = `"id", "title", "content", "category", "published", "createdAt", "updatedAt"`

func scanArticle(row pgx.Row) (Article, error) {
	var a Article
	err := row.Scan(&a.ID, &a.Title, &a.Content, &a.Category, &a.Published, &a.CreatedAt.Time, &a.UpdatedAt.Time)
	return a, err
}

func (r pgRepo) List(ctx context.Context, publishedOnly bool) ([]Article, error) {
	q := `SELECT ` + columns + ` FROM "Article"`
	if publishedOnly {
		q += ` WHERE "published" = true`
	}
	q += ` ORDER BY "createdAt" DESC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Article{}
	for rows.Next() {
		a, err := scanArticle(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r pgRepo) GetPublished(ctx context.Context, id string) (Article, error) {
	a, err := scanArticle(r.db.QueryRow(ctx,
		`SELECT `+columns+` FROM "Article" WHERE "id" = $1 AND "published" = true`, id))
	return a, db.NotFound(err)
}

func (r pgRepo) Create(ctx context.Context, in Input) (Article, error) {
	return scanArticle(r.db.QueryRow(ctx,
		`INSERT INTO "Article" ("id", "title", "content", "category", "published", "createdAt", "updatedAt")
		 VALUES (gen_random_uuid()::text, $1, $2, $3, $4, now(), now())
		 RETURNING `+columns,
		in.Title, in.Content, in.Category, in.Published))
}

func (r pgRepo) Update(ctx context.Context, id string, in Input) (Article, error) {
	a, err := scanArticle(r.db.QueryRow(ctx,
		`UPDATE "Article" SET "title" = $2, "content" = $3, "category" = $4, "published" = $5, "updatedAt" = now()
		 WHERE "id" = $1
		 RETURNING `+columns,
		id, in.Title, in.Content, in.Category, in.Published))
	return a, db.NotFound(err)
}

func (r pgRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM "Article" WHERE "id" = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return db.ErrNotFound
	}
	return nil
}
