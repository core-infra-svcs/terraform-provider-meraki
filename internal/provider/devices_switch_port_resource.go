package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &DevicesSwitchPortResource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &DevicesSwitchPortResource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &DevicesSwitchPortResource{} // Interface for resources with import state functionality
)

func NewDevicesSwitchPortResource() resource.Resource {
	return &DevicesSwitchPortResource{}
}

// DevicesSwitchPortResource struct defines the structure for this resource.
// It includes an APIClient field for making requests to the Meraki API.
type DevicesSwitchPortResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The DevicesSwitchPortResourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type DevicesSwitchPortResourceModel struct {
	Id                          jsontypes.String `tfsdk:"id"`
	Serial                      jsontypes.String `tfsdk:"serial" json:"serial"`
	PortId                      jsontypes.String `tfsdk:"port_id" json:"portId"`
	Name                        jsontypes.String `tfsdk:"name" json:"name"`
	Tags                        types.Set        `tfsdk:"tags" json:"tags"`
	Enabled                     jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	PoeEnabled                  jsontypes.Bool   `tfsdk:"poe_enabled" json:"poeEnabled"`
	Type                        jsontypes.String `tfsdk:"type" json:"type"`
	Vlan                        jsontypes.Int64  `tfsdk:"vlan" json:"vlan"`
	VoiceVlan                   jsontypes.Int64  `tfsdk:"voice_vlan" json:"voiceVlan"`
	AllowedVlans                jsontypes.String `tfsdk:"allowed_vlans" json:"allowedVlans"`
	AccessPolicyNumber          jsontypes.Int64  `tfsdk:"access_policy_number" json:"accessPolicyNumber"`
	AccessPolicyType            jsontypes.String `tfsdk:"access_policy_type" json:"accessPolicyType"`
	PortScheduleId              jsontypes.String `tfsdk:"port_schedule_id" json:"portScheduleId"`
	StickyMacAllowListLimit     jsontypes.Int64  `tfsdk:"sticky_mac_allow_list_limit" json:"stickyMacWhitelistLimit"`
	MacAllowList                types.Set        `tfsdk:"mac_allow_list" json:"macWhitelist"`
	StickyMacAllowList          types.Set        `tfsdk:"sticky_mac_allow_list" json:"stickyMacWhitelist"`
	StormControlEnabled         jsontypes.Bool   `tfsdk:"storm_control_enabled" json:"stormControlEnabled"`
	AdaptivePolicyGroupId       jsontypes.String `tfsdk:"adaptive_policy_group_id" json:"adaptivePolicyGroupId"`
	PeerSgtCapable              jsontypes.Bool   `tfsdk:"peer_sgt_capable" json:"peerSgtCapable"`
	FlexibleStackingEnabled     jsontypes.Bool   `tfsdk:"flexible_stacking_enabled" json:"flexibleStackingEnabled"`
	DaiTrusted                  jsontypes.Bool   `tfsdk:"dai_trusted" json:"daiTrusted"`
	IsolationEnabled            jsontypes.Bool   `tfsdk:"isolation_enabled" json:"isolationEnabled"`
	RstpEnabled                 jsontypes.Bool   `tfsdk:"rstp_enabled" json:"rstpEnabled"`
	StpGuard                    jsontypes.String `tfsdk:"stp_guard" json:"stpGuard"`
	LinkNegotiation             jsontypes.String `tfsdk:"link_negotiation" json:"linkNegotiation"`
	LinkNegotiationCapabilities types.List       `tfsdk:"link_negotiation_capabilities" json:"linkNegotiationCapabilities"`
	Udld                        jsontypes.String `tfsdk:"udld" json:"udld"`
	Profile                     types.Object     `tfsdk:"profile" json:"profile"`
}

