package templates

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// DataSourceTemplate provides a reusable base implementation for data sources.
type DataSourceTemplate struct {
	Client *openApiClient.APIClient
}

// Configure sets up the API client for the data source.
func (t *DataSourceTemplate) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configuration Type",
			fmt.Sprintf("Expected *openApiClient.APIClient, got: %T", req.ProviderData),
		)
		return
	}

	t.Client = client
}

// Metadata sets the data source type name.
func Metadata(req datasource.MetadataRequest, resp *datasource.MetadataResponse, suffix string) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, suffix)
}

// WithRetry provides retry logic for API calls.
func (t *DataSourceTemplate) WithRetry(ctx context.Context, operation func() (interface{}, *http.Response, error)) (interface{}, *http.Response, error) {
	maxRetries := t.Client.GetConfig().MaximumRetries
	retryDelay := time.Duration(t.Client.GetConfig().Retry4xxErrorWaitTime)

	return utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, operation)
}
