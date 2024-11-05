package organizations_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccAdaptivePolicyAclResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_adaptive_policy_acl"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_adaptive_policy_acl"),
			},

			// Create and Read testing (ACL)
			{
				Config: AdaptivePolicyAclResourceConfigCreate(),
				Check:  AdaptivePolicyAclTestChecks("Block sensitive web traffic", "Blocks sensitive web traffic", "ipv6", 1),
			},

			// Update testing (ACL)
			{
				Config: AdaptivePolicyAclResourceConfigUpdate(),
				Check:  AdaptivePolicyAclTestUpdateChecks(2),
			},

			// Import State testing
			{
				ResourceName:            "meraki_organizations_adaptive_policy_acl.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

// AdaptivePolicyAclResourceConfigCreate returns a configuration string to create an adaptive policy ACL resource
func AdaptivePolicyAclResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
	resource "meraki_organizations_adaptive_policy_acl" "test" {
		organization_id = meraki_organization.test.organization_id
		name            = "Block sensitive web traffic"
		description     = "Blocks sensitive web traffic"
		ip_version      = "ipv6"

		// Inline list of rules as a list of maps (not strings)
		rules = [
			{
				policy   = "deny"
				protocol = "tcp"
				src_port = "1,33"
				dst_port = "22-30"
			}
		]
	}
	`,
		utils.CreateOrganizationConfig("test_acc_meraki_organizations_adaptive_policy_acl"),
	)
}

// AdaptivePolicyAclResourceConfigUpdate returns a configuration string for updating an adaptive policy ACL resource
func AdaptivePolicyAclResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
	resource "meraki_organizations_adaptive_policy_acl" "test" {
		organization_id = meraki_organization.test.organization_id
		name            = "Block sensitive web traffic"
		description     = "Blocks sensitive web traffic"
		ip_version      = "ipv6"

		rules = [
			{
				policy   = "deny"
				protocol = "tcp"
				src_port = "1,33"
				dst_port = "22-30"
			},
			{
				policy   = "allow"
				protocol = "any"
				src_port = "any"
				dst_port = "any"
			}
		]
	}
	`,
		utils.CreateOrganizationConfig("test_acc_meraki_organizations_adaptive_policy_acl"),
	)
}

// AdaptivePolicyAclTestChecks returns the aggregated test check functions for an adaptive policy ACL resource
func AdaptivePolicyAclTestChecks(name, description, ipVersion string, ruleCount int) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":             name,
		"description":      description,
		"ip_version":       ipVersion,
		"rules.#":          fmt.Sprintf("%d", ruleCount),
		"rules.0.policy":   "deny",
		"rules.0.protocol": "tcp",
		"rules.0.src_port": "1,33",
		"rules.0.dst_port": "22-30",
	}

	return utils.ResourceTestCheck("meraki_organizations_adaptive_policy_acl.test", expectedAttrs)
}

// AdaptivePolicyAclTestUpdateChecks returns the aggregated test check functions for an adaptive policy ACL resource after an update
func AdaptivePolicyAclTestUpdateChecks(ruleCount int) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"rules.#":          fmt.Sprintf("%d", ruleCount),
		"rules.0.policy":   "deny",
		"rules.0.protocol": "tcp",
		"rules.0.src_port": "1,33",
		"rules.0.dst_port": "22-30",
		"rules.1.policy":   "allow",
		"rules.1.protocol": "any",
		"rules.1.src_port": "any",
		"rules.1.dst_port": "any",
	}

	return utils.ResourceTestCheck("meraki_organizations_adaptive_policy_acl.test", expectedAttrs)
}
