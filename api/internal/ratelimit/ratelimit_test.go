package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

func newClock() *fakeClock { return &fakeClock{t: time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)} }

func TestWindowLimit(t *testing.T) {
	clk := newClock()
	l := New(5, 15*time.Minute, WithClock(clk.Now))
	for i := 1; i <= 5; i++ {
		res := l.Allow("1.2.3.4")
		if !res.Allowed || res.Remaining != 5-i {
			t.Fatalf("hit %d: %+v", i, res)
		}
		clk.Advance(time.Minute)
	}
	if res := l.Allow("1.2.3.4"); res.Allowed || res.Remaining != 0 {
		t.Fatalf("6th hit allowed: %+v", res)
	}
}

func TestFixedWindowReset(t *testing.T) {
	clk := newClock()
	l := New(5, 15*time.Minute, WithClock(clk.Now))
	start := clk.Now()
	for range 6 {
		l.Allow("ip")
	}
	// Still blocked just before the window (opened by the first hit) ends.
	clk.Advance(15*time.Minute - time.Millisecond)
	if res := l.Allow("ip"); res.Allowed || !res.Reset.Equal(start.Add(15*time.Minute)) {
		t.Fatalf("expected block within window: %+v", res)
	}
	clk.Advance(time.Millisecond)
	res := l.Allow("ip")
	if !res.Allowed || res.Remaining != 4 || !res.Reset.Equal(clk.Now().Add(15*time.Minute)) {
		t.Fatalf("expected fresh window: %+v", res)
	}
}

func TestPerKeyIsolation(t *testing.T) {
	clk := newClock()
	l := New(2, time.Minute, WithClock(clk.Now))
	l.Allow("a")
	l.Allow("a")
	if l.Allow("a").Allowed {
		t.Fatal("a should be blocked")
	}
	if !l.Allow("b").Allowed {
		t.Fatal("b must not be affected by a")
	}
}

func TestCleanup(t *testing.T) {
	clk := newClock()
	l := New(5, 15*time.Minute, WithClock(clk.Now))
	for _, ip := range []string{"a", "b", "c"} {
		l.Allow(ip)
	}
	clk.Advance(10 * time.Minute)
	l.Allow("d")
	l.Sweep()
	if l.Len() != 4 {
		t.Fatalf("nothing expired yet, len = %d", l.Len())
	}
	clk.Advance(5 * time.Minute) // a, b, c expire; d has 10 minutes left
	l.Sweep()
	if l.Len() != 1 {
		t.Fatalf("len after sweep = %d, want 1", l.Len())
	}

	// Automatic sweep: once a window length has passed since the last sweep, Allow cleans up.
	l2 := New(5, time.Minute, WithClock(clk.Now))
	for i := range 100 {
		l2.Allow(strings.Repeat("x", i+1))
	}
	clk.Advance(time.Minute)
	l2.Allow("new")
	if l2.Len() != 1 {
		t.Fatalf("auto sweep: len = %d, want 1", l2.Len())
	}
}

func TestMiddleware(t *testing.T) {
	clk := newClock()
	l := New(2, 15*time.Minute, WithClock(clk.Now))
	h := l.Middleware(func(r *http.Request) string { return r.RemoteAddr }, "slow down",
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))

	hit := func(addr string) *httptest.ResponseRecorder {
		r := httptest.NewRequest("POST", "/", nil)
		r.RemoteAddr = addr
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}
	if w := hit("x"); w.Code != 204 || w.Header().Get("RateLimit") != "limit=2, remaining=1, reset=900" ||
		w.Header().Get("RateLimit-Policy") != "2;w=900" {
		t.Fatalf("first: %d %v", w.Code, w.Header())
	}
	hit("x")
	clk.Advance(time.Minute)
	w := hit("x")
	if w.Code != 429 || strings.TrimSpace(w.Body.String()) != `{"message":"slow down"}` || w.Header().Get("Retry-After") != "840" {
		t.Fatalf("third: %d %s %v", w.Code, w.Body, w.Header())
	}
	if w := hit("y"); w.Code != 204 {
		t.Fatalf("other key: %d", w.Code)
	}
}
