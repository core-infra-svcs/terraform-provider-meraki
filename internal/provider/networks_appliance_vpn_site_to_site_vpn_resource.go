package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NetworksApplianceVpnSiteToSiteVpnResource{}
var _ resource.ResourceWithImportState = &NetworksApplianceVpnSiteToSiteVpnResource{}

func NewNetworksApplianceVpnSiteToSiteVpnResource() resource.Resource {
	return &NetworksApplianceVpnSiteToSiteVpnResource{}
}

// NetworksApplianceVpnSiteToSiteVpnResource defines the resource implementation.
type NetworksApplianceVpnSiteToSiteVpnResource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceVpnSiteToSiteVpnResourceModel describes the resource data model.
type NetworksApplianceVpnSiteToSiteVpnResourceModel struct {
	Id        jsontypes.String                                    `tfsdk:"id"`
	NetworkId jsontypes.String                                    `tfsdk:"network_id" json:"network_id"`
	Mode      jsontypes.String                                    `tfsdk:"mode" json:"mode"`
	Hubs      []NetworksNetworkIdApplianceVpnSiteToSiteVpnHubs    `tfsdk:"hubs" json:"hubs"`
	Subnets   []NetworksNetworkIdApplianceVpnSiteToSiteVpnSubnets `tfsdk:"subnets" json:"subnets"`
}

type NetworksNetworkIdApplianceVpnSiteToSiteVpnHubs struct {
	HubId           jsontypes.String `tfsdk:"hub_id" json:"hubId"`
	UseDefaultRoute jsontypes.Bool   `tfsdk:"use_default_route" json:"useDefaultRoute"`
}

type NetworksNetworkIdApplianceVpnSiteToSiteVpnSubnets struct {
	LocalSubnet jsontypes.String `tfsdk:"local_subnet" json:"localSubnet"`
	UseVpn      jsontypes.Bool   `tfsdk:"use_vpn" json:"useVpn"`
}

func (r *NetworksApplianceVpnSiteToSiteVpnResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vpn_site_to_site_vpn"
}

func (r *NetworksApplianceVpnSiteToSiteVpnResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Networks appliance vpn site to site vpn resource for updating networks appliance vpn site to site vpn. Only valid for MX networks in NAT mode.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "The site-to-site VPN mode.",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"hubs": schema.SetNestedAttribute{
				Description: "The list of VPN hubs, in order of preference.",
				Required:    true,
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
			"subnets": schema.SetNestedAttribute{
				Description: "The list of subnets and their VPN presence.",
				Required:    true,
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

func (r *NetworksApplianceVpnSiteToSiteVpnResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceVpnSiteToSiteVpnResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksApplianceVpnSiteToSiteVpnResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceVpnSiteToSiteVpn := *openApiClient.NewUpdateNetworkApplianceVpnSiteToSiteVpnRequest(data.Mode.ValueString())

	// For mode value "none" hubs and subnets should be empty
	if data.Mode.ValueString() != "none" {

		// Hubs
		var hubs []openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestHubsInner
		for _, attribute := range data.Hubs {
			var hubData openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestHubsInner
			hubData.SetHubId(attribute.HubId.ValueString())
			hubData.SetUseDefaultRoute(attribute.UseDefaultRoute.ValueBool())
			hubs = append(hubs, hubData)
		}
		updateNetworkApplianceVpnSiteToSiteVpn.SetHubs(hubs)

		// Subnets
		var subnets []openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestSubnetsInner
		for _, attribute := range data.Subnets {
			var subnetData openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestSubnetsInner
			subnetData.SetLocalSubnet(attribute.LocalSubnet.ValueString())
			subnetData.SetUseVpn(attribute.UseVpn.ValueBool())
			subnets = append(subnets, subnetData)
		}
		updateNetworkApplianceVpnSiteToSiteVpn.SetSubnets(subnets)
	}

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceVpnSiteToSiteVpn(ctx, data.NetworkId.ValueString()).UpdateNetworkApplianceVpnSiteToSiteVpnRequest(updateNetworkApplianceVpnSiteToSiteVpn).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksApplianceVpnSiteToSiteVpnResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceVpnSiteToSiteVpnResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVpnSiteToSiteVpn(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksApplianceVpnSiteToSiteVpnResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksApplianceVpnSiteToSiteVpnResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceVpnSiteToSiteVpn := *openApiClient.NewUpdateNetworkApplianceVpnSiteToSiteVpnRequest(data.Mode.ValueString())

	// For mode value "none" hubs and subnets should be empty
	if data.Mode.ValueString() != "none" {

		// Hubs
		var hubs []openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestHubsInner
		for _, attribute := range data.Hubs {
			var hubData openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestHubsInner
			hubData.SetHubId(attribute.HubId.ValueString())
			hubData.SetUseDefaultRoute(attribute.UseDefaultRoute.ValueBool())
			hubs = append(hubs, hubData)
		}
		updateNetworkApplianceVpnSiteToSiteVpn.SetHubs(hubs)

		// Subnets
		var subnets []openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestSubnetsInner
		for _, attribute := range data.Subnets {
			var subnetData openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestSubnetsInner
			subnetData.SetLocalSubnet(attribute.LocalSubnet.ValueString())
			subnetData.SetUseVpn(attribute.UseVpn.ValueBool())
			subnets = append(subnets, subnetData)
		}
		updateNetworkApplianceVpnSiteToSiteVpn.SetSubnets(subnets)
	}

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceVpnSiteToSiteVpn(ctx, data.NetworkId.ValueString()).UpdateNetworkApplianceVpnSiteToSiteVpnRequest(updateNetworkApplianceVpnSiteToSiteVpn).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceVpnSiteToSiteVpnResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksApplianceVpnSiteToSiteVpnResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVpnSiteToSiteVpn(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksApplianceVpnSiteToSiteVpnResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
