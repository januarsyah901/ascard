//go:build contract

// Package contract is a black-box test suite for the public HTTP contract of the
// Bongkar Ilusi server. It talks to a running server over HTTP, so the same suite runs
// against the Express reference (server/index.js) and the Go rewrite (cmd/server).
// Both must pass it before the Go server may replace Express.
//
// Run it with scripts/contract.sh, or by hand:
//
//	BASE_URL=http://127.0.0.1:5001 DATABASE_URL=postgres://... JWT_SECRET=... \
//	  go test -tags=contract -count=1 ./contract/
//
// The suite empties the "Article" and "Admin" tables of DATABASE_URL: never point it at
// a database you care about.
package contract

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	adminUsername = "contract-admin"
	adminPassword = "contract-password-123"
	// Produced by bcryptjs (the Express dependency): bcrypt.hashSync(adminPassword, 10).
	// Logging in with it proves the Go server accepts the hashes already in production.
	adminHash = "$2b$10$Tsu.YZ6xMspZT6vhuKbjh.svQpWdBu5R8NPTspHZLTbio0EyykSZ."
)

var (
	baseURL   string
	jwtSecret string
	adminID   string
	client    = &http.Client{Timeout: 10 * time.Second}
	isoMillis = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`)
)

func TestMain(m *testing.M) {
	baseURL = strings.TrimRight(os.Getenv("BASE_URL"), "/")
	jwtSecret = os.Getenv("JWT_SECRET")
	dbURL := os.Getenv("DATABASE_URL")
	if baseURL == "" || jwtSecret == "" || dbURL == "" {
		fmt.Fprintln(os.Stderr, "contract: BASE_URL, DATABASE_URL and JWT_SECRET are required")
		os.Exit(2)
	}
	if err := resetDatabase(dbURL); err != nil {
		fmt.Fprintln(os.Stderr, "contract: reset database:", err)
		os.Exit(2)
	}
	os.Exit(m.Run())
}

func resetDatabase(url string) error {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, strings.Replace(url, "?schema=public", "", 1))
	if err != nil {
		return err
	}
	defer pool.Close()
	if _, err := pool.Exec(ctx, `TRUNCATE "Article", "Admin"`); err != nil {
		return err
	}
	return pool.QueryRow(ctx,
		`INSERT INTO "Admin" ("id", "username", "password") VALUES (gen_random_uuid()::text, $1, $2) RETURNING "id"`,
		adminUsername, adminHash).Scan(&adminID)
}

// ---------- helpers ----------

type response struct {
	Status int
	Header http.Header
	Body   []byte
}

func (r response) json(t *testing.T) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(r.Body, &m); err != nil {
		t.Fatalf("body is not a JSON object: %v\n%s", err, r.Body)
	}
	return m
}

func (r response) list(t *testing.T) []map[string]any {
	t.Helper()
	var l []map[string]any
	if err := json.Unmarshal(r.Body, &l); err != nil {
		t.Fatalf("body is not a JSON array: %v\n%s", err, r.Body)
	}
	return l
}

type req struct {
	method, path string
	body         any    // marshalled as JSON when not a string
	rawBody      string // sent as-is when set
	contentType  string
	token        string
	headers      map[string]string
}

func do(t *testing.T, q req) response {
	t.Helper()
	var body io.Reader
	ctype := q.contentType
	switch {
	case q.rawBody != "":
		body = strings.NewReader(q.rawBody)
	case q.body != nil:
		b, _ := json.Marshal(q.body)
		body = bytes.NewReader(b)
		if ctype == "" {
			ctype = "application/json"
		}
	}
	r, err := http.NewRequest(q.method, baseURL+q.path, body)
	if err != nil {
		t.Fatal(err)
	}
	if ctype != "" {
		r.Header.Set("Content-Type", ctype)
	}
	if q.token != "" {
		r.Header.Set("Authorization", "Bearer "+q.token)
	}
	// Every request gets its own client IP so the login rate limit never leaks between tests.
	r.Header.Set("X-Forwarded-For", randomIP())
	for k, v := range q.headers {
		r.Header.Set(k, v)
	}
	res, err := client.Do(r)
	if err != nil {
		t.Fatalf("%s %s: %v", q.method, q.path, err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return response{Status: res.StatusCode, Header: res.Header, Body: b}
}

func randomIP() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return fmt.Sprintf("10.%d.%d.%d", b[0], b[1], b[2]|1)
}

func wantStatus(t *testing.T, r response, status int) {
	t.Helper()
	if r.Status != status {
		t.Fatalf("status = %d, want %d\nbody: %s", r.Status, status, r.Body)
	}
}

func wantMessage(t *testing.T, r response, status int, msg string) {
	t.Helper()
	wantStatus(t, r, status)
	if got := r.json(t)["message"]; got != msg {
		t.Fatalf("message = %q, want %q", got, msg)
	}
}

func wantJSONContentType(t *testing.T, r response) {
	t.Helper()
	if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
}

func login(t *testing.T) string {
	t.Helper()
	r := do(t, req{method: "POST", path: "/api/auth/login",
		body: map[string]string{"username": adminUsername, "password": adminPassword}})
	wantStatus(t, r, 200)
	tok, _ := r.json(t)["token"].(string)
	if tok == "" {
		t.Fatalf("no token in %s", r.Body)
	}
	return tok
}

func signToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func uniqueTitle(prefix string) string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return prefix + " " + hex.EncodeToString(b)
}

func createArticle(t *testing.T, tok string, body map[string]any) map[string]any {
	t.Helper()
	r := do(t, req{method: "POST", path: "/api/articles", token: tok, body: body})
	wantStatus(t, r, 201)
	return r.json(t)
}

var articleKeys = []string{"id", "title", "content", "category", "published", "createdAt", "updatedAt"}

func checkArticleShape(t *testing.T, a map[string]any) {
	t.Helper()
	if len(a) != len(articleKeys) {
		t.Fatalf("article has keys %v, want exactly %v", keys(a), articleKeys)
	}
	for _, k := range articleKeys {
		if _, ok := a[k]; !ok {
			t.Fatalf("article missing %q: %v", k, a)
		}
	}
	for _, k := range []string{"createdAt", "updatedAt"} {
		if s, _ := a[k].(string); !isoMillis.MatchString(s) {
			t.Fatalf("%s = %v, want JS ISO string like 2026-11-11T09:21:00.123Z", k, a[k])
		}
	}
	if _, ok := a["published"].(bool); !ok {
		t.Fatalf("published must be a boolean: %v", a["published"])
	}
}

func keys(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// ---------- platform ----------

func TestHealth(t *testing.T) {
	r := do(t, req{method: "GET", path: "/api/health"})
	wantStatus(t, r, 200)
	wantJSONContentType(t, r)
	m := r.json(t)
	if m["status"] != "ok" || m["database"] != "ok" {
		t.Fatalf("health = %v", m)
	}
}

func TestUnknownAPIPathIsJSON404(t *testing.T) {
	for _, p := range []string{"/api/nope", "/api/auth/seed", "/api/articles/x/y/z"} {
		r := do(t, req{method: "GET", path: p})
		wantMessage(t, r, 404, "Not found")
	}
}

func TestSeedEndpointIsGone(t *testing.T) {
	r := do(t, req{method: "POST", path: "/api/auth/seed",
		body: map[string]string{"username": "evil", "password": "evil-password-123"}})
	wantStatus(t, r, 404)
}

func TestSPAServesIndexForClientRoutes(t *testing.T) {
	for _, p := range []string{"/", "/simulation", "/knowledge-base", "/article/abc", "/admin/dashboard"} {
		r := do(t, req{method: "GET", path: p})
		wantStatus(t, r, 200)
		if ct := r.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
			t.Fatalf("%s: Content-Type = %q", p, ct)
		}
		if !bytes.Contains(r.Body, []byte(`<div id="root">`)) {
			t.Fatalf("%s: not the SPA index.html", p)
		}
		if cc := r.Header.Get("Cache-Control"); cc != "no-cache" {
			t.Fatalf("%s: index.html Cache-Control = %q, want no-cache", p, cc)
		}
	}
}

func TestSPAStaticAssets(t *testing.T) {
	index := do(t, req{method: "GET", path: "/"})
	m := regexp.MustCompile(`src="(/assets/[^"]+\.js)"`).FindSubmatch(index.Body)
	if m == nil {
		t.Fatalf("no /assets/*.js script in index.html")
	}
	r := do(t, req{method: "GET", path: string(m[1])})
	wantStatus(t, r, 200)
	if ct := r.Header.Get("Content-Type"); !strings.Contains(ct, "javascript") {
		t.Fatalf("asset Content-Type = %q", ct)
	}
	cc := r.Header.Get("Cache-Control")
	if !strings.Contains(cc, "max-age=31536000") || !strings.Contains(cc, "immutable") {
		t.Fatalf("asset Cache-Control = %q, want long-lived immutable", cc)
	}
}

