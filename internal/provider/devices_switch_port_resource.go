package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
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
	Id                 jsontypes.String   `tfsdk:"id"`
	Serial             jsontypes.String   `tfsdk:"serial" json:"serial"`
	PortId             jsontypes.String   `tfsdk:"port_id" json:"portId"`
	Name               jsontypes.String   `tfsdk:"name" json:"name"`
	Tags               []jsontypes.String `tfsdk:"tags" json:"tags"`
	Enabled            jsontypes.Bool     `tfsdk:"enabled" json:"enabled"`
	PoeEnabled         jsontypes.Bool     `tfsdk:"poe_enabled" json:"poeEnabled"`
	Type               jsontypes.String   `tfsdk:"type" json:"type"`
	Vlan               jsontypes.Int64    `tfsdk:"vlan" json:"vlan"`
	VoiceVlan          jsontypes.Int64    `tfsdk:"voice_vlan" json:"voiceVlan"`
	AllowedVlans       jsontypes.String   `tfsdk:"allowed_vlans" json:"allowedVlans"`
	IsolationEnabled   jsontypes.Bool     `tfsdk:"isolation_enabled" json:"isolationEnabled"`
	RstpEnabled        jsontypes.Bool     `tfsdk:"rstp_enabled" json:"rstpEnabled"`
	StpGuard           jsontypes.String   `tfsdk:"stp_guard" json:"stpGuard"`
	AccessPolicyNumber jsontypes.Int64    `tfsdk:"access_policy_number" json:"accessPolicyNumber"`
	AccessPolicyType   jsontypes.String   `tfsdk:"access_policy_type" json:"accessPolicyType"`
	LinkNegotiation    jsontypes.String   `tfsdk:"link_negotiation" json:"linkNegotiation"`
	//LinkNegotiationCapabilities types.List                      `tfsdk:"link_negotiation_capabilities" json:"linkNegotiationCapabilities"`
	PortScheduleId          jsontypes.String                `tfsdk:"port_schedule_id" json:"portScheduleId"`
	Udld                    jsontypes.String                `tfsdk:"udld" json:"udld"`
	StickyMacAllowListLimit jsontypes.Int64                 `tfsdk:"sticky_mac_allow_list_limit" json:"stickyMacWhitelistLimit"`
	StormControlEnabled     jsontypes.Bool                  `tfsdk:"storm_control_enabled" json:"stormControlEnabled"`
	MacAllowList            jsontypes.Set[jsontypes.String] `tfsdk:"mac_allow_list" json:"macWhitelist"`
	StickyMacAllowList      jsontypes.Set[jsontypes.String] `tfsdk:"sticky_mac_allow_list" json:"stickyMacWhitelist"`
	AdaptivePolicyGroupId   jsontypes.String                `tfsdk:"adaptive_policy_group_id" json:"adaptivePolicyGroupId"`
	PeerSgtCapable          jsontypes.Bool                  `tfsdk:"peer_sgt_capable" json:"peerSgtCapable"`
	FlexibleStackingEnabled jsontypes.Bool                  `tfsdk:"flexible_stacking_enabled" json:"flexibleStackingEnabled"`
	DaiTrusted              jsontypes.Bool                  `tfsdk:"dai_trusted" json:"daiTrusted"`
	Profile                 DevicesSerialSwitchPortProfile  `tfsdk:"profile" json:"profile"`
}

type DevicesSerialSwitchPortProfile struct {
	Enabled jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	Id      jsontypes.String `tfsdk:"id" json:"id"`
	Iname   jsontypes.String `tfsdk:"iname" json:"iname"`
}

var devicesSerialSwitchPortProfile map[string]attr.Type = map[string]attr.Type{
	"enabled": types.BoolType,
	"id":      types.StringType,
	"iname":   types.StringType,
}

func MyModelGoToTerraform(ctx context.Context, my sdk.My) (types.Object, diag.Diagnostics) {
	return types.ObjectValueFrom(ctx, devicesSerialSwitchPortProfile, DevicesSerialSwitchPortProfile{
		Enabled: types.BoolValue(my.Enabled),
		Id:      types.StringValue(my.Id),
		Iname:   types.StringValue(my.Iname),
	})
}

