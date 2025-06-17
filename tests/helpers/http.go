package helpers

import (
	"bytes"
	"fmt"
	"io"
)

// FormatStatusCode creates a test-friendly name for HTTP status codes.
func FormatStatusCode(statusCode int) string {
	switch statusCode {
	case 400:
		return "400_bad_request"
	case 401:
		return "401_unauthorized"
	case 403:
		return "403_forbidden"
	case 404:
		return "404_not_found"
	case 429:
		return "429_rate_limit"
	case 500:
		return "500_internal_server_error"
	case 502:
		return "502_bad_gateway"
	case 503:
		return "503_service_unavailable"
	case 504:
		return "504_gateway_timeout"
	default:
		return fmt.Sprintf("%d_status", statusCode)
	}
}

// CreateResponseBody creates an io.ReadCloser from byte slice for HTTP responses.
func CreateResponseBody(data []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(data))
}