func TestSecurityHeaders(t *testing.T) {
	for _, p := range []string{"/", "/api/health"} {
		r := do(t, req{method: "GET", path: p})
		want := map[string]string{
			"X-Content-Type-Options":     "nosniff",
			"X-Frame-Options":            "SAMEORIGIN",
			"Referrer-Policy":            "no-referrer",
			"Strict-Transport-Security":  "max-age=31536000; includeSubDomains",
			"Cross-Origin-Opener-Policy": "same-origin",
		}
		for h, v := range want {
			if got := r.Header.Get(h); got != v {
				t.Fatalf("%s: %s = %q, want %q", p, h, got, v)
			}
		}
		csp := r.Header.Get("Content-Security-Policy")
		for _, directive := range []string{
			"default-src 'self'",
			"script-src 'self'",
			"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
			"font-src 'self' https://fonts.gstatic.com",
			"img-src 'self' data: https:",
			"connect-src 'self'",
			"object-src 'none'",
			"frame-ancestors 'self'",
		} {
			if !strings.Contains(csp, directive) {
				t.Fatalf("%s: CSP %q lacks %q", p, csp, directive)
			}
		}
		if r.Header.Get("X-Powered-By") != "" {
			t.Fatalf("%s: X-Powered-By must not be sent", p)
		}
	}
}

