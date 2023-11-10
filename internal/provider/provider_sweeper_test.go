package provider

import (
	"crypto/tls"
	"github.com/hashicorp/go-retryablehttp"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// SweeperHTTPClient returns a common provider client configured for the specified region
func SweeperHTTPClient() (*openApiClient.APIClient, error) {

	// Get http retryClient variables and default values
	configuration := openApiClient.NewConfiguration()

	// UserAgent
	configuration.UserAgent = configuration.UserAgent + " terraform/dev-sweeper"

	// Custom Transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	// Set Bearer Token in transport
	authenticatedTransport := &bearerAuthTransport{
		Transport: transport,
	}

	authenticatedTransport.Token = os.Getenv("MERAKI_DASHBOARD_API_KEY")

	// New custom retryable retryClient
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = authenticatedTransport
	configuration.HTTPClient = retryClient.HTTPClient

	client := openApiClient.NewAPIClient(configuration)

	return client, nil
}
