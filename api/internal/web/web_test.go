package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/januarsyah901/ascard/api/internal/app"
	"github.com/januarsyah901/ascard/api/internal/config"
)

const indexHTML = `<!doctype html><html><body><div id="root"></div><script src="/assets/index-abc.js"></script></body></html>`

// newWebDir builds <tmp>/secret.txt next to <tmp>/dist with a small SPA inside.
func newWebDir(t *testing.T) string {
	t.Helper()
	base := t.TempDir()
	dist := filepath.Join(base, "dist")
	files := map[string]string{
		"index.html":           indexHTML,
		"assets/index-abc.js":  "console.log(1)",
		"assets/index-abc.css": "body{}",
		"icon.png":             "\x89PNG\r\n\x1a\n",
		".env":                 "SECRET=1",
		"sub/dir/keep.txt":     "x",
	}
	for name, body := range files {
		p := filepath.Join(dist, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(base, "secret.txt"), []byte("TOP SECRET"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(base, "secret.txt"), filepath.Join(dist, "escape.txt")); err != nil {
		t.Fatal(err)
	}
	return dist
}

func newMux(dir string) *http.ServeMux {
	mux := http.NewServeMux()
	Register(mux, app.Deps{Config: config.Config{WebDir: dir}})
	return mux
}

func serve(h http.Handler, method, target string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func wantJSON404(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want JSON", ct)
	}
	var m map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil || m["message"] != "Not found" {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func wantIndex(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("Content-Type = %q", ct)
	}
	if cc := rec.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Fatalf("Cache-Control = %q, want no-cache", cc)
	}
	if !strings.Contains(rec.Body.String(), `<div id="root">`) {
		t.Fatalf("body is not index.html: %q", rec.Body.String())
	}
}

func TestSPAIndexFallback(t *testing.T) {
	mux := newMux(newWebDir(t))
	for _, p := range []string{"/", "/simulation", "/article/abc", "/admin/dashboard", "/index.html", "/nope.png"} {
		t.Run(p, func(t *testing.T) { wantIndex(t, serve(mux, "GET", p)) })
	}
}

func TestSPARealFilesAreImmutable(t *testing.T) {
	mux := newMux(newWebDir(t))
	cases := map[string]string{
		"/assets/index-abc.js":  "javascript",
		"/assets/index-abc.css": "text/css",
		"/icon.png":             "image/png",
	}
	for p, ctype := range cases {
		for _, method := range []string{"GET", "HEAD"} {
			rec := serve(mux, method, p)
			if rec.Code != http.StatusOK {
				t.Fatalf("%s %s: status = %d", method, p, rec.Code)
			}
			if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, ctype) {
				t.Fatalf("%s %s: Content-Type = %q, want %s", method, p, ct, ctype)
			}
			if cc := rec.Header().Get("Cache-Control"); cc != "public, max-age=31536000, immutable" {
				t.Fatalf("%s %s: Cache-Control = %q", method, p, cc)
			}
			if method == "HEAD" && rec.Body.Len() != 0 {
				t.Fatalf("HEAD %s returned a body", p)
			}
		}
	}
}

func TestSPAMissingAssetIs404(t *testing.T) {
	mux := newMux(newWebDir(t))
	for _, p := range []string{"/assets/index-old.js", "/assets/", "/assets", "/assets/.hidden"} {
		for _, method := range []string{"GET", "HEAD"} {
			rec := serve(mux, method, p)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("%s %s: status = %d, want 404", method, p, rec.Code)
			}
		}
	}
}

func TestSPATraversalRefused(t *testing.T) {
	spa := NewSPA(newWebDir(t))
	// Call the handler directly so ServeMux's own path cleaning cannot help.
	for _, p := range []string{"/../secret.txt", "/assets/../../secret.txt", "/escape.txt", "/.env", "/%2e%2e/secret.txt"} {
		req := httptest.NewRequest("GET", "/", nil)
		req.URL.Path = p
		if strings.Contains(p, "%") {
			req = httptest.NewRequest("GET", p, nil)
		}
		rec := httptest.NewRecorder()
		spa.ServeHTTP(rec, req)
		if strings.Contains(rec.Body.String(), "TOP SECRET") || strings.Contains(rec.Body.String(), "SECRET=1") {
			t.Fatalf("%s leaked a file outside the SPA: %q", p, rec.Body.String())
		}
		wantIndex(t, rec)
	}
}

func TestSPANoDirectoryListing(t *testing.T) {
	mux := newMux(newWebDir(t))
	for _, p := range []string{"/sub", "/sub/", "/sub/dir/"} {
		rec := serve(mux, "GET", p)
		if strings.Contains(rec.Body.String(), "keep.txt") {
			t.Fatalf("%s listed a directory: %q", p, rec.Body.String())
		}
		wantIndex(t, rec)
	}
	if rec := serve(mux, "GET", "/sub/dir/keep.txt"); rec.Code != http.StatusOK || rec.Body.String() != "x" {
		t.Fatalf("real nested file: %d %q", rec.Code, rec.Body.String())
	}
}

func TestSPANonGETIsJSON404(t *testing.T) {
	mux := newMux(newWebDir(t))
	for _, method := range []string{"POST", "PUT", "PATCH", "DELETE", "OPTIONS"} {
		for _, p := range []string{"/", "/assets/index-abc.js", "/simulation"} {
			wantJSON404(t, serve(mux, method, p))
		}
	}
}

