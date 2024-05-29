package tools

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"log"
	"net/http"
	"time"
)

// HttpDiagnostics - responsible for gathering and logging HTTP driven events
func HttpDiagnostics(httpResp *http.Response) string {
	if httpResp != nil {
		defer httpResp.Body.Close()

		// Read the response body, so we can include it in the diagnostics message.
		bodyBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			log.Fatal(err)
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

	return "No HTTP Response to Diagnose (Check Internet Connectivity)"
}

// RetryOn4xx Helper function for retrying API calls
func RetryOn4xx(ctx context.Context, maxRetries int, retryDelay time.Duration, apiCall func() (map[string]interface{}, *http.Response, error)) (map[string]interface{}, *http.Response, error) {
	retries := 0
	for retries < maxRetries {
		result, httpResp, err := apiCall()
		if err == nil {
			return result, httpResp, nil
		}
		if httpResp != nil && httpResp.StatusCode >= 400 && httpResp.StatusCode < 500 {
			tflog.Warn(ctx, "Retrying API call due to 4xx error", map[string]interface{}{
				"maxRetries":        maxRetries,
				"retryDelay":        retryDelay,
				"remainingAttempts": maxRetries - retries - 1,
				"httpStatusCode":    httpResp.StatusCode,
				"httpBody":          httpResp.Body,
			})
			time.Sleep(retryDelay * time.Second)
			retries++
		} else {
			return nil, httpResp, err
		}
	}
	return nil, nil, fmt.Errorf("max retries reached")
}
