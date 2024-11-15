package provider

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices"
	switchDevices "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance"
	applianceFirewall "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall"
	appliancePorts "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/ports"
	applianceVlans "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans"
	applianceVpn "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vpn"
	networkGroupPolicy "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/groupPolicy"
	merakiSwitch "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless"
	wirelessSsids "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssid"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *CiscoMerakiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		devices.NewDevicesCellularSimsResource,
		devices.NewDevicesResource,
		switchDevices.NewDevicesSwitchPortResource,
		switchDevices.NewDevicesSwitchPortsCycleResource,

		networks.NewNetworksCellularGatewaySubnetPoolResource,
		networks.NewNetworksCellularGatewayUplinkResource,
		networks.NewNetworksDevicesClaimResource,
		networks.NewNetworksNetflowResource,
		networks.NewNetworkResource,
		networks.NewNetworksSettingsResource,
		networks.NewNetworksSnmpResource,
		networks.NewNetworksStormControlResource,
		networks.NewNetworksSyslogServersResource,
		networks.NewNetworksTrafficAnalysisResource,
		networkGroupPolicy.NewNetworksGroupPolicyResource,

		appliancePorts.NewNetworksAppliancePortsResource,
		appliance.NewNetworksApplianceSettingsResource,
		appliance.NewNetworkApplianceStaticRoutesResource,
		applianceVpn.NewNetworksApplianceVpnSiteToSiteVpnResource,
		appliance.NewNetworksApplianceTrafficShapingUplinkBandWidthResource,
		applianceFirewall.NewNetworksApplianceFirewallL3FirewallRulesResource,
		applianceFirewall.NewNetworksApplianceFirewallL7FirewallRulesResource,
		applianceFirewall.NewNetworksApplianceFirewallSettingsResource,
		applianceVlans.NewNetworksApplianceVLANResource,
		applianceVlans.NewNetworksApplianceVlansSettingsResource,

		merakiSwitch.NewNetworksSwitchDscpToCosMappingsResource,
		merakiSwitch.NewNetworksSwitchMtuResource,
		merakiSwitch.NewNetworksSwitchQosRuleResource,
		merakiSwitch.NewNetworksSwitchSettingsResource,

		wireless.NewNetworksWirelessSsidsFirewallL3FirewallRulesResource,
		wireless.NewNetworksWirelessSsidsFirewallL7FirewallRulesResource,
		wireless.NewNetworksWirelessSsidsSplashSettingsResource,
		wirelessSsids.NewNetworksWirelessSsidsResource,

		organizations.NewAdaptivePolicyAclResource,
		organizations.NewOrganizationsAdminResource,
		organizations.NewOrganizationsApplianceVpnVpnFirewallRulesResource,
		organizations.NewOrganizationsClaimResource,
		organizations.NewOrganizationsLicenseResource,
		organizations.NewOrganizationsSamlIdpResource,
		organizations.NewOrganizationSamlResource,
		organizations.NewOrganizationsSamlRolesResource,
		organizations.NewOrganizationsSnmpResource,
		organizations.NewOrganizationResource,
		organizations.NewOrganizationPolicyObjectResource,
	}
}