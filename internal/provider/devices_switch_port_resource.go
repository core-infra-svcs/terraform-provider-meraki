package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"strings"
	"time"
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
type DevicesSwitchPortResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The DevicesSwitchPortResourceModel structure describes the data model.
type DevicesSwitchPortResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	Serial                      types.String `tfsdk:"serial" json:"serial"`
	PortId                      types.String `tfsdk:"port_id" json:"portId"`
	Name                        types.String `tfsdk:"name" json:"name"`
	Tags                        types.Set    `tfsdk:"tags" json:"tags"`
	Enabled                     types.Bool   `tfsdk:"enabled" json:"enabled"`
	PoeEnabled                  types.Bool   `tfsdk:"poe_enabled" json:"poeEnabled"`
	Type                        types.String `tfsdk:"type" json:"type"`
	Vlan                        types.Int64  `tfsdk:"vlan" json:"vlan"`
	VoiceVlan                   types.Int64  `tfsdk:"voice_vlan" json:"voiceVlan"`
	AllowedVlans                types.String `tfsdk:"allowed_vlans" json:"allowedVlans"`
	AccessPolicyNumber          types.Int64  `tfsdk:"access_policy_number" json:"accessPolicyNumber"`
	AccessPolicyType            types.String `tfsdk:"access_policy_type" json:"accessPolicyType"`
	PortScheduleId              types.String `tfsdk:"port_schedule_id" json:"portScheduleId"`
	StickyMacAllowListLimit     types.Int64  `tfsdk:"sticky_mac_allow_list_limit" json:"stickyMacWhitelistLimit"`
	MacAllowList                types.Set    `tfsdk:"mac_allow_list" json:"macWhitelist"`
	StickyMacAllowList          types.Set    `tfsdk:"sticky_mac_allow_list" json:"stickyMacWhitelist"`
	StormControlEnabled         types.Bool   `tfsdk:"storm_control_enabled" json:"stormControlEnabled"`
	AdaptivePolicyGroupId       types.String `tfsdk:"adaptive_policy_group_id" json:"adaptivePolicyGroupId"`
	PeerSgtCapable              types.Bool   `tfsdk:"peer_sgt_capable" json:"peerSgtCapable"`
	FlexibleStackingEnabled     types.Bool   `tfsdk:"flexible_stacking_enabled" json:"flexibleStackingEnabled"`
	DaiTrusted                  types.Bool   `tfsdk:"dai_trusted" json:"daiTrusted"`
	IsolationEnabled            types.Bool   `tfsdk:"isolation_enabled" json:"isolationEnabled"`
	RstpEnabled                 types.Bool   `tfsdk:"rstp_enabled" json:"rstpEnabled"`
	StpGuard                    types.String `tfsdk:"stp_guard" json:"stpGuard"`
	LinkNegotiation             types.String `tfsdk:"link_negotiation" json:"linkNegotiation"`
	LinkNegotiationCapabilities types.List   `tfsdk:"link_negotiation_capabilities" json:"linkNegotiationCapabilities"`
	Udld                        types.String `tfsdk:"udld" json:"udld"`
	Profile                     types.Object `tfsdk:"profile" json:"profile"`
}

type DevicesSwitchPortResourceModelProfile struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled"`
	Id      types.String `tfsdk:"id" json:"id"`
	Iname   types.String `tfsdk:"iname" json:"iname"`
}

// Metadata provides a way to define information about the resource.
func (r *DevicesSwitchPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_devices_switch_port"
}

