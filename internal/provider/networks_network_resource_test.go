package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"strings"
	"testing"
	"time"
)

func init() {
	resource.AddTestSweepers("meraki_network", &resource.Sweeper{
		Name: "meraki_network",

		// The Organization used in acceptance tests
		F: func(organization string) error {
			client, err := SweeperHTTPClient()
			if err != nil {
				return fmt.Errorf("error getting http client: %s", err)
			}

			// HTTP DELETE METHOD does not leverage the retry-after header and throws 400 errors.
			retries := 3
			wait := 1
			var deletedFromMerakiPortal bool
			deletedFromMerakiPortal = false

			// Search test organization for networks
			perPage := int32(100000)
			inlineResp, _, err := client.NetworksApi.GetOrganizationNetworks(nil, organization).PerPage(perPage).Execute()
			if err != nil {
				return fmt.Errorf("error getting network list from organization:%s \nerror: %s", organization, err)
			}

			for _, merakiNetwork := range inlineResp {

				// match on networks starting with "test_acc" in name
				if strings.HasPrefix(*merakiNetwork.Name, "test_acc") {
					fmt.Println(fmt.Sprintf("deleting network: %s, id: %s", *merakiNetwork.Name, *merakiNetwork.Id))

					for retries > 0 {
						// Delete test network
						httpResp, err2 := client.NetworksApi.DeleteNetwork(nil, *merakiNetwork.Id).Execute()
						if err2 != nil {
							return fmt.Errorf("error deleting network from organization:%s \nerror: %s", organization, err2)
						}

						if httpResp.StatusCode == 204 {
							fmt.Println(fmt.Sprintf("Successfully deleted network: %s, id: %s", *merakiNetwork.Name, *merakiNetwork.Id))
							deletedFromMerakiPortal = true

							// escape loop
							break

						} else {

							// decrement retry counter
							retries -= 1

							// exponential wait
							time.Sleep(time.Duration(wait) * time.Second)
							wait += 1
						}

						if !deletedFromMerakiPortal {
							fmt.Println(fmt.Sprintf("Failed to delete network: %s, id: %s", *merakiNetwork.Name, *merakiNetwork.Id))
							fmt.Println(fmt.Sprintf("HTTP response: \n%v", httpResp))
						}

					}

				}

			}
			return nil
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
