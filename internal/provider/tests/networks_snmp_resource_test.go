package tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSnmpSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: testAccOrganizationsSnmpSettingsResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_snmp_settings"),
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

			// Create
			{
				Config: testAccNetworkSnmpSettingsResourceConfigUpdateNetworkSnmpSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_snmp.test", "access", "community"),
					resource.TestCheckResourceAttr("meraki_networks_snmp.test", "community_string", "public"),
				),
			},

			{
				Config: testAccNetworkSnmpSettingsResourceConfigUpdateSNMPSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_snmp.test", "access", "users"),
					resource.TestCheckResourceAttr("meraki_networks_snmp.test", "users.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_snmp.test", "users.0.username", "snmp_user"),
					resource.TestCheckResourceAttr("meraki_networks_snmp.test", "users.0.passphrase", "snmp_passphrase"),
				),
			},
			{
				ResourceName:      "meraki_networks_snmp.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources["meraki_networks_snmp.test"]
					if !ok {
						return "", fmt.Errorf("not found: %s", "meraki_networks_snmp.test")
					}
					return rs.Primary.ID, nil
				},
			},
		},
	})
}

func testAccOrganizationsSnmpSettingsResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_snmp_settings"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`, orgId)
	return result
}

const testAccNetworkSnmpSettingsResourceConfigUpdateNetworkSnmpSettings = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]	
}
resource "meraki_organizations_snmp" "test" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	access = "community"
	community_string = "public"
}
`

const testAccNetworkSnmpSettingsResourceConfigUpdateSNMPSettings = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]	
}
resource "meraki_organizations_snmp" "test" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	access = "users"
	users = [{
		username = "snmp_user"
		passphrase = "snmp_passphrase"
	}]
}
`
