package authapi

import (
	"bytes"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/januarsyah901/ascard/api/internal/auth"
	"github.com/januarsyah901/ascard/api/internal/db"
	"github.com/januarsyah901/ascard/api/internal/httpx"
)

// Messages identical to server/routes/auth.js.
const (
	MsgCredentialsRequired = "Username and password are required"
	MsgInvalidCredentials  = "Invalid credentials"
	MsgTooManyAttempts     = "Too many login attempts. Try again in 15 minutes."
)

// Signer issues tokens for an admin id (auth.JWT in production).
type Signer interface {
	Sign(adminID string) (string, error)
}

// Handler serves the auth endpoints.
type Handler struct {
	Repo   Repo
	Tokens Signer
}

// adminJSON is the public view of an admin: {"id","username"}.
type adminJSON struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type loginResponse struct {
	Token string    `json:"token"`
	Admin adminJSON `json:"admin"`
}

type meResponse struct {
	Admin adminJSON `json:"admin"`
}

// Login handles POST /api/auth/login.
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	fields, ok := httpx.DecodeObject(w, r)
	if !ok {
		return
	}
	username, okU := jsString(fields, "username")
	password, okP := jsString(fields, "password")
	if !okU || !okP {
		httpx.Error(w, http.StatusBadRequest, MsgCredentialsRequired)
		return
	}

	admin, err := h.Repo.AdminByUsername(r.Context(), username)
	if errors.Is(err, db.ErrNotFound) {
		httpx.Error(w, http.StatusUnauthorized, MsgInvalidCredentials)
		return
	}
	if err != nil {
		httpx.ServerError(w, r, "login: load admin", err)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		httpx.Error(w, http.StatusUnauthorized, MsgInvalidCredentials)
		return
	}

	token, err := h.Tokens.Sign(admin.ID)
	if err != nil {
		httpx.ServerError(w, r, "login: sign token", err)
		return
	}
	httpx.JSON(w, http.StatusOK, loginResponse{Token: token, Admin: adminJSON{ID: admin.ID, Username: admin.Username}})
}

// jsString is fields.String that also rejects JSON null, matching `typeof v === 'string'`.
// (httpx.Fields.String accepts null as "" because json.Unmarshal of null into a string is a no-op.)
func jsString(fields httpx.Fields, key string) (string, bool) {
	if string(bytes.TrimSpace(fields[key])) == "null" {
		return "", false
	}
	return fields.String(key)
}

// Me handles GET /api/auth/me (behind JWT.Guard).
func (h Handler) Me(w http.ResponseWriter, r *http.Request) {
	admin, err := h.Repo.AdminByID(r.Context(), auth.AdminID(r.Context()))
	if errors.Is(err, db.ErrNotFound) {
		httpx.Error(w, http.StatusUnauthorized, auth.MsgInvalidToken)
		return
	}
	if err != nil {
		httpx.ServerError(w, r, "me: load admin", err)
		return
	}
	httpx.JSON(w, http.StatusOK, meResponse{Admin: adminJSON{ID: admin.ID, Username: admin.Username}})
}
