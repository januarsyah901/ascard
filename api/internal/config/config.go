// Package config loads and validates runtime configuration from the environment.
package config

import (
	"errors"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
)

// Config is the validated runtime configuration.
type Config struct {
	Port           int
	DatabaseURL    string
	JWTSecret      string
	CORSOrigins    []string
	WebDir         string         // directory holding the built SPA (Vite dist/)
	TrustedProxies []netip.Prefix // peers whose X-Forwarded-For is believed
	LogFormat      string         // "text" or "json"
	AppEnv         string
}

// MinJWTSecretLen matches the Express server's startup check.
const MinJWTSecretLen = 32

// DefaultTrustedProxies covers loopback and private ranges: Caddy reaches the app over the
// Docker gateway-net (172.16.0.0/12) and local tooling uses loopback.
const DefaultTrustedProxies = "127.0.0.0/8,::1/128,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16"

// Load reads the environment. It fails on anything the server cannot run safely without.
func Load() (Config, error) {
	return load(os.Getenv)
}

func load(getenv func(string) string) (Config, error) {
	cfg := Config{
		DatabaseURL: getenv("DATABASE_URL"),
		JWTSecret:   getenv("JWT_SECRET"),
		WebDir:      or(getenv("WEB_DIR"), "dist"),
		LogFormat:   or(getenv("LOG_FORMAT"), "text"),
		AppEnv:      or(getenv("APP_ENV"), "development"),
	}

	var errs []error
	if cfg.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL is required"))
	}
	if cfg.JWTSecret == "" {
		errs = append(errs, errors.New("JWT_SECRET is required"))
	} else if len(cfg.JWTSecret) < MinJWTSecretLen {
		errs = append(errs, fmt.Errorf("JWT_SECRET must be at least %d characters long", MinJWTSecretLen))
	}

	port, err := strconv.Atoi(or(getenv("PORT"), "5001"))
	if err != nil || port < 1 || port > 65535 {
		errs = append(errs, fmt.Errorf("PORT must be a number between 1 and 65535"))
	}
	cfg.Port = port

	cfg.CORSOrigins = splitList(getenv("CORS_ORIGINS"))

	for _, s := range splitList(or(getenv("TRUSTED_PROXIES"), DefaultTrustedProxies)) {
		p, err := netip.ParsePrefix(s)
		if err != nil {
			errs = append(errs, fmt.Errorf("TRUSTED_PROXIES: invalid CIDR %q", s))
			continue
		}
		cfg.TrustedProxies = append(cfg.TrustedProxies, p)
	}

	if cfg.LogFormat != "text" && cfg.LogFormat != "json" {
		errs = append(errs, errors.New(`LOG_FORMAT must be "text" or "json"`))
	}

	return cfg, errors.Join(errs...)
}

func or(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func splitList(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, part)
		}
	}
	return out
}
