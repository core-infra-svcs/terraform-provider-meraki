package provider

import (
	merakiAdministered "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/administered"
	merakiDevices "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices"
	merakiNetworks "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks"
	merakiOrganizations "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/organizations"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var MerakiResources []func() resource.Resource
var MerakiDataSources []func() datasource.DataSource

func init() {

	// Initialize all provider resources //
	MerakiResources = append(MerakiResources, merakiAdministered.Resources()...)
	MerakiResources = append(MerakiResources, merakiDevices.Resources()...)
	MerakiResources = append(MerakiResources, merakiOrganizations.Resources()...)
	MerakiResources = append(MerakiResources, merakiNetworks.Resources()...)

	// Initialize all provider data sources //
	MerakiDataSources = append(MerakiDataSources, merakiAdministered.DataSources()...)
	MerakiDataSources = append(MerakiDataSources, merakiDevices.DataSources()...)
	MerakiDataSources = append(MerakiDataSources, merakiOrganizations.DataSources()...)
	MerakiDataSources = append(MerakiDataSources, merakiNetworks.DataSources()...)

}
