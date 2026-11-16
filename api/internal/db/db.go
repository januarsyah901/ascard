// Package db opens the PostgreSQL pool, ensures the schema exists and offers helpers
// shared by repositories.
//
// The schema was created by Prisma in the Express version; table and column names are
// PascalCase/camelCase and must be double-quoted in SQL, e.g. SELECT "createdAt" FROM "Article".
// Prisma generated ids and "updatedAt" on the client side, so the columns have no
// defaults for them: always pass gen_random_uuid()::text and now() yourself.
package db

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the subset of pgxpool.Pool / pgx.Tx that repositories use.
type DB interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Open creates a connection pool and verifies it with a ping.
func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}
	// Prisma stores TIMESTAMP(3) WITHOUT TIME ZONE values in UTC. Pin the session time
	// zone so now() and CURRENT_TIMESTAMP agree with existing rows on any server.
	cfg.ConnConfig.RuntimeParams["timezone"] = "UTC"
	// Accept Prisma-style URLs (?schema=public): Postgres has no "schema" parameter.
	if schema, ok := cfg.ConnConfig.RuntimeParams["schema"]; ok {
		delete(cfg.ConnConfig.RuntimeParams, "schema")
		if schema != "" {
			cfg.ConnConfig.RuntimeParams["search_path"] = schema
		}
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

//go:embed schema.sql
var schemaSQL string

// EnsureSchema creates the tables when they do not exist yet. It is idempotent and a
// no-op on a database that Prisma already migrated, so the Go server can take over the
// live database without a data migration.
func EnsureSchema(ctx context.Context, d DB) error {
	if _, err := d.Exec(ctx, schemaSQL); err != nil {
		return fmt.Errorf("ensure schema: %w", err)
	}
	return nil
}

// ErrNotFound is returned by repositories when a row does not exist.
var ErrNotFound = errors.New("not found")

// IsUniqueViolation reports whether err is a PostgreSQL unique-constraint error (SQLSTATE 23505).
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// NotFound maps pgx.ErrNoRows to ErrNotFound and passes other errors through.
func NotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
