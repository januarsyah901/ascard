// Package authapi serves POST /api/auth/login and GET /api/auth/me.
//
// TODO(go/auth): implement per docs/SPEC.md §3.1. This stub only exists so the server compiles.
package authapi

import (
	"net/http"

	"github.com/januarsyah901/ascard/api/internal/app"
)

// Register mounts the auth routes.
func Register(mux *http.ServeMux, d app.Deps) {
	_ = mux
	_ = d
}
