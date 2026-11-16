package authapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/januarsyah901/ascard/api/internal/app"
	"github.com/januarsyah901/ascard/api/internal/auth"
	"github.com/januarsyah901/ascard/api/internal/db"
	"github.com/januarsyah901/ascard/api/internal/ratelimit"
)

const (
	testSecret   = "test-jwt-secret-0123456789abcdef0123456789"
	testUser     = "contract-admin"
	testPassword = "contract-password-123"
	// bcryptjs output (bcrypt.hashSync(testPassword, 10)), as stored by the Express version.
	testHash = "$2b$10$Tsu.YZ6xMspZT6vhuKbjh.svQpWdBu5R8NPTspHZLTbio0EyykSZ."
	testID   = "8d0f5f8e-5d7a-4c1c-9f63-2f0a4b8d1e11"
)

type fakeRepo struct {
	admins map[string]Admin // by id
	err    error
}

func (f *fakeRepo) AdminByUsername(_ context.Context, username string) (Admin, error) {
	if f.err != nil {
		return Admin{}, f.err
	}
	for _, a := range f.admins {
		if a.Username == username {
			return a, nil
		}
	}
	return Admin{}, db.ErrNotFound
}

func (f *fakeRepo) AdminByID(_ context.Context, id string) (Admin, error) {
	if f.err != nil {
		return Admin{}, f.err
	}
	a, ok := f.admins[id]
	if !ok {
		return Admin{}, db.ErrNotFound
	}
	return a, nil
}

type fixture struct {
	mux  *http.ServeMux
	repo *fakeRepo
	jwt  *auth.JWT
}

func newFixture(limit int) fixture {
	repo := &fakeRepo{admins: map[string]Admin{testID: {ID: testID, Username: testUser, PasswordHash: testHash}}}
	j := auth.NewJWT(testSecret)
	mux := http.NewServeMux()
	d := app.Deps{JWT: j}
	mount(mux, d, Handler{Repo: repo, Tokens: j}, ratelimit.New(limit, LoginWindow))
	return fixture{mux: mux, repo: repo, jwt: j}
}

func (f fixture) do(method, path, ctype, body, token, remote string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if ctype != "" {
		r.Header.Set("Content-Type", ctype)
	}
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	if remote != "" {
		r.RemoteAddr = remote
	}
	w := httptest.NewRecorder()
	f.mux.ServeHTTP(w, r)
	return w
}

func decode(t *testing.T, w *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("body not JSON: %v: %s", err, w.Body)
	}
	return m
}

func wantMessage(t *testing.T, w *httptest.ResponseRecorder, status int, msg string) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("status = %d, want %d (%s)", w.Code, status, w.Body)
	}
	if got := decode(t, w)["message"]; got != msg {
		t.Fatalf("message = %v, want %q", got, msg)
	}
}

func TestLoginValidationAndCredentials(t *testing.T) {
	cases := []struct {
		name, ctype, body string
		status            int
		msg               string
	}{
		{"missing password", "application/json", `{"username":"contract-admin"}`, 400, MsgCredentialsRequired},
		{"missing username", "application/json", `{"password":"x"}`, 400, MsgCredentialsRequired},
		{"non-string username", "application/json", `{"username":42,"password":"x"}`, 400, MsgCredentialsRequired},
		{"null password", "application/json", `{"username":"a","password":null}`, 400, MsgCredentialsRequired},
		{"array body", "application/json", `["a","b"]`, 400, MsgCredentialsRequired},
		{"empty body", "application/json", ``, 400, MsgCredentialsRequired},
		{"form body", "application/x-www-form-urlencoded", `username=a&password=b`, 400, MsgCredentialsRequired},
		{"malformed JSON", "application/json", `{bad`, 400, "Invalid JSON"},
		{"unknown user", "application/json", `{"username":"nobody","password":"contract-password-123"}`, 401, MsgInvalidCredentials},
		{"wrong password", "application/json", `{"username":"contract-admin","password":"wrong"}`, 401, MsgInvalidCredentials},
		{"empty password", "application/json", `{"username":"contract-admin","password":""}`, 401, MsgInvalidCredentials},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := newFixture(LoginLimit)
			wantMessage(t, f.do("POST", "/api/auth/login", c.ctype, c.body, "", ""), c.status, c.msg)
		})
	}
}