/*
func (p *DevicesSerialSwitchPortProfile) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*p = DevicesSerialSwitchPortProfile{
		Enabled: jsontypes.BoolNull(),
		Id:      jsontypes.StringNull(),
		Iname:   jsontypes.StringNull(),
	}
	return nil
}
*/

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
				Optional:   true,
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
			"link_negotiation": schema.StringAttribute{
				MarkdownDescription: "The link speed for the switch port.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			/*
				"link_negotiation_capabilities": schema.ListAttribute{
						MarkdownDescription: "The link speeds for the switch port.",
						ElementType:         jsontypes.StringType,
						Optional:            true,
						Computed:            true,
					},
			*/
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
			"sticky_mac_allow_list_limit": schema.Int64Attribute{
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"mac_allow_list": schema.SetAttribute{
				MarkdownDescription: "Only devices with MAC addresses specified in this list will have access to this port. Up to 20 MAC addresses can be defined. Only applicable when 'accessPolicyType' is 'MAC allow list'.",
				CustomType:          jsontypes.SetType[jsontypes.String](),
				ElementType:         jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"sticky_mac_allow_list": schema.SetAttribute{
				MarkdownDescription: "The initial list of MAC addresses for sticky Mac allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
				CustomType:          jsontypes.SetType[jsontypes.String](),
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
			"profile": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "When enabled, override this port's configuration with a port profile.",
						Optional:            true,
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

	payload, err := DevicesSwitchPortResourcePayload(data)

	/*

		if !data.Profile.IsNull() && !data.Profile.IsUnknown() {

			profilePayload := openApiClient.NewGetDeviceSwitchPorts200ResponseInnerProfile()
			//profilePayload.SetIname(data.Profile.Iname.ValueString())
			//profilePayload.SetEnabled(data.Profile.Enabled.ValueBool())
			payload.SetProfile(*profilePayload)

		} else if !data.AdaptivePolicyGroupId.IsNull() && !data.AdaptivePolicyGroupId.IsUnknown() {

			// AdaptivePolicyGroupId
			payload.SetAdaptivePolicyGroupId(data.AdaptivePolicyGroupId.ValueString())
		}

	*/
	// Initialize provider client and make API call
	_, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
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

	newData, err := DevicesSwitchPortResourceResponse(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create JSON Decode issue",
			fmt.Sprintf("%v", err.Error()),
		)
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)

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
	_, httpResp, err := r.client.SwitchApi.GetDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).Execute()
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

	payload, err := DevicesSwitchPortResourcePayload(data)

	/*
		if !data.Profile.IsNull() && !data.Profile.IsUnknown() {

				profilePayload := openApiClient.NewGetDeviceSwitchPorts200ResponseInnerProfile()
				//profilePayload.SetIname(data.Profile.Iname.ValueString())
				//profilePayload.SetEnabled(data.Profile.Enabled.ValueBool())
				payload.SetProfile(*profilePayload)

			} else if !data.AdaptivePolicyGroupId.IsNull() && !data.AdaptivePolicyGroupId.IsUnknown() {

				// AdaptivePolicyGroupId
				payload.SetAdaptivePolicyGroupId(data.AdaptivePolicyGroupId.ValueString())
			}

	*/

	// Initialize provider client and make API call
	_, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"UPDATE HTTP Client Failure",
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

	newData, err := DevicesSwitchPortResourceResponse(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", err.Error()),
		)
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &newData)...)

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

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("serial"), req, resp)

}

