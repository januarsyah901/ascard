// Package articles serves /api/articles (SPEC §3.2).
package articles

import (
	"net/http"

	"github.com/januarsyah901/ascard/api/internal/app"
)

// Register mounts the article routes with the PostgreSQL repository.
func Register(mux *http.ServeMux, d app.Deps) {
	var repo Repo
	if d.Pool != nil {
		repo = NewRepo(d.Pool)
	}
	Mount(mux, d, Handler{Repo: repo})
}

// Mount registers h's endpoints. "GET /api/articles/admin/all" is more specific than
// "GET /api/articles/{id}", so ServeMux routes it first regardless of order.
func Mount(mux *http.ServeMux, d app.Deps, h Handler) {
	guard := func(f http.HandlerFunc) http.Handler { return d.JWT.Guard(f) }

	mux.HandleFunc("GET /api/articles", h.listPublished)
	mux.HandleFunc("GET /api/articles/{$}", h.listPublished) // Express routing is not strict about the trailing slash
	mux.Handle("GET /api/articles/admin/all", guard(h.listAll))
	mux.HandleFunc("GET /api/articles/{id}", h.get)
	mux.Handle("POST /api/articles", guard(h.create))
	mux.Handle("POST /api/articles/{$}", guard(h.create))
	mux.Handle("PUT /api/articles/{id}", guard(h.update))
	mux.Handle("DELETE /api/articles/{id}", guard(h.delete))
}
