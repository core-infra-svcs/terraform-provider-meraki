package vpn

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworksApplianceVpnSiteToSiteVpnDatasource{}

func NewNetworksApplianceVpnSiteToSiteVpnDatasource() datasource.DataSource {
	return &NetworksApplianceVpnSiteToSiteVpnDatasource{}
}

// NetworksApplianceVpnSiteToSiteVpnDatasource defines the resource implementation.
type NetworksApplianceVpnSiteToSiteVpnDatasource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceVpnSiteToSiteVpnDatasourceModel NetworksApplianceVpnSiteToSiteVpnResourceModel describes the resource data model.
type NetworksApplianceVpnSiteToSiteVpnDatasourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id" json:"network_id"`
	Mode      jsontypes.String `tfsdk:"mode" json:"mode"`
	Hubs      types.List       `tfsdk:"hubs" json:"hubs"`
	Subnets   types.List       `tfsdk:"subnets" json:"subnets"`
}

type NetworksApplianceVpnSiteToSiteVpnDatasourceModelHubs struct {
	HubId           jsontypes.String `tfsdk:"hub_id" json:"hubId"`
	UseDefaultRoute jsontypes.Bool   `tfsdk:"use_default_route" json:"useDefaultRoute"`
}

type NetworksApplianceVpnSiteToSiteVpnDatasourceModelSubnets struct {
	LocalSubnet jsontypes.String `tfsdk:"local_subnet" json:"localSubnet"`
	UseVpn      jsontypes.Bool   `tfsdk:"use_vpn" json:"useVpn"`
}

func (r *NetworksApplianceVpnSiteToSiteVpnDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vpn_site_to_site_vpn"
}

func (r *NetworksApplianceVpnSiteToSiteVpnDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage networks appliance vpn site to site vpn. Only valid for MX networks in NAT mode.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "The site-to-site VPN mode.",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"hubs": schema.ListNestedAttribute{
				Description: "The list of VPN hubs, in order of preference.",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hub_id": schema.StringAttribute{
							MarkdownDescription: "The network ID of the hub",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"use_default_route": schema.BoolAttribute{
							MarkdownDescription: "Indicates whether default route traffic should be sent to this hub.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
			"subnets": schema.ListNestedAttribute{
				Description: "The list of subnets and their VPN presence.",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"local_subnet": schema.StringAttribute{
							MarkdownDescription: "The CIDR notation subnet used within the VPN",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"use_vpn": schema.BoolAttribute{
							MarkdownDescription: "Indicates the presence of the subnet in the VPN.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
		},
	}
}

func (r *NetworksApplianceVpnSiteToSiteVpnDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NetworksApplianceVpnSiteToSiteVpnDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *NetworksApplianceVpnSiteToSiteVpnDatasourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVpnSiteToSiteVpn(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Read HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	} else {
		resp.Diagnostics.Append()
	}

	var diags diag.Diagnostics

	data.Mode = jsontypes.StringValue(response.GetMode())

	// Hubs
	var hubs []NetworksApplianceVpnSiteToSiteVpnResourceModelHubs
	for _, element := range response.GetHubs() {
		var hub NetworksApplianceVpnSiteToSiteVpnResourceModelHubs
		hub.UseDefaultRoute = jsontypes.BoolValue(element.GetUseDefaultRoute())
		hub.HubId = jsontypes.StringValue(element.GetHubId())
		hubs = append(hubs, hub)

	}

	hubAttributes := map[string]attr.Type{
		"use_default_route": types.BoolType,
		"hub_id":            types.StringType,
	}

	hubSchema := types.ObjectType{
		AttrTypes: hubAttributes,
	}

	data.Hubs, diags = types.ListValueFrom(ctx, hubSchema, hubs)
	if diags.HasError() {
		return
	}

	// Subnets
	var subnets []NetworksApplianceVpnSiteToSiteVpnResourceModelSubnets
	for _, element := range response.GetSubnets() {
		var subnet NetworksApplianceVpnSiteToSiteVpnResourceModelSubnets
		subnet.UseVpn = jsontypes.BoolValue(element.GetUseVpn())
		subnet.LocalSubnet = jsontypes.StringValue(element.GetLocalSubnet())

		subnets = append(subnets, subnet)

	}

	subnetAttributes := map[string]attr.Type{
		"use_vpn":      types.BoolType,
		"local_subnet": types.StringType,
	}

	subnetSchema := types.ObjectType{
		AttrTypes: subnetAttributes,
	}

	data.Subnets, diags = types.ListValueFrom(ctx, subnetSchema, subnets)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