type DevicesSwitchPortResourceModelProfile struct {
	Enabled jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	Id      jsontypes.String `tfsdk:"id" json:"id"`
	Iname   jsontypes.String `tfsdk:"iname" json:"iname"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *DevicesSwitchPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	// The TypeName attribute is important as it provides the user-friendly name for the resource/data source.
	// This is the name users will use to reference the resource/data source and it's also used in the acceptance tests.
	resp.TypeName = req.ProviderTypeName + "_devices_switch_port"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *DevicesSwitchPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "Network Devices resource. This only works for devices associated with a network.",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "The devices serial number",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
				CustomType: jsontypes.StringType,
			},
			"port_id": schema.StringAttribute{
				MarkdownDescription: "The identifier of the switch port.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the switch port.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "The list of tags of the switch port.",
				ElementType:         jsontypes.StringType,
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The status of the switch port.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
			},
			"poe_enabled": schema.BoolAttribute{
				MarkdownDescription: "The PoE status of the switch port.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the switch port ('trunk' or 'access').",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"vlan": schema.Int64Attribute{
				MarkdownDescription: "The VLAN of the switch port. A null value will clear the value set for trunk ports.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"voice_vlan": schema.Int64Attribute{
				MarkdownDescription: "The voice VLAN of the switch port. Only applicable to access ports.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"allowed_vlans": schema.StringAttribute{
				MarkdownDescription: "The VLANs allowed on the switch port. Only applicable to trunk ports.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"access_policy_type": schema.StringAttribute{
				MarkdownDescription: "The type of the access policy of the switch port. Only applicable to access ports. Can be one of 'Open', 'Custom access policy', 'MAC allow list' or 'Sticky MAC allow list'.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Open", "Custom access policy", "MAC allow list", "Sticky MAC allow list"),
				},
			},
			"access_policy_number": schema.Int64Attribute{
				MarkdownDescription: "The number of a custom access policy to configure on the switch port. Only applicable when 'accessPolicyType' is 'Custom access policy'.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"port_schedule_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the port schedule. A value of null will clear the port schedule.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"sticky_mac_allow_list_limit": schema.Int64Attribute{
				MarkdownDescription: "The maximum number of MAC addresses for sticky MAC allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"mac_allow_list": schema.SetAttribute{
				MarkdownDescription: "Only devices with MAC addresses specified in this list will have access to this port. Up to 20 MAC addresses can be defined. Only applicable when 'accessPolicyType' is 'MAC allow list'.",
				ElementType:         jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"sticky_mac_allow_list": schema.SetAttribute{
				MarkdownDescription: "The initial list of MAC addresses for sticky Mac allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
				ElementType:         jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"storm_control_enabled": schema.BoolAttribute{
				MarkdownDescription: "The storm control status of the switch port.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"adaptive_policy_group_id": schema.StringAttribute{
				MarkdownDescription: "The adaptive policy group ID that will be used to tag traffic through this switch port. This ID must pre-exist during the configuration, else needs to be created using adaptivePolicy/groups API. Cannot be applied to a port on a switch bound to profile.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"peer_sgt_capable": schema.BoolAttribute{
				MarkdownDescription: "If true, Peer SGT is enabled for traffic through this switch port. Applicable to trunk port only, not access port. Cannot be applied to a port on a switch bound to profile.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"flexible_stacking_enabled": schema.BoolAttribute{
				MarkdownDescription: "For supported switches (e.g. MS420/MS425), whether or not the port has flexible stacking enabled.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dai_trusted": schema.BoolAttribute{
				MarkdownDescription: "If true, ARP packets for this port will be considered trusted, and Dynamic ARP Inspection will allow the traffic.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"isolation_enabled": schema.BoolAttribute{
				MarkdownDescription: "The isolation status of the switch port.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"rstp_enabled": schema.BoolAttribute{
				MarkdownDescription: "The rapid spanning tree protocol status.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"stp_guard": schema.StringAttribute{
				MarkdownDescription: "The state of the STP guard ('disabled', 'root guard', 'bpdu guard' or 'loop guard').",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "root guard", "bpdu guard", "loop guard"),
				},
			},
			"link_negotiation": schema.StringAttribute{
				MarkdownDescription: "The link speed for the switch port.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"link_negotiation_capabilities": schema.ListAttribute{
				MarkdownDescription: "The link speeds for the switch port.",
				Computed:            true,
				ElementType:         jsontypes.StringType,
			},
			"udld": schema.StringAttribute{
				MarkdownDescription: "The action to take when Unidirectional Link is detected (Alert only, Enforce). Default configuration is Alert only.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"profile": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "When enabled, override this port's configuration with a port profile.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "When enabled, the ID of the port profile used to override the port's configuration.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"iname": schema.StringAttribute{
						MarkdownDescription: "When enabled, the IName of the profile.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
				},
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
func (r *DevicesSwitchPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// The provider must be properly configured before it can be used.
	if req.ProviderData == nil {
		return
	}

	// Here we expect the provider data to be of type *openApiClient.APIClient.
	client, ok := req.ProviderData.(*openApiClient.APIClient)

	// This is a fatal error and the provider cannot proceed without it.
	// If you see this error, it means there is an issue with the provider setup.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	// This allows the resource to use the configured provider for any API calls it needs to make.
	r.client = client
}

// Create method is responsible for creating a new resource.
// It takes a CreateRequest containing the planned state of the new resource and returns a CreateResponse
// with the final state of the new resource or an error.
func (r *DevicesSwitchPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesSwitchPortResourceModel
	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diag := DevicesSwitchPortResourcePayload(context.Background(), data)
	if diag.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", diag))
		return
	}

	// Initialize provider client and make API call
	response, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Create HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data, diag = DevicesSwitchPortResourceResponse(ctx, response, data)
	if diag.HasError() {
		resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", diag))
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
// It takes a ReadRequest and returns a ReadResponse with the current state of the resource or an error.
func (r *DevicesSwitchPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesSwitchPortResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize provider client and make API call
	response, httpResp, err := r.client.SwitchApi.GetDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).Execute()
	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Read HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	data, diag := DevicesSwitchPortResourceResponse(ctx, response, data)
	if diag.HasError() {
		resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", diag))
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *DevicesSwitchPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *DevicesSwitchPortResourceModel
	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diag := DevicesSwitchPortResourcePayload(context.Background(), data)
	if diag.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", diag))
		return
	}

	// Initialize provider client and make API call
	response, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Create HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	data, diag = DevicesSwitchPortResourceResponse(ctx, response, data)
	if diag.HasError() {
		resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", diag))
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")

}

// Delete function is responsible for deleting a resource.
// It uses a DeleteRequest and responds with a DeleteResponse which contains the updated state of the resource or an error.
func (r *DevicesSwitchPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *DevicesSwitchPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there was an error reading the state, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	payload := *openApiClient.NewUpdateDeviceSwitchPortRequest()
	payload.SetName("")
	payload.SetTags([]string{})
	payload.SetEnabled(false)
	payload.SetPoeEnabled(false)
	payload.SetType("trunk")
	payload.SetVlan(1)
	payload.SetVoiceVlan(1)
	payload.SetAllowedVlans("1")
	payload.SetAccessPolicyType("Open")

	// Initialize provider client and make API call
	_, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
// The function expects an ImportStateRequest and responds with an ImportStateResponse which contains
// the new state of the resource or an error.
func (r *DevicesSwitchPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: serial, port_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("port_id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func DevicesSwitchPortResourcePayload(ctx context.Context, data *DevicesSwitchPortResourceModel) (openApiClient.UpdateDeviceSwitchPortRequest, diag.Diagnostics) {

	// Create HTTP request body
	payload := *openApiClient.NewUpdateDeviceSwitchPortRequest()

	if !data.Enabled.IsUnknown() && !data.Enabled.IsNull() {
		payload.SetEnabled(data.Enabled.ValueBool())
	}
	if !data.Name.IsUnknown() && !data.Name.IsNull() {
		payload.SetName(data.Name.ValueString())
	}
	if !data.Tags.IsUnknown() && !data.Tags.IsNull() {
		var tags []string
		diags := data.Tags.ElementsAs(ctx, &tags, false)
		if diags.HasError() {
			return payload, diags

		}
		payload.SetTags(tags)
	}
	if !data.PoeEnabled.IsUnknown() && !data.PoeEnabled.IsNull() {
		payload.SetEnabled(data.PoeEnabled.ValueBool())
	}
	if !data.Type.IsUnknown() && !data.Type.IsNull() {
		payload.SetType(data.Type.ValueString())
	}
	if !data.IsolationEnabled.IsUnknown() && !data.IsolationEnabled.IsNull() {
		payload.SetIsolationEnabled(data.IsolationEnabled.ValueBool())
	}
	if !data.RstpEnabled.IsUnknown() && !data.RstpEnabled.IsNull() {
		payload.SetRstpEnabled(data.RstpEnabled.ValueBool())
	}
	if !data.StpGuard.IsUnknown() && !data.StpGuard.IsNull() {
		payload.SetStpGuard(data.StpGuard.ValueString())
	}
	if !data.LinkNegotiation.IsUnknown() && !data.LinkNegotiation.IsNull() {
		payload.SetLinkNegotiation(data.LinkNegotiation.ValueString())
	}
	if !data.Udld.IsUnknown() && !data.Udld.IsNull() {
		payload.SetUdld(data.Udld.ValueString())
	}

	if !data.Vlan.IsUnknown() && !data.Vlan.IsNull() {
		payload.SetVlan(int32(data.Vlan.ValueInt64()))
	}

	if !data.VoiceVlan.IsUnknown() && !data.VoiceVlan.IsNull() {
		payload.SetVoiceVlan(int32(data.VoiceVlan.ValueInt64()))
	}

	if !data.AllowedVlans.IsUnknown() && !data.AllowedVlans.IsNull() {
		payload.SetAllowedVlans(data.AllowedVlans.ValueString())
	}

	if !data.AccessPolicyNumber.IsUnknown() && !data.AccessPolicyNumber.IsNull() {
		payload.SetAccessPolicyNumber(int32(data.AccessPolicyNumber.ValueInt64()))
	}

	if !data.AccessPolicyType.IsUnknown() && !data.AccessPolicyType.IsNull() {
		payload.SetAccessPolicyType(data.AccessPolicyType.ValueString())
	}

	if !data.AllowedVlans.IsUnknown() && !data.AllowedVlans.IsNull() {
		payload.SetAllowedVlans(data.AllowedVlans.ValueString())
	}

	if !data.DaiTrusted.IsUnknown() && !data.DaiTrusted.IsNull() {
		payload.SetDaiTrusted(data.DaiTrusted.ValueBool())
	}

	if !data.FlexibleStackingEnabled.IsUnknown() && !data.FlexibleStackingEnabled.IsNull() {
		payload.SetFlexibleStackingEnabled(data.FlexibleStackingEnabled.ValueBool())
	}

	if !data.PeerSgtCapable.IsUnknown() && !data.PeerSgtCapable.IsNull() {
		payload.SetPeerSgtCapable(data.PeerSgtCapable.ValueBool())
	}

	if !data.AdaptivePolicyGroupId.IsUnknown() && !data.AdaptivePolicyGroupId.IsNull() {
		payload.SetAdaptivePolicyGroupId(data.AdaptivePolicyGroupId.ValueString())
	}

	if !data.StickyMacAllowListLimit.IsUnknown() && !data.StickyMacAllowListLimit.IsNull() {
		payload.SetStickyMacAllowListLimit(int32(data.StickyMacAllowListLimit.ValueInt64()))
	}

	if !data.MacAllowList.IsUnknown() && !data.MacAllowList.IsNull() {
		var macAllowList []string
		diags := data.MacAllowList.ElementsAs(ctx, &macAllowList, false)
		if diags.HasError() {
			return payload, diags

		}
		payload.SetMacAllowList(macAllowList)
	}

	if !data.StickyMacAllowList.IsUnknown() && !data.StickyMacAllowList.IsNull() {
		var stickyMacAllowList []string
		diags := data.StickyMacAllowList.ElementsAs(ctx, &stickyMacAllowList, false)
		if diags.HasError() {
			return payload, diags

		}
		payload.SetStickyMacAllowList(stickyMacAllowList)
	}

	if !data.StormControlEnabled.IsUnknown() && !data.StormControlEnabled.IsNull() {
		payload.SetStormControlEnabled(data.StormControlEnabled.ValueBool())
	}

	if !data.PortScheduleId.IsUnknown() && !data.PortScheduleId.IsNull() {
		payload.SetPortScheduleId(data.PortScheduleId.ValueString())
	}

	if !data.Profile.IsUnknown() && !data.Profile.IsNull() {
		var profile openApiClient.GetDeviceSwitchPorts200ResponseInnerProfile
		var profileData DevicesSwitchPortResourceModelProfile
		data.Profile.As(ctx, &profileData, basetypes.ObjectAsOptions{})
		profile.SetEnabled(profileData.Enabled.ValueBool())
		profile.SetId(profileData.Id.ValueString())
		profile.SetIname(profileData.Iname.ValueString())
		payload.SetProfile(profile)
	}

	return payload, nil
}

func DevicesSwitchPortResourceResponse(ctx context.Context, response *openApiClient.GetDeviceSwitchPorts200ResponseInner, data *DevicesSwitchPortResourceModel) (*DevicesSwitchPortResourceModel, diag.Diagnostics) {

	profileobjectMap := map[string]attr.Type{
		"enabled": jsontypes.BoolType,
		"id":      jsontypes.StringType,
		"iname":   jsontypes.StringType,
	}
	var profileData DevicesSwitchPortResourceModelProfile
	profileData.Enabled = jsontypes.BoolValue(response.Profile.GetEnabled())
	profileData.Id = jsontypes.StringValue(response.Profile.GetId())
	profileData.Iname = jsontypes.StringValue(response.Profile.GetIname())
	profileObjectValue, diags := types.ObjectValueFrom(ctx, profileobjectMap, profileData)
	if diags.HasError() {
		return data, diags
	}
	data.Profile = profileObjectValue
	data.Enabled = jsontypes.BoolValue(response.GetEnabled())
	data.IsolationEnabled = jsontypes.BoolValue(response.GetIsolationEnabled())
	data.Name = jsontypes.StringValue(response.GetName())
	var tags []jsontypes.String
	for _, element := range response.GetTags() {
		tags = append(tags, jsontypes.StringValue(element))
	}
	tagValues, diags := types.SetValueFrom(ctx, types.StringType, tags)
	if diags.HasError() {
		return data, diags
	}
	data.Tags = tagValues
	data.PoeEnabled = jsontypes.BoolValue(response.GetPoeEnabled())
	data.Type = jsontypes.StringValue(response.GetType())
	data.StpGuard = jsontypes.StringValue(response.GetStpGuard())
	data.RstpEnabled = jsontypes.BoolValue(response.GetRstpEnabled())
	data.LinkNegotiation = jsontypes.StringValue(response.GetLinkNegotiation())
	data.Udld = jsontypes.StringValue(response.GetUdld())
	data.AllowedVlans = jsontypes.StringValue(response.GetAllowedVlans())
	data.Vlan = jsontypes.Int64Value(int64(response.GetVlan()))
	data.VoiceVlan = jsontypes.Int64Value(int64(response.GetVoiceVlan()))
	data.AccessPolicyNumber = jsontypes.Int64Value(int64(response.GetAccessPolicyNumber()))
	data.AccessPolicyType = jsontypes.StringValue(response.GetAccessPolicyType())
	data.DaiTrusted = jsontypes.BoolValue(response.GetDaiTrusted())
	data.StickyMacAllowListLimit = jsontypes.Int64Value(int64(response.GetStickyMacAllowListLimit()))
	var macAllowList []jsontypes.String
	for _, element := range response.GetMacAllowList() {
		macAllowList = append(macAllowList, jsontypes.StringValue(element))
	}
	macAllowListValues, diags := types.SetValueFrom(ctx, types.StringType, macAllowList)
	if diags.HasError() {
		return data, diags
	}
	data.MacAllowList = macAllowListValues
	var stickyMacAllowList []jsontypes.String
	for _, element := range response.GetStickyMacAllowList() {
		stickyMacAllowList = append(stickyMacAllowList, jsontypes.StringValue(element))
	}
	stickyMacAllowListValues, diags := types.SetValueFrom(ctx, types.StringType, stickyMacAllowList)
	if diags.HasError() {
		return data, diags
	}
	data.StickyMacAllowList = stickyMacAllowListValues
	data.StormControlEnabled = jsontypes.BoolValue(response.GetStormControlEnabled())
	data.FlexibleStackingEnabled = jsontypes.BoolValue(response.GetFlexibleStackingEnabled())
	data.PeerSgtCapable = jsontypes.BoolValue(response.GetPeerSgtCapable())
	data.AdaptivePolicyGroupId = jsontypes.StringValue(response.GetAdaptivePolicyGroupId())
	data.PortScheduleId = jsontypes.StringValue(response.GetPortScheduleId())
	var linkNegotiationCapabilities []jsontypes.String
	for _, element := range response.GetLinkNegotiationCapabilities() {
		linkNegotiationCapabilities = append(linkNegotiationCapabilities, jsontypes.StringValue(element))
	}
	listValue, diags := types.ListValueFrom(ctx, types.StringType, linkNegotiationCapabilities)
	if diags.HasError() {
		return data, diags
	}
	data.LinkNegotiationCapabilities = listValue
	data.Id = jsontypes.StringValue(data.Serial.ValueString() + "," + data.PortId.ValueString())

	return data, nil
}
