// Package app defines the dependencies handed to every module's Register function.
package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/januarsyah901/ascard/api/internal/auth"
	"github.com/januarsyah901/ascard/api/internal/config"
	"github.com/januarsyah901/ascard/api/internal/httpx"
)

// Deps is everything a module needs to mount its routes.
type Deps struct {
	Config   config.Config
	Pool     *pgxpool.Pool // nil in handler tests; modules take repositories via interfaces
	JWT      *auth.JWT     // JWT.Sign for login, JWT.Guard to protect admin routes
	ClientIP httpx.ClientIP
}
