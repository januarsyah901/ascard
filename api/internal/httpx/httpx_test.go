package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"
	"time"
)

func TestTimeMatchesJavaScriptISOString(t *testing.T) {
	ts := Time{time.Date(2026, 11, 11, 9, 21, 0, 123456789, time.FixedZone("WIB", 7*3600))}
	b, err := json.Marshal(ts)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(b), `"2026-11-11T02:21:00.123Z"`; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestDecodeObject(t *testing.T) {
	cases := []struct {
		name, ctype, body string
		wantOK            bool
		wantStatus        int
		wantKeys          int
	}{
		{"object", "application/json", `{"a":"b"}`, true, 0, 1},
		{"charset", "application/json; charset=utf-8", `{"a":"b"}`, true, 0, 1},
		{"not json content type is ignored", "text/plain", `{"a":"b"}`, true, 0, 0},
		{"empty body", "application/json", ``, true, 0, 0},
		{"array becomes empty", "application/json", `[1,2]`, true, 0, 0},
		{"malformed", "application/json", `{bad`, false, 400, 0},
		{"too large", "application/json", `{"a":"` + strings.Repeat("x", MaxJSONBody) + `"}`, false, 413, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(c.body))
			r.Header.Set("Content-Type", c.ctype)
			w := httptest.NewRecorder()
			f, ok := DecodeObject(w, r)
			if ok != c.wantOK {
				t.Fatalf("ok = %v, want %v", ok, c.wantOK)
			}
			if !ok && w.Code != c.wantStatus {
				t.Fatalf("status = %d, want %d", w.Code, c.wantStatus)
			}
			if ok && len(f) != c.wantKeys {
				t.Fatalf("keys = %d, want %d", len(f), c.wantKeys)
			}
		})
	}
}

func TestFields(t *testing.T) {
	f := Fields{}
	_ = json.Unmarshal([]byte(`{"s":"x","n":1,"t":true,"ts":"true","z":null}`), &f)
	if _, ok := f.String("z"); ok {
		t.Fatal("null must not be a string")
	}
	if s, ok := f.String("s"); !ok || s != "x" {
		t.Fatal("String(s)")
	}
	if _, ok := f.String("n"); ok {
		t.Fatal("number must not be a string")
	}
	if !f.IsTrue("t") || f.IsTrue("ts") || f.IsTrue("missing") {
		t.Fatal("IsTrue must only accept the literal true")
	}
}

func TestClientIP(t *testing.T) {
	c := ClientIP{Trusted: []netip.Prefix{netip.MustParsePrefix("172.16.0.0/12")}}
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.RemoteAddr = "172.18.0.5:1234"
	r.Header.Set("X-Forwarded-For", "1.1.1.1, 9.9.9.9")
	if got := c.Resolve(r); got != "9.9.9.9" {
		t.Fatalf("trusted peer: got %s", got)
	}

	r.RemoteAddr = "8.8.8.8:1234"
	if got := c.Resolve(r); got != "8.8.8.8" {
		t.Fatalf("untrusted peer must ignore XFF: got %s", got)
	}
}
