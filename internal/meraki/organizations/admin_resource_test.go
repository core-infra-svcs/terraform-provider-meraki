package organizations_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsAdminResource(t *testing.T) {
	admins := 2                            // number of admins
	timestamp := utils.GenerateTimestamp() // Use utils for timestamp generation

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read testing (Organization)
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_admin"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_admin"),
			},

			// Create test network
			{
				Config: utils.CreateNetworkConfig("test_acc_meraki_organizations_admin", "test_acc_organizations_admin"),
				Check:  utils.NetworkTestChecks("test_acc_organizations_admin"),
			},

			// Create and Read testing (admin)
			{
				Config: AdminResourceConfig(timestamp),
				Check:  AdminResourceTestChecks(timestamp),
			},

			// Update testing (admin)
			{
				Config: AdminResourceConfigUpdate(timestamp),
				Check:  AdminUpdateTestChecks(timestamp),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organizations_admin.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},

			// Test the creation of multiple group policies
			{
				Config: AdminMultiplePoliciesConfig(admins, timestamp),
				Check:  AdminMultiplePoliciesTestChecks(admins, timestamp),
			},
		},
	})
}

// AdminResourceConfig returns a configuration string to create an admin resource
func AdminResourceConfig(timestamp string) string {
	return fmt.Sprintf(`
	%s
	resource "meraki_organizations_admin" "test" {
		depends_on = ["meraki_organization.test", "meraki_network.test"]
		organization_id = resource.meraki_organization.test.organization_id
		name        = "test_acc_admin"
		email       = "test_acc_meraki_organizations_admin_test_%s@example.com"
		org_access  = "read-only"
		authentication_method = "Email"
		tags = [
			{
				tag = "west"
				access = "read-only"
			}
		]
		networks = [{
			id = resource.meraki_network.test.network_id
			access = "read-only"
		}]
	}
	`,
		utils.CreateNetworkConfig("test_acc_meraki_organizations_admin",
			"test_acc_meraki_organizations_admin"), timestamp)
}

// AdminResourceTestChecks returns the aggregated test check functions for an admin resource
func AdminResourceTestChecks(timestamp string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "name", "test_acc_admin"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "email", fmt.Sprintf("test_acc_meraki_organizations_admin_test_%s@example.com", timestamp)),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "org_access", "read-only"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "authentication_method", "Email"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "has_api_key", "false"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.tag", "west"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.access", "read-only"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.0.access", "read-only"),
	)
}

// AdminResourceConfigUpdate returns a configuration string for updating an admin resource
func AdminResourceConfigUpdate(timestamp string) string {
	return fmt.Sprintf(`
	%s
	resource "meraki_organizations_admin" "test" {
		depends_on = ["meraki_organization.test", "meraki_network.test"]
		organization_id = resource.meraki_organization.test.organization_id
		name        = "test_acc_admin"
		email       = "test_acc_meraki_organizations_admin_test_%s@example.com"
		org_access  = "read-only"
		authentication_method = "Email"
		tags = [
			{
				tag = "east"
				access = "read-only"
			}
		]
		networks = [{
			id = resource.meraki_network.test.network_id
			access = "read-only"
		}]
	}
	`,
		utils.CreateNetworkConfig("test_acc_meraki_organizations_admin", "test_acc_meraki_organizations_admin"), timestamp)
}

// AdminUpdateTestChecks returns the aggregated test check functions for an admin resource after an update
func AdminUpdateTestChecks(timestamp string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "name", "test_acc_admin"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "email", fmt.Sprintf("test_acc_meraki_organizations_admin_test_%s@example.com", timestamp)),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "org_access", "read-only"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "authentication_method", "Email"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.tag", "east"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "tags.0.access", "read-only"),
		resource.TestCheckResourceAttr("meraki_organizations_admin.test", "networks.0.access", "read-only"),
	)
}

// AdminMultiplePoliciesConfig returns a configuration string for creating multiple admin resources
func AdminMultiplePoliciesConfig(admins int, timestamp string) string {
	config := fmt.Sprintf(`
	%s`,
		utils.CreateNetworkConfig("test_acc_meraki_organizations_admin", "test_acc_meraki_organizations_admin"))

	// Append each admin configuration
	for i := 1; i <= admins; i++ {
		config += fmt.Sprintf(`
	resource "meraki_organizations_admin" "test%d" {
		depends_on = ["meraki_network.test"]
		organization_id = resource.meraki_organization.test.organization_id
		name        = "test_acc_admin_%d"
		email       = "test_acc_meraki_organizations_admin_test_%s_%d@example.com"
		org_access  = "read-only"
		authentication_method = "Email"
		tags = [
			{
				tag = "west"
				access = "read-only"
			}
		]
		networks = [{
			id = resource.meraki_network.test.network_id
			access = "read-only"
		}]
	}
	`, i, i, timestamp, i)
	}
	return config
}

// AdminMultiplePoliciesTestChecks returns the test check functions for multiple admin policies
func AdminMultiplePoliciesTestChecks(admins int, timestamp string) resource.TestCheckFunc {
	var checks []resource.TestCheckFunc
	for i := 1; i <= admins; i++ {
		resourceName := fmt.Sprintf("meraki_organizations_admin.test%d", i)
		checks = append(checks,
			resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test_acc_admin_%d", i)),
			resource.TestCheckResourceAttr(resourceName, "email", fmt.Sprintf("test_acc_meraki_organizations_admin_test_%s_%d@example.com", timestamp, i)),
			resource.TestCheckResourceAttr(resourceName, "org_access", "read-only"),
			resource.TestCheckResourceAttr(resourceName, "authentication_method", "Email"),
			resource.TestCheckResourceAttr(resourceName, "tags.0.tag", "west"),
			resource.TestCheckResourceAttr(resourceName, "tags.0.access", "read-only"),
			resource.TestCheckResourceAttr(resourceName, "networks.0.access", "read-only"),
		)
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}
