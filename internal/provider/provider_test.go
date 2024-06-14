package provider

import (
	"context"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"meraki": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	// Check environmental variables
	ev := []string{"MERAKI_DASHBOARD_API_KEY", "TF_ACC_MERAKI_MX_LICENCE", "TF_ACC_MERAKI_MX_SERIAL",
		"TF_ACC_MERAKI_MS_SERIAL", "TF_ACC_MERAKI_MG_SERIAL", "TF_ACC_MERAKI_ORDER_NUMBER", "TF_ACC_MERAKI_ORGANIZATION_ID"}
	for _, v := range ev {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for acceptance tests", v)
		}
	}

	// Sweep for test networks prefixed with "acc_test" to delete

}

func TestHttpClientRetryLogic(t *testing.T) {
	t.Run("Retry with Retry-After header", func(t *testing.T) {
		mockServer := createMockServer(1, 3, 10) // 404 condition not relevant for this test
		defer mockServer.Close()

		client := configureClient(mockServer.URL, 5, 2*time.Second)
		ctx := context.Background()

		req, err := http.NewRequest("GET", mockServer.URL+"/test", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.GetConfig().HTTPClient.Do(req.WithContext(ctx))
		assert.NoError(t, err, "Expected no error after retries")
		if resp != nil {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK after retries")
		} else {
			t.Error("Response is nil")
		}
	})

	t.Run("Immediate success without retries", func(t *testing.T) {
		mockServer := createMockServer(1, 1, 10) // 404 condition not relevant for this test
		defer mockServer.Close()

		client := configureClient(mockServer.URL, 3, 2*time.Second)
		ctx := context.Background()

		req, err := http.NewRequest("GET", mockServer.URL+"/test", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		resp, err := client.GetConfig().HTTPClient.Do(req.WithContext(ctx))
		assert.NoError(t, err, "Expected no error without retries")
		if resp != nil {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK without retries")
		} else {
			t.Error("Response is nil")
		}
	})

}

func createMockServer(retryAfter int, successAfter int, notFoundAfter int) *httptest.Server {
	var count int64

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		currentCount := atomic.AddInt64(&count, 1)
		if currentCount < int64(successAfter) {
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "too many requests"}`))
		} else if currentCount == int64(notFoundAfter) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "not found"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": "true"}`))
		}
	})

	mux.HandleFunc("/networks", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": "L_1234567890",
				"name": "Network 1",
				"timeZone": "America/Los_Angeles",
				"tags": " tag1 tag2 "
			},
			{
				"id": "L_0987654321",
				"name": "Network 2",
				"timeZone": "America/New_York",
				"tags": " tag3 tag4 "
			}
		]`))
	})

	mux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"id": "123456",
				"name": "Organization 1"
			},
			{
				"id": "654321",
				"name": "Organization 2"
			}
		]`))
	})

	mux.HandleFunc("/devices", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{
				"serial": "Q234-ABCD-5678",
				"mac": "00:11:22:33:44:55",
				"networkId": "L_1234567890",
				"model": "MR34",
				"name": "My AP"
			}
		]`))
	})

	return httptest.NewServer(mux)
}

func configureClient(baseURL string, retries int, retryWaitMax time.Duration) *openApiClient.APIClient {
	configuration := openApiClient.NewConfiguration()
	configuration.Servers = openApiClient.ServerConfigurations{
		{
			URL: baseURL,
		},
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = retries
	retryClient.RetryWaitMax = retryWaitMax
	//retryClient.CheckRetry = customRetryPolicy
	retryClient.HTTPClient.Transport = &bearerAuthTransport{
		Transport: http.DefaultTransport.(*http.Transport),
		Token:     "dummy_api_key",
	}

	configuration.HTTPClient = retryClient.StandardClient()
	return openApiClient.NewAPIClient(configuration)
}

/*
func customRetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			retryAfterSeconds, convErr := strconv.Atoi(retryAfter)
			if convErr == nil {
				select {
				case <-ctx.Done():
					return false, ctx.Err()
				case <-time.After(time.Duration(retryAfterSeconds) * time.Second):
				}
			}
		}
		return true, nil
	}
	return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
}
*/