func TestNoCORSForUnlistedOrigin(t *testing.T) {
	r := do(t, req{method: "GET", path: "/api/articles", headers: map[string]string{"Origin": "https://evil.example"}})
	if v := r.Header.Get("Access-Control-Allow-Origin"); v != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want none", v)
	}
}

// ---------- auth ----------

func TestLoginSuccess(t *testing.T) {
	r := do(t, req{method: "POST", path: "/api/auth/login",
		body: map[string]string{"username": adminUsername, "password": adminPassword}})
	wantStatus(t, r, 200)
	wantJSONContentType(t, r)
	m := r.json(t)
	admin, _ := m["admin"].(map[string]any)
	if admin["id"] != adminID || admin["username"] != adminUsername || len(admin) != 2 {
		t.Fatalf("admin = %v", m["admin"])
	}
	tok, _ := m["token"].(string)
	parsed, err := jwt.Parse(tok, func(*jwt.Token) (any, error) { return []byte(jwtSecret), nil },
		jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		t.Fatalf("token does not verify with JWT_SECRET as HS256: %v", err)
	}
	claims := parsed.Claims.(jwt.MapClaims)
	if claims["id"] != adminID {
		t.Fatalf("token id claim = %v", claims["id"])
	}
	exp, _ := claims.GetExpirationTime()
	if exp == nil || exp.Sub(time.Now()) < 23*time.Hour || exp.Sub(time.Now()) > 25*time.Hour {
		t.Fatalf("token must expire in ~1 day, exp = %v", exp)
	}
}

