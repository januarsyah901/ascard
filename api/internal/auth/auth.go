// Package auth signs and verifies admin JWTs and provides the guard middleware.
//
// Tokens are compatible with the Express version (jsonwebtoken): HS256, payload
// {"id": "<admin uuid>", "iat": ..., "exp": ...}, valid for one day. A token issued by
// the Express server stays valid after the cutover and vice versa.
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/januarsyah901/ascard/api/internal/httpx"
)

// TokenTTL matches jwt.sign(..., { expiresIn: '1d' }).
const TokenTTL = 24 * time.Hour

// Messages identical to the Express middleware.
const (
	MsgNoToken      = "No token provided"
	MsgInvalidToken = "Invalid token"
)

// JWT signs and verifies admin tokens.
type JWT struct {
	secret []byte
	now    func() time.Time
}

// NewJWT returns a signer/verifier for the given secret.
func NewJWT(secret string) *JWT {
	return &JWT{secret: []byte(secret), now: time.Now}
}

type claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

// Sign returns a token for the admin id.
func (j *JWT) Sign(adminID string) (string, error) {
	now := j.now()
	c := claims{
		ID: adminID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenTTL)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(j.secret)
}

// ErrInvalidToken is returned for any token that must be rejected.
var ErrInvalidToken = errors.New("invalid token")

// Verify returns the admin id carried by a valid token.
func (j *JWT) Verify(token string) (string, error) {
	var c claims
	_, err := jwt.ParseWithClaims(token, &c, func(*jwt.Token) (any, error) { return j.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithTimeFunc(j.now),
	)
	if err != nil || c.ID == "" {
		return "", ErrInvalidToken
	}
	return c.ID, nil
}

type ctxKey struct{}

// AdminID returns the id of the authenticated admin, set by Guard.
func AdminID(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}

// WithAdminID stores an admin id in ctx. Handler tests use it to fake an authenticated request.
func WithAdminID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// Guard returns middleware that only lets requests with a valid Bearer token through,
// answering 401 {"message":"No token provided"} or {"message":"Invalid token"} otherwise.
func (j *JWT) Guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			httpx.Error(w, http.StatusUnauthorized, MsgNoToken)
			return
		}
		// Express used authHeader.split(' ')[1]: the token is the second space-separated part.
		token := strings.Split(header, " ")[1]
		id, err := j.Verify(token)
		if err != nil {
			httpx.Error(w, http.StatusUnauthorized, MsgInvalidToken)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithAdminID(r.Context(), id)))
	})
}
