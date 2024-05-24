package base

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/httplog/v2"
	"github.com/unexpectedtoken/recipes/view/layout"
)

type BaseHandler struct {
}

func New() *BaseHandler {
	return &BaseHandler{}
}

func (h *BaseHandler) LoggerFromRContext(r *http.Request) slog.Logger {
	return httplog.LogEntry(r.Context())
}

func (h *BaseHandler) HandleClientError(w http.ResponseWriter, r *http.Request, err error) {
	logger := h.LoggerFromRContext(r)

	logger.Error("client error: %s", err)

	w.WriteHeader(http.StatusBadRequest)
}

func (h *BaseHandler) HandleServerError(w http.ResponseWriter, r *http.Request, err error) {
	logger := h.LoggerFromRContext(r)

	logger.Error("server error: %s", err)

	w.WriteHeader(http.StatusInternalServerError)
}

func (h *BaseHandler) RenderHTMXWithLayout(w http.ResponseWriter, r *http.Request, comp templ.Component) {
	err := layout.Base(comp).Render(r.Context(), w)

	if err != nil {
		h.HandleServerError(w, r, err)
	}
}

func (h *BaseHandler) RenderHTMX(w http.ResponseWriter, r *http.Request, comp templ.Component) {
	err := comp.Render(r.Context(), w)

	if err != nil {
		h.HandleServerError(w, r, err)
	}
}
