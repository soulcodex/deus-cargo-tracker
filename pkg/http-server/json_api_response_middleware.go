package httpserver

import (
	"context"
	"net/http"

	"github.com/google/jsonapi"

	"github.com/soulcodex/deus-cargo-tracker/pkg/logger"
)

type JSONAPIResponseMiddleware struct {
	logger logger.ZerologLogger
}

func NewJSONAPIResponseMiddleware(logger logger.ZerologLogger) *JSONAPIResponseMiddleware {
	return &JSONAPIResponseMiddleware{logger: logger}
}

func (jrm *JSONAPIResponseMiddleware) WriteErrorResponse(
	ctx context.Context,
	writer http.ResponseWriter,
	errors []*jsonapi.ErrorObject,
	statusCode int,
	previous error,
) {
	writer.Header().Set("Content-Type", jsonapi.MediaType)
	writer.WriteHeader(statusCode)

	jrm.logError(ctx, previous, statusCode)

	if err := jsonapi.MarshalErrors(writer, errors); err != nil {
		jrm.logger.Error().
			Ctx(ctx).
			Err(err).
			Msg("unexpected error marshalling json api response error")
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (jrm *JSONAPIResponseMiddleware) WriteResponse(
	ctx context.Context,
	writer http.ResponseWriter,
	payload interface{},
	statusCode int,
) {
	writer.Header().Set("Content-Type", jsonapi.MediaType)
	writer.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	if err := jsonapi.MarshalPayload(writer, payload); err != nil {
		jrm.logger.Error().
			Ctx(ctx).
			Err(err).
			Msg("unexpected error marshalling json api response")
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (jrm *JSONAPIResponseMiddleware) logError(ctx context.Context, err error, statusCode int) {
	if err == nil {
		return
	}

	if statusCode >= http.StatusInternalServerError {
		jrm.logger.Error().
			Ctx(ctx).
			Err(err).
			Msg("unexpected error processing request")
		return
	}

	jrm.logger.Warn().
		Ctx(ctx).
		Err(err).
		Msg("unexpected error processing request")
}
