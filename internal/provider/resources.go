package provider

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/cellular"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/device"
	_interface "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/management/interface"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/port"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/ports/cycle"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance"
	applianceFirewallL3Rules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/l3/firewall/rules"
	applianceFirewallL7Rules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/l7/firewall/rules"
	applianceFirewallSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/settings"
	appliancePorts "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/ports"
	applianceVlansSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans/settings"
	applianceVlansVlan "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans/vlan"
	applianceVpn "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vpn"
	networkGroupPolicy "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/groupPolicy"
	networkSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/settings"
	merakiSwitch "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless"
	wirelessSsids "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssid"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *CiscoMerakiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		cellular.NewDevicesCellularSimsResource,
		device.NewDevicesResource,
		port.NewDevicesSwitchPortResource,
		cycle.NewDevicesSwitchPortsCycleResource,
		_interface.NewManagementInterfaceResource,

		networks.NewNetworksCellularGatewaySubnetPoolResource,
		networks.NewNetworksCellularGatewayUplinkResource,
		networks.NewNetworksDevicesClaimResource,
		networks.NewNetworksNetflowResource,
		networks.NewNetworkResource,
		networkSettings.NewNetworksSettingsResource,
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
		applianceFirewallL3Rules.NewNetworksApplianceFirewallL3FirewallRulesResource,
		applianceFirewallL7Rules.NewNetworksApplianceFirewallL7FirewallRulesResource,
		applianceFirewallSettings.NewNetworksApplianceFirewallSettingsResource,
		applianceVlansVlan.NewNetworksApplianceVLANResource,
		applianceVlansSettings.NewNetworksApplianceVlansSettingsResource,

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
