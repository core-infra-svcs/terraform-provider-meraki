package administered

import (
	"context"
	"fmt"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider-defined types fully satisfy framework interfaces
var _ datasource.DataSource = &DataSource{}

// NewDataSource initializes a new Administered Identities Me data source.
func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

// DataSource implements the Terraform data source for retrieving the current user's identity.
type DataSource struct {
	client *openApiClient.APIClient
}

// Metadata sets the data source type name.
func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_administered_identities_me", req.ProviderTypeName)
}

// Schema sets the data source schema.
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = identitiesMeSchema
}

// Configure initializes the API client for the data source.
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Ensure the provider has been configured
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

// Read retrieves the current user's identity and populates the Terraform state.
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config dataSourceModel

	// Read Terraform configuration data into the model
	tflog.Debug(ctx, "[identities_me] Reading Terraform configuration data")
	if err := req.Config.Get(ctx, &config); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	// Call the API to retrieve the current user's identity
	tflog.Debug(ctx, "[identities_me] Calling API to retrieve identity")
	apiResponse, httpResp, err := d.client.AdministeredApi.GetAdministeredIdentitiesMe(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"API Request Failure",
			fmt.Sprintf("Failed to call API: %s\nDiagnostics: %s", err.Error(), utils.HttpDiagnostics(httpResp)),
		)
		return
	}

	// Handle unexpected HTTP status codes
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Status Code",
			fmt.Sprintf("Received HTTP status code: %d", httpResp.StatusCode),
		)
		return
	}

	// Map the API response to the Terraform state model
	tflog.Debug(ctx, "[identities_me] Mapping API response to Terraform state")
	state, diags := mapAPIResponseToState(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"State Mapping Failure",
			"Failed to map API response to Terraform state. Check diagnostics for details.",
		)
		return
	}

	// Save the mapped state to Terraform
	tflog.Debug(ctx, "[identities_me] Saving state to Terraform")
	if err := resp.State.Set(ctx, &state); err != nil {
		resp.Diagnostics.Append(err...)
	}
}
