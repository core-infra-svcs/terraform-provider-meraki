package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
	"net/url"
	"os"
	"time"
)

// BearerAuthTransport Custom transport to add bearer token in the Authorization header
type BearerAuthTransport struct {
	Transport *http.Transport
	Token     string
}

func (t *BearerAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add the bearer token to the Authorization header
	req.Header.Set("Authorization", "Bearer "+t.Token)
	// Use the underlying transport to perform the actual request
	return t.Transport.RoundTrip(req)
}

func (p *CiscoMerakiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CiscoMerakiProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get http retryClient variables and default values
	configuration := openApiClient.NewConfiguration()

	// Debug
	if p.version == "dev" {
		// always enable debug for provider development
		configuration.Debug = true
	} else if data.LoggingEnabled.ValueBool() {
		// check if user enabled debug in the provider
		configuration.Debug = data.LoggingEnabled.ValueBool()
	}

	// MERAKI BASE URL
	if !data.BaseUrl.IsNull() {
		baseUrl, err := url.Parse(data.BaseUrl.ValueString())
		if err == nil {
			configuration.Servers = openApiClient.ServerConfigurations{
				{
					URL:         baseUrl.String() + "/{basePath}",
					Description: "No description provided",
					Variables: map[string]openApiClient.ServerVariable{
						"basePath": {
							Description:  "Meraki API Go Client",
							DefaultValue: data.BasePath.ValueString(),
						},
					},
				},
			}
		}
	}

	// UserAgent
	configuration.UserAgent = configuration.UserAgent + " terraform/" + p.version

	// Set certificate path
	if !data.CertificatePath.IsNull() {
		configuration.CertificatePath = data.CertificatePath.ValueString()
	}

	// Proxy
	if !data.Proxy.IsNull() {
		configuration.RequestsProxy = data.Proxy.ValueString()
	}

	// SingleRequestTimeout
	if !data.SingleRequestTimeout.IsNull() {
		configuration.SingleRequestTimeout = int(data.SingleRequestTimeout.ValueInt64())
	}

	// MaximumRetries
	if !data.MaximumRetries.IsNull() {
		configuration.MaximumRetries = int(data.MaximumRetries.ValueInt64())
	}

	// Nginx429RetryWaitTime
	if !data.Nginx429RetryWaitTime.IsNull() {
		configuration.Nginx429RetryWaitTime = int(data.Nginx429RetryWaitTime.ValueInt64())
	}

	// New custom retryable retryClient
	retryClient := retryablehttp.NewClient()

	// Custom Transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	// add certificate to retryClient if certificate path isn't empty
	if configuration.CertificatePath != "" {
		// Load the certificate file
		certFile := configuration.CertificatePath
		cert, err := os.ReadFile(certFile)
		if err != nil {
			e := fmt.Sprintf("%v", err.Error())
			tflog.Error(ctx, e)
		}

		// Create a certificate pool and add the certificate
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(cert)

		// Create a custom Cert pool with the certificate and add TLS configuration to transport
		transport.TLSClientConfig.RootCAs = certPool
	}

	if configuration.RequestsProxy != "" {
		proxyUrl, err := url.Parse(configuration.RequestsProxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	// Set single request timeout in transport
	retryClient.HTTPClient.Timeout = time.Duration(configuration.SingleRequestTimeout) * time.Second

	retryClient.RetryMax = configuration.MaximumRetries
	retryClient.RetryWaitMax = time.Duration(configuration.Nginx429RetryWaitTime) * time.Second

	configuration.UserAgent = configuration.UserAgent + "terraform" + p.version

	// Set Bearer Token in transport
	authenticatedTransport := &BearerAuthTransport{
		Transport: transport,
	}

	// MERAKI DASHBOARD API KEY
	if !data.ApiKey.IsNull() {
		authenticatedTransport.Token = data.ApiKey.ValueString()
	} else {
		authenticatedTransport.Token = os.Getenv("MERAKI_DASHBOARD_API_KEY")
	}
	retryClient.HTTPClient.Transport = authenticatedTransport
	configuration.HTTPClient = retryClient.HTTPClient

	client := openApiClient.NewAPIClient(configuration)

	if client == nil {
		tflog.Error(ctx, "Error creating API Client")
		return
	}

	// Pass the encryption key to resources and data sources
	resp.DataSourceData = client
	resp.ResourceData = client
}
