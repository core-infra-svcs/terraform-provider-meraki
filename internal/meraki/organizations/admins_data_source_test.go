package organizations_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsAdminsDataSource(t *testing.T) {

	timestamp := utils.GenerateTimestamp() // Using the utility function to generate a timestamp

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing (Organization)
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_admins"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_admins"),
			},

			// Create and Read testing (network)
			{
				Config: utils.CreateNetworkConfig("test_acc_meraki_organizations_admins", "test_acc_organizations_admins"),
				Check:  utils.NetworkTestChecks("test_acc_organizations_admins"),
			},

			// Create and Read testing (admin)
			{
				Config: AdminResourceConfig(timestamp),
				Check:  AdminResourceTestChecks(timestamp),
			},

			// Read Testing (admins)
			{
				Config: AdminsDataSourceConfigRead(),
				Check:  AdminsDataSourceTestChecks(timestamp),
			},
		},
	})
}

// AdminsDataSourceConfigRead returns the configuration string to read multiple admins from a data source
func AdminsDataSourceConfigRead() string {
	return fmt.Sprintf(`
	%s
	data "meraki_organizations_admins" "test" {
		depends_on = [resource.meraki_organization.test]
		organization_id = resource.meraki_organization.test.organization_id
	}
	`,
		AdminResourceConfig(utils.GenerateTimestamp()),
	)
}

// AdminsDataSourceTestChecks returns the aggregated test check functions for reading multiple admins
func AdminsDataSourceTestChecks(timestamp string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.#", "2"),
		resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.name", "test_acc_admin"),
		resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.email", fmt.Sprintf("test_acc_meraki_organizations_admin_test_%s@example.com", timestamp)),
		resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.org_access", "read-only"),
		resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.has_api_key", "false"),
		resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.tags.0.tag", "west"),
		resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.tags.0.access", "read-only"),
		resource.TestCheckResourceAttr("data.meraki_organizations_admins.test", "list.1.authentication_method", "Email"),
	)
}
