package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/logger"
)

type JSONResponseWriter struct {
	logger logger.ZerologLogger
}

func NewJSONResponseMiddleware(logger logger.ZerologLogger) *JSONResponseWriter {
	return &JSONResponseWriter{logger: logger}
}

func (jrm *JSONResponseWriter) WriteErrorResponse(
	ctx context.Context,
	writer http.ResponseWriter,
	errors []string,
	statusCode int,
) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	jsonBody, _ := json.Marshal(map[string]interface{}{
		"errors": errors,
	})

	if _, err := writer.Write(jsonBody); err != nil {
		jrm.logger.Error().
			Ctx(ctx).
			Err(err).
			Msg("unexpected error writing json response error")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (jrm *JSONResponseWriter) WriteResponse(
	ctx context.Context,
	writer http.ResponseWriter,
	payload interface{},
	statusCode int,
) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	jsonBody, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		jrm.logger.Error().
			Ctx(ctx).
			Err(marshalErr).
			Msg("unexpected error marshalling json response")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := writer.Write(jsonBody); err != nil {
		jrm.logger.Error().
			Ctx(ctx).
			Err(err).
			Msg("unexpected error writing json response error")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
