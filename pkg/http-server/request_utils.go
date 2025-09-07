package httpserver

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const reverseProxyForwardedByHeader = "X-Forwarded-For"

func FetchBoolQueryParamValue(values url.Values, param string, defaultVal bool) bool {
	if queryParamVal := values.Get(param); queryParamVal != "" {
		if unescapedVal, unescapeErr := url.QueryUnescape(queryParamVal); unescapeErr == nil {
			if parsed, parseErr := strconv.ParseBool(unescapedVal); parseErr == nil {
				return parsed
			}

			return defaultVal
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
