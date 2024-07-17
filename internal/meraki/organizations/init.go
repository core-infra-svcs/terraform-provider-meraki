package organizations

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var organizationResources []func() resource.Resource
var organizationDataSources []func() datasource.DataSource

func init() {

	// Initialize with all organization resources //

	//organizationResources = append(organizationResources, subFolder.Resources()...)

	// package level resources
	organizationResources = append(organizationResources,
		NewOrganizationsAdaptivePolicyAclResource,
		NewOrganizationsAdminResource,
		NewOrganizationsApplianceVpnVpnFirewallRulesResource,
		NewOrganizationsClaimResource,
		NewOrganizationsSamlIdpResource,
		NewOrganizationSamlResource,
		NewOrganizationsSamlRolesResource,
		NewOrganizationsSnmpResource,
		NewOrganizationResource,
		NewOrganizationPolicyObjectResource,
	)

	// Initialize with all organization data sources //

	//organizationDataSources = append(organizationDataSources, subFolder.DataSources()...)

	// package level data sources
	organizationDataSources = append(organizationDataSources,
		NewOrganizationsAdaptivePolicyAclsDataSource,
		NewOrganizationsAdminsDataSource,
		NewOrganizationsCellularGatewayUplinkStatusesDataSource,
		NewOrganizationsDataSource,
		NewOrganizationsSamlIdpsDataSource,
		NewOrganizationsSamlRolesDataSource,
	)

}

func Resources() []func() resource.Resource {
	return organizationResources
}

func DataSources() []func() datasource.DataSource {
	return organizationDataSources
}
