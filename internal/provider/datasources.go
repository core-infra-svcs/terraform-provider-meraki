package provider

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/administered"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/appliance"
	_switch "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/switch"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssids"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func (p *CiscoMerakiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		administered.NewAdministeredIdentitiesMeDataSource,

		devices.NewNetworkDevicesDataSource,
		devices.NewDevicesManagementInterfaceDatasource,
		devices.NewDevicesSwitchPortsStatusesDataSource,
		devices.NewDevicesApplianceDhcpSubnetsDataSource,

		networks.NewNetworkGroupPoliciesDataSource,
		networks.NewNetworksSwitchStormControlDataSource,

		appliance.NewNetworksAppliancePortsDataSource,
		appliance.NewNetworksApplianceVLANsDatasource,
		appliance.NewNetworksApplianceVlansSettingsDatasource,
		appliance.NewNetworksApplianceVpnSiteToSiteVpnDatasource,

		_switch.NewNetworksSwitchMtuDataSource,
		_switch.NewNetworksSwitchQosRulesDataSource,

		ssids.NewNetworksWirelessSsidsDataSource,

		organizations.NewOrganizationsAdaptivePolicyAclsDataSource,
		organizations.NewOrganizationsAdminsDataSource,
		organizations.NewOrganizationsCellularGatewayUplinkStatusesDataSource,
		organizations.NewOrganizationsDataSource,
		organizations.NewOrganizationsSamlIdpsDataSource,
		organizations.NewOrganizationsSamlRolesDataSource,
		organizations.NewOrganizationsInventoryDevicesDataSource,
		organizations.NewOrganizationsNetworksDataSource,
	}
}
