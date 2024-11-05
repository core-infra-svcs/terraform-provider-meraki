package provider

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance"
	_switch "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssids"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *CiscoMerakiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		devices.NewDevicesCellularSimsResource,
		devices.NewDevicesResource,
		devices.NewDevicesSwitchPortResource,
		devices.NewDevicesSwitchPortsCycleResource,

		networks.NewNetworksCellularGatewaySubnetPoolResource,
		networks.NewNetworksCellularGatewayUplinkResource,
		networks.NewNetworksDevicesClaimResource,
		networks.NewNetworksGroupPolicyResource,
		networks.NewNetworksNetflowResource,
		networks.NewNetworkResource,
		networks.NewNetworksSettingsResource,
		networks.NewNetworksSnmpResource,
		networks.NewNetworksStormControlResource,
		networks.NewNetworksSyslogServersResource,
		networks.NewNetworksTrafficAnalysisResource,

		appliance.NewNetworksApplianceFirewallL3FirewallRulesResource,
		appliance.NewNetworksApplianceFirewallL7FirewallRulesResource,
		appliance.NewNetworksApplianceFirewallSettingsResource,
		appliance.NewNetworksAppliancePortsResource,
		appliance.NewNetworksApplianceSettingsResource,
		appliance.NewNetworkApplianceStaticRoutesResource,
		appliance.NewNetworksApplianceTrafficShapingUplinkBandWidthResource,
		appliance.NewNetworksApplianceVLANResource,
		appliance.NewNetworksApplianceVlansSettingsResource,
		appliance.NewNetworksApplianceVpnSiteToSiteVpnResource,

		_switch.NewNetworksSwitchDscpToCosMappingsResource,
		_switch.NewNetworksSwitchMtuResource,
		_switch.NewNetworksSwitchQosRuleResource,
		_switch.NewNetworksSwitchSettingsResource,

		wireless.NewNetworksWirelessSsidsFirewallL3FirewallRulesResource,
		wireless.NewNetworksWirelessSsidsFirewallL7FirewallRulesResource,
		wireless.NewNetworksWirelessSsidsSplashSettingsResource,

		ssids.NewNetworksWirelessSsidsResource,

		organizations.NewAdaptivePolicyAclResource,
		organizations.NewOrganizationsAdminResource,
		organizations.NewOrganizationsApplianceVpnVpnFirewallRulesResource,
		organizations.NewOrganizationsClaimResource,
		organizations.NewOrganizationsSamlIdpResource,
		organizations.NewOrganizationSamlResource,
		organizations.NewOrganizationsSamlRolesResource,
		organizations.NewOrganizationsSnmpResource,
		organizations.NewOrganizationResource,
		organizations.NewOrganizationPolicyObjectResource,
	}
}