func TestSPAWithoutIndexIs404(t *testing.T) {
	empty := t.TempDir()
	wantJSON404(t, serve(newMux(empty), "GET", "/simulation"))
	wantJSON404(t, serve(newMux(filepath.Join(empty, "missing")), "GET", "/"))
}

func TestSecurityHeadersOnEveryResponse(t *testing.T) {
	want := map[string]string{
		"Content-Security-Policy":           "default-src 'self';base-uri 'self';font-src 'self' https://fonts.gstatic.com;form-action 'self';frame-ancestors 'self';img-src 'self' data: https:;object-src 'none';script-src 'self';script-src-attr 'none';style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;upgrade-insecure-requests;connect-src 'self'",
		"Cross-Origin-Opener-Policy":        "same-origin",
		"Cross-Origin-Resource-Policy":      "same-origin",
		"Origin-Agent-Cluster":              "?1",
		"Referrer-Policy":                   "no-referrer",
		"Strict-Transport-Security":         "max-age=31536000; includeSubDomains",
		"X-Content-Type-Options":            "nosniff",
		"X-DNS-Prefetch-Control":            "off",
		"X-Download-Options":                "noopen",
		"X-Frame-Options":                   "SAMEORIGIN",
		"X-Permitted-Cross-Domain-Policies": "none",
		"X-XSS-Protection":                  "0",
	}
	h := SecurityHeaders(newMux(newWebDir(t)))
	for _, tc := range []struct{ method, path string }{
		{"GET", "/"}, {"GET", "/assets/index-abc.js"}, {"GET", "/assets/missing.js"}, {"POST", "/"},
	} {
		rec := serve(h, tc.method, tc.path)
		for k, v := range want {
			if got := rec.Header().Get(k); got != v {
				t.Fatalf("%s %s: %s = %q, want %q", tc.method, tc.path, k, got, v)
			}
		}
		if rec.Header().Get("X-Powered-By") != "" {
			t.Fatalf("X-Powered-By must not be sent")
		}
	}
}

func TestCORS(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	const good = "http://localhost:5173"
	h := CORS([]string{good, "https://other.example"})(next)

	do := func(method, origin string, extra map[string]string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(method, "/api/articles", nil)
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		for k, v := range extra {
			req.Header.Set(k, v)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec
	}

	t.Run("allowed origin", func(t *testing.T) {
		rec := do("GET", good, nil)
		if rec.Code != http.StatusTeapot {
			t.Fatalf("status = %d, next not called", rec.Code)
		}
		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != good {
			t.Fatalf("ACAO = %q", got)
		}
		if got := rec.Header().Values("Vary"); len(got) == 0 || got[0] != "Origin" {
			t.Fatalf("Vary = %q", got)
		}
		if rec.Header().Get("Access-Control-Allow-Methods") != "" {
			t.Fatalf("non-preflight must not send Allow-Methods")
		}
	})

	t.Run("denied origin and no origin", func(t *testing.T) {
		for _, o := range []string{"https://evil.example", ""} {
			rec := do("GET", o, nil)
			if rec.Code != http.StatusTeapot {
				t.Fatalf("status = %d", rec.Code)
			}
			for k := range rec.Header() {
				if strings.HasPrefix(k, "Access-Control-") || k == "Vary" {
					t.Fatalf("origin %q got CORS header %s", o, k)
				}
			}
		}
		// A preflight from an unlisted origin is not answered by CORS either.
		if rec := do("OPTIONS", "https://evil.example", map[string]string{"Access-Control-Request-Method": "POST"}); rec.Code != http.StatusTeapot || rec.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Fatalf("unlisted preflight: %d %v", rec.Code, rec.Header())
		}
	})

	t.Run("preflight", func(t *testing.T) {
		rec := do("OPTIONS", good, map[string]string{
			"Access-Control-Request-Method":  "PUT",
			"Access-Control-Request-Headers": "authorization,content-type",
		})
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want 204", rec.Code)
		}
		hd := rec.Header()
		if hd.Get("Access-Control-Allow-Origin") != good ||
			hd.Get("Access-Control-Allow-Methods") != "GET,HEAD,PUT,PATCH,POST,DELETE" ||
			hd.Get("Access-Control-Allow-Headers") != "authorization,content-type" {
			t.Fatalf("preflight headers = %v", hd)
		}
		if rec.Body.Len() != 0 {
			t.Fatalf("preflight body = %q", rec.Body.String())
		}
	})

	t.Run("plain OPTIONS is not a preflight", func(t *testing.T) {
		rec := do("OPTIONS", good, nil)
		if rec.Code != http.StatusTeapot || rec.Header().Get("Access-Control-Allow-Origin") != good {
			t.Fatalf("plain OPTIONS: %d %v", rec.Code, rec.Header())
		}
	})

	t.Run("empty list adds nothing", func(t *testing.T) {
		h := CORS(nil)(next)
		req := httptest.NewRequest("OPTIONS", "/", nil)
		req.Header.Set("Origin", good)
		req.Header.Set("Access-Control-Request-Method", "GET")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusTeapot || rec.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Fatalf("empty list: %d %v", rec.Code, rec.Header())
		}
	})
}
