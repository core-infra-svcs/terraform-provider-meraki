package provider

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/administered"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/dhcp/subnets"
	_interface "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/management/interface"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/ports"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks"
	applianceFirewallL3Rules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/l3/firewall/rules"
	appliancePorts "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/ports"
	//applianceFirewallL7Rules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/l7/firewall/rules"
	applianceVlansSettings "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans/settings"
	applianceVlansVlan "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans/vlan"
	applianceVpn "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vpn"
	networkGroupPolicy "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/group/policy"
	networksSwitchMtu "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/mtu"
	networksSwitchQosRules "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch/qos/rules"
	wirelessSsids "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssid"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (p *CiscoMerakiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		administered.NewAdministeredIdentitiesMeDataSource,

		devices.NewNetworkDevicesDataSource,
		_interface.NewDevicesManagementInterfaceDataSource,
		ports.NewDevicesSwitchPortsStatusesDataSource,
		subnets.NewDevicesApplianceDhcpSubnetsDataSource,

		networkGroupPolicy.NewDataSource,
		networks.NewNetworksSwitchStormControlDataSource,

		appliancePorts.NewNetworksAppliancePortsDataSource,
		applianceVlansVlan.NewNetworksApplianceVLANsDatasource,
		applianceVlansSettings.NewNetworksApplianceVlansSettingsDatasource,
		applianceVpn.NewNetworksApplianceVpnSiteToSiteVpnDatasource,
		applianceFirewallL3Rules.NewNetworksApplianceFirewallL3FirewallRulesDataSource,
		//applianceFirewallL7Rules.NewNetworksApplianceFirewallL7FirewallRulesDataSource,
		networksSwitchMtu.NewDataSource,
		networksSwitchQosRules.NewDataSource,

		wirelessSsids.NewDataSource,

		organizations.NewOrganizationsAdaptivePolicyAclsDataSource,
		organizations.NewOrganizationsAdminsDataSource,
		organizations.NewOrganizationsLicensesDataSource,
		organizations.NewOrganizationsCellularGatewayUplinkStatusesDataSource,
		organizations.NewOrganizationsDataSource,
		organizations.NewOrganizationsSamlIdpsDataSource,
		organizations.NewOrganizationsSamlRolesDataSource,
		organizations.NewOrganizationsInventoryDevicesDataSource,
		organizations.NewOrganizationsNetworksDataSource,
	}
}
