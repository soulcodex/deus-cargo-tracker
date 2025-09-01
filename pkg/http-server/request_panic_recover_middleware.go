package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/soulcodex/deus-cargo-tracker/pkg/logger"
)

type PanicRecoverMiddleware struct {
	logger logger.ZerologLogger
}

func NewPanicRecoverMiddleware(l logger.ZerologLogger) *PanicRecoverMiddleware {
	return &PanicRecoverMiddleware{logger: l}
}

func (prm *PanicRecoverMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				prm.logger.Error().Ctx(r.Context()).Msg("unhandled Error")
				jsonBody, _ := json.Marshal(map[string]interface{}{
					"errors": []map[string]string{
						{
							"title": "Internal Server Error",
						},
					},
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if _, writeErr := w.Write(jsonBody); writeErr != nil {
					prm.logger.Error().
						Ctx(r.Context()).
						Msg("couldn't write the response of the panic error")
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
