package ports

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

func mapSwitchDataToPort(switchData openApiClient.GetDeviceSwitchPorts200ResponseInner) PortResourceModel {
	var devicesSwitchPortData PortResourceModel
	devicesSwitchPortData.Name = types.StringValue(switchData.GetName())
	devicesSwitchPortData.PortId = types.StringValue(switchData.GetPortId())
	devicesSwitchPortData.Enabled = types.BoolValue(switchData.GetEnabled())
	devicesSwitchPortData.PoeEnabled = types.BoolValue(switchData.GetPoeEnabled())
	devicesSwitchPortData.Type = types.StringValue(switchData.GetType())
	devicesSwitchPortData.Vlan = types.Int64Value(int64(switchData.GetVlan()))
	devicesSwitchPortData.VoiceVlan = types.Int64Value(int64(switchData.GetVoiceVlan()))
	devicesSwitchPortData.AllowedVlans = types.StringValue(switchData.GetAllowedVlans())
	devicesSwitchPortData.IsolationEnabled = types.BoolValue(switchData.GetIsolationEnabled())
	devicesSwitchPortData.RstpEnabled = types.BoolValue(switchData.GetRstpEnabled())
	devicesSwitchPortData.StpGuard = types.StringValue(switchData.GetStpGuard())
	devicesSwitchPortData.AccessPolicyNumber = types.Int64Value(int64(switchData.GetAccessPolicyNumber()))
	devicesSwitchPortData.AccessPolicyType = types.StringValue(switchData.GetAccessPolicyType())
	devicesSwitchPortData.LinkNegotiation = types.StringValue(switchData.GetLinkNegotiation())
	devicesSwitchPortData.PortScheduleId = types.StringValue(switchData.GetPortScheduleId())
	devicesSwitchPortData.Udld = types.StringValue(switchData.GetUdld())
	devicesSwitchPortData.StickyMacWhitelistLimit = types.Int64Value(int64(switchData.GetStickyMacAllowListLimit()))
	devicesSwitchPortData.StormControlEnabled = types.BoolValue(switchData.GetStormControlEnabled())
	devicesSwitchPortData.AdaptivePolicyGroupId = types.StringValue(switchData.GetAdaptivePolicyGroupId())
	devicesSwitchPortData.PeerSgtCapable = types.BoolValue(switchData.GetPeerSgtCapable())
	devicesSwitchPortData.FlexibleStackingEnabled = types.BoolValue(switchData.GetFlexibleStackingEnabled())
	devicesSwitchPortData.DaiTrusted = types.BoolValue(switchData.GetDaiTrusted())
	devicesSwitchPortData.Profile.Enabled = types.BoolValue(switchData.Profile.GetEnabled())
	devicesSwitchPortData.Profile.Id = types.StringValue(switchData.Profile.GetId())
	devicesSwitchPortData.Profile.Iname = types.StringValue(switchData.Profile.GetIname())

	for _, attribute := range switchData.GetStickyMacAllowList() {
		devicesSwitchPortData.StickyMacWhitelist = append(devicesSwitchPortData.StickyMacWhitelist, types.StringValue(attribute))
	}
	for _, attribute := range switchData.GetTags() {
		devicesSwitchPortData.Tags = append(devicesSwitchPortData.Tags, types.StringValue(attribute))
	}
	for _, attribute := range switchData.GetMacAllowList() {
		devicesSwitchPortData.MacWhitelist = append(devicesSwitchPortData.MacWhitelist, types.StringValue(attribute))
	}

	return devicesSwitchPortData
}
