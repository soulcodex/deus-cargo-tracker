package testutils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
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

func CheckResponse(
	t *testing.T,
	expectedStatusCode int,
	expectedResponse string,
	response *httptest.ResponseRecorder,
	formats ...interface{},
) {
	ja := jsonassert.New(t)
	CheckResponseCode(t, expectedStatusCode, response.Code)

	receivedResponse := response.Body.String()
	if receivedResponse == "" {
		assert.Equal(t, expectedResponse, receivedResponse)
		return
	}
	if formats != nil {
		ja.Assertf(receivedResponse, expectedResponse, formats...)
	} else {
		ja.Assertf(receivedResponse, "%s", expectedResponse)
	}
}

func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func EmptyHTTPHeaders() map[string]string {
	return map[string]string{}
}
