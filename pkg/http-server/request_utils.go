package httpserver

import (
	"net/http"
	"net/url"
	"strings"
)

const reverseProxyForwardedByHeader = "X-Forwarded-For"

func FetchStringQueryParamValue(values url.Values, param string, defaultVal string) string {
	if queryParamVal := values.Get(param); queryParamVal != "" {
		if unescapedVal, unescapeErr := url.QueryUnescape(queryParamVal); unescapeErr == nil {
			return unescapedVal
		}
	}
	return defaultVal
}

func ClientIP(req *http.Request) string {
	ipAddress := req.RemoteAddr
	fwdAddress := req.Header.Get(reverseProxyForwardedByHeader)
	if fwdAddress != "" {
		ipAddress = fwdAddress

		ips := strings.Split(fwdAddress, ", ")
		if len(ips) > 1 {
			ipAddress = ips[0]
		}
	}

	return ipAddress
}
