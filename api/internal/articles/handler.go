package articles

import (
	"errors"
	"net/http"

	"github.com/januarsyah901/ascard/api/internal/db"
	"github.com/januarsyah901/ascard/api/internal/httpx"
)

// Messages identical to server/routes/articles.js.
const (
	MsgNotFound  = "Article not found"
	MsgRequired  = "Title and content are required"
	MsgDuplicate = "An article with this title already exists"
	MsgDeleted   = "Article deleted successfully"
)

// Handler serves the article endpoints.
type Handler struct{ Repo Repo }

func (h Handler) listPublished(w http.ResponseWriter, r *http.Request) { h.list(w, r, true) }
func (h Handler) listAll(w http.ResponseWriter, r *http.Request)       { h.list(w, r, false) }

func (h Handler) list(w http.ResponseWriter, r *http.Request, publishedOnly bool) {
	list, err := h.Repo.List(r.Context(), publishedOnly)
	if err != nil {
		httpx.ServerError(w, r, "list articles", err)
		return
	}
	if list == nil {
		list = []Article{}
	}
	httpx.JSON(w, http.StatusOK, list)
}

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	a, err := h.Repo.GetPublished(r.Context(), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, r, "get article", err)
		return
	}
	httpx.JSON(w, http.StatusOK, a)
}

func (h Handler) create(w http.ResponseWriter, r *http.Request) {
	in, ok := decode(w, r)
	if !ok {
		return
	}
	a, err := h.Repo.Create(r.Context(), in)
	if err != nil {
		h.writeErr(w, r, "create article", err)
		return
	}
	httpx.JSON(w, http.StatusCreated, a)
}

func (h Handler) update(w http.ResponseWriter, r *http.Request) {
	in, ok := decode(w, r)
	if !ok {
		return
	}
	a, err := h.Repo.Update(r.Context(), r.PathValue("id"), in)
	if err != nil {
		h.writeErr(w, r, "update article", err)
		return
	}
	httpx.JSON(w, http.StatusOK, a)
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request) {
	if err := h.Repo.Delete(r.Context(), r.PathValue("id")); err != nil {
		h.writeErr(w, r, "delete article", err)
		return
	}
	httpx.Error(w, http.StatusOK, MsgDeleted)
}

// decode reads and validates the body; it has written the response when ok is false.
func decode(w http.ResponseWriter, r *http.Request) (Input, bool) {
	fields, ok := httpx.DecodeObject(w, r)
	if !ok {
		return Input{}, false
	}
	in, ok := validate(fields)
	if !ok {
		httpx.Error(w, http.StatusBadRequest, MsgRequired)
	}
	return in, ok
}

func (h Handler) writeErr(w http.ResponseWriter, r *http.Request, label string, err error) {
	switch {
	case errors.Is(err, db.ErrNotFound):
		httpx.Error(w, http.StatusNotFound, MsgNotFound)
	case db.IsUniqueViolation(err):
		httpx.Error(w, http.StatusConflict, MsgDuplicate)
	default:
		httpx.ServerError(w, r, label, err)
	}
}
