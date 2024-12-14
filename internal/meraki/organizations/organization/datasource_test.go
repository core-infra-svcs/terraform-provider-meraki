package organization_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations"),
			},

			// Read Organizations DataSource
			{
				Config: testAccOrganizationsDataSourceConfigRead(),
				Check:  OrganizationsDataSourceTestChecks(),
			},
		},
	})
}

// testAccOrganizationsDataSourceConfigRead returns the configuration for reading the organizations data source
func testAccOrganizationsDataSourceConfigRead() string {
	return `
	data "meraki_organizations" "test" {
	}
	`
}

// OrganizationsDataSourceTestChecks returns the test check functions for verifying the organizations data source
func OrganizationsDataSourceTestChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		// Add your expected checks here, e.g.
		// "list.#": "2",
		// "list.1.name": "test_acc_meraki_organizations",
		// "list.1.api_enabled": "true",
	}

	return utils.ResourceTestCheck("data.meraki_organizations.test", expectedAttrs)
}
