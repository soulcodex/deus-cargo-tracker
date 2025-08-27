package httpserver

import (
	"context"
	"net/http"

	"github.com/soulcodex/deus-cargoes-tracker/pkg/utils"
)

type ctxKeyRequestID struct{}

const HeaderRequestID = "X-Request-Id"

// IDGenerator defines a function that produces a new request ID.
type IDGenerator func() string

// RequestIDMiddleware injects or propagates a request ID via headers and context.
type RequestIDMiddleware struct {
	gen IDGenerator
}

// NewRequestIDMiddleware creates a RequestIDMiddleware. If gen is nil, it uses DefaultIDGenerator.
func NewRequestIDMiddleware(gen IDGenerator) *RequestIDMiddleware {
	if gen == nil {
		gen = defaultIDGenerator()
	}
	return &RequestIDMiddleware{gen: gen}
}

// Handler wraps next and ensures every request has X-Request-Id set on both
// the ResponseWriter and the incoming *http.Request's context.
func (m *RequestIDMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// reuse existing header if present, otherwise generate a new one
		reqID := r.Header.Get(HeaderRequestID)
		if reqID == "" {
			reqID = m.gen()
		}

		// set on response
		w.Header().Set(HeaderRequestID, reqID)
		// propagate into context and request header for downstream handlers
		r = r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID{}, reqID))
		r.Header.Set(HeaderRequestID, reqID)

		next.ServeHTTP(w, r)
	})
}

// GetRequestID retrieves the request ID from context.
func GetRequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ctxKeyRequestID{}).(string)
	return id, ok
}

// defaultIDGenerator returns a generator that uses uuid.NewString().
func defaultIDGenerator() IDGenerator {
	provider := utils.NewRandomUUIDProvider()
	return func() string {
		return provider.New().String()
	}
}
