// Package web serves the built React SPA and the global security/CORS middleware.
//
// TODO(go/platform): implement per docs/SPEC.md §3.3. These stubs only exist so the server compiles.
package web

import (
	"net/http"

	"github.com/januarsyah901/ascard/api/internal/app"
	"github.com/januarsyah901/ascard/api/internal/httpx"
)

// Register mounts the SPA on "/".
func Register(mux *http.ServeMux, d app.Deps) {
	_ = d
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.StatusNotFound, httpx.MsgNotFound)
	})
}

// SecurityHeaders sets helmet-equivalent headers on every response.
func SecurityHeaders(next http.Handler) http.Handler { return next }

// CORS allows cross-origin requests only from the configured origins.
func CORS(origins []string) httpx.Middleware {
	_ = origins
	return func(next http.Handler) http.Handler { return next }
}
