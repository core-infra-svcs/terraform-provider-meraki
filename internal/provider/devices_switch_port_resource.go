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

	// The Id field is mandatory for all resources. It's used for resource identification and is required
	// for the acceptance tests to run.
	Id                      jsontypes.String               `tfsdk:"id"`
	Serial                  jsontypes.String               `tfsdk:"serial"`
	PortId                  jsontypes.String               `tfsdk:"port_id"`
	Name                    jsontypes.String               `tfsdk:"name"`
	Tags                    []jsontypes.String             `tfsdk:"tags"`
	Enabled                 jsontypes.Bool                 `tfsdk:"enabled"`
	PoeEnabled              jsontypes.Bool                 `tfsdk:"poe_enabled"`
	Type                    jsontypes.String               `tfsdk:"type"`
	Vlan                    jsontypes.Int64                `tfsdk:"vlan"`
	VoiceVlan               jsontypes.Int64                `tfsdk:"voice_vlan"`
	AllowedVlans            jsontypes.String               `tfsdk:"allowed_vlans"`
	IsolationEnabled        jsontypes.Bool                 `tfsdk:"isolation_enabled"`
	RstpEnabled             jsontypes.Bool                 `tfsdk:"rstp_enabled"`
	StpGuard                jsontypes.String               `tfsdk:"stp_guard"`
	AccessPolicyNumber      jsontypes.Int64                `tfsdk:"access_policy_number"`
	AccessPolicyType        jsontypes.String               `tfsdk:"access_policy_type"`
	LinkNegotiation         jsontypes.String               `tfsdk:"link_negotiation"`
	PortScheduleId          jsontypes.String               `tfsdk:"port_schedule_id"`
	Udld                    jsontypes.String               `tfsdk:"udld"`
	StickyMacWhitelistLimit jsontypes.Int64                `tfsdk:"sticky_mac_white_list_limit"`
	StormControlEnabled     jsontypes.Bool                 `tfsdk:"storm_control_enabled"`
	MacWhitelist            []jsontypes.String             `tfsdk:"mac_white_list"`
	StickyMacWhitelist      []jsontypes.String             `tfsdk:"sticky_mac_white_list"`
	AdaptivePolicyGroupId   jsontypes.String               `tfsdk:"adaptive_policy_group_id"`
	PeerSgtCapable          jsontypes.Bool                 `tfsdk:"peer_sgt_capable"`
	FlexibleStackingEnabled jsontypes.Bool                 `tfsdk:"flexible_stacking_enabled"`
	DaiTrusted              jsontypes.Bool                 `tfsdk:"dai_trusted"`
	Profile                 DevicesSerialSwitchPortProfile `tfsdk:"profile"`
}

