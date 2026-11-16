package articles

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/januarsyah901/ascard/api/internal/app"
	"github.com/januarsyah901/ascard/api/internal/auth"
	"github.com/januarsyah901/ascard/api/internal/db"
)

// fakeRepo is an in-memory Repo with the same semantics as pgRepo.
type fakeRepo struct {
	mu    sync.Mutex
	rows  map[string]Article
	seq   int
	clock time.Time
	fail  error // returned by every call when set
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{rows: map[string]Article{}, clock: time.Date(2026, 11, 11, 9, 21, 0, 123_000_000, time.UTC)}
}

func (f *fakeRepo) tick() time.Time { f.clock = f.clock.Add(time.Second); return f.clock }

func (f *fakeRepo) titleTaken(title, except string) bool {
	for id, a := range f.rows {
		if a.Title == title && id != except {
			return true
		}
	}
	return false
}

var errUnique = &pgconn.PgError{Code: "23505"}

func (f *fakeRepo) List(_ context.Context, publishedOnly bool) ([]Article, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fail != nil {
		return nil, f.fail
	}
	var list []Article // nil on purpose: the handler must still answer []
	for _, a := range f.rows {
		if !publishedOnly || a.Published {
			list = append(list, a)
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt.Time) })
	return list, nil
}

func (f *fakeRepo) GetPublished(_ context.Context, id string) (Article, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fail != nil {
		return Article{}, f.fail
	}
	a, ok := f.rows[id]
	if !ok || !a.Published {
		return Article{}, db.ErrNotFound
	}
	return a, nil
}

func (f *fakeRepo) Create(_ context.Context, in Input) (Article, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fail != nil {
		return Article{}, f.fail
	}
	if f.titleTaken(in.Title, "") {
		return Article{}, errUnique
	}
	f.seq++
	now := f.tick()
	a := Article{ID: fmt.Sprintf("00000000-0000-0000-0000-%012d", f.seq), Title: in.Title, Content: in.Content,
		Category: in.Category, Published: in.Published}
	a.CreatedAt.Time, a.UpdatedAt.Time = now, now
	f.rows[a.ID] = a
	return a, nil
}

func (f *fakeRepo) Update(_ context.Context, id string, in Input) (Article, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fail != nil {
		return Article{}, f.fail
	}
	a, ok := f.rows[id]
	if !ok {
		return Article{}, db.ErrNotFound
	}
	if f.titleTaken(in.Title, id) {
		return Article{}, errUnique
	}
	a.Title, a.Content, a.Category, a.Published = in.Title, in.Content, in.Category, in.Published
	a.UpdatedAt.Time = f.tick()
	f.rows[id] = a
	return a, nil
}

func (f *fakeRepo) Delete(_ context.Context, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.fail != nil {
		return f.fail
	}
	if _, ok := f.rows[id]; !ok {
		return db.ErrNotFound
	}
	delete(f.rows, id)
	return nil
}

const testSecret = "handler-test-secret-0123456789abcdef0123"

type env struct {
	t    *testing.T
	srv  http.Handler
	repo *fakeRepo
	tok  string
}

func newEnv(t *testing.T) *env {
	t.Helper()
	jwt := auth.NewJWT(testSecret)
	tok, err := jwt.Sign("admin-1")
	if err != nil {
		t.Fatal(err)
	}
	repo := newFakeRepo()
	mux := http.NewServeMux()
	Mount(mux, app.Deps{JWT: jwt}, Handler{Repo: repo})
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot) // marks "fell through to the foundation's 404"
	})
	return &env{t: t, srv: mux, repo: repo, tok: tok}
}

func (e *env) do(method, path, body string, authed bool) *httptest.ResponseRecorder {
	e.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if authed {
		req.Header.Set("Authorization", "Bearer "+e.tok)
	}
	rec := httptest.NewRecorder()
	e.srv.ServeHTTP(rec, req)
	return rec
}

func (e *env) create(body string) map[string]any {
	e.t.Helper()
	rec := e.do("POST", "/api/articles", body, true)
	if rec.Code != http.StatusCreated {
		e.t.Fatalf("create: %d %s", rec.Code, rec.Body)
	}
	return decodeObj(e.t, rec)
}

func decodeObj(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &m); err != nil {
		t.Fatalf("not a JSON object: %s", rec.Body)
	}
	return m
}

