package provider

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/administered"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/dhcp/subnets"
	_interface "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/management/interface"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/switch/ports"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/firewall/l3/firewall/rules"
	appliancePorts "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/ports"
	applianceVlans "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vlans"
	applianceVpn "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance/vpn"
	networkGroupPolicy "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/groupPolicy"
	merakiSwitch "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch"
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

		networkGroupPolicy.NewNetworkGroupPoliciesDataSource,
		networks.NewNetworksSwitchStormControlDataSource,

		appliancePorts.NewNetworksAppliancePortsDataSource,
		applianceVlans.NewNetworksApplianceVLANsDatasource,
		applianceVlans.NewNetworksApplianceVlansSettingsDatasource,
		applianceVpn.NewNetworksApplianceVpnSiteToSiteVpnDatasource,
		rules.NewNetworksApplianceFirewallL3FirewallRulesDataSource,

		merakiSwitch.NewNetworksSwitchMtuDataSource,
		merakiSwitch.NewNetworksSwitchQosRulesDataSource,

		wirelessSsids.NewNetworksWirelessSsidsDataSource,

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
