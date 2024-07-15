package appliance

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworksApplianceFirewallL3FirewallRulesResource,
		NewNetworksApplianceFirewallL7FirewallRulesResource,
		NewNetworksApplianceFirewallSettingsResource,
		NewNetworksAppliancePortsResource,
		NewNetworksApplianceSettingsResource,
		NewNetworkApplianceStaticRoutesResource,
		NewNetworksApplianceTrafficShapingUplinkBandWidthResource,
		NewNetworksApplianceVLANResource,
		NewNetworksApplianceVlansSettingsResource,
		NewNetworksApplianceVpnSiteToSiteVpnResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDevicesApplianceDhcpSubnetsDataSource,
		NewNetworksAppliancePortsDataSource,
		NewNetworksApplianceVLANsDatasource,
		NewNetworksApplianceVlansSettingsDatasource,
		NewNetworksApplianceVpnSiteToSiteVpnDatasource,
	}
}