func wantMsg(t *testing.T, rec *httptest.ResponseRecorder, status int, msg string) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status = %d, want %d (%s)", rec.Code, status, rec.Body)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != fmt.Sprintf(`{"message":%q}`, msg) {
		t.Fatalf("body = %s, want message %q", got, msg)
	}
}

var timeRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`)

// TestArticleJSONShape checks the exact key order and timestamp format.
func TestArticleJSONShape(t *testing.T) {
	e := newEnv(t)
	rec := e.do("POST", "/api/articles", `{"title":"A","content":"B","published":true}`, true)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content type %q", ct)
	}
	want := regexp.MustCompile(`^\{"id":"[^"]+","title":"A","content":"B","category":null,"published":true,` +
		`"createdAt":"2026-11-11T09:21:01\.123Z","updatedAt":"2026-11-11T09:21:01\.123Z"\}$`)
	if body := strings.TrimSpace(rec.Body.String()); !want.MatchString(body) {
		t.Fatalf("body = %s", body)
	}
	m := decodeObj(t, rec)
	if len(m) != 7 {
		t.Fatalf("want 7 keys, got %v", m)
	}
	for _, k := range []string{"createdAt", "updatedAt"} {
		if !timeRe.MatchString(m[k].(string)) {
			t.Fatalf("%s = %v", k, m[k])
		}
	}
}

func TestTimeIsUTC(t *testing.T) {
	a := Article{Category: ptr("c")}
	a.CreatedAt.Time = time.Date(2026, 1, 2, 10, 0, 0, 5_000_000, time.FixedZone("WIB", 7*3600))
	b, _ := json.Marshal(a)
	if !strings.Contains(string(b), `"createdAt":"2026-01-02T03:00:00.005Z"`) || !strings.Contains(string(b), `"category":"c"`) {
		t.Fatalf("got %s", b)
	}
}

func TestListNeverNull(t *testing.T) {
	e := newEnv(t)
	for _, path := range []string{"/api/articles", "/api/articles/"} {
		rec := e.do("GET", path, "", false)
		if rec.Code != 200 || strings.TrimSpace(rec.Body.String()) != "[]" {
			t.Fatalf("%s: %d %s", path, rec.Code, rec.Body)
		}
	}
	rec := e.do("GET", "/api/articles/admin/all", "", true)
	if rec.Code != 200 || strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("admin: %d %s", rec.Code, rec.Body)
	}
}

func TestListsAndAdminRoute(t *testing.T) {
	e := newEnv(t)
	older := e.create(`{"title":"older","content":"x","published":true}`)
	draft := e.create(`{"title":"draft","content":"x"}`)
	newer := e.create(`{"title":"newer","content":"x","published":true}`)

	ids := func(rec *httptest.ResponseRecorder) []any {
		t.Helper()
		if rec.Code != 200 {
			t.Fatalf("status %d %s", rec.Code, rec.Body)
		}
		var list []map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
			t.Fatal(err)
		}
		var out []any
		for _, a := range list {
			out = append(out, a["id"])
		}
		return out
	}
	if got := fmt.Sprint(ids(e.do("GET", "/api/articles", "", false))); got != fmt.Sprint([]any{newer["id"], older["id"]}) {
		t.Fatalf("public list = %s", got)
	}
	// admin/all wins over {id} and includes drafts
	if got := fmt.Sprint(ids(e.do("GET", "/api/articles/admin/all", "", true))); got != fmt.Sprint([]any{newer["id"], draft["id"], older["id"]}) {
		t.Fatalf("admin list = %s", got)
	}
	wantMsg(t, e.do("GET", "/api/articles/admin/all", "", false), 401, auth.MsgNoToken)
}

func TestGet(t *testing.T) {
	e := newEnv(t)
	pub := e.create(`{"title":"pub","content":"x","published":true}`)
	draft := e.create(`{"title":"draft","content":"x"}`)

	rec := e.do("GET", "/api/articles/"+pub["id"].(string), "", false)
	if rec.Code != 200 || decodeObj(t, rec)["id"] != pub["id"] {
		t.Fatalf("get published: %d %s", rec.Code, rec.Body)
	}
	wantMsg(t, e.do("GET", "/api/articles/"+draft["id"].(string), "", false), 404, MsgNotFound)
	wantMsg(t, e.do("GET", "/api/articles/does-not-exist", "", false), 404, MsgNotFound)
}

