package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("meraki_organization", &resource.Sweeper{
		Name: "meraki_organization",

		// The Organization used in acceptance tests
		F: func(organization string) error {
			client, err := SweeperHTTPClient()
			if err != nil {
				return fmt.Errorf("error getting http client: %s", err)
			}

			// Search test organization for networks
			inlineResp, _, err := client.OrganizationsApi.GetOrganizations(context.Background()).Execute()
			if err != nil {
				return fmt.Errorf("error getting organizations list from Meraki API:%s\n", err)
			}

			// List of Organizations
			for _, merakiOrganization := range inlineResp {

				// match on organizations starting with "test_acc" in name
				if strings.HasPrefix(*merakiOrganization.Name, "test") {
					fmt.Printf("deleting organization: %s, id: %s\n", *merakiOrganization.Name, *merakiOrganization.Id)

					// Check for Networks
					perPage := int32(100000)
					inlineRespNetwork, _, err1 := client.NetworksApi.GetOrganizationNetworks(context.Background(), *merakiOrganization.Id).PerPage(perPage).Execute()
					if err1 != nil {
						return fmt.Errorf("error getting network list from organization:%s error: %s\n", *merakiOrganization.Id, err1)
					}

					for _, merakiNetwork := range inlineRespNetwork {
						fmt.Printf("deleting network: %s, id: %s\n", *merakiNetwork.Name, *merakiNetwork.Id)

						// Delete test network
						networkHttpResp, err2 := client.NetworksApi.DeleteNetwork(context.Background(), *merakiNetwork.Id).Execute()
						if err2 != nil {
							fmt.Printf("%v\n", networkHttpResp)
						}

						if networkHttpResp.StatusCode == 204 {
							fmt.Printf("Successfully deleted network: %s, id: %s\n", *merakiNetwork.Name, *merakiNetwork.Id)
						}
					}

					// TODO: Delete admins

					// TODO: Org Inventory Devices (No api endpoint)

					// Delete test Organization
					httpRespOrg, err3 := client.OrganizationsApi.DeleteOrganization(context.Background(), *merakiOrganization.Id).Execute()
					if err3 != nil {
						fmt.Printf("%v\n", httpRespOrg.Body)
					}
					if httpRespOrg.StatusCode == 204 {
						fmt.Printf("Successfully deleted organization: %s, id: %s\n", *merakiOrganization.Name, *merakiOrganization.Id)
					}

				}
			}

			return nil
		}})
}

func TestAccOrganizationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organization"),
					resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),
				),
			},

			// Update testing
			{
				Config: testAccOrganizationResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organization_update"),
					resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organization.test", "management_details_name", "MSP ID"),
					resource.TestCheckResourceAttr("meraki_organization.test", "management_details_value", "123456"),
				),
			},

			// Clone Organization testing
			{
				Config: testAccOrganizationResourceConfigClone,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.clone", "name", "test_acc_meraki_organization_clone"),
					resource.TestCheckResourceAttr("meraki_organization.clone", "api_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organization.clone", "management_details_name", "MSP ID"),
					resource.TestCheckResourceAttr("meraki_organization.clone", "management_details_value", "123456"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationResourceConfig = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organization"
	api_enabled = true
}
`

const testAccOrganizationResourceConfigUpdate = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_organization_update"
	api_enabled = true
	management_details_name = "MSP ID"
	management_details_value = "123456"
}
`

const testAccOrganizationResourceConfigClone = `
resource "meraki_organization" "test" {}

resource "meraki_organization" "clone" {
	depends_on = [meraki_organization.test]
	clone_organization_id = resource.meraki_organization.test.organization_id
	name = "test_acc_meraki_organization_clone"
	api_enabled = true
	
}
`
