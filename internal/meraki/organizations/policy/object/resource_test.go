package object_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"log"
	"time"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationPolicyObjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create Organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_policy_object"),
				Check: resource.ComposeAggregateTestCheckFunc(
					utils.OrganizationTestChecks("test_acc_meraki_organizations_policy_object"),
				),
			},

			// Create and Verify Policy Object
			{
				Config: testAccOrganizationPolicyObjectResourceConfigCreate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					OrganizationPolicyObjectTestChecks("test_acc_meraki_organizations_policy_object"),
					resource.TestCheckResourceAttrSet("meraki_organizations_policy_object.test", "id"),
				),
			},

			// Import State Testing
			{
				ResourceName:            "meraki_organizations_policy_object.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ip", "cidr"},
			},

			// Pause to Ensure Resource is Created in API
			{
				PreConfig: func() {
					log.Println("Waiting before update to ensure API sync...")
					time.Sleep(5 * time.Second)
				},
				Config: testAccOrganizationPolicyObjectResourceConfigUpdate(),
				Check:  OrganizationPolicyObjectUpdateTestChecks("test_acc_meraki_organizations_policy_object"),
			},
		},
	})
}

// Create Configuration
func testAccOrganizationPolicyObjectResourceConfigCreate() string {
	return `
	resource "meraki_organization" "test" {
		name = "test_acc_meraki_organization_policy_object"
	}

	resource "meraki_organizations_policy_object" "test" {
		depends_on = [meraki_organization.test]
		organization_id = meraki_organization.test.organization_id
		name = "test_acc_meraki_organizations_policy_object"
		category = "network"
		type = "cidr"
		cidr = "10.0.0.0/24"
		ip = "1.2.3.4"
		group_ids = []
	}
	`
}

// Update Configuration
func testAccOrganizationPolicyObjectResourceConfigUpdate() string {
	return `
	resource "meraki_organization" "test" {
		name = "test_acc_meraki_organization_policy_object"
	}

	resource "meraki_organizations_policy_object" "test" {
		depends_on = [meraki_organization.test]
		organization_id = meraki_organization.test.organization_id
		name = "test_acc_meraki_organizations_policy_object"
		category = "network"
		type = "cidr"
		cidr = "10.0.0.0/24"
		ip = "1.2.3.5"
		group_ids = []
	}
	`
}

func OrganizationPolicyObjectTestChecks(name string) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":     name,
		"category": "network",
		"type":     "cidr",
		"cidr":     "10.0.0.0/24",
		"ip":       "1.2.3.4",
	}

	fmt.Printf("Verifying initial policy object attributes: %+v\n", expectedAttrs)
	return utils.ResourceTestCheck("meraki_organizations_policy_object.test", expectedAttrs)
}

func OrganizationPolicyObjectUpdateTestChecks(name string) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":     name,
		"category": "network",
		"type":     "cidr",
		"cidr":     "10.0.0.0/24",
		"ip":       "1.2.3.5",
	}

	fmt.Printf("Verifying updated policy object attributes: %+v\n", expectedAttrs)
	return utils.ResourceTestCheck("meraki_organizations_policy_object.test", expectedAttrs)
}
