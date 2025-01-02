package utils

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"
)

// HandleAPIError processes API errors and maps them to Terraform diagnostics.
func HandleAPIError(ctx context.Context, resp *http.Response, err error, diags *diag.Diagnostics) error {
	if err != nil {
		diags.AddError("API Error", fmt.Sprintf("Error during API call: %s", err.Error()))
		return err
	}

	if resp.StatusCode >= 400 {
		diags.AddError("API Response Error", fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// ExtractResponseToMap reads an HTTP response body and unmarshals the JSON content into a map[string]interface{}.
// It returns the map along with any error that occurs during the read or unmarshal process.
func ExtractResponseToMap(resp *http.Response) (map[string]interface{}, error) {
	if resp == nil {
		return nil, fmt.Errorf("received nil http response")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response to map: %v", err)
	}

	return result, nil
}

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

// NewHttpDiagnostics - responsible for gathering and logging HTTP driven events
func NewHttpDiagnostics(httpResp *http.Response, bodyContent string) string {
	if httpResp == nil {
		return "No HTTP Response to Diagnose (Check Internet Connectivity)"
	}

	sanitizedHeaders := make(http.Header)
	for key, values := range httpResp.Request.Header {
		if key == "Authorization" {
			sanitizedHeaders[key] = []string{"**REDACTED**"}
		} else {
			sanitizedHeaders[key] = values
		}
	}

	return fmt.Sprintf("HTTP Method: %s\nRequest URL: %s\nRequest Headers: %v\nStatus Code: %d\nResponse Body: %s",
		httpResp.Request.Method, httpResp.Request.URL, sanitizedHeaders, httpResp.StatusCode, bodyContent)
}

// ReadAndCloseBody Define a helper function to read and close the HTTP response body
func ReadAndCloseBody(httpResp *http.Response) (string, error) {
	if httpResp == nil {
		return "", nil
	}
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	httpResp.Body.Close() // Ensure the body is closed after reading
	return string(bodyBytes), nil
}

// CustomHttpRequestRetry Helper function for retrying API calls. This is to recover from backend congestion errors which manifest as 4XX response codes
func CustomHttpRequestRetry[T any](ctx context.Context, maxRetries int, initialRetryDelay time.Duration, apiCall func() (T, *http.Response, error)) (T, *http.Response, error) {
	var result T
	var lastResponse *http.Response
	var lastError error

	n, _ := rand.Int(rand.Reader, big.NewInt(5))
	retryDelay := time.Duration(n.Int64()+1) * time.Second

	for i := 0; i < maxRetries; i++ {
		tflog.Info(ctx, fmt.Sprintf("Attempt %d/%d", i+1, maxRetries))

		result, httpResp, err := apiCall()
		if httpResp != nil && httpResp.StatusCode >= 200 && httpResp.StatusCode <= 299 {
			return result, httpResp, nil
		}

		// Log error message before retry
		if err != nil {
			var responseBody string
			if httpResp != nil {
				if httpResp.Body != nil {
					bodyBytes, readErr := io.ReadAll(httpResp.Body)
					if readErr == nil {
						responseBody = string(bodyBytes)
					} else {
						responseBody = fmt.Sprintf("Failed to read response body: %s", readErr)
					}
					closeErr := httpResp.Body.Close()
					if closeErr != nil {
						tflog.Error(ctx, fmt.Sprintf("Error closing HTTP response body: %s", closeErr))
					} // Close the body to free up the connection
				} else {
					responseBody = "No response body"
				}

				// Check if the error message matches the criteria for retrying
				if httpResp.StatusCode >= 400 && httpResp.StatusCode <= 599 {
					tflog.Warn(ctx, fmt.Sprintf("Retry %d/%d due to specific error: %s", i+1, maxRetries, responseBody))
				} else {
					// If not, break the loop and return the error
					return result, httpResp, fmt.Errorf("API call failed with status %d: %s", httpResp.StatusCode, responseBody)
				}

				tflog.Warn(ctx, fmt.Sprintf("Retry %d/%d", i+1, maxRetries))
				tflog.Trace(ctx, fmt.Sprintf("API call failed with status %d: %s", httpResp.StatusCode, responseBody))
				lastResponse = httpResp
				lastError = fmt.Errorf("%w: %s", err, responseBody)
			} else {
				tflog.Warn(ctx, fmt.Sprintf("Retry %d/%d", i+1, maxRetries))
				tflog.Trace(ctx, fmt.Sprintf("API call failed: %s", err))
				lastError = err
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("Retry %d/%d due to no response or no error", i+1, maxRetries))
		}

		if i < maxRetries-1 {
			tflog.Info(ctx, fmt.Sprintf("Sleeping for %s before next retry", retryDelay))
			time.Sleep(retryDelay)
		}
	}

	return result, lastResponse, fmt.Errorf("after %d retries, last error: %w", maxRetries, lastError)
}

// CustomHttpRequestRetryStronglyTyped is a generic function that leverages CustomHttpRequestRetry
func CustomHttpRequestRetryStronglyTyped[T any](ctx context.Context, maxRetries int, retryDelay time.Duration, apiCall func() (T, *http.Response, error, diag.Diagnostics)) (T, *http.Response, error, diag.Diagnostics) {
	var diags diag.Diagnostics
	var result T
	var lastResponse *http.Response
	var lastError error

	for i := 0; i < maxRetries; i++ {
		tflog.Info(ctx, fmt.Sprintf("Attempt %d/%d", i+1, maxRetries))
		result, httpResp, err, newDiags := apiCall()
		diags = append(diags, newDiags...) // Accumulate diagnostics from each attempt

		if httpResp != nil && httpResp.StatusCode >= 200 && httpResp.StatusCode <= 299 {
			return result, httpResp, nil, diags
		}

		// Append diagnostics about retry
		if httpResp != nil {
			var responseBody string
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				} else {
					responseBody = fmt.Sprintf("Failed to read response body: %s", readErr)
				}
				closeErr := httpResp.Body.Close()
				if closeErr != nil {
					diags.AddError("HTTP Body Close Error", fmt.Sprintf("Error closing HTTP response body: %s", closeErr))
				} // Close the body to free up the connection
			} else {
				responseBody = "No response body"
			}

			// Check if the error message matches the criteria for retrying
			if httpResp.StatusCode >= 400 && httpResp.StatusCode <= 599 {
				tflog.Warn(ctx, fmt.Sprintf("Retry %d/%d due to specific error: %s", i+1, maxRetries, responseBody))
				diags = append(diags, diag.NewErrorDiagnostic(
					fmt.Sprintf("Retry %d/%d", i+1, maxRetries),
					fmt.Sprintf("API call failed with status %d: %s", httpResp.StatusCode, responseBody),
				))
			} else {
				// If not, break the loop and return the error
				return result, httpResp, fmt.Errorf("API call failed with status %d: %s", httpResp.StatusCode, responseBody), diags
			}

			lastResponse = httpResp
			lastError = err
		} else {
			diags = append(diags, diag.NewErrorDiagnostic(
				fmt.Sprintf("Retry %d/%d", i+1, maxRetries),
				fmt.Sprintf("API call failed: %s", err),
			))
			lastError = err
		}

		if i < maxRetries-1 {
			// Ensure retryDelay is in milliseconds
			tflog.Info(ctx, fmt.Sprintf("Sleeping for %s before next retry", retryDelay*time.Millisecond))
			time.Sleep(retryDelay * time.Millisecond)
			// Exponential backoff: Increase retry delay for next attempt
			retryDelay *= 2
		}
	}

	// Return the last error and response if all retries are exhausted
	return result, lastResponse, fmt.Errorf("after %d retries, last error: %w", maxRetries, lastError), diags
}
