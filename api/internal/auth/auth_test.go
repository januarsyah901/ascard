package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secret = "0123456789abcdef0123456789abcdef"

func TestSignVerifyRoundTrip(t *testing.T) {
	j := NewJWT(secret)
	tok, err := j.Sign("admin-1")
	if err != nil {
		t.Fatal(err)
	}
	id, err := j.Verify(tok)
	if err != nil || id != "admin-1" {
		t.Fatalf("id=%q err=%v", id, err)
	}
}

// A token shaped like the one jsonwebtoken produced in the Express server.
func TestVerifyAcceptsExpressToken(t *testing.T) {
	now := time.Now()
	tok, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": "abc", "iat": now.Unix(), "exp": now.Add(time.Hour).Unix(),
	}).SignedString([]byte(secret))
	if id, err := NewJWT(secret).Verify(tok); err != nil || id != "abc" {
		t.Fatalf("id=%q err=%v", id, err)
	}
}

func TestVerifyRejects(t *testing.T) {
	now := time.Now()
	sign := func(m jwt.Claims, method jwt.SigningMethod, key any) string {
		s, err := jwt.NewWithClaims(method, m).SignedString(key)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	cases := map[string]string{
		"expired":      sign(jwt.MapClaims{"id": "a", "exp": now.Add(-time.Minute).Unix()}, jwt.SigningMethodHS256, []byte(secret)),
		"wrong secret": sign(jwt.MapClaims{"id": "a", "exp": now.Add(time.Hour).Unix()}, jwt.SigningMethodHS256, []byte(strings.Repeat("x", 32))),
		"other alg":    sign(jwt.MapClaims{"id": "a", "exp": now.Add(time.Hour).Unix()}, jwt.SigningMethodHS512, []byte(secret)),
		"no id":        sign(jwt.MapClaims{"exp": now.Add(time.Hour).Unix()}, jwt.SigningMethodHS256, []byte(secret)),
		"garbage":      "not.a.token",
	}
	for name, tok := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := NewJWT(secret).Verify(tok); err == nil {
				t.Fatal("expected rejection")
			}
		})
	}
}

func TestGuard(t *testing.T) {
	j := NewJWT(secret)
	good, _ := j.Sign("admin-1")
	var seen string
	h := j.Guard(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = AdminID(r.Context())
	}))

	cases := []struct {
		header, wantMsg string
		wantStatus      int
	}{
		{"", MsgNoToken, 401},
		{"Basic abc", MsgNoToken, 401},
		{"Bearer nope", MsgInvalidToken, 401},
		{"Bearer " + good, "", 200},
	}
	for _, c := range cases {
		seen = ""
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		if c.header != "" {
			r.Header.Set("Authorization", c.header)
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != c.wantStatus || !strings.Contains(w.Body.String(), c.wantMsg) {
			t.Fatalf("%q: status=%d body=%s", c.header, w.Code, w.Body)
		}
		if c.wantStatus == 200 && seen != "admin-1" {
			t.Fatalf("admin id not in context: %q", seen)
		}
	}
}
