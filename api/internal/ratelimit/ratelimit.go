// Package ratelimit provides an in-memory fixed-window limiter keyed by client IP.
//
// TODO(go/auth): implement per docs/SPEC.md §3.1 (login: 5 requests / 15 minutes / IP).
package ratelimit
