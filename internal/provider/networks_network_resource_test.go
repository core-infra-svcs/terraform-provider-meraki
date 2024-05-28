package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func init() {
	resource.AddTestSweepers("meraki_network", &resource.Sweeper{
		Name: "meraki_network",
		F: func(organization string) error {
			return sweepMerakiNetwork(organization)
		},
	})
}

func TestAccOrganizationsNetworkResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing (network).
			{
				Config: testAccOrganizationsNetworkResourceConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network"),
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

			// Update testing
			{
				Config: testAccOrganizationsNetworkResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_network.test", "enrollment_string", "my-enrollment-string"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.1", "tag2"),
				),
			},

			/* TODO: Need OrganizationConfigTemplate resource in order to test...
			// Bind Network Test
				{
					Config: testAccOrganizationsNetworkResourceConfigBind(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
					Check: resource.ComposeAggregateTestCheckFunc(
						//resource.TestCheckResourceAttr("meraki_network.test", "enrollment_string", "my-enrollment-string"),
						resource.TestCheckResourceAttr("meraki_network.bind", "product_types.#", "3"),
						resource.TestCheckResourceAttr("meraki_network.bind", "product_types.0", "appliance"),
						resource.TestCheckResourceAttr("meraki_network.bind", "product_types.1", "switch"),
						resource.TestCheckResourceAttr("meraki_network.bind", "product_types.2", "wireless"),
						resource.TestCheckResourceAttr("meraki_network.bind", "tags.#", "2"),
						resource.TestCheckResourceAttr("meraki_network.bind", "tags.0", "tag1"),
						resource.TestCheckResourceAttr("meraki_network.bind", "tags.1", "tag2"),
						resource.TestCheckResourceAttr("meraki_network.bind", "auto_bind", "true"),
					),
				},

				// Unbind Network Test
				{
					Config: testAccOrganizationsNetworkResourceConfigUnBind(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
					Check: resource.ComposeAggregateTestCheckFunc(
						//resource.TestCheckResourceAttr("meraki_network.test", "enrollment_string", "my-enrollment-string"),
						resource.TestCheckResourceAttr("meraki_network.bind", "product_types.#", "3"),
						resource.TestCheckResourceAttr("meraki_network.bind", "product_types.0", "appliance"),
						resource.TestCheckResourceAttr("meraki_network.bind", "product_types.1", "switch"),
						resource.TestCheckResourceAttr("meraki_network.bind", "product_types.2", "wireless"),
						resource.TestCheckResourceAttr("meraki_network.bind", "tags.#", "2"),
						resource.TestCheckResourceAttr("meraki_network.bind", "tags.0", "tag1"),
						resource.TestCheckResourceAttr("meraki_network.bind", "tags.1", "tag2"),
						resource.TestCheckResourceAttr("meraki_network.bind", "auto_bind", "false"),
					),
				},


			*/

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccOrganizationsNetworkResourceConfig(orgId string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_network"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccOrganizationsNetworkResourceConfigUpdate = `

resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1", "tag2"]
	name = "test_acc_network-2"
	timezone = "America/Chicago"
	notes = "Additional description of the network-2"
}
`

/* TODO: Need OrganizationConfigTemplate resource in order to test...
func testAccOrganizationsNetworkResourceConfigBind(orgId string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_network" "bind" {
	depends_on = [resource.meraki_network.test]
	auto_bind = true
    config_template_id = resource.meraki_network.test.network_id
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_network_bind"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId, orgId)
	return result
}

func testAccOrganizationsNetworkResourceConfigUnBind(orgId string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
}

resource "meraki_network" "bind" {
	depends_on = [resource.meraki_network.test]
	auto_bind = false
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_network_bind"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId, orgId)
	return result
}
*/
