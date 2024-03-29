package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsAdminResource(t *testing.T) {
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	admins := 13 // number of admins

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: testAccOrganizationsAdminResourceConfigCreateOrg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_admin"),
				),
			},

			// Create and Read testing (network)
			{
				Config: testAccOrganizationsAdminResourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_organizations_admin"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Create and Read testing (admin)
			{
				Config: testAccOrganizationsAdminResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "name", "test_acc_admin"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "email", "meraki_organizations_admin_test_2023_06_05@example.com"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "org_access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "authentication_method", "Email"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "has_api_key", "false"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.tag", "west"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.0.access", "read-only"),
				),
			},

			// Update testing
			{
				Config: testUpdatedAccOrganizationsAdminResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.tag", "east"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.access", "read-only"),
					resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.0.access", "read-only"),
				),
			},

			// Test the creation of multiple group policies.
			{
				Config: testAccOrganizationsAdminResourceConfigMultiplePolicies(orgId, admins),
				Check: func(s *terraform.State) error {
					var checks []resource.TestCheckFunc
					// Dynamically generate checks for each group policy
					for i := 1; i <= admins; i++ {
						resourceName := fmt.Sprintf("meraki_organizations_admin.test%d", i)
						checks = append(checks,
							resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test_acc_admin_%d", i)),
							resource.TestCheckResourceAttr(resourceName, "email", fmt.Sprintf("meraki_organizations_admin_test_%d@example.com", i)),
							resource.TestCheckResourceAttr(resourceName, "org_access", "read-only"),
							resource.TestCheckResourceAttr(resourceName, "authentication_method", "Email"),
							resource.TestCheckResourceAttr(resourceName, "tags.0.tag", "west"),
							resource.TestCheckResourceAttr(resourceName, "tags.0.access", "read-only"),
							resource.TestCheckResourceAttr(resourceName, "networks.0.access", "read-only"),
						)
					}
					return resource.ComposeAggregateTestCheckFunc(checks...)(s)
				},
			},

			// Delete testing automatically occurs in TestCase
			// TODO - This test can result in orphaned resources as organizations cannot be deleted with admins still present.
		},
	})
}

const testAccOrganizationsAdminResourceConfigCreateOrg = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_admin"
 	api_enabled = true
 }
 `

const testAccOrganizationsAdminResourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	depends_on = ["meraki_organization.test"]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_organizations_admin"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

const testAccOrganizationsAdminResourceConfig = `
resource "meraki_organization" "test" {}

resource "meraki_network" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_organizations_admin" "test" {
	depends_on = ["meraki_organization.test", "meraki_network.test"]
	organization_id = resource.meraki_organization.test.organization_id
	name        = "test_acc_admin"
	email       = "meraki_organizations_admin_test_2023_06_05@example.com"
	org_access   = "read-only"
	authentication_method = "Email"
    tags = [
			  {
			   tag = "west"
			   access = "read-only"
			  }]
    networks    = [{
                  id = resource.meraki_network.test.network_id
                  access = "read-only"
                }]
}
`

const testUpdatedAccOrganizationsAdminResourceConfig = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_organizations_admin" "test" {
	depends_on = ["meraki_organization.test", "meraki_network.test"]
	organization_id = resource.meraki_organization.test.organization_id
	name        = "test_acc_admin"
	email       = "meraki_organizations_admin_test_2023_06_05@example.com"
	org_access   = "read-only"
	authentication_method = "Email"
    tags = [
			{
				tag = "east"
				access = "read-only"
			}]
    networks    = [{
                  id = resource.meraki_network.test.network_id
                  access = "read-only"
                }]
}
`

func testAccOrganizationsAdminResourceConfigMultiplePolicies(orgId string, admins int) string {
	config := fmt.Sprintf(`
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_admin_2"
 	api_enabled = true
 }

resource "meraki_network" "test" {
	depends_on = ["meraki_organization.test"]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_organizations_admin"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`)

	// Append each admin configuration
	for i := 1; i <= admins; i++ {
		config += fmt.Sprintf(`

resource "meraki_organizations_admin" "test%d" {
	depends_on = ["meraki_network.test"]
	organization_id = resource.meraki_organization.test.organization_id
	name        = "test_acc_admin_%d"
	email       = "meraki_organizations_admin_test_%d@example.com"
	org_access   = "read-only"
	authentication_method = "Email"
    tags = [
			  {
			   tag = "west"
			   access = "read-only"
			  }]
    networks    = [{
                  id = resource.meraki_network.test.network_id
                  access = "read-only"
                }]
}


`, i, i, i)
	}

	return config
}
