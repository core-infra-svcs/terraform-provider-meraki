package ports

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"strings"

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
var _ resource.Resource = &NetworksAppliancePortsResource{}
var _ resource.ResourceWithImportState = &NetworksAppliancePortsResource{}

func NewNetworksAppliancePortsResource() resource.Resource {
	return &NetworksAppliancePortsResource{}
}

// NetworksAppliancePortsResource defines the resource implementation.
type NetworksAppliancePortsResource struct {
	client *openApiClient.APIClient
}

// NetworksAppliancePortsResourceModel describes the resource data model.
type NetworksAppliancePortsResourceModel struct {
	Id                  jsontypes.String `tfsdk:"id"`
	NetworkId           jsontypes.String `tfsdk:"network_id"`
	PortId              jsontypes.String `tfsdk:"port_id"`
	Accesspolicy        jsontypes.String `tfsdk:"access_policy" json:"access_policy"`
	Allowedvlans        jsontypes.String `tfsdk:"allowed_vlans" json:"allowed_vlans"`
	Dropuntaggedtraffic jsontypes.Bool   `tfsdk:"drop_untagged_traffic" json:"drop_untagged_traffic"`
	Enabled             jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	Number              jsontypes.Int64  `tfsdk:"number" json:"number"`
	Type                jsontypes.String `tfsdk:"type" json:"type"`
	Vlan                jsontypes.Int64  `tfsdk:"vlan" json:"vlan"`
}

func (r *NetworksAppliancePortsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_ports"
}

func (r *NetworksAppliancePortsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksAppliancePorts resource for updating Network Appliance Firewall L3 Firewall Rules.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"port_id": schema.StringAttribute{
				MarkdownDescription: "Port ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
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
	}
}