func TestLoginFailures(t *testing.T) {
	cases := []struct {
		name   string
		q      req
		status int
		msg    string
	}{
		{"wrong password", req{body: map[string]string{"username": adminUsername, "password": "wrong-password"}}, 401, "Invalid credentials"},
		{"unknown user", req{body: map[string]string{"username": "nobody", "password": adminPassword}}, 401, "Invalid credentials"},
		{"missing password", req{body: map[string]string{"username": adminUsername}}, 400, "Username and password are required"},
		{"non-string username", req{body: map[string]any{"username": 42, "password": adminPassword}}, 400, "Username and password are required"},
		{"empty object", req{body: map[string]any{}}, 400, "Username and password are required"},
		{"not JSON content type", req{rawBody: "username=a&password=b", contentType: "application/x-www-form-urlencoded"}, 400, "Username and password are required"},
		{"malformed JSON", req{rawBody: "{bad", contentType: "application/json"}, 400, "Invalid JSON"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.q.method, c.q.path = "POST", "/api/auth/login"
			wantMessage(t, do(t, c.q), c.status, c.msg)
		})
	}
}

func TestLoginRateLimit(t *testing.T) {
	ip := randomIP()
	body := map[string]string{"username": adminUsername, "password": "wrong-password"}
	for i := 1; i <= 5; i++ {
		r := do(t, req{method: "POST", path: "/api/auth/login", body: body, headers: map[string]string{"X-Forwarded-For": ip}})
		wantStatus(t, r, 401)
	}
	r := do(t, req{method: "POST", path: "/api/auth/login", body: body, headers: map[string]string{"X-Forwarded-For": ip}})
	wantMessage(t, r, 429, "Too many login attempts. Try again in 15 minutes.")

	// Even correct credentials are refused for that IP until the window ends...
	r = do(t, req{method: "POST", path: "/api/auth/login", headers: map[string]string{"X-Forwarded-For": ip},
		body: map[string]string{"username": adminUsername, "password": adminPassword}})
	wantStatus(t, r, 429)

	// ...while other clients are unaffected.
	login(t)
}

func TestMe(t *testing.T) {
	tok := login(t)
	r := do(t, req{method: "GET", path: "/api/auth/me", token: tok})
	wantStatus(t, r, 200)
	admin, _ := r.json(t)["admin"].(map[string]any)
	if admin["id"] != adminID || admin["username"] != adminUsername || len(admin) != 2 {
		t.Fatalf("me = %s", r.Body)
	}
}

// A token minted the way jsonwebtoken did in Express must keep working after the cutover.
func TestExpressStyleTokenAccepted(t *testing.T) {
	now := time.Now()
	tok := signToken(t, jwtSecret, jwt.MapClaims{"id": adminID, "iat": now.Unix(), "exp": now.Add(24 * time.Hour).Unix()})
	wantStatus(t, do(t, req{method: "GET", path: "/api/auth/me", token: tok}), 200)
}

func TestProtectedRoutesRejectBadTokens(t *testing.T) {
	now := time.Now()
	tokens := map[string]string{
		"garbage":              "not-a-jwt",
		"wrong secret":         signToken(t, strings.Repeat("x", 40), jwt.MapClaims{"id": adminID, "exp": now.Add(time.Hour).Unix()}),
		"expired":              signToken(t, jwtSecret, jwt.MapClaims{"id": adminID, "exp": now.Add(-time.Minute).Unix()}),
		"admin does not exist": signToken(t, jwtSecret, jwt.MapClaims{"id": "00000000-0000-0000-0000-000000000000", "exp": now.Add(time.Hour).Unix()}),
	}
	for name, tok := range tokens {
		t.Run(name, func(t *testing.T) {
			wantMessage(t, do(t, req{method: "GET", path: "/api/auth/me", token: tok}), 401, "Invalid token")
			if name != "admin does not exist" { // only /me looks the admin up
				wantMessage(t, do(t, req{method: "GET", path: "/api/articles/admin/all", token: tok}), 401, "Invalid token")
			}
		})
	}
	routes := []req{
		{method: "GET", path: "/api/auth/me"},
		{method: "GET", path: "/api/articles/admin/all"},
		{method: "POST", path: "/api/articles", body: map[string]any{"title": "x", "content": "y"}},
		{method: "PUT", path: "/api/articles/some-id", body: map[string]any{"title": "x", "content": "y"}},
		{method: "DELETE", path: "/api/articles/some-id"},
	}
	for _, q := range routes {
		wantMessage(t, do(t, q), 401, "No token provided")
		q.headers = map[string]string{"Authorization": "Basic abc"}
		wantMessage(t, do(t, q), 401, "No token provided")
	}
}

