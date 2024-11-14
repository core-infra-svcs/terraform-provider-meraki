package ports

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
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

type NetworksAppliancePortsDataSourceModel struct {
	Id        jsontypes.String                            `tfsdk:"id"`
	NetworkId jsontypes.String                            `tfsdk:"network_id"`
	List      []NetworksAppliancePortsDataSourceModelList `tfsdk:"list"`
}

// NetworksAppliancePortsDataSourceModelList describes the data source data model.
type NetworksAppliancePortsDataSourceModelList struct {
	Accesspolicy        jsontypes.String `tfsdk:"access_policy" json:"access_policy"`
	Allowedvlans        jsontypes.String `tfsdk:"allowed_vlans" json:"allowed_vlans"`
	Dropuntaggedtraffic jsontypes.Bool   `tfsdk:"drop_untagged_traffic" json:"drop_untagged_traffic"`
	Enabled             jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	Number              jsontypes.Int64  `tfsdk:"number" json:"number"`
	Type                jsontypes.String `tfsdk:"type" json:"type"`
	Vlan                jsontypes.Int64  `tfsdk:"vlan" json:"vlan"`
}

func (d *NetworksAppliancePortsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_ports"
}

func (d *NetworksAppliancePortsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get appliance ports",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Description:         "Example identifier",
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"list": schema.SetNestedAttribute{
				MarkdownDescription: "Ports of Network Appliance Ports",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_policy": schema.StringAttribute{
							MarkdownDescription: "The name of the policy. Only applicable to Access ports.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"allowed_vlans": schema.StringAttribute{
							MarkdownDescription: "Comma-delimited list of the VLAN ID's allowed on the port, or 'all' to permit all VLAN's on the port.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"drop_untagged_traffic": schema.BoolAttribute{
							MarkdownDescription: "Whether the trunk port can drop all untagged traffic.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"enabled": schema.BoolAttribute{
							Description:         "The status of the port",
							MarkdownDescription: "The status of the port",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"number": schema.Int64Attribute{
							MarkdownDescription: "SsidNumber of the port",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the port: 'access' or 'trunk'.",
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"access", "trunk"}...),
								stringvalidator.LengthAtLeast(4),
							},
							CustomType: jsontypes.StringType,
						},
						"vlan": schema.Int64Attribute{
							MarkdownDescription: "Native VLAN when the port is in Trunk mode. Access VLAN when the port is in Access mode.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
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
	var data NetworksAppliancePortsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := d.client.ApplianceApi.GetNetworkAppliancePorts(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	data.Id = jsontypes.StringValue("example-id")
	for _, appliance_port := range inlineResp {
		var result NetworksAppliancePortsDataSourceModelList
		result.Accesspolicy = jsontypes.StringValue(appliance_port.GetAccessPolicy())
		result.Allowedvlans = jsontypes.StringValue(appliance_port.GetAllowedVlans())
		result.Dropuntaggedtraffic = jsontypes.BoolValue(appliance_port.GetDropUntaggedTraffic())
		result.Type = jsontypes.StringValue(appliance_port.GetType())
		result.Number = jsontypes.Int64Value((int64(appliance_port.GetNumber())))
		result.Vlan = jsontypes.Int64Value(int64(appliance_port.GetVlan()))
		result.Enabled = jsontypes.BoolValue(appliance_port.GetEnabled())
		result.Number = jsontypes.Int64Value(int64(appliance_port.GetNumber()))
		data.List = append(data.List, result)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
