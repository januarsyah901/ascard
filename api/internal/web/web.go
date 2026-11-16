// Package web serves the built React SPA and the global security/CORS middleware.
package web

import (
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/januarsyah901/ascard/api/internal/app"
	"github.com/januarsyah901/ascard/api/internal/httpx"
)

const (
	cacheImmutable = "public, max-age=31536000, immutable"
	cacheNoCache   = "no-cache"
	indexFile      = "index.html"
)

// Register mounts the SPA on "/", serving files from d.Config.WebDir.
func Register(mux *http.ServeMux, d app.Deps) {
	mux.Handle("/", NewSPA(d.Config.WebDir))
}

// SPA serves a built single-page app: real files with long-lived caching, index.html for
// client-side routes, and a hard 404 for missing hashed assets.
type SPA struct {
	root *os.Root // nil when the directory could not be opened: everything is 404
}

// NewSPA opens dir. A missing directory or index.html is logged, not fatal: the API still runs.
func NewSPA(dir string) *SPA {
	root, err := os.OpenRoot(dir)
	if err != nil {
		slog.Warn("web dir not available, SPA disabled", "web_dir", dir, "err", err)
		return &SPA{}
	}
	if fi, err := root.Stat(indexFile); err != nil || !fi.Mode().IsRegular() {
		slog.Warn("web dir has no index.html, SPA routes will answer 404", "web_dir", dir)
	}
	return &SPA{root: root}
}

func (s *SPA) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		httpx.Error(w, http.StatusNotFound, httpx.MsgNotFound)
		return
	}
	if s.root == nil {
		httpx.Error(w, http.StatusNotFound, httpx.MsgNotFound)
		return
	}

	// Rooted Clean drops any ".." that would climb above the web dir; os.Root additionally
	// refuses symlinks that escape it.
	name := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
	if name != "" && name != indexFile && servable(name) && s.serveFile(w, r, name, cacheImmutable, "") {
		return
	}
	if name == "assets" || strings.HasPrefix(name, "assets/") {
		httpx.Error(w, http.StatusNotFound, httpx.MsgNotFound)
		return
	}
	if !s.serveFile(w, r, indexFile, cacheNoCache, "text/html; charset=utf-8") {
		httpx.Error(w, http.StatusNotFound, httpx.MsgNotFound)
	}
}

// servable rejects hidden path segments (like express.static's dotfiles: "ignore").
func servable(name string) bool {
	for _, seg := range strings.Split(name, "/") {
		if strings.HasPrefix(seg, ".") {
			return false
		}
	}
	return true
}

// serveFile writes the regular file name if it exists. It reports false, having written
// nothing, when the file is missing or is a directory (no listings).
func (s *SPA) serveFile(w http.ResponseWriter, r *http.Request, name, cacheControl, contentType string) bool {
	f, err := s.root.Open(name)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			slog.DebugContext(r.Context(), "spa open", "name", name, "err", err)
		}
		return false
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil || !fi.Mode().IsRegular() {
		return false
	}
	h := w.Header()
	h.Set("Cache-Control", cacheControl)
	if contentType != "" {
		h.Set("Content-Type", contentType)
	}
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
	return true
}

// securityHeaders are helmet's defaults plus the CSP the Express server configures.
var securityHeaders = [][2]string{
	{"Content-Security-Policy", "default-src 'self';base-uri 'self';font-src 'self' https://fonts.gstatic.com;" +
		"form-action 'self';frame-ancestors 'self';img-src 'self' data: https:;object-src 'none';" +
		"script-src 'self';script-src-attr 'none';style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;" +
		"upgrade-insecure-requests;connect-src 'self'"},
	{"Cross-Origin-Opener-Policy", "same-origin"},
	{"Cross-Origin-Resource-Policy", "same-origin"},
	{"Origin-Agent-Cluster", "?1"},
	{"Referrer-Policy", "no-referrer"},
	{"Strict-Transport-Security", "max-age=31536000; includeSubDomains"},
	{"X-Content-Type-Options", "nosniff"},
	{"X-DNS-Prefetch-Control", "off"},
	{"X-Download-Options", "noopen"},
	{"X-Frame-Options", "SAMEORIGIN"},
	{"X-Permitted-Cross-Domain-Policies", "none"},
	{"X-XSS-Protection", "0"},
}

// SecurityHeaders sets helmet-equivalent headers on every response.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		for _, kv := range securityHeaders {
			h.Set(kv[0], kv[1])
		}
		next.ServeHTTP(w, r)
	})
}

const corsMethods = "GET,HEAD,PUT,PATCH,POST,DELETE"

// CORS allows cross-origin requests only from the configured origins. Requests from any
// other origin (or with no Origin) pass through without CORS headers.
func CORS(origins []string) httpx.Middleware {
	allowed := make(map[string]bool, len(origins))
	for _, o := range origins {
		allowed[o] = true
	}
	return func(next http.Handler) http.Handler {
		if len(allowed) == 0 {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" || !allowed[origin] {
				next.ServeHTTP(w, r)
				return
			}
			h := w.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Add("Vary", "Origin")
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				h.Set("Access-Control-Allow-Methods", corsMethods)
				if reqHeaders := r.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
					h.Set("Access-Control-Allow-Headers", reqHeaders)
				}
				h.Add("Vary", "Access-Control-Request-Headers")
				h.Set("Content-Length", "0")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