func (r *NetworksAppliancePortsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksAppliancePortsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data *NetworksAppliancePortsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkAppliancePortRequest()

	if !data.Vlan.IsUnknown() && !data.Vlan.IsNull() && data.Vlan != jsontypes.Int64Value(0) {
		var vlan = int32(data.Vlan.ValueInt64())
		payload.Vlan = &vlan
	}

	if !data.Type.IsUnknown() && !data.Type.IsNull() && data.Type != jsontypes.StringValue("") {
		payload.Type = data.Type.ValueStringPointer()
	}

	if !data.Enabled.IsUnknown() && !data.Enabled.IsNull() {
		payload.Enabled = data.Enabled.ValueBoolPointer()
	}

	if !data.Accesspolicy.IsUnknown() && !data.Accesspolicy.IsNull() && data.Accesspolicy != jsontypes.StringValue("") {
		payload.AccessPolicy = data.Accesspolicy.ValueStringPointer()
	}
	if !data.Allowedvlans.IsUnknown() && !data.Allowedvlans.IsNull() && data.Allowedvlans != jsontypes.StringValue("") {
		payload.AllowedVlans = data.Allowedvlans.ValueStringPointer()
	}
	if !data.Dropuntaggedtraffic.IsUnknown() && !data.Dropuntaggedtraffic.IsNull() {
		payload.DropUntaggedTraffic = data.Dropuntaggedtraffic.ValueBoolPointer()
	}

	response, httpResp, err := r.client.ApplianceApi.UpdateNetworkAppliancePort(context.Background(), data.NetworkId.ValueString(), data.PortId.ValueString()).UpdateNetworkAppliancePortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

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

	data.Enabled = jsontypes.BoolValue(response.GetEnabled())
	data.Type = jsontypes.StringValue(response.GetType())
	data.Allowedvlans = jsontypes.StringValue(response.GetAllowedVlans())
	data.Dropuntaggedtraffic = jsontypes.BoolValue(response.GetDropUntaggedTraffic())
	data.Vlan = jsontypes.Int64Value(int64(response.GetVlan()))
	data.Accesspolicy = jsontypes.StringValue(response.GetAccessPolicy())
	data.Number = jsontypes.Int64Value(int64(response.GetNumber()))
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksAppliancePortsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksAppliancePortsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, httpResp, err := r.client.ApplianceApi.GetNetworkAppliancePort(context.Background(), data.NetworkId.ValueString(), data.PortId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
	}

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

	data.Enabled = jsontypes.BoolValue(response.GetEnabled())
	data.Type = jsontypes.StringValue(response.GetType())
	data.Allowedvlans = jsontypes.StringValue(response.GetAllowedVlans())
	data.Dropuntaggedtraffic = jsontypes.BoolValue(response.GetDropUntaggedTraffic())
	data.Vlan = jsontypes.Int64Value(int64(response.GetVlan()))
	data.Accesspolicy = jsontypes.StringValue(response.GetAccessPolicy())
	data.Number = jsontypes.Int64Value(int64(response.GetNumber()))
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksAppliancePortsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksAppliancePortsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkAppliancePortRequest()

	if !data.Vlan.IsUnknown() && !data.Vlan.IsNull() && data.Vlan != jsontypes.Int64Value(0) {
		var vlan = int32(data.Vlan.ValueInt64())
		payload.Vlan = &vlan
	}

	if !data.Type.IsUnknown() && !data.Type.IsNull() && data.Type != jsontypes.StringValue("") {
		payload.Type = data.Type.ValueStringPointer()
	}

	if !data.Enabled.IsUnknown() && !data.Enabled.IsNull() {
		payload.Enabled = data.Enabled.ValueBoolPointer()
	}

	if !data.Accesspolicy.IsUnknown() && !data.Accesspolicy.IsNull() && data.Accesspolicy != jsontypes.StringValue("") {
		payload.AccessPolicy = data.Accesspolicy.ValueStringPointer()
	}
	if !data.Allowedvlans.IsUnknown() && !data.Allowedvlans.IsNull() && data.Allowedvlans != jsontypes.StringValue("") {
		payload.AllowedVlans = data.Allowedvlans.ValueStringPointer()
	}
	if !data.Dropuntaggedtraffic.IsUnknown() && !data.Dropuntaggedtraffic.IsNull() {
		payload.DropUntaggedTraffic = data.Dropuntaggedtraffic.ValueBoolPointer()
	}

	response, httpResp, err := r.client.ApplianceApi.UpdateNetworkAppliancePort(context.Background(), data.NetworkId.ValueString(), data.PortId.ValueString()).UpdateNetworkAppliancePortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

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

	data.Enabled = jsontypes.BoolValue(response.GetEnabled())
	data.Type = jsontypes.StringValue(response.GetType())
	data.Allowedvlans = jsontypes.StringValue(response.GetAllowedVlans())
	data.Dropuntaggedtraffic = jsontypes.BoolValue(response.GetDropUntaggedTraffic())
	data.Vlan = jsontypes.Int64Value(int64(response.GetVlan()))
	data.Accesspolicy = jsontypes.StringValue(response.GetAccessPolicy())
	data.Number = jsontypes.Int64Value(int64(response.GetNumber()))
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksAppliancePortsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksAppliancePortsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkAppliancePortRequest()

	if !data.Vlan.IsUnknown() && !data.Vlan.IsNull() && data.Vlan != jsontypes.Int64Value(0) {
		var vlan = int32(data.Vlan.ValueInt64())
		payload.Vlan = &vlan
	}

	if !data.Type.IsUnknown() && !data.Type.IsNull() && data.Type != jsontypes.StringValue("") {
		payload.Type = data.Type.ValueStringPointer()
	}

	if !data.Enabled.IsUnknown() && !data.Enabled.IsNull() {
		payload.Enabled = data.Enabled.ValueBoolPointer()
	}

	if !data.Accesspolicy.IsUnknown() && !data.Accesspolicy.IsNull() && data.Accesspolicy != jsontypes.StringValue("") {
		payload.AccessPolicy = data.Accesspolicy.ValueStringPointer()
	}
	if !data.Allowedvlans.IsUnknown() && !data.Allowedvlans.IsNull() && data.Allowedvlans != jsontypes.StringValue("") {
		payload.AllowedVlans = data.Allowedvlans.ValueStringPointer()
	}
	if !data.Dropuntaggedtraffic.IsUnknown() && !data.Dropuntaggedtraffic.IsNull() {
		payload.DropUntaggedTraffic = data.Dropuntaggedtraffic.ValueBoolPointer()
	}

	response, httpResp, err := r.client.ApplianceApi.UpdateNetworkAppliancePort(context.Background(), data.NetworkId.ValueString(), data.PortId.ValueString()).UpdateNetworkAppliancePortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

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

	data.Enabled = jsontypes.BoolValue(response.GetEnabled())
	data.Type = jsontypes.StringValue(response.GetType())
	data.Allowedvlans = jsontypes.StringValue(response.GetAllowedVlans())
	data.Dropuntaggedtraffic = jsontypes.BoolValue(response.GetDropUntaggedTraffic())
	data.Vlan = jsontypes.Int64Value(int64(response.GetVlan()))
	data.Accesspolicy = jsontypes.StringValue(response.GetAccessPolicy())
	data.Number = jsontypes.Int64Value(int64(response.GetNumber()))
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksAppliancePortsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, network_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("port_id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}
}
