package _interface

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure the provider-defined types fully satisfy the framework interfaces
var _ datasource.DataSource = &DataSource{}

// NewDataSource initializes a new Management Interface data source.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource represents the data source for management interface settings.
type DataSource struct {
	client *openApiClient.APIClient
}

// Metadata sets the data source type name.
func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_devices_management_interface", req.ProviderTypeName)
}

// Schema defines the schema for the data source.
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = GetDatasourceSchema
}

// Configure initializes the API client for the data source.
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DataSourceModel

	// Read Terraform configuration data into the model
	tflog.Debug(ctx, "[management_interface] Reading Terraform configuration data")
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the required 'serial' attribute
	if config.Serial.IsNull() || config.Serial.IsUnknown() || config.Serial.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Serial Number",
			"The 'serial' attribute must be specified in the data source configuration.",
		)
		return
	}

	// Call the READ API (GET)
	tflog.Debug(ctx, "[management_interface] Calling API to retrieve management interface settings")
	apiResponse, httpResp, err := CallReadAPI(ctx, d.client, config.Serial.ValueString())
	if err := utils.HandleAPIError(ctx, httpResp, err, &resp.Diagnostics); err != nil {
		resp.Diagnostics.AddError("API Call Failure", fmt.Sprintf("Error details: %s", err.Error()))
		return
	}

	// Handle unexpected HTTP status codes
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// Marshal API response into Terraform state
	state, diags := MarshalStateFromAPI(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state for Terraform
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log successful operation
	tflog.Debug(ctx, "[management_interface] Successfully completed Read operation", map[string]interface{}{
		"serial": config.Serial.ValueString(),
	})
}
