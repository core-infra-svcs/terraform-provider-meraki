package provider

import (
	"context"
	devicesCellular "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/cellular"
	devicesDevice "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/device"
	devicesManagementInterface "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/management/interface"
	devicesSwitchPort "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/port"
	devicesSwitchPortsCycle "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/ports/cycle"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks"
	networksApplianceFirewallL3Rules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/l3/firewall/rules"
	networksApplianceFirewallL7Rules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/l7/firewall/rules"
	networksApplianceFirewallSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/settings"
	networksAppliancePorts "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/ports"
	networksApplianceSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/settings"
	networksApplianceStaticRoutes "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/static/routes"
	networksApplianceTrafficShapingUplinkBandWidth "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/traffic/shaping/uplink/bandwidth"
	networksApplianceVlansSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans/settings"
	networksApplianceVlansVlan "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans/vlan"
	networksApplianceVpn "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vpn"
	networksGroupPolicy "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/group/policy"
	networksSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/settings"
	networksSwitchDscpToCosMappings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/dscp/to/cos/mappings"
	networksSwitchMtu "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/mtu"
	networksSwitchQosRules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/qos/rules"
	networksSwitchSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/settings"
	networksWireless "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless"
	networksWirelessSsids "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssid"
	networksWirelessSsidsFirewallL3FirewallRules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssids/firewall/l3/firewall/rules"
	networksWirelessSsidsFirewallL7FirewallRules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssids/firewall/l7/firewall/rules"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *CiscoMerakiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		devicesCellular.NewDevicesCellularSimsResource,
		devicesDevice.NewDevicesResource,
		devicesSwitchPort.NewDevicesSwitchPortResource,
		devicesSwitchPortsCycle.NewDevicesSwitchPortsCycleResource,
		devicesManagementInterface.NewManagementInterfaceResource,

		networks.NewNetworksCellularGatewaySubnetPoolResource,
		networks.NewNetworksCellularGatewayUplinkResource,
		networks.NewNetworksDevicesClaimResource,
		networks.NewNetworksNetflowResource,
		networks.NewNetworkResource,
		networksSettings.NewResource,
		networks.NewNetworksSnmpResource,
		networks.NewNetworksStormControlResource,
		networks.NewNetworksSyslogServersResource,
		networks.NewNetworksTrafficAnalysisResource,
		networksGroupPolicy.NewResource,

		networksAppliancePorts.NewNetworksAppliancePortsResource,
		networksApplianceSettings.NewNetworksApplianceSettingsResource,
		networksApplianceStaticRoutes.NewNetworkApplianceStaticRoutesResource,
		networksApplianceTrafficShapingUplinkBandWidth.NewResource, // this pattern fyi
		networksApplianceVpn.NewNetworksApplianceVpnSiteToSiteVpnResource,
		networksApplianceFirewallL3Rules.NewNetworksApplianceFirewallL3FirewallRulesResource,
		networksApplianceFirewallL7Rules.NewNetworksApplianceFirewallL7FirewallRulesResource,
		networksApplianceFirewallSettings.NewNetworksApplianceFirewallSettingsResource,
		networksApplianceVlansVlan.NewNetworksApplianceVLANResource,
		networksApplianceVlansSettings.NewNetworksApplianceVlansSettingsResource,

		networksSwitchDscpToCosMappings.NewResource,
		networksSwitchMtu.NewResource,
		networksSwitchQosRules.NewResource,
		networksSwitchSettings.NewResource,

		networksWirelessSsidsFirewallL3FirewallRules.NewResource,
		networksWirelessSsidsFirewallL7FirewallRules.NewResource,
		networksWireless.NewNetworksWirelessSsidsSplashSettingsResource,
		networksWirelessSsids.NewResource,

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
