package provider

import (
	"context"
	"fmt"

	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworksAppliancePortsDataSource{}

func NewNetworksAppliancePortsDataSource() datasource.DataSource {
	return &NetworksAppliancePortsDataSource{}
}

// NetworksAppliancePortsDataSource defines the data source implementation.
type NetworksAppliancePortsDataSource struct {
	client *openApiClient.APIClient
}

type NetworksAppliancePortsDataSourceListModel struct {
	Id        types.String                            `tfsdk:"id"`
	NetworkId types.String                            `tfsdk:"network_id"`
	List      []NetworksAppliancePortsDataSourceModel `tfsdk:"list"`
}

// NetworksAppliancePortsDataSourceModel describes the data source data model.
type NetworksAppliancePortsDataSourceModel struct {
	Accesspolicy        types.String `tfsdk:"access_policy"`
	Allowedvlans        types.String `tfsdk:"allowed_vlans"`
	Dropuntaggedtraffic types.Bool   `tfsdk:"drop_untagged_traffic"`
	Enabled             types.Bool   `tfsdk:"enabled"`
	Number              types.Int64  `tfsdk:"number"`
	Type                types.String `tfsdk:"type"`
	Vlan                types.Int64  `tfsdk:"vlan"`
}

func (d *NetworksAppliancePortsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_ports"
}

func (d *NetworksAppliancePortsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksAppliancePorts data source for listing appliance ports",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Description:         "Example identifier",
				MarkdownDescription: "Example identifier",
				Optional:            false,
				Computed:            true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"list": schema.SetNestedAttribute{
				MarkdownDescription: "List of Network Appliance Ports",
				Optional:            true,
				Computed:            true,
				Description:         "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_policy": schema.StringAttribute{
							MarkdownDescription: "The name of the policy. Only applicable to Access ports.",
							Optional:            true,
							Computed:            true,
						},
						"allowed_vlans": schema.StringAttribute{
							MarkdownDescription: "Comma-delimited list of the VLAN ID's allowed on the port, or 'all' to permit all VLAN's on the port.",
							Optional:            true,
							Computed:            true,
						},
						"drop_untagged_traffic": schema.BoolAttribute{
							MarkdownDescription: "Whether the trunk port can drop all untagged traffic.",
							Optional:            true,
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							Description:         "The status of the port",
							MarkdownDescription: "The status of the port",
							Optional:            true,
							Computed:            true,
						},
						"number": schema.Int64Attribute{
							MarkdownDescription: "Number of the port",
							Optional:            true,
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the port: 'access' or 'trunk'.",
							Optional:            true,
							Computed:            true,
						},
						"vlan": schema.Int64Attribute{
							MarkdownDescription: "Native VLAN when the port is in Trunk mode. Access VLAN when the port is in Access mode.",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *NetworksAppliancePortsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworksAppliancePortsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NetworksAppliancePortsDataSourceListModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if len(data.NetworkId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing NetworkId", fmt.Sprintf("Value: %s", data.NetworkId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := d.client.ApplianceApi.GetNetworkAppliancePorts(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read datasource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	data.Id = types.StringValue("example-id")
	for _, appliance_port := range inlineResp {
		var result NetworksAppliancePortsDataSourceModel
		result.Accesspolicy = types.StringValue(appliance_port.GetAccessPolicy())
		result.Allowedvlans = types.StringValue(appliance_port.GetAllowedVlans())
		result.Dropuntaggedtraffic = types.BoolValue(appliance_port.GetDropUntaggedTraffic())
		result.Type = types.StringValue(appliance_port.GetType())
		result.Number = types.Int64Value(int64(appliance_port.GetNumber()))
		result.Vlan = types.Int64Value(int64(appliance_port.GetVlan()))
		result.Enabled = types.BoolValue(appliance_port.GetEnabled())
		data.List = append(data.List, result)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
