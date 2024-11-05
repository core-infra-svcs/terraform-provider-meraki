package networks_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksSyslogServersResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_syslog_servers.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network.
			{
				Config: testAccNetworksSyslogServersResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_syslog_servers"),
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

			// Update and Read Networks Syslog Servers.
			{
				Config: testAccSyslogServersResourceConfigUpdateNetworkSyslogServersSettings,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_syslog_servers.test", "servers.#", "1"),
					resource.TestCheckResourceAttr("meraki_networks_syslog_servers.test", "servers.0.host", "1.2.3.67"),
					resource.TestCheckResourceAttr("meraki_networks_syslog_servers.test", "servers.0.port", "443"),
					resource.TestCheckResourceAttr("meraki_networks_syslog_servers.test", "servers.0.roles.0", "URLs"),
				),
			},
		},
	})
}

func testAccNetworksSyslogServersResourceConfigCreateNetwork(orgId string) string {
	result := fmt.Sprintf(`
 resource "meraki_network" "test" {
	organization_id = %s
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_networks_syslog_servers"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
 `, orgId)
	return result
}

const testAccSyslogServersResourceConfigUpdateNetworkSyslogServersSettings = `
resource "meraki_network" "test" {
	product_types = ["appliance", "switch", "wireless"]
	
}
resource "meraki_networks_syslog_servers" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  servers = [{
		host = "1.2.3.67"
		port = "443"
		roles = ["URLs"]
	}] 
}
`
