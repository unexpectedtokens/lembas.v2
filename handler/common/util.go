package handler_util

import (
	"context"
	"io"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RenderPage(comp templ.Component, w io.Writer, ctx context.Context) error {

	return comp.Render(ctx, w)
}

func RenderComp(ctx context.Context, comp templ.Component, w http.ResponseWriter) {}

func GetLoggerFromReqContext(r *http.Request) slog.Logger {
	return httplog.LogEntry(r.Context())
}

func LogErrorWithMessage(r *http.Request, msg string, err error) {
	logger := GetLoggerFromReqContext(r)

	logger.ErrorContext(r.Context(), msg, "err", err)
}

func ObjectIDFromR(r *http.Request, name string) (primitive.ObjectID, error) {
	id := chi.URLParam(r, name)

	objectId, err := primitive.ObjectIDFromHex(id)

	return objectId, err
}

func HandleClientErr() {}

func HandleServerErr() {}

// TODO: Implement generic error handlers