func TestCreate(t *testing.T) {
	e := newEnv(t)
	a := e.create(`{"title":"  T  ","content":" c ","category":"  Cat  ","published":"true"}`)
	if a["title"] != "T" || a["content"] != " c " || a["category"] != "Cat" || a["published"] != false {
		t.Fatalf("created = %v", a)
	}
	wantMsg(t, e.do("POST", "/api/articles", `{"title":"T","content":"y"}`, true), 409, MsgDuplicate)
	wantMsg(t, e.do("POST", "/api/articles", `{"title":" ","content":"y"}`, true), 400, MsgRequired)
	wantMsg(t, e.do("POST", "/api/articles", `{bad`, true), 400, "Invalid JSON")
	wantMsg(t, e.do("POST", "/api/articles", `{"title":"a","content":"b"}`, false), 401, auth.MsgNoToken)
	if len(e.repo.rows) != 1 {
		t.Fatalf("rows = %d", len(e.repo.rows))
	}
}

func TestUpdate(t *testing.T) {
	e := newEnv(t)
	a := e.create(`{"title":"a","content":"x","category":"c","published":true}`)
	e.create(`{"title":"b","content":"x"}`)
	path := "/api/articles/" + a["id"].(string)

	rec := e.do("PUT", path, `{"title":" a2 ","content":"y"}`, true)
	if rec.Code != 200 {
		t.Fatalf("update: %d %s", rec.Code, rec.Body)
	}
	u := decodeObj(t, rec)
	if u["title"] != "a2" || u["content"] != "y" || u["category"] != nil || u["published"] != false {
		t.Fatalf("update must replace every field: %v", u)
	}
	if u["createdAt"] != a["createdAt"] || u["updatedAt"] == a["updatedAt"] {
		t.Fatalf("timestamps: before %v, after %v", a, u)
	}
	wantMsg(t, e.do("PUT", path, `{"title":"b","content":"y"}`, true), 409, MsgDuplicate)
	wantMsg(t, e.do("PUT", "/api/articles/missing", `{"title":"z","content":"y"}`, true), 404, MsgNotFound)
	// validation happens before any DB access, so a missing id with a bad body is 400
	wantMsg(t, e.do("PUT", "/api/articles/missing", `{"title":"z"}`, true), 400, MsgRequired)
	wantMsg(t, e.do("PUT", path, `{"title":"z","content":"y"}`, false), 401, auth.MsgNoToken)
}

func TestDelete(t *testing.T) {
	e := newEnv(t)
	a := e.create(`{"title":"a","content":"x"}`)
	path := "/api/articles/" + a["id"].(string)
	wantMsg(t, e.do("DELETE", path, "", false), 401, auth.MsgNoToken)
	wantMsg(t, e.do("DELETE", path, "", true), 200, MsgDeleted)
	wantMsg(t, e.do("DELETE", path, "", true), 404, MsgNotFound)
}

func TestRepoErrorsAre500(t *testing.T) {
	e := newEnv(t)
	e.repo.fail = errors.New("boom")
	cases := []struct {
		method, path, body string
		authed             bool
	}{
		{"GET", "/api/articles", "", false},
		{"GET", "/api/articles/admin/all", "", true},
		{"GET", "/api/articles/x", "", false},
		{"POST", "/api/articles", `{"title":"a","content":"b"}`, true},
		{"PUT", "/api/articles/x", `{"title":"a","content":"b"}`, true},
		{"DELETE", "/api/articles/x", "", true},
	}
	for _, c := range cases {
		wantMsg(t, e.do(c.method, c.path, c.body, c.authed), 500, "Internal server error")
	}
}

func TestUnknownPathsFallThrough(t *testing.T) {
	e := newEnv(t)
	for _, c := range [][2]string{
		{"GET", "/api/articles/a/b"},
		{"PATCH", "/api/articles/x"},
		{"POST", "/api/articles/admin/all"},
		{"PUT", "/api/articles"},
	} {
		if rec := e.do(c[0], c[1], "", true); rec.Code != http.StatusTeapot {
			t.Fatalf("%s %s = %d, want fallthrough to /api/ 404", c[0], c[1], rec.Code)
		}
	}
}
