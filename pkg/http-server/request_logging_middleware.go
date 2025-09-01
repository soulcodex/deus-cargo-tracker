package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/soulcodex/deus-cargo-tracker/pkg/logger"
	"github.com/soulcodex/deus-cargo-tracker/pkg/utils"
)

type RequestLoggingMiddleware struct {
	logger       logger.ZerologLogger
	timeProvider utils.DateTimeProvider
}

func NewRequestLoggingMiddleware(logger logger.ZerologLogger, tp utils.DateTimeProvider) *RequestLoggingMiddleware {
	return &RequestLoggingMiddleware{logger: logger, timeProvider: tp}
}

func (rlm *RequestLoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		recorder, ipAddress, startedAt := NewStatusRecorder(w), ClientIP(req), rlm.timeProvider.Now()

		defer func() {
			totalTime := time.Since(startedAt).Milliseconds()

			rlm.logger.Info().
				Ctx(req.Context()).
				Str("remote_addr_ip", ipAddress).
				Int64("request_time_d", totalTime).
				Int("status", recorder.Status()).
				Str("request", req.RequestURI).
				Time("request_started_at", startedAt).
				Str("request_method", req.Method).
				Str("http_referrer", req.Referer()).
				Str("http_user_agent", req.Header.Get("User-Agent")).
				Str("response_content_type", w.Header().Get("Content-Type")).
				Str("correlation_id", w.Header().Get(HeaderRequestID)).
				Msg(
					fmt.Sprintf(
						"%s %s %d %s %s %d %s",
						req.Method,
						req.RequestURI,
						recorder.Status(),
						rlm.timeProvider.Now().Format(time.RFC3339),
						ipAddress,
						totalTime,
						req.Referer(),
					),
				)
		}()

		next.ServeHTTP(recorder, req)
	})
}