type DevicesSerialSwitchPortProfile struct {
	Enabled jsontypes.Bool   `tfsdk:"enabled"`
	Id      jsontypes.String `tfsdk:"id"`
	Iname   jsontypes.String `tfsdk:"iname"`
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
				Computed:            true,
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
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The status of the switch port.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
				Computed:            true,
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
			},
			"access_policy_type": schema.StringAttribute{
				MarkdownDescription: "The type of the access policy of the switch port. Only applicable to access ports. Can be one of 'Open', 'Custom access policy', 'MAC allow list' or 'Sticky MAC allow list'.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"access_policy_number": schema.Int64Attribute{
				MarkdownDescription: "The number of a custom access policy to configure on the switch port. Only applicable when 'accessPolicyType' is 'Custom access policy'.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"link_negotiation": schema.StringAttribute{
				MarkdownDescription: "The link speed for the switch port.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"port_schedule_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the port schedule. A value of null will clear the port schedule.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"udld": schema.StringAttribute{
				MarkdownDescription: "The action to take when Unidirectional Link is detected (Alert only, Enforce). Default configuration is Alert only.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"sticky_mac_white_list_limit": schema.Int64Attribute{
				MarkdownDescription: "The maximum number of MAC addresses for sticky MAC allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"storm_control_enabled": schema.BoolAttribute{
				MarkdownDescription: "The storm control status of the switch port.",
				CustomType:          jsontypes.BoolType,
				Optional:            true,
				Computed:            true,
			},
			"mac_white_list": schema.SetAttribute{
				MarkdownDescription: "Only devices with MAC addresses specified in this list will have access to this port. Up to 20 MAC addresses can be defined. Only applicable when 'accessPolicyType' is 'MAC allow list'.",
				ElementType:         jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"sticky_mac_white_list": schema.SetAttribute{
				MarkdownDescription: "The initial list of MAC addresses for sticky Mac allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
				ElementType:         jsontypes.StringType,
				Optional:            true,
				Computed:            true,
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
			},
			"dai_trusted": schema.BoolAttribute{
				MarkdownDescription: "If true, ARP packets for this port will be considered trusted, and Dynamic ARP Inspection will allow the traffic.",
				CustomType:          jsontypes.BoolType,
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

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Create HTTP request body
	payload := *openApiClient.NewUpdateDeviceSwitchPortRequest()

	var tags []string
	for _, tag := range data.Tags {
		tags = append(tags, tag.ValueString())
	}
	payload.Tags = tags

	payload.Enabled = data.Enabled.ValueBoolPointer()
	payload.PoeEnabled = data.PoeEnabled.ValueBoolPointer()
	payload.Type = data.Type.ValueStringPointer()

	vlan := int32(data.Vlan.ValueInt64())
	payload.Vlan = &vlan

	voiceVlan := int32(data.VoiceVlan.ValueInt64())
	payload.VoiceVlan = &voiceVlan

	payload.AllowedVlans = data.AllowedVlans.ValueStringPointer()
	payload.IsolationEnabled = data.IsolationEnabled.ValueBoolPointer()
	payload.RstpEnabled = data.RstpEnabled.ValueBoolPointer()
	payload.StpGuard = data.StpGuard.ValueStringPointer()
	payload.LinkNegotiation = data.LinkNegotiation.ValueStringPointer()
	payload.PortScheduleId = data.PortScheduleId.ValueStringPointer()
	payload.Udld = data.Udld.ValueStringPointer()
	payload.AccessPolicyType = data.AccessPolicyType.ValueStringPointer()

	accessPolicyNumber := int32(data.AccessPolicyNumber.ValueInt64())
	payload.AccessPolicyNumber = &accessPolicyNumber

	var macAllowList []string
	for _, mac := range data.MacWhitelist {
		macAllowList = append(macAllowList, mac.ValueString())
	}
	payload.MacAllowList = macAllowList

	var stickyMacAllowList []string
	for _, mac := range data.StickyMacWhitelist {
		stickyMacAllowList = append(stickyMacAllowList, mac.ValueString())
	}
	payload.StickyMacAllowList = stickyMacAllowList

	stickyMacAllowListLimit := int32(data.StickyMacWhitelistLimit.ValueInt64())
	payload.StickyMacAllowListLimit = &stickyMacAllowListLimit

	payload.StormControlEnabled = data.StormControlEnabled.ValueBoolPointer()
	payload.AdaptivePolicyGroupId = data.AdaptivePolicyGroupId.ValueStringPointer()
	payload.PeerSgtCapable = data.PeerSgtCapable.ValueBoolPointer()
	payload.FlexibleStackingEnabled = data.FlexibleStackingEnabled.ValueBoolPointer()
	payload.DaiTrusted = data.DaiTrusted.ValueBoolPointer()

	profile := openApiClient.NewGetDeviceSwitchPorts200ResponseInnerProfile()
	profile.Id = data.Profile.Id.ValueStringPointer()
	profile.Iname = data.Profile.Iname.ValueStringPointer()
	profile.Enabled = data.Profile.Enabled.ValueBoolPointer()
	payload.Profile = profile

	// Initialize provider client and make API call
	_, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Initialize provider client and make API call
	_, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).Execute()
	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *DevicesSwitchPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *DevicesSwitchPortResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Create HTTP request body
	payload := *openApiClient.NewUpdateDeviceSwitchPortRequest()

	var tags []string
	for _, tag := range data.Tags {
		tags = append(tags, tag.ValueString())
	}
	payload.Tags = tags

	payload.Enabled = data.Enabled.ValueBoolPointer()
	payload.PoeEnabled = data.PoeEnabled.ValueBoolPointer()
	payload.Type = data.Type.ValueStringPointer()

	vlan := int32(data.Vlan.ValueInt64())
	payload.Vlan = &vlan

	voiceVlan := int32(data.VoiceVlan.ValueInt64())
	payload.VoiceVlan = &voiceVlan

	payload.AllowedVlans = data.AllowedVlans.ValueStringPointer()
	payload.IsolationEnabled = data.IsolationEnabled.ValueBoolPointer()
	payload.RstpEnabled = data.RstpEnabled.ValueBoolPointer()
	payload.StpGuard = data.StpGuard.ValueStringPointer()
	payload.LinkNegotiation = data.LinkNegotiation.ValueStringPointer()
	payload.PortScheduleId = data.PortScheduleId.ValueStringPointer()
	payload.Udld = data.Udld.ValueStringPointer()
	payload.AccessPolicyType = data.AccessPolicyType.ValueStringPointer()

	accessPolicyNumber := int32(data.AccessPolicyNumber.ValueInt64())
	payload.AccessPolicyNumber = &accessPolicyNumber

	var macAllowList []string
	for _, mac := range data.MacWhitelist {
		macAllowList = append(macAllowList, mac.ValueString())
	}
	payload.MacAllowList = macAllowList

	var stickyMacAllowList []string
	for _, mac := range data.StickyMacWhitelist {
		stickyMacAllowList = append(stickyMacAllowList, mac.ValueString())
	}
	payload.StickyMacAllowList = stickyMacAllowList

	stickyMacAllowListLimit := int32(data.StickyMacWhitelistLimit.ValueInt64())
	payload.StickyMacAllowListLimit = &stickyMacAllowListLimit

	payload.StormControlEnabled = data.StormControlEnabled.ValueBoolPointer()
	payload.AdaptivePolicyGroupId = data.AdaptivePolicyGroupId.ValueStringPointer()
	payload.PeerSgtCapable = data.PeerSgtCapable.ValueBoolPointer()
	payload.FlexibleStackingEnabled = data.FlexibleStackingEnabled.ValueBoolPointer()
	payload.DaiTrusted = data.DaiTrusted.ValueBoolPointer()

	profile := openApiClient.NewGetDeviceSwitchPorts200ResponseInnerProfile()
	profile.Id = data.Profile.Id.ValueStringPointer()
	profile.Iname = data.Profile.Iname.ValueStringPointer()
	profile.Enabled = data.Profile.Enabled.ValueBoolPointer()
	payload.Profile = profile

	// Initialize provider client and make API call
	_, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

	var tags []string
	payload.Tags = tags

	var macAllowList []string
	payload.MacAllowList = macAllowList

	var stickyMacAllowList []string
	payload.StickyMacAllowList = stickyMacAllowList

	profile := openApiClient.NewGetDeviceSwitchPorts200ResponseInnerProfile()
	payload.Profile = profile

	// Initialize provider client and make API call
	_, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
// The function expects an ImportStateRequest and responds with an ImportStateResponse which contains
// the new state of the resource or an error.
func (r *DevicesSwitchPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("serial"), req, resp)

}
