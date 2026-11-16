// Package ratelimit provides an in-memory fixed-window limiter keyed by client IP.
//
// It mirrors express-rate-limit's default MemoryStore: the first hit from a key opens a
// window of fixed length; every hit inside it counts; once the window ends the count
// starts over. Expired windows are swept periodically so the map cannot grow without bound.
package ratelimit

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/januarsyah901/ascard/api/internal/httpx"
)

// Limiter counts hits per key in fixed windows.
type Limiter struct {
	limit  int
	window time.Duration
	now    func() time.Time

	mu        sync.Mutex
	hits      map[string]*entry
	nextSweep time.Time
}

type entry struct {
	count int
	reset time.Time // end of the window
}

// Option configures a Limiter.
type Option func(*Limiter)

// WithClock replaces time.Now (tests).
func WithClock(now func() time.Time) Option {
	return func(l *Limiter) { l.now = now }
}

// New returns a limiter allowing limit hits per key in each window.
func New(limit int, window time.Duration, opts ...Option) *Limiter {
	l := &Limiter{limit: limit, window: window, now: time.Now, hits: map[string]*entry{}}
	for _, o := range opts {
		o(l)
	}
	l.nextSweep = l.now().Add(window)
	return l
}

// Result describes the state of a key after a hit.
type Result struct {
	Allowed   bool
	Limit     int
	Remaining int
	Reset     time.Time // when the current window ends
}

// Allow records a hit for key and reports whether it is within the limit.
func (l *Limiter) Allow(key string) Result {
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()

	if !now.Before(l.nextSweep) {
		l.sweepLocked(now)
		l.nextSweep = now.Add(l.window)
	}

	e, ok := l.hits[key]
	if !ok || !now.Before(e.reset) {
		e = &entry{reset: now.Add(l.window)}
		l.hits[key] = e
	}
	e.count++
	return Result{
		Allowed:   e.count <= l.limit,
		Limit:     l.limit,
		Remaining: max(l.limit-e.count, 0),
		Reset:     e.reset,
	}
}

// Len returns the number of tracked keys (tests and metrics).
func (l *Limiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.hits)
}

// Sweep drops every expired window now.
func (l *Limiter) Sweep() {
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sweepLocked(now)
}

func (l *Limiter) sweepLocked(now time.Time) {
	for k, e := range l.hits {
		if !now.Before(e.reset) {
			delete(l.hits, k)
		}
	}
}

// Middleware limits next per key(r). Over the limit it answers status 429 with
// {"message": msg}. It sets the draft-7 RateLimit / RateLimit-Policy headers that
// express-rate-limit sent (standardHeaders: 'draft-7') and Retry-After on 429.
func (l *Limiter) Middleware(key func(*http.Request) string, msg string, next http.Handler) http.Handler {
	policy := fmt.Sprintf("%d;w=%d", l.limit, int(l.window/time.Second))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := l.Allow(key(r))
		resetSecs := int(math.Ceil(res.Reset.Sub(l.now()).Seconds()))
		if resetSecs < 0 {
			resetSecs = 0
		}
		h := w.Header()
		h.Set("RateLimit-Policy", policy)
		h.Set("RateLimit", fmt.Sprintf("limit=%d, remaining=%d, reset=%d", res.Limit, res.Remaining, resetSecs))
		if !res.Allowed {
			h.Set("Retry-After", strconv.Itoa(resetSecs))
			httpx.Error(w, http.StatusTooManyRequests, msg)
			return
		}
		next.ServeHTTP(w, r)
	})
}
