package administered

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &IdentitiesMeDataSource{}

func NewAdministeredIdentitiesMeDataSource() datasource.DataSource {
	return &IdentitiesMeDataSource{}
}

// IdentitiesMeDataSource defines the data source implementation.
type IdentitiesMeDataSource struct {
	client *openApiClient.APIClient
}

func (d *IdentitiesMeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_administered_identities_me"
}

func (d *IdentitiesMeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = identitiesMeSchema
}

func (d *IdentitiesMeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *IdentitiesMeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data identitiesMeAttrModel

	// Read Terraform configuration data into the model
	tflog.Trace(ctx, "Reading Terraform configuration data")
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API call to retrieve identity
	tflog.Trace(ctx, "Calling API to retrieve identity")
	inlineResp, httpResp, err := d.client.AdministeredApi.GetAdministeredIdentitiesMe(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			fmt.Sprintf("Error: %s\nDiagnostics: %s", err.Error(), utils.HttpDiagnostics(httpResp)),
		)
		return
	}

	// Handle unexpected status codes
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("Received HTTP status code: %v", httpResp.StatusCode),
		)
		return
	}

	// Marshal API response into Terraform state
	tflog.Trace(ctx, "Marshaling API response into Terraform state")
	marshaledData, diags := marshalState(ctx, inlineResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"Data Model Error",
			"Failed to marshal HTTP response into Terraform state. Check diagnostics for details.",
		)
		return
	}

	// Save marshaled data into Terraform state
	tflog.Trace(ctx, "Saving marshaled data into Terraform state")
	resp.Diagnostics.Append(resp.State.Set(ctx, &marshaledData)...)
}