// Schema provides a way to define the structure of the resource data.
func (r *DevicesSwitchPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "Network Devices resource. This only works for devices associated with a network.",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed: true,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "The devices serial number",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(14, 14),
				},
			},
			"port_id": schema.StringAttribute{
				MarkdownDescription: "The identifier of the switch port.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the switch port.",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "The list of tags of the switch port.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "The status of the switch port.",
				Optional:            true,
				Computed:            true,
			},
			"poe_enabled": schema.BoolAttribute{
				MarkdownDescription: "The PoE status of the switch port.",
				Optional:            true,
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the switch port ('trunk' or 'access').",
				Optional:            true,
				Computed:            true,
			},
			"vlan": schema.Int64Attribute{
				MarkdownDescription: "The VLAN of the switch port. A null value will clear the value set for trunk ports.",
				Optional:            true,
				Computed:            true,
			},
			"voice_vlan": schema.Int64Attribute{
				MarkdownDescription: "The voice VLAN of the switch port. Only applicable to access ports.",
				Optional:            true,
				Computed:            true,
			},
			"allowed_vlans": schema.StringAttribute{
				MarkdownDescription: "The VLANs allowed on the switch port. Only applicable to trunk ports.",
				Optional:            true,
				Computed:            true,
			},
			"access_policy_type": schema.StringAttribute{
				MarkdownDescription: "The type of the access policy of the switch port. Only applicable to access ports. Can be one of 'Open', 'Custom access policy', 'MAC allow list' or 'Sticky MAC allow list'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("Open", "Custom access policy", "MAC allow list", "Sticky MAC allow list"),
				},
			},
			"access_policy_number": schema.Int64Attribute{
				MarkdownDescription: "The number of a custom access policy to configure on the switch port. Only applicable when 'accessPolicyType' is 'Custom access policy'.",
				Optional:            true,
				Computed:            true,
			},
			"port_schedule_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the port schedule. A value of null will clear the port schedule.",
				Optional:            true,
				Computed:            true,
			},
			"sticky_mac_allow_list_limit": schema.Int64Attribute{
				MarkdownDescription: "The maximum number of MAC addresses for sticky MAC allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
				Optional:            true,
				Computed:            true,
			},
			"mac_allow_list": schema.SetAttribute{
				MarkdownDescription: "Only devices with MAC addresses specified in this list will have access to this port. Up to 20 MAC addresses can be defined. Only applicable when 'accessPolicyType' is 'MAC allow list'.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"sticky_mac_allow_list": schema.SetAttribute{
				MarkdownDescription: "The initial list of MAC addresses for sticky Mac allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
			},
			"storm_control_enabled": schema.BoolAttribute{
				MarkdownDescription: "The storm control status of the switch port.",
				Optional:            true,
				Computed:            true,
			},
			"adaptive_policy_group_id": schema.StringAttribute{
				MarkdownDescription: "The adaptive policy group ID that will be used to tag traffic through this switch port. This ID must pre-exist during the configuration, else needs to be created using adaptivePolicy/groups API. Cannot be applied to a port on a switch bound to profile.",
				Optional:            true,
				Computed:            true,
			},
			"peer_sgt_capable": schema.BoolAttribute{
				MarkdownDescription: "If true, Peer SGT is enabled for traffic through this switch port. Applicable to trunk port only, not access port. Cannot be applied to a port on a switch bound to profile.",
				Optional:            true,
				Computed:            true,
			},
			"flexible_stacking_enabled": schema.BoolAttribute{
				MarkdownDescription: "For supported switches (e.g. MS420/MS425), whether or not the port has flexible stacking enabled.",
				Optional:            true,
				Computed:            true,
			},
			"dai_trusted": schema.BoolAttribute{
				MarkdownDescription: "If true, ARP packets for this port will be considered trusted, and Dynamic ARP Inspection will allow the traffic.",
				Optional:            true,
				Computed:            true,
			},
			"isolation_enabled": schema.BoolAttribute{
				MarkdownDescription: "The isolation status of the switch port.",
				Optional:            true,
				Computed:            true,
			},
			"rstp_enabled": schema.BoolAttribute{
				MarkdownDescription: "The rapid spanning tree protocol status.",
				Optional:            true,
				Computed:            true,
			},
			"stp_guard": schema.StringAttribute{
				MarkdownDescription: "The state of the STP guard ('disabled', 'root guard', 'bpdu guard' or 'loop guard').",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "root guard", "bpdu guard", "loop guard"),
				},
			},
			"link_negotiation": schema.StringAttribute{
				MarkdownDescription: "The link speed for the switch port.",
				Optional:            true,
				Computed:            true,
			},
			"link_negotiation_capabilities": schema.ListAttribute{
				MarkdownDescription: "The link speeds for the switch port.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"udld": schema.StringAttribute{
				MarkdownDescription: "The action to take when Unidirectional Link is detected (Alert only, Enforce). Default configuration is Alert only.",
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
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "When enabled, the ID of the port profile used to override the port's configuration.",
						Optional:            true,
						Computed:            true,
					},
					"iname": schema.StringAttribute{
						MarkdownDescription: "When enabled, the IName of the profile.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
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

func DevicesSwitchPortResourcePayload(ctx context.Context, plan *DevicesSwitchPortResourceModel) (openApiClient.UpdateDeviceSwitchPortRequest, diag.Diagnostics) {

	// Create HTTP request body
	payload := *openApiClient.NewUpdateDeviceSwitchPortRequest()

	//  Name
	if !plan.Name.IsUnknown() && !plan.Name.IsNull() {
		payload.SetName(plan.Name.ValueString())
	}

	//  Tags
	if !plan.Tags.IsUnknown() && !plan.Tags.IsNull() {
		var tags []string
		diags := plan.Tags.ElementsAs(ctx, &tags, false)
		if diags.HasError() {
			return payload, diags

		}
		payload.SetTags(tags)
	}

	//    Enabled
	if !plan.Enabled.IsUnknown() && !plan.Enabled.IsNull() {
		payload.SetEnabled(plan.Enabled.ValueBool())
	}

	//    PoeEnabled
	if !plan.PoeEnabled.IsUnknown() && !plan.PoeEnabled.IsNull() {
		payload.SetEnabled(plan.PoeEnabled.ValueBool())
	}

	//    Type
	if !plan.Type.IsUnknown() && !plan.Type.IsNull() {
		payload.SetType(plan.Type.ValueString())
	}

	//    Vlan
	if !plan.Vlan.IsUnknown() && !plan.Vlan.IsNull() {
		payload.SetVlan(int32(plan.Vlan.ValueInt64()))
	}

	//    VoiceVlan
	if !plan.VoiceVlan.IsUnknown() && !plan.VoiceVlan.IsNull() {
		payload.SetVoiceVlan(int32(plan.VoiceVlan.ValueInt64()))
	}

	//    AllowedVlans
	if !plan.AllowedVlans.IsUnknown() && !plan.AllowedVlans.IsNull() {
		payload.SetAllowedVlans(plan.AllowedVlans.ValueString())
	}

	//    IsolationEnabled
	if !plan.IsolationEnabled.IsUnknown() && !plan.IsolationEnabled.IsNull() {
		payload.SetIsolationEnabled(plan.IsolationEnabled.ValueBool())
	}

	//    RstpEnabled
	if !plan.RstpEnabled.IsUnknown() && !plan.RstpEnabled.IsNull() {
		payload.SetRstpEnabled(plan.RstpEnabled.ValueBool())
	}

	//    StpGuard
	if !plan.StpGuard.IsUnknown() && !plan.StpGuard.IsNull() {
		payload.SetStpGuard(plan.StpGuard.ValueString())
	}

	//    LinkNegotiation
	if !plan.LinkNegotiation.IsUnknown() && !plan.LinkNegotiation.IsNull() {
		payload.SetLinkNegotiation(plan.LinkNegotiation.ValueString())
	}

	//    PortScheduleId
	if !plan.PortScheduleId.IsUnknown() && !plan.PortScheduleId.IsNull() {
		payload.SetPortScheduleId(plan.PortScheduleId.ValueString())
	}

	//    Udld
	if !plan.Udld.IsUnknown() && !plan.Udld.IsNull() {
		payload.SetUdld(plan.Udld.ValueString())
	}

	//    AccessPolicyType
	if !plan.AccessPolicyType.IsUnknown() && !plan.AccessPolicyType.IsNull() {
		payload.SetAccessPolicyType(plan.AccessPolicyType.ValueString())
	}

	//    AccessPolicyNumber
	if !plan.AccessPolicyNumber.IsUnknown() && !plan.AccessPolicyNumber.IsNull() {
		payload.SetAccessPolicyNumber(int32(plan.AccessPolicyNumber.ValueInt64()))
	}

	//    MacAllowList
	if !plan.MacAllowList.IsUnknown() && !plan.MacAllowList.IsNull() {
		var macAllowList []string
		diags := plan.MacAllowList.ElementsAs(ctx, &macAllowList, false)
		if diags.HasError() {
			return payload, diags

		}
		payload.SetMacAllowList(macAllowList)
	}

	//    StickyMacAllowList
	if !plan.StickyMacAllowList.IsUnknown() && !plan.StickyMacAllowList.IsNull() {
		var stickyMacAllowList []string
		diags := plan.StickyMacAllowList.ElementsAs(ctx, &stickyMacAllowList, false)
		if diags.HasError() {
			return payload, diags

		}
		payload.SetStickyMacAllowList(stickyMacAllowList)
	}

	//    StickyMacAllowListLimit
	if !plan.StickyMacAllowListLimit.IsUnknown() && !plan.StickyMacAllowListLimit.IsNull() {
		payload.SetStickyMacAllowListLimit(int32(plan.StickyMacAllowListLimit.ValueInt64()))
	}

	//    StormControlEnabled
	if !plan.StormControlEnabled.IsUnknown() && !plan.StormControlEnabled.IsNull() {
		payload.SetStormControlEnabled(plan.StormControlEnabled.ValueBool())
	}

	//    AdaptivePolicyGroupId
	if !plan.AdaptivePolicyGroupId.IsUnknown() && !plan.AdaptivePolicyGroupId.IsNull() {
		payload.SetAdaptivePolicyGroupId(plan.AdaptivePolicyGroupId.ValueString())
	}

	//    PeerSgtCapable
	if !plan.PeerSgtCapable.IsUnknown() && !plan.PeerSgtCapable.IsNull() {
		payload.SetPeerSgtCapable(plan.PeerSgtCapable.ValueBool())
	}

	//    FlexibleStackingEnabled
	if !plan.FlexibleStackingEnabled.IsUnknown() && !plan.FlexibleStackingEnabled.IsNull() {
		payload.SetFlexibleStackingEnabled(plan.FlexibleStackingEnabled.ValueBool())
	}

	//    DaiTrusted
	if !plan.DaiTrusted.IsUnknown() && !plan.DaiTrusted.IsNull() {
		payload.SetDaiTrusted(plan.DaiTrusted.ValueBool())
	}

	//    Profile
	if !plan.Profile.IsUnknown() && !plan.Profile.IsNull() {

		var profile openApiClient.GetDeviceSwitchPorts200ResponseInnerProfile
		var profileData DevicesSwitchPortResourceModelProfile

		plan.Profile.As(ctx, &profileData, basetypes.ObjectAsOptions{})

		profile.SetEnabled(profileData.Enabled.ValueBool())
		profile.SetId(profileData.Id.ValueString())
		profile.SetIname(profileData.Iname.ValueString())

		payload.SetProfile(profile)
	}

	return payload, nil
}

func DevicesSwitchPortResourceResponse(ctx context.Context, response *openApiClient.GetDeviceSwitchPorts200ResponseInner, state *DevicesSwitchPortResourceModel) (*DevicesSwitchPortResourceModel, diag.Diagnostics) {

	//   import Id
	if !state.PortId.IsNull() || !state.PortId.IsUnknown() && !state.Serial.IsNull() || !state.Serial.IsUnknown() {
		state.Id = types.StringValue(state.Serial.ValueString() + "," + state.PortId.ValueString())
	} else {
		state.Id = types.StringNull()
	}

	//    Name
	if state.Name.IsNull() || state.Name.IsUnknown() {
		state.Name = types.StringValue(response.GetName())
	}

	//    Tags
	if state.Tags.IsNull() || state.Tags.IsUnknown() {
		var tags []types.String
		for _, element := range response.GetTags() {
			tags = append(tags, types.StringValue(element))
		}
		tagValues, diags := types.SetValueFrom(ctx, types.StringType, tags)
		if diags.HasError() {
			return state, diags
		}
		state.Tags = tagValues
	}

	//    Enabled
	if state.Enabled.IsNull() || state.Enabled.IsUnknown() {
		state.Enabled = types.BoolValue(response.GetEnabled())
	}

	//    PoeEnabled
	if state.PoeEnabled.IsNull() || state.PoeEnabled.IsUnknown() {

		state.PoeEnabled = types.BoolValue(response.GetPoeEnabled())
	}

	//    Type
	if state.Type.IsNull() || state.Type.IsUnknown() {
		state.Type = types.StringValue(response.GetType())
	}

	//    Vlan
	if state.Vlan.IsNull() || state.Vlan.IsUnknown() {
		state.Vlan = types.Int64Value(int64(response.GetVlan()))
	}

	//    VoiceVlan
	if state.VoiceVlan.IsNull() || state.VoiceVlan.IsUnknown() {
		state.VoiceVlan = types.Int64Value(int64(response.GetVoiceVlan()))
	}

	//    AllowedVlans
	if state.AllowedVlans.IsNull() || state.AllowedVlans.IsUnknown() {
		state.AllowedVlans = types.StringValue(response.GetAllowedVlans())
	}

	//    IsolationEnabled
	if state.IsolationEnabled.IsNull() || state.IsolationEnabled.IsUnknown() {
		state.IsolationEnabled = types.BoolValue(response.GetIsolationEnabled())
	}

	//    RstpEnabled
	if state.RstpEnabled.IsNull() || state.RstpEnabled.IsUnknown() {
		state.RstpEnabled = types.BoolValue(response.GetRstpEnabled())
	}

	//    StpGuard
	if state.StpGuard.IsNull() || state.StpGuard.IsUnknown() {
		state.StpGuard = types.StringValue(response.GetStpGuard())
	}

	//    LinkNegotiation
	if state.LinkNegotiation.IsNull() || state.LinkNegotiation.IsUnknown() {
		state.LinkNegotiation = types.StringValue(response.GetLinkNegotiation())
	}

	//    LinkNegotiationCapabilities
	if state.LinkNegotiationCapabilities.IsNull() || state.LinkNegotiationCapabilities.IsUnknown() {
		var linkNegotiationCapabilities []types.String
		for _, element := range response.GetLinkNegotiationCapabilities() {
			linkNegotiationCapabilities = append(linkNegotiationCapabilities, types.StringValue(element))
		}
		listValue, diags := types.ListValueFrom(ctx, types.StringType, linkNegotiationCapabilities)
		if diags.HasError() {
			return state, diags
		}
		state.LinkNegotiationCapabilities = listValue
	}

	//    PortScheduleId
	if state.PortScheduleId.IsNull() || state.PortScheduleId.IsUnknown() {
		state.PortScheduleId = types.StringValue(response.GetPortScheduleId())
	}

	//    Udld
	if state.Udld.IsNull() || state.Udld.IsUnknown() {
		state.Udld = types.StringValue(response.GetUdld())
	}

	//    AccessPolicyType
	if state.AccessPolicyType.IsNull() || state.AccessPolicyType.IsUnknown() {
		state.AccessPolicyType = types.StringValue(response.GetAccessPolicyType())
	}

	//    AccessPolicyNumber
	if state.AccessPolicyNumber.IsNull() || state.AccessPolicyNumber.IsUnknown() {
		state.AccessPolicyNumber = types.Int64Value(int64(response.GetAccessPolicyNumber()))
	}

	//    MacAllowList
	if state.MacAllowList.IsNull() || state.MacAllowList.IsUnknown() {

		var macAllowList []types.String
		for _, element := range response.GetMacAllowList() {
			macAllowList = append(macAllowList, types.StringValue(element))
		}
		macAllowListValues, diags := types.SetValueFrom(ctx, types.StringType, macAllowList)
		if diags.HasError() {
			return state, diags
		}
		state.MacAllowList = macAllowListValues

	}

	//    StickyMacAllowList
	if state.StickyMacAllowList.IsNull() || state.StickyMacAllowListLimit.IsUnknown() {
		var stickyMacAllowList []types.String
		for _, element := range response.GetStickyMacAllowList() {
			stickyMacAllowList = append(stickyMacAllowList, types.StringValue(element))
		}
		stickyMacAllowListValues, diags := types.SetValueFrom(ctx, types.StringType, stickyMacAllowList)
		if diags.HasError() {
			return state, diags
		}
		state.StickyMacAllowList = stickyMacAllowListValues
	}

	//    StickyMacAllowListLimit
	if state.StickyMacAllowListLimit.IsNull() || state.StickyMacAllowListLimit.IsUnknown() {
		state.StickyMacAllowListLimit = types.Int64Value(int64(response.GetStickyMacAllowListLimit()))
	}

	//    StormControlEnabled
	if state.StormControlEnabled.IsNull() || state.StormControlEnabled.IsUnknown() {
		state.StormControlEnabled = types.BoolValue(response.GetStormControlEnabled())
	}

	//    AdaptivePolicyGroupId
	if state.AdaptivePolicyGroupId.IsNull() || state.AdaptivePolicyGroupId.IsUnknown() {
		state.AdaptivePolicyGroupId = types.StringValue(response.GetAdaptivePolicyGroupId())
	}

	//    PeerSgtCapable
	if state.PeerSgtCapable.IsNull() || state.PeerSgtCapable.IsUnknown() {
		state.PeerSgtCapable = types.BoolValue(response.GetPeerSgtCapable())
	}

	//    FlexibleStackingEnabled
	if state.FlexibleStackingEnabled.IsNull() || state.FlexibleStackingEnabled.IsUnknown() {
		state.FlexibleStackingEnabled = types.BoolValue(response.GetFlexibleStackingEnabled())
	}

	//    DaiTrusted
	if state.DaiTrusted.IsNull() || state.DaiTrusted.IsUnknown() {
		state.DaiTrusted = types.BoolValue(response.GetDaiTrusted())
	}

	//    Profile
	if state.Profile.IsNull() || state.Profile.IsUnknown() {
		profileObjectMap := map[string]attr.Type{
			"enabled": types.BoolType,
			"id":      types.StringType,
			"iname":   types.StringType,
		}

		var profileData DevicesSwitchPortResourceModelProfile

		profileData.Enabled = types.BoolValue(response.Profile.GetEnabled())
		profileData.Id = types.StringValue(response.Profile.GetId())
		profileData.Iname = types.StringValue(response.Profile.GetIname())

		profileObjectValue, diags := types.ObjectValueFrom(ctx, profileObjectMap, profileData)
		if diags.HasError() {
			return state, diags
		}

		state.Profile = profileObjectValue

	}

	return state, nil
}

// Create method is responsible for creating a new resource.
func (r *DevicesSwitchPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesSwitchPortResourceModel
	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := DevicesSwitchPortResourcePayload(context.Background(), data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", diags))
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error) {
		inline, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
		return inline, httpResp, err
	}

	apiResp, httpResp, err := tools.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating switch port config",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

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

	data, diags = DevicesSwitchPortResourceResponse(ctx, apiResp, data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", diags))
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
func (r *DevicesSwitchPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesSwitchPortResourceModel
	var diags diag.Diagnostics
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// usage of CustomHttpRequestRetry with a strongly typed struct
	apiCall := func() (*openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error, diag.Diagnostics) {
		inline, httpResp, err := r.client.SwitchApi.GetDeviceSwitchPort(ctx, data.Serial.ValueString(), data.PortId.ValueString()).Execute()

		return inline, httpResp, err, diags
	}

	inlineResp, httpResp, err, tfDiags := tools.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		if tfDiags.HasError() {
			resp.Diagnostics.AddError("Diagnostics Errors", fmt.Sprintf(" %s", tfDiags.Errors()))
		}
		resp.Diagnostics.AddError("Error reading device switch port", fmt.Sprintf(" %s", err))

		if httpResp != nil {
			var responseBody string
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				} else {
					responseBody = fmt.Sprintf("Failed to read response body: %s", readErr)
				}
			} else {
				responseBody = "No response body"
			}
			resp.Diagnostics.AddError("Failed to create resource.",
				fmt.Sprintf("HTTP Status Code: %d, Response Body: %s\n", httpResp.StatusCode, responseBody))
		} else {
			resp.Diagnostics.AddError("HTTP Response is nil", "")
		}

		return
	}

	// Ensure inlineResp is not nil before dereferencing it
	if inlineResp == nil {
		fmt.Printf("Received nil response for device switch port: %s, port ID: %s\n", data.Serial.ValueString(), data.PortId.ValueString())
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Use typedApiResp with the correct type for further processing
	data, diags = DevicesSwitchPortResourceResponse(ctx, inlineResp, data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", diags))
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
func (r *DevicesSwitchPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *DevicesSwitchPortResourceModel
	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := DevicesSwitchPortResourcePayload(context.Background(), data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", diags))
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error) {
		inline, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
		return inline, httpResp, err
	}

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating switch port config",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
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

	data, diags = DevicesSwitchPortResourceResponse(ctx, inlineResp, data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", diags))
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")

}

// Delete function is responsible for deleting a resource.
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

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error) {
		inline, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
		return inline, httpResp, err
	}

	_, httpResp, err := tools.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating switch port config",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
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
