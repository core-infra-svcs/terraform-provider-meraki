package network_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccOrganizationsNetworkResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing (network).
			{
				Config: testAccOrganizationsNetworkResourceConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Update testing
			{
				Config: testAccOrganizationsNetworkResourceConfigUpdate(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_update"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Chicago"),
					//resource.TestCheckResourceAttr("meraki_network.test", "enrollment_string", "test_acc_string"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network update"),
				),
			},

			{
				ResourceName:      "meraki_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOrganizationsNetworkResourceConfig(orgId string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
	organization_id = "%s"
	name = "test_acc_network"
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

func testAccOrganizationsNetworkResourceConfigUpdate(orgId string) string {
	result := fmt.Sprintf(`

resource "meraki_network" "test" {
	organization_id = "%s"
 	product_types = ["appliance", "switch", "wireless"]
	name = "test_acc_network_update"
	timezone = "America/Chicago"
	tags = ["tag1", "tag2"]
	// enrollment_string = "test_acc_string"
	notes = "Additional description of the network update"
}
`, orgId)
	return result
}
