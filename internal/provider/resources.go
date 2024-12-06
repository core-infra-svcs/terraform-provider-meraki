package provider

import (
	"context"
	devicesCellular "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/cellular"
	devicesDevice "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/device"
	devicesManagementInterface "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/management/interface"
	devicesSwitchPort "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/port"
	devicesSwitchPortsCycle "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/ports/cycle"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks"
	networksAppliance "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance"
	networksApplianceFirewallL3Rules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/l3/firewall/rules"
	networksApplianceFirewallL7Rules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/l7/firewall/rules"
	networksApplianceFirewallSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/settings"
	networksAppliancePorts "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/ports"
	networksApplianceSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/settings"
	networksApplianceVlansSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans/settings"
	networksApplianceVlansVlan "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans/vlan"
	networksApplianceVpn "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vpn"
	networksGroupPolicy "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/groupPolicy"
	networksSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/settings"
	networksSwitch "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch"
	networksWireless "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless"
	networksWirelessSsids "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssid"
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
		networksSettings.NewNetworksSettingsResource,
		networks.NewNetworksSnmpResource,
		networks.NewNetworksStormControlResource,
		networks.NewNetworksSyslogServersResource,
		networks.NewNetworksTrafficAnalysisResource,
		networksGroupPolicy.NewNetworksGroupPolicyResource,

		networksAppliancePorts.NewNetworksAppliancePortsResource,
		networksApplianceSettings.NewNetworksApplianceSettingsResource,
		networksAppliance.NewNetworkApplianceStaticRoutesResource,
		networksAppliance.NewNetworksApplianceTrafficShapingUplinkBandWidthResource,
		networksApplianceVpn.NewNetworksApplianceVpnSiteToSiteVpnResource,
		networksApplianceFirewallL3Rules.NewNetworksApplianceFirewallL3FirewallRulesResource,
		networksApplianceFirewallL7Rules.NewNetworksApplianceFirewallL7FirewallRulesResource,
		networksApplianceFirewallSettings.NewNetworksApplianceFirewallSettingsResource,
		networksApplianceVlansVlan.NewNetworksApplianceVLANResource,
		networksApplianceVlansSettings.NewNetworksApplianceVlansSettingsResource,

		networksSwitch.NewNetworksSwitchDscpToCosMappingsResource,
		networksSwitch.NewNetworksSwitchMtuResource,
		networksSwitch.NewNetworksSwitchQosRuleResource,
		networksSwitch.NewNetworksSwitchSettingsResource,

		networksWireless.NewNetworksWirelessSsidsFirewallL3FirewallRulesResource,
		networksWireless.NewNetworksWirelessSsidsFirewallL7FirewallRulesResource,
		networksWireless.NewNetworksWirelessSsidsSplashSettingsResource,
		networksWirelessSsids.NewNetworksWirelessSsidsResource,

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
