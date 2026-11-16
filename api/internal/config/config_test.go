package config

import (
	"strings"
	"testing"
)

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := load(env(map[string]string{
		"DATABASE_URL": "postgres://x",
		"JWT_SECRET":   strings.Repeat("s", 32),
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 5001 || cfg.WebDir != "dist" || cfg.LogFormat != "text" {
		t.Fatalf("bad defaults: %+v", cfg)
	}
	if len(cfg.CORSOrigins) != 0 {
		t.Fatalf("CORS should default to none, got %v", cfg.CORSOrigins)
	}
	if len(cfg.TrustedProxies) == 0 {
		t.Fatal("trusted proxies should have defaults")
	}
}

func TestLoadRejectsUnsafeConfig(t *testing.T) {
	cases := map[string]map[string]string{
		"missing everything": {},
		"short secret":       {"DATABASE_URL": "postgres://x", "JWT_SECRET": "short"},
		"bad port":           {"DATABASE_URL": "postgres://x", "JWT_SECRET": strings.Repeat("s", 32), "PORT": "abc"},
		"bad proxy":          {"DATABASE_URL": "postgres://x", "JWT_SECRET": strings.Repeat("s", 32), "TRUSTED_PROXIES": "nope"},
	}
	for name, m := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := load(env(m)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestSplitList(t *testing.T) {
	got := splitList(" https://a.test , ,https://b.test")
	if len(got) != 2 || got[0] != "https://a.test" || got[1] != "https://b.test" {
		t.Fatalf("got %v", got)
	}
}
