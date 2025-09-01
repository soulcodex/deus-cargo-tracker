package testutils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	httpserver "github.com/soulcodex/deus-cargo-tracker/pkg/http-server"
)

func ExecuteJSONRequest(
	t *testing.T,
	router *httpserver.Router,
	verb,
	path string,
	body []byte,
) *httptest.ResponseRecorder {
	req, err := http.NewRequestWithContext(t.Context(), verb, path, bytes.NewBuffer(body))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	httpRecorder := httptest.NewRecorder()
	router.GetMuxRouter().ServeHTTP(httpRecorder, req)
	return httpRecorder
}
