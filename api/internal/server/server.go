// Package server assembles the HTTP handler: global middleware, health check and every module.
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/januarsyah901/ascard/api/internal/app"
	"github.com/januarsyah901/ascard/api/internal/articles"
	"github.com/januarsyah901/ascard/api/internal/authapi"
	"github.com/januarsyah901/ascard/api/internal/httpx"
	"github.com/januarsyah901/ascard/api/internal/web"
)

// New returns the fully wired handler: the JSON API under /api and the SPA everywhere else.
func New(d app.Deps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		if d.Pool == nil {
			httpx.JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "error", "database": "unreachable"})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := d.Pool.Ping(ctx); err != nil {
			httpx.JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "error", "database": "unreachable"})
			return
		}
		httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok", "database": "ok"})
	})

	authapi.Register(mux, d)  // /api/auth/*      — worker go/auth
	articles.Register(mux, d) // /api/articles/*  — worker go/articles

	// Unknown API paths get a JSON 404 instead of the SPA.
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.StatusNotFound, httpx.MsgNotFound)
	})

	web.Register(mux, d) // "/" — built SPA + fallback — worker go/platform

	return httpx.Chain(mux,
		httpx.WithRequestID,
		httpx.AccessLog,
		httpx.Recover,
		web.SecurityHeaders,            // worker go/platform
		web.CORS(d.Config.CORSOrigins), // worker go/platform
	)
}
