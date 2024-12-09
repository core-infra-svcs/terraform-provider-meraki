package _interface

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
	"time"
)

// Ensure the provider-defined types fully satisfy the framework interfaces
var _ datasource.DataSource = &ManagementInterfaceDataSource{}

// NewDataSource initializes a new Management Interface data source.
func NewDataSource() datasource.DataSource {
	return &ManagementInterfaceDataSource{}
}

// ManagementInterfaceDataSource represents the data source for management interface settings.
type ManagementInterfaceDataSource struct {
	client *openApiClient.APIClient
}

// Metadata sets the data source type name.
func (d *ManagementInterfaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_devices_management_interface", req.ProviderTypeName)
}

// Schema defines the schema for the data source.
func (d *ManagementInterfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GetDatasourceSchema
}

// Configure initializes the API client for the data source.
func (d *ManagementInterfaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Ensure the provider has been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configuration",
			fmt.Sprintf("Expected *openApiClient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Read retrieves the management interface settings and updates the Terraform state.
func (d *ManagementInterfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config resourceModel

	// Read Terraform configuration data into the model
	tflog.Debug(ctx, "[management_interface] Reading Terraform configuration data")
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retry mechanism for API call
	maxRetries := d.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(d.client.GetConfig().Retry4xxErrorWaitTime)

	// Execute API call with retry logic
	tflog.Debug(ctx, "[management_interface] Calling API to retrieve management interface settings")
	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
		return d.client.DevicesApi.GetDeviceManagementInterface(ctx, config.Serial.ValueString()).Execute()
	})
	if err != nil {
		tflog.Error(ctx, "[management_interface] API call failed", map[string]interface{}{
			"error": err.Error(),
		})
		resp.Diagnostics.AddError("API Call Failure", fmt.Sprintf("Error details: %s", err.Error()))
		return
	}

	// Handle unexpected HTTP status codes
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response",
			utils.HttpDiagnostics(httpResp),
		)

	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// HTTP 400 counts as an error so moving this here
		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				utils.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", config))
			return
		}

	}

	diags = updateDatasourceState(ctx, &config, inlineResp, httpResp)

	config.Id = types.StringValue(config.Serial.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)

	// Write logs using the tflog package
	tflog.Debug(ctx, "[management_interface] Successfully completed Read operation")
	tflog.Trace(ctx, "read resource")
}
