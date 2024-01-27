package logging

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

func getLoggerFromReqContext(r *http.Request) slog.Logger {
	return httplog.LogEntry(r.Context())
}

func LogRequestError(r *http.Request, err error) {
	logger := getLoggerFromReqContext(r)

	logger.ErrorContext(r.Context(), "err", err)
}