func DevicesSwitchPortResourcePayload(data *DevicesSwitchPortResourceModel) (openApiClient.UpdateDeviceSwitchPortRequest, error) {

	// Create HTTP request body
	payload := *openApiClient.NewUpdateDeviceSwitchPortRequest()

	payload.SetEnabled(data.Enabled.ValueBool())
	payload.SetName(data.Name.ValueString())
	payload.SetPoeEnabled(data.PoeEnabled.ValueBool())
	payload.SetType(data.Type.ValueString())
	payload.SetIsolationEnabled(data.IsolationEnabled.ValueBool())
	payload.SetRstpEnabled(data.RstpEnabled.ValueBool())
	payload.SetStpGuard(data.StpGuard.ValueString())
	payload.SetLinkNegotiation(data.LinkNegotiation.ValueString())
	payload.SetUdld(data.Udld.ValueString())

	if !data.AccessPolicyType.IsUnknown() && !data.AccessPolicyType.IsNull() {
		payload.SetAccessPolicyType(data.AccessPolicyType.ValueString())
	}

	payload.SetDaiTrusted(data.DaiTrusted.ValueBool())

	if !data.Vlan.IsUnknown() && !data.Vlan.IsNull() {
		payload.SetVlan(int32(data.Vlan.ValueInt64()))
	}

	// PortScheduleId
	if !data.PortScheduleId.IsNull() && !data.PortScheduleId.IsUnknown() {
		payload.SetPortScheduleId(data.PortScheduleId.ValueString())
	}

	// StormControlEnabled
	if data.StormControlEnabled.ValueBool() == true {
		payload.SetStormControlEnabled(data.StormControlEnabled.ValueBool())
	}

	// FlexibleStackingEnabled
	if data.FlexibleStackingEnabled.ValueBool() == true {
		payload.SetFlexibleStackingEnabled(data.FlexibleStackingEnabled.ValueBool())
	}

	// Tags
	var tags []string
	for _, tag := range data.Tags {
		tags = append(tags, tag.ValueString())
	}
	payload.SetTags(tags)

	// Trunk Port Settings
	if data.Type.ValueString() == "trunk" {

		// Allowed VLANS
		if !data.AllowedVlans.IsUnknown() && !data.AllowedVlans.IsNull() {
			payload.SetAllowedVlans(data.AllowedVlans.ValueString())
		}

		// PeerSgtCapable
		if !data.PeerSgtCapable.IsNull() && !data.PeerSgtCapable.IsUnknown() {
			payload.SetPeerSgtCapable(data.PeerSgtCapable.ValueBool())
		}

	}

	// Access Port Settings
	if data.Type.ValueString() == "access" {

		// Voice VLAN
		if !data.VoiceVlan.IsNull() {
			payload.SetVoiceVlan(int32(data.VoiceVlan.ValueInt64()))
		}

	}

	// AccessPolicyType
	if data.AccessPolicyType.ValueString() == "Custom access policy" {

		// AccessPolicyNumber
		payload.SetAccessPolicyNumber(int32(data.AccessPolicyNumber.ValueInt64()))

	} else if data.AccessPolicyType.ValueString() == "MAC allow list" {

		// MacAllowList
		var macAllowList []string
		for _, mac := range data.MacAllowList.Elements() {
			macAllowList = append(macAllowList, mac.String())
		}
		payload.SetMacAllowList(macAllowList)

	} else if data.AccessPolicyType.ValueString() == "Sticky MAC allow list" {

		// StickyMacAllowList
		var stickyMacAllowList []string
		for _, mac := range data.StickyMacAllowList.Elements() {
			stickyMacAllowList = append(stickyMacAllowList, mac.String())
		}
		payload.SetStickyMacAllowList(stickyMacAllowList)

		// StickyMacAllowListLimit
		payload.SetStickyMacAllowListLimit(int32(data.StickyMacAllowListLimit.ValueInt64()))
	}
	return payload, nil
}

func DevicesSwitchPortResourceResponse(httpRespBody io.ReadCloser) (DevicesSwitchPortResourceModel, error) {
	var data DevicesSwitchPortResourceModel
	var profile DevicesSerialSwitchPortProfile

	data.Profile = profile

	if err := json.NewDecoder(httpRespBody).Decode(&data); err != nil {
		return data, err
	}

	if data.Profile.Enabled.IsUnknown() || data.Profile.Enabled.IsNull() {
		data.Profile.Enabled = jsontypes.BoolNull()
	}

	if data.Profile.Id.IsUnknown() || data.Profile.Id.IsNull() {
		data.Profile.Id = jsontypes.StringNull()
	}

	if data.Profile.Iname.IsUnknown() || data.Profile.Iname.IsNull() {
		data.Profile.Iname = jsontypes.StringNull()
	}

	if data.AccessPolicyNumber.IsUnknown() {
		data.AccessPolicyNumber = jsontypes.Int64Null()
	}

	if data.AdaptivePolicyGroupId.IsUnknown() {
		data.AdaptivePolicyGroupId = jsontypes.StringNull()
	}

	if data.FlexibleStackingEnabled.IsUnknown() {
		data.FlexibleStackingEnabled = jsontypes.BoolNull()
	}

	if data.PeerSgtCapable.IsUnknown() {
		data.PeerSgtCapable = jsontypes.BoolNull()
	}

	if data.StickyMacAllowList.IsUnknown() {
		var stickyMacAllowList jsontypes.Set[jsontypes.String]
		data.StickyMacAllowList = stickyMacAllowList
	}

	if data.MacAllowList.IsUnknown() {
		var macAllowList jsontypes.Set[jsontypes.String]
		data.MacAllowList = macAllowList
	}

	if data.StickyMacAllowListLimit.IsUnknown() {
		data.StickyMacAllowListLimit = jsontypes.Int64Null()
	}

	if data.StormControlEnabled.IsUnknown() {
		data.StormControlEnabled = jsontypes.BoolNull()
	}

	data.Id = jsontypes.StringValue("example-id")

	return data, nil
}
