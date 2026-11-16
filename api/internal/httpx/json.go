// Package httpx holds HTTP helpers shared by every module: JSON I/O, error responses,
// timestamps and middleware. Error bodies always have the shape {"message": "..."} with
// the exact English messages of the Express version.
package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"time"
)

// Shared messages, identical to the Express version.
const (
	MsgServerError  = "Internal server error"
	MsgNotFound     = "Not found"
	MsgInvalidJSON  = "Invalid JSON"
	MsgBodyTooLarge = "Request body too large"
)

// MaxJSONBody is the largest JSON request body accepted (1 MB, like express.json({limit:'1mb'})).
const MaxJSONBody = 1 << 20

// Message is the standard error/info body.
type Message struct {
	Message string `json:"message"`
}

// JSON writes v as JSON with the given status.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("write json response", "err", err)
	}
}

// Error writes {"message": msg} with the given status.
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, Message{Message: msg})
}

// ServerError logs err with a label and writes a 500 with the standard message.
func ServerError(w http.ResponseWriter, r *http.Request, label string, err error) {
	slog.ErrorContext(r.Context(), label, "err", err, "request_id", RequestID(r.Context()))
	Error(w, http.StatusInternalServerError, MsgServerError)
}

// Fields is a decoded JSON object whose values keep their raw JSON so handlers can check
// each field's type the way the Express code did with typeof.
type Fields map[string]json.RawMessage

// String returns the field as a string; ok is false when it is missing or not a JSON string.
func (f Fields) String(key string) (string, bool) {
	raw, present := f[key]
	// json.Unmarshal of null into a string is a silent no-op, but typeof null !== 'string'.
	if !present || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return "", false
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false
	}
	return s, true
}

// IsTrue reports whether the field is exactly the JSON literal true (JS `=== true`).
func (f Fields) IsTrue(key string) bool {
	var b bool
	raw, present := f[key]
	return present && json.Unmarshal(raw, &b) == nil && b
}

// DecodeObject reads a JSON request body into Fields, mirroring express.json():
//   - a body that is not declared as application/json is ignored (empty Fields, ok);
//   - an empty body or a JSON value that is not an object gives empty Fields;
//   - malformed JSON writes 400 {"message":"Invalid JSON"} and returns ok=false;
//   - a body over 1 MB writes 413 and returns ok=false.
func DecodeObject(w http.ResponseWriter, r *http.Request) (Fields, bool) {
	fields := Fields{}
	mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mt != "application/json" {
		return fields, true
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, MaxJSONBody))
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			Error(w, http.StatusRequestEntityTooLarge, MsgBodyTooLarge)
			return nil, false
		}
		Error(w, http.StatusBadRequest, MsgInvalidJSON)
		return nil, false
	}
	if len(body) == 0 {
		return fields, true
	}
	var anyVal any
	if err := json.Unmarshal(body, &anyVal); err != nil {
		Error(w, http.StatusBadRequest, MsgInvalidJSON)
		return nil, false
	}
	if _, isObj := anyVal.(map[string]any); !isObj {
		return fields, true
	}
	if err := json.Unmarshal(body, &fields); err != nil {
		Error(w, http.StatusBadRequest, MsgInvalidJSON)
		return nil, false
	}
	return fields, true
}

// Time serialises like a JavaScript Date in JSON: UTC with millisecond precision,
// e.g. "2026-11-11T09:21:00.123Z". Use it for every timestamp in a response.
type Time struct{ time.Time }

// JSTimeLayout is the layout of JavaScript's Date.prototype.toISOString.
const JSTimeLayout = "2006-01-02T15:04:05.000Z"

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UTC().Format(JSTimeLayout))
}
