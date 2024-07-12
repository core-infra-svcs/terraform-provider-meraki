package networks

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/networks/appliance"
	switch_ "github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/networks/switch"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/networks/wireless"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var networkResources []func() resource.Resource
var networkDataSources []func() datasource.DataSource

func init() {

	// Initialize with all network resources //

	networkResources = append(networkResources, appliance.Resources()...)
	networkResources = append(networkResources, switch_.Resources()...)
	networkResources = append(networkResources, wireless.Resources()...)

	// package level resources
	networkResources = append(networkResources,
		NewNetworksCellularGatewaySubnetPoolResource,
		NewNetworksCellularGatewayUplinkResource,
		NewNetworksDevicesClaimResource,
		NewNetworksGroupPolicyResource,
		NewNetworksNetflowResource,
		NewNetworkResource,
		NewNetworksSettingsResource,
		NewNetworksSnmpResource,
		NewNetworksStormControlResource,
		NewNetworksSyslogServersResource,
		NewNetworksTrafficAnalysisResource,
	)

	// Initialize with all network data sources //

	networkDataSources = append(networkDataSources, appliance.DataSources()...)
	networkDataSources = append(networkDataSources, switch_.DataSources()...)
	networkDataSources = append(networkDataSources, wireless.DataSources()...)

	// package level data sources
	networkDataSources = append(networkDataSources,
		NewNetworkGroupPoliciesDataSource,
		NewNetworksSwitchStormControlDataSource,
	)

}

func Resources() []func() resource.Resource {
	return networkResources
}

func DataSources() []func() datasource.DataSource {
	return networkDataSources
}
