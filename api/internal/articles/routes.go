// Package articles serves /api/articles.
//
// TODO(go/articles): implement per docs/SPEC.md §3.2. This stub only exists so the server compiles.
package articles

import (
	"net/http"

	"github.com/januarsyah901/ascard/api/internal/app"
)

// Register mounts the article routes.
func Register(mux *http.ServeMux, d app.Deps) {
	_ = mux
	_ = d
}
