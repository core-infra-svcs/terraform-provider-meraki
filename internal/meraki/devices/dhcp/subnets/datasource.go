package subnets

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

type ApplianceDhcpSubnetsDataSource struct {
	client *openApiClient.APIClient
}

// NewDevicesApplianceDhcpSubnetsDataSource initializes the data source.
func NewDevicesApplianceDhcpSubnetsDataSource() datasource.DataSource {
	return &ApplianceDhcpSubnetsDataSource{}
}

// Metadata provides metadata for the data source.
func (d *ApplianceDhcpSubnetsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_appliance_dhcp_subnets"
}

// SchemaResource returns the schema definition.
func (d *ApplianceDhcpSubnetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = SchemaDataSource()
}

// Configure configures the data source with the API client.
func (d *ApplianceDhcpSubnetsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := req.ProviderData.(*openApiClient.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *openApiClient.APIClient, got: %T", req.ProviderData),
		)
		return
	}
	d.client = client
}

// Read fetches data from the API and sets the state.
func (d *ApplianceDhcpSubnetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := d.client.SubnetsApi.GetDeviceApplianceDhcpSubnets(ctx, data.Serial.ValueString())
	apiResp, httpResp, err := apiReq.Execute()
	if err != nil || httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Failed to fetch DHCP subnets",
			fmt.Sprintf("Error: %v, HTTP Response: %v", err, httpResp),
		)
		return
	}

	resp.Diagnostics.Append(mapApiResponseToModel(apiResp, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(data.Serial.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Debug(ctx, "Read appliance DHCP subnets", map[string]interface{}{"serial": data.Serial.ValueString()})
}
