package httpserver

import (
	"fmt"
	"net/http"
)

// StatusRecorder wraps an http.ResponseWriter and records
// the HTTP status code sent to the client.
type StatusRecorder struct {
	Writer      http.ResponseWriter
	StatusCode  int  // recorded status code
	wroteHeader bool // true once WriteHeader has been called
}

// NewStatusRecorder creates a new StatusRecorder.
// StatusCode remains zero until a header is written or data is sent.
func NewStatusRecorder(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		Writer:      w,
		StatusCode:  0,
		wroteHeader: false,
	}
}

// WriteHeader records the status code on the first call
// and forwards it to the underlying ResponseWriter.
func (r *StatusRecorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.StatusCode = code
		r.wroteHeader = true
		r.Writer.WriteHeader(code)
	}
}

// Write sends data to the client. If WriteHeader was never called,
// it records and sends StatusOK by default.
func (r *StatusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	bytesWritten, err := r.Writer.Write(b)
	if err != nil {
		return bytesWritten, fmt.Errorf("failed to write response: %w", err)
	}
	return bytesWritten, nil
}

// Header returns the header map that will be sent by WriteHeader.
func (r *StatusRecorder) Header() http.Header {
	return r.Writer.Header()
}

// Status returns the recorded status code.
func (r *StatusRecorder) Status() int {
	return r.StatusCode
}