func TestLoginSuccess(t *testing.T) {
	f := newFixture(LoginLimit)
	w := f.do("POST", "/api/auth/login", "application/json",
		`{"username":"contract-admin","password":"contract-password-123"}`, "", "")
	if w.Code != 200 {
		t.Fatalf("status = %d: %s", w.Code, w.Body)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content-type = %q", ct)
	}
	var body struct {
		Token string          `json:"token"`
		Admin json.RawMessage `json:"admin"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if want := `{"id":"` + testID + `","username":"contract-admin"}`; string(body.Admin) != want {
		t.Fatalf("admin = %s, want %s", body.Admin, want)
	}
	id, err := f.jwt.Verify(body.Token)
	if err != nil || id != testID {
		t.Fatalf("token: id=%q err=%v", id, err)
	}
}

func TestLoginRepoError(t *testing.T) {
	f := newFixture(LoginLimit)
	f.repo.err = errors.New("db down")
	wantMessage(t, f.do("POST", "/api/auth/login", "application/json", `{"username":"a","password":"b"}`, "", ""),
		500, "Internal server error")
}

func TestLoginRateLimitCountsEveryRequest(t *testing.T) {
	f := newFixture(LoginLimit)
	ip := "203.0.113.7:1234"
	// 400, 401, malformed JSON and a success all count towards the limit.
	bodies := []string{`{}`, `{"username":"x","password":"y"}`, `{bad`,
		`{"username":"contract-admin","password":"contract-password-123"}`, `{}`}
	for _, b := range bodies {
		if w := f.do("POST", "/api/auth/login", "application/json", b, "", ip); w.Code == 429 {
			t.Fatalf("limited too early: %s", b)
		}
	}
	wantMessage(t, f.do("POST", "/api/auth/login", "application/json",
		`{"username":"contract-admin","password":"contract-password-123"}`, "", ip), 429, MsgTooManyAttempts)
	if w := f.do("POST", "/api/auth/login", "application/json",
		`{"username":"contract-admin","password":"contract-password-123"}`, "", "203.0.113.8:1234"); w.Code != 200 {
		t.Fatalf("other IP: %d", w.Code)
	}
}

func TestMe(t *testing.T) {
	f := newFixture(LoginLimit)
	tok, err := f.jwt.Sign(testID)
	if err != nil {
		t.Fatal(err)
	}
	w := f.do("GET", "/api/auth/me", "", "", tok, "")
	if w.Code != 200 {
		t.Fatalf("status = %d: %s", w.Code, w.Body)
	}
	if want := `{"admin":{"id":"` + testID + `","username":"contract-admin"}}`; strings.TrimSpace(w.Body.String()) != want {
		t.Fatalf("body = %s", w.Body)
	}

	// Admin deleted after the token was issued.
	delete(f.repo.admins, testID)
	wantMessage(t, f.do("GET", "/api/auth/me", "", "", tok, ""), 401, auth.MsgInvalidToken)
}

func TestMeGuarded(t *testing.T) {
	f := newFixture(LoginLimit)
	wantMessage(t, f.do("GET", "/api/auth/me", "", "", "", ""), 401, auth.MsgNoToken)
	wantMessage(t, f.do("GET", "/api/auth/me", "", "", "garbage", ""), 401, auth.MsgInvalidToken)
}

func TestMeRepoError(t *testing.T) {
	f := newFixture(LoginLimit)
	tok, _ := f.jwt.Sign(testID)
	f.repo.err = errors.New("db down")
	wantMessage(t, f.do("GET", "/api/auth/me", "", "", tok, ""), 500, "Internal server error")
}