// ---------- articles ----------

func TestArticleLifecycle(t *testing.T) {
	tok := login(t)
	title := uniqueTitle("Contract Published")

	created := createArticle(t, tok, map[string]any{
		"title": "  " + title + "  ", "content": "Isi artikel", "category": "  Matematika  ", "published": true,
	})
	checkArticleShape(t, created)
	if created["title"] != title || created["category"] != "Matematika" || created["published"] != true {
		t.Fatalf("create must trim title/category: %v", created)
	}
	id := created["id"].(string)
	if !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`).MatchString(id) {
		t.Fatalf("id %q is not a UUID", id)
	}

	// public read
	r := do(t, req{method: "GET", path: "/api/articles/" + id})
	wantStatus(t, r, 200)
	got := r.json(t)
	checkArticleShape(t, got)
	if got["id"] != id || got["content"] != "Isi artikel" {
		t.Fatalf("get = %v", got)
	}

	// update replaces every field
	r = do(t, req{method: "PUT", path: "/api/articles/" + id, token: tok,
		body: map[string]any{"title": title + " v2", "content": "Baru", "category": "", "published": "true"}})
	wantStatus(t, r, 200)
	upd := r.json(t)
	checkArticleShape(t, upd)
	if upd["title"] != title+" v2" || upd["content"] != "Baru" || upd["category"] != nil || upd["published"] != false {
		t.Fatalf("update: empty category must become null and published must be exactly true: %v", upd)
	}
	if upd["createdAt"] != created["createdAt"] {
		t.Fatalf("update must keep createdAt")
	}
	if upd["updatedAt"].(string) < created["updatedAt"].(string) {
		t.Fatalf("updatedAt must not go backwards")
	}

	// now a draft: hidden from the public
	wantMessage(t, do(t, req{method: "GET", path: "/api/articles/" + id}), 404, "Article not found")

	// delete
	r = do(t, req{method: "DELETE", path: "/api/articles/" + id, token: tok})
	wantMessage(t, r, 200, "Article deleted successfully")
	wantMessage(t, do(t, req{method: "DELETE", path: "/api/articles/" + id, token: tok}), 404, "Article not found")
	wantMessage(t, do(t, req{method: "PUT", path: "/api/articles/" + id, token: tok,
		body: map[string]any{"title": "x", "content": "y"}}), 404, "Article not found")
}

func TestPublicListHidesDraftsAndSortsNewestFirst(t *testing.T) {
	tok := login(t)
	older := createArticle(t, tok, map[string]any{"title": uniqueTitle("Older"), "content": "a", "published": true})
	time.Sleep(20 * time.Millisecond)
	draft := createArticle(t, tok, map[string]any{"title": uniqueTitle("Draft"), "content": "b", "published": false})
	time.Sleep(20 * time.Millisecond)
	newer := createArticle(t, tok, map[string]any{"title": uniqueTitle("Newer"), "content": "c", "published": true})

	r := do(t, req{method: "GET", path: "/api/articles"})
	wantStatus(t, r, 200)
	wantJSONContentType(t, r)
	list := r.list(t)
	pos := map[string]int{}
	for i, a := range list {
		checkArticleShape(t, a)
		if a["published"] != true {
			t.Fatalf("public list contains a draft: %v", a)
		}
		pos[a["id"].(string)] = i
	}
	if _, ok := pos[draft["id"].(string)]; ok {
		t.Fatal("draft must not be listed publicly")
	}
	if pos[newer["id"].(string)] > pos[older["id"].(string)] {
		t.Fatal("list must be ordered by createdAt descending")
	}
	wantMessage(t, do(t, req{method: "GET", path: "/api/articles/" + draft["id"].(string)}), 404, "Article not found")

	r = do(t, req{method: "GET", path: "/api/articles/admin/all", token: tok})
	wantStatus(t, r, 200)
	found := false
	for _, a := range r.list(t) {
		checkArticleShape(t, a)
		if a["id"] == draft["id"] {
			found = true
		}
	}
	if !found {
		t.Fatal("admin list must include drafts")
	}
}

func TestEmptyListIsArray(t *testing.T) {
	r := do(t, req{method: "GET", path: "/api/articles"})
	wantStatus(t, r, 200)
	if b := bytes.TrimSpace(r.Body); len(b) == 0 || b[0] != '[' {
		t.Fatalf("list must be a JSON array, got %s", r.Body)
	}
}

func TestArticleValidation(t *testing.T) {
	tok := login(t)
	const msg = "Title and content are required"
	cases := map[string]map[string]any{
		"missing title":      {"content": "x"},
		"blank title":        {"title": "   ", "content": "x"},
		"numeric title":      {"title": 5, "content": "x"},
		"title too long":     {"title": strings.Repeat("t", 201), "content": "x"},
		"missing content":    {"title": "x"},
		"blank content":      {"title": "x", "content": " \n "},
		"content too long":   {"title": "x", "content": strings.Repeat("c", 100001)},
		"content not string": {"title": "x", "content": []string{"a"}},
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			wantMessage(t, do(t, req{method: "POST", path: "/api/articles", token: tok, body: body}), 400, msg)
			wantMessage(t, do(t, req{method: "PUT", path: "/api/articles/any-id", token: tok, body: body}), 400, msg)
		})
	}
	wantMessage(t, do(t, req{method: "POST", path: "/api/articles", token: tok,
		rawBody: "{bad", contentType: "application/json"}), 400, "Invalid JSON")
}

func TestArticleLimitsAndDefaults(t *testing.T) {
	tok := login(t)
	title := uniqueTitle(strings.Repeat("T", 180))
	if len(title) > 200 {
		title = title[:200]
	}
	a := createArticle(t, tok, map[string]any{
		"title": title, "content": "x", "category": strings.Repeat("k", 60),
	})
	if a["category"] != strings.Repeat("k", 50) {
		t.Fatalf("category must be cut to 50 characters, got %v", a["category"])
	}
	if a["published"] != false {
		t.Fatalf("published must default to false, got %v", a["published"])
	}
	b := createArticle(t, tok, map[string]any{"title": uniqueTitle("No category"), "content": "x", "category": 7})
	if b["category"] != nil {
		t.Fatalf("non-string category must become null, got %v", b["category"])
	}
}

func TestDuplicateTitleConflict(t *testing.T) {
	tok := login(t)
	title := uniqueTitle("Duplicate")
	createArticle(t, tok, map[string]any{"title": title, "content": "x"})
	wantMessage(t, do(t, req{method: "POST", path: "/api/articles", token: tok,
		body: map[string]any{"title": title, "content": "y"}}), 409, "An article with this title already exists")

	other := createArticle(t, tok, map[string]any{"title": uniqueTitle("Other"), "content": "x"})
	wantMessage(t, do(t, req{method: "PUT", path: "/api/articles/" + other["id"].(string), token: tok,
		body: map[string]any{"title": title, "content": "y"}}), 409, "An article with this title already exists")
}

func TestUnknownArticle(t *testing.T) {
	wantMessage(t, do(t, req{method: "GET", path: "/api/articles/does-not-exist"}), 404, "Article not found")
	wantMessage(t, do(t, req{method: "GET", path: "/api/articles/00000000-0000-0000-0000-000000000000"}), 404, "Article not found")
}
