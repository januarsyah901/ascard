// Package authapi serves POST /api/auth/login and GET /api/auth/me (docs/SPEC.md §3.1).
package authapi

import (
	"net/http"
	"time"

	"github.com/januarsyah901/ascard/api/internal/app"
	"github.com/januarsyah901/ascard/api/internal/ratelimit"
)

// Login rate limit: 5 requests per 15 minutes per client IP, like express-rate-limit in Express.
const (
	LoginLimit  = 5
	LoginWindow = 15 * time.Minute
)

// Register mounts the auth routes.
func Register(mux *http.ServeMux, d app.Deps) {
	h := Handler{Repo: NewRepo(d.Pool), Tokens: d.JWT}
	mount(mux, d, h, ratelimit.New(LoginLimit, LoginWindow))
}

func mount(mux *http.ServeMux, d app.Deps, h Handler, limiter *ratelimit.Limiter) {
	mux.Handle("POST /api/auth/login",
		limiter.Middleware(d.ClientIP.Resolve, MsgTooManyAttempts, http.HandlerFunc(h.Login)))
	mux.Handle("GET /api/auth/me", d.JWT.Guard(http.HandlerFunc(h.Me)))
}
