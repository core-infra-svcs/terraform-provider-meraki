package provider

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/administered"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices"
	devicesCellular "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/cellular"
	devicesDevice "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/device"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/dhcp/subnets"
	devicesManagementInterface "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/management/interface"
	devicesSwitchPort "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/port"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/ports"
	devicesSwitchPortsCycle "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/ports/cycle"
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
	networksCellularGatewaySubnetPool "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/cellular/gateway/subnet/pool"
	networksCellularGatewayUplink "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/cellular/gateway/uplink"
	networksDevicesClaim "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/devices/claim"
	networksGroupPolicy "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/group/policy"
	networksNetflow "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/netflow"
	networksNetwork "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/network"
	networksSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/settings"
	networksSnmp "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/snmp"
	networksStormControl "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/storm/control"
	networksSwitchDscpToCosMappings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/dscp/to/cos/mappings"
	networksSwitchMtu "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/mtu"
	networksSwitchQosRules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/qos/rules"
	networksSwitchSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/settings"
	networksSyslogServers "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/syslog/servers"
	networksTrafficAnalysis "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/traffic/analysis"
	networksWirelessSsids "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssid"
	networksWirelessSsidsFirewallL3FirewallRules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssids/firewall/l3/firewall/rules"
	networksWirelessSsidsFirewallL7FirewallRules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssids/firewall/l7/firewall/rules"
	networksWirelessSsidsSplashSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssids/splash/settings"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations"
	organizationsAdaptivePolicyAcls "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/adaptive/policy/acls"
	organizationsAdmins "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/admins"
	organizationsApplianceVpnFirewallRules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/appliance/vpn/firewall/rules"
	organizationsCellularGatewayUplinkStatuses "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/cellular/gateway/uplink/statuses"
	organizationsClaim "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/claim"
	organizationsInventoryDevices "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/inventory/devices"
	organizationsLicences "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/licences"
	organizationsLicencesMove "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/licences/move"
	organizationsNetworks "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/networks"
	organizationsOrganization "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/organization"
	organizationsPolicyObject "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/policy/object"
	organizationsSamlIdps "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations/saml/idps"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (p *CiscoMerakiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		devicesCellular.NewResource,
		devicesDevice.NewResource,
		devicesSwitchPort.NewResource,
		devicesSwitchPortsCycle.NewResource,
		devicesManagementInterface.NewResource,
		networksCellularGatewaySubnetPool.NewResource,
		networksCellularGatewayUplink.NewResource,
		networksDevicesClaim.NewResource,
		networksNetflow.NewResource,
		networksNetwork.NewResource,
		networksSettings.NewResource,
		networksSnmp.NewResource,
		networksStormControl.NewResource,
		networksSyslogServers.NewResource,
		networksTrafficAnalysis.NewResource,
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
		networksWirelessSsidsSplashSettings.NewResource,
		networksWirelessSsids.NewResource,
		organizationsAdaptivePolicyAcls.NewResource,
		organizationsAdmins.NewResource,
		organizationsApplianceVpnFirewallRules.NewResource,
		organizationsClaim.NewResource,
		organizationsLicencesMove.NewResource,
		organizationsSamlIdps.NewResource,
		organizations.NewOrganizationSamlResource,
		organizations.NewOrganizationsSamlRolesResource,
		organizations.NewOrganizationsSnmpResource,
		organizationsOrganization.NewResource,
		organizationsPolicyObject.NewResource,
	}
}

func (p *CiscoMerakiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		administered.NewAdministeredIdentitiesMeDataSource,

		devices.NewNetworkDevicesDataSource,
		devicesManagementInterface.NewDevicesManagementInterfaceDataSource,
		ports.NewDevicesSwitchPortsStatusesDataSource,
		subnets.NewDevicesApplianceDhcpSubnetsDataSource,
		networksGroupPolicy.NewDataSource,
		networksStormControl.NewDataSource,
		networksAppliancePorts.NewNetworksAppliancePortsDataSource,
		networksApplianceVlansVlan.NewNetworksApplianceVLANsDatasource,
		networksApplianceVlansSettings.NewNetworksApplianceVlansSettingsDatasource,
		networksApplianceVpn.NewNetworksApplianceVpnSiteToSiteVpnDatasource,
		networksApplianceFirewallL3Rules.NewNetworksApplianceFirewallL3FirewallRulesDataSource,
		//networksApplianceFirewallL3Rules.NewNetworksApplianceFirewallL7FirewallRulesDataSource,
		networksSwitchMtu.NewDataSource,
		networksSwitchQosRules.NewDataSource,
		networksWirelessSsids.NewDataSource,
		organizationsAdaptivePolicyAcls.NewDataSource,
		organizationsAdmins.NewDataSource,
		organizationsLicences.NewDataSource,
		organizationsCellularGatewayUplinkStatuses.NewDataSource,
		organizationsOrganization.NewDataSource,
		organizationsSamlIdps.NewDataSource,
		organizations.NewOrganizationsSamlRolesDataSource,
		organizationsInventoryDevices.NewDataSource,
		organizationsNetworks.NewDataSource,
	}
}
