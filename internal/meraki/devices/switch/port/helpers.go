package port

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

func PortResourcePayload(ctx context.Context, plan *PortResourceModel) (openApiClient.UpdateDeviceSwitchPortRequest, diag.Diagnostics) {

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
		payload.SetPoeEnabled(plan.PoeEnabled.ValueBool())
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
		var profileData PortProfileModel

		plan.Profile.As(ctx, &profileData, basetypes.ObjectAsOptions{})

		profile.SetEnabled(profileData.Enabled.ValueBool())
		profile.SetId(profileData.Id.ValueString())
		profile.SetIname(profileData.Iname.ValueString())

		payload.SetProfile(profile)
	}

	return payload, nil
}

func PortResourceState(ctx context.Context, response *openApiClient.GetDeviceSwitchPorts200ResponseInner, state *PortResourceModel) (*PortResourceModel, diag.Diagnostics) {

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

		var profileData PortProfileModel

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
