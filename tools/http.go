package tools

import (
	"fmt"
	"io"
	"net/http"
)

// HttpDiagnostics - responsible for gathering and logging HTTP driven events
func HttpDiagnostics(httpResp *http.Response) string {
	defer httpResp.Body.Close()

	// Read the response body, so we can include it in the diagnostics message.
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		// Handle error if needed.
	}

	// Make a copy of the request headers and exclude the 'Authorization' header
	requestHeaders := make(http.Header, len(httpResp.Request.Header))

	for key, values := range httpResp.Request.Header {
		if key == "Authorization" {
			requestHeaders["Authorization"] = append([]string{}, "**REDACTED**")
		} else {
			requestHeaders[key] = append([]string{}, values...)
		}
	}

	results := fmt.Sprintf(
		"HTTP Method: %v\n\nRequest URL: %v\n\nRequest Headers: %v\n\nRequest Payload: %v\n\n"+
			"Response Headers: %v\n\nResponse Time: %v\n\nStatus Code: %d\n\nResponse Body: %s\n\n",
		httpResp.Request.Method, httpResp.Request.URL, requestHeaders, httpResp.Request.Body,
		httpResp.Header, httpResp.Header.Get("Date"), httpResp.StatusCode, string(bodyBytes),
	)

	return results
}
