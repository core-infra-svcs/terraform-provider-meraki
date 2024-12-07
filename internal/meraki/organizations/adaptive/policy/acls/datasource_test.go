package acls_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsAdaptivePolicyAclsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_meraki_organizations_adaptive_policy_acls"),
				Check:  utils.OrganizationTestChecks("test_acc_meraki_organizations_adaptive_policy_acls"),
			},

			// Create OrganizationsAdaptivePolicyAcl
			{
				Config: AdaptivePolicyAclResourceConfigCreate(),
				Check:  AdaptivePolicyAclTestChecks("Block sensitive web traffic", "Blocks sensitive web traffic", "ipv6", 1),
			},

			// Read OrganizationsAdaptivePolicyAcls
			{
				Config: AdaptivePolicyAclsDataSourceConfigRead(),
				Check:  AdaptivePolicyAclsDataSourceTestChecks(),
			},
		},
	})
}

// AdaptivePolicyAclsDataSourceConfigRead returns the configuration string for reading adaptive policy ACLs from a data source
func AdaptivePolicyAclsDataSourceConfigRead() string {
	return fmt.Sprintf(`
	%s

	// Ensure the ACL resource is created
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
			}
		]
	}

	// Data source to read adaptive policy ACLs
	data "meraki_organizations_adaptive_policy_acls" "test" {
		organization_id = meraki_organization.test.organization_id
	}
	`,
		utils.CreateOrganizationConfig("test_acc_meraki_organizations_adaptive_policy_acls"),
	)
}

// AdaptivePolicyAclsDataSourceTestChecks returns the test check functions for reading adaptive policy ACLs from a data source
func AdaptivePolicyAclsDataSourceTestChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.#":                  "1", // Expecting one ACL in the list
		"list.0.name":             "Block sensitive web traffic",
		"list.0.description":      "Blocks sensitive web traffic",
		"list.0.ip_version":       "ipv6",
		"list.0.rules.0.policy":   "deny",
		"list.0.rules.0.protocol": "tcp",
		"list.0.rules.0.src_port": "1,33",
		"list.0.rules.0.dst_port": "22-30",
	}

	return utils.ResourceTestCheck("data.meraki_organizations_adaptive_policy_acls.test", expectedAttrs)
}
