package main

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func env(m map[string]string) func(string) string { return func(k string) string { return m[k] } }

func TestLoadSettings(t *testing.T) {
	cases := []struct {
		name    string
		env     map[string]string
		wantErr string
		user    string
	}{
		{"missing password", map[string]string{"DATABASE_URL": "postgres://x"}, "ADMIN_PASSWORD", ""},
		{"short password", map[string]string{"DATABASE_URL": "postgres://x", "ADMIN_PASSWORD": "elevenchars"}, "ADMIN_PASSWORD", ""},
		{"missing database url", map[string]string{"ADMIN_PASSWORD": "twelve-chars"}, "DATABASE_URL", ""},
		{"default username", map[string]string{"DATABASE_URL": "postgres://x", "ADMIN_PASSWORD": "twelve-chars"}, "", "admin"},
		{"custom username", map[string]string{"DATABASE_URL": "postgres://x", "ADMIN_PASSWORD": "twelve-chars", "ADMIN_USERNAME": "root"}, "", "root"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := loadSettings(env(c.env))
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("err = %v, want mention of %s", err, c.wantErr)
				}
				return
			}
			if err != nil || s.username != c.user {
				t.Fatalf("settings = %+v, err = %v", s, err)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	h, err := hashPassword("correct horse battery")
	if err != nil {
		t.Fatal(err)
	}
	if cost, _ := bcrypt.Cost([]byte(h)); cost != bcryptCost {
		t.Fatalf("cost = %d", cost)
	}
	if bcrypt.CompareHashAndPassword([]byte(h), []byte("correct horse battery")) != nil {
		t.Fatal("hash does not verify")
	}
	long := strings.Repeat("a", 100)
	h, err = hashPassword(long)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(h), []byte(long)) != nil {
		t.Fatalf("long password: %v", err)
	}
}
