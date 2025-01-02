package networks_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsNetworksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_networks"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_networks"),
			},

			// Create and Read Network
			{
				Config: utils.CreateNetworkConfig("test_acc_meraki_organizations_networks", "test_acc_network"),
				Check:  utils.NetworkTestChecks("test_acc_network"),
			},

			// Read Organizations Networks
			{
				Config: testAccOrganizationsNetworksDataSourceConfigRead(),
				Check:  OrganizationsNetworksDataSourceTestChecks(),
			},
		},
	})
}

// Ensure that you are creating the network with 2 tags
func testAccOrganizationsNetworksDataSourceConfigRead() string {
	return fmt.Sprintf(`
	%s

	resource "meraki_network" "test" {
		name         = "test_acc_network"
		organization_id = resource.meraki_organization.test.organization_id
		tags         = ["tag1"]
		timezone     = "America/Los_Angeles"
		product_types = ["appliance", "switch", "wireless"]
		notes        = "Additional description of the network"
	}

	data "meraki_organizations_networks" "test" {
		organization_id = resource.meraki_organization.test.organization_id
	}
	`,
		utils.CreateOrganizationConfig("test_acc_meraki_organizations_networks"),
	)
}

// OrganizationsNetworksDataSourceTestChecks returns the test check functions for verifying networks data source
func OrganizationsNetworksDataSourceTestChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.#":                 "1",
		"list.0.name":            "test_acc_network",
		"list.0.timezone":        "America/Los_Angeles",
		"list.0.tags.#":          "1",
		"list.0.tags.0":          "tag1",
		"list.0.product_types.#": "3",
		"list.0.product_types.0": "appliance",
		"list.0.product_types.1": "switch",
		"list.0.product_types.2": "wireless",
		"list.0.notes":           "Additional description of the network",
	}

	return utils.ResourceTestCheck("data.meraki_organizations_networks.test", expectedAttrs)
}
