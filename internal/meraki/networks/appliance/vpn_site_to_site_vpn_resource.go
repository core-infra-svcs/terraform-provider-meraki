package appliance

import (
	"context"
	"fmt"
	jsontypes2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NetworksApplianceVpnSiteToSiteVpnResource{}
var _ resource.ResourceWithConfigure = &NetworksApplianceVpnSiteToSiteVpnResource{}
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
	Id        jsontypes2.String `tfsdk:"id"`
	NetworkId jsontypes2.String `tfsdk:"network_id" json:"network_id"`
	Mode      jsontypes2.String `tfsdk:"mode" json:"mode"`
	Hubs      types.List        `tfsdk:"hubs" json:"hubs"`
	Subnets   types.List        `tfsdk:"subnets" json:"subnets"`
}

type NetworksApplianceVpnSiteToSiteVpnResourceModelHubs struct {
	HubId           jsontypes2.String `tfsdk:"hub_id" json:"hubId"`
	UseDefaultRoute jsontypes2.Bool   `tfsdk:"use_default_route" json:"useDefaultRoute"`
}

type NetworksApplianceVpnSiteToSiteVpnResourceModelSubnets struct {
	LocalSubnet jsontypes2.String `tfsdk:"local_subnet" json:"localSubnet"`
	UseVpn      jsontypes2.Bool   `tfsdk:"use_vpn" json:"useVpn"`
}

func (r *NetworksApplianceVpnSiteToSiteVpnResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vpn_site_to_site_vpn"
}

func (r *NetworksApplianceVpnSiteToSiteVpnResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage networks appliance vpn site to site vpn. Only valid for MX networks in NAT mode.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes2.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes2.StringType,
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
				CustomType:          jsontypes2.StringType,
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
							CustomType:          jsontypes2.StringType,
						},
						"use_default_route": schema.BoolAttribute{
							MarkdownDescription: "Indicates whether default route traffic should be sent to this hub.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes2.BoolType,
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
							CustomType:          jsontypes2.StringType,
						},
						"use_vpn": schema.BoolAttribute{
							MarkdownDescription: "Indicates the presence of the subnet in the VPN.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes2.BoolType,
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

	payload, diags := NetworkApplianceVpnSiteToSiteVpnResourcePayload(context.Background(), data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Create Payload Error", fmt.Sprintf("\n%v", diags))
		return
	}

	response, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceVpnSiteToSiteVpn(ctx, data.NetworkId.ValueString()).UpdateNetworkApplianceVpnSiteToSiteVpnRequest(*payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Create HTTP Client Failure",
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Save data into Terraform state
	data, diags = NetworksApplianceVpnSiteToSiteVpnResourceResponse(ctx, response, data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Create Response Error", fmt.Sprintf("\n%v", diags))
		return
	}

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

	// Save data into Terraform state
	data, diags := NetworksApplianceVpnSiteToSiteVpnResourceResponse(ctx, response, data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Read Response Error", fmt.Sprintf("\n%v", diags))
		return
	}

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

	payload, diags := NetworkApplianceVpnSiteToSiteVpnResourcePayload(context.Background(), data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Update Payload Error", fmt.Sprintf("\n%v", diags))
		return
	}

	response, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceVpnSiteToSiteVpn(ctx, data.NetworkId.ValueString()).UpdateNetworkApplianceVpnSiteToSiteVpnRequest(*payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Update HTTP Client Failure",
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Save data into Terraform state
	data, diags = NetworksApplianceVpnSiteToSiteVpnResourceResponse(ctx, response, data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Update Response Error", fmt.Sprintf("\n%v", diags))
		return
	}

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

	updateNetworkApplianceVpnSiteToSiteVpn := *openApiClient.NewUpdateNetworkApplianceVpnSiteToSiteVpnRequest(data.Mode.ValueString())

	updateNetworkApplianceVpnSiteToSiteVpn.SetMode("none")

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceVpnSiteToSiteVpn(ctx, data.NetworkId.ValueString()).UpdateNetworkApplianceVpnSiteToSiteVpnRequest(updateNetworkApplianceVpnSiteToSiteVpn).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete HTTP Client Failure",
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

func NetworkApplianceVpnSiteToSiteVpnResourcePayload(ctx context.Context, data *NetworksApplianceVpnSiteToSiteVpnResourceModel) (*openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequest, diag.Diagnostics) {

	payload := *openApiClient.NewUpdateNetworkApplianceVpnSiteToSiteVpnRequest(data.Mode.ValueString())

	// For mode value "none" hubs and subnets should be empty
	if data.Mode.ValueString() != "none" {

		// Hubs
		var hubs []openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestHubsInner

		if !data.Hubs.IsUnknown() && !data.Hubs.IsNull() {
			var hubsPayload []NetworksApplianceVpnSiteToSiteVpnResourceModelHubs
			diags := data.Hubs.ElementsAs(ctx, &hubsPayload, false)
			if diags.HasError() {
				return nil, diags
			}

			for _, hubValue := range hubsPayload {
				var hubData openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestHubsInner

				hubData.SetHubId(hubValue.HubId.ValueString())
				hubData.SetUseDefaultRoute(hubValue.UseDefaultRoute.ValueBool())
				hubs = append(hubs, hubData)
			}

			payload.SetHubs(hubs)
		} else {
			payload.SetHubs(nil)
		}

		if !data.Subnets.IsUnknown() && !data.Subnets.IsNull() {
			// Subnets
			var subnets []openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestSubnetsInner
			var subnetsPayload []NetworksApplianceVpnSiteToSiteVpnResourceModelSubnets
			diags := data.Subnets.ElementsAs(ctx, &subnetsPayload, false)
			if diags.HasError() {
				return nil, diags
			}

			for _, subnetValue := range subnetsPayload {
				var subnetData openApiClient.UpdateNetworkApplianceVpnSiteToSiteVpnRequestSubnetsInner
				subnetData.SetLocalSubnet(subnetValue.LocalSubnet.ValueString())
				subnetData.SetUseVpn(subnetValue.UseVpn.ValueBool())
				subnets = append(subnets, subnetData)
			}

			payload.SetSubnets(subnets)
		}
	} else {
		payload.SetSubnets(nil)
	}

	data.Id = jsontypes2.StringValue("example-id")

	return &payload, nil

}

func NetworksApplianceVpnSiteToSiteVpnResourceResponse(ctx context.Context, response *openApiClient.GetNetworkApplianceVpnSiteToSiteVpn200Response, data *NetworksApplianceVpnSiteToSiteVpnResourceModel) (*NetworksApplianceVpnSiteToSiteVpnResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	data.Mode = jsontypes2.StringValue(response.GetMode())

	// Hubs
	var hubs []NetworksApplianceVpnSiteToSiteVpnResourceModelHubs
	for _, element := range response.GetHubs() {
		var hub NetworksApplianceVpnSiteToSiteVpnResourceModelHubs
		hub.UseDefaultRoute = jsontypes2.BoolValue(element.GetUseDefaultRoute())
		hub.HubId = jsontypes2.StringValue(element.GetHubId())
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
		return data, diags
	}

	// Subnets
	var subnets []NetworksApplianceVpnSiteToSiteVpnResourceModelSubnets
	for _, element := range response.GetSubnets() {
		var subnet NetworksApplianceVpnSiteToSiteVpnResourceModelSubnets
		subnet.UseVpn = jsontypes2.BoolValue(element.GetUseVpn())
		subnet.LocalSubnet = jsontypes2.StringValue(element.GetLocalSubnet())

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
		return data, diags
	}
	return data, nil
}
