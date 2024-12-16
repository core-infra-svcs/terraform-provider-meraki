package snmp_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccOrganizationsSnmpSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create Organization
			{
				Config: utils.CreateOrganizationConfig("test_acc_organizations_snmp"),
				Check:  utils.OrganizationTestChecks("test_acc_organizations_snmp"),
			},

			// Create and Read SNMP settings
			{
				Config: testAccOrganizationsSnmpSettingsResourceConfigCreate(),
				Check:  OrganizationsSnmpSettingsTestChecks([]string{"1.1.1.1"}, "SHA", "AES128"),
			},

			// Update and Read SNMP settings
			{
				Config: testAccOrganizationsSnmpSettingsResourceConfigUpdate(),
				Check:  OrganizationsSnmpSettingsTestChecks([]string{"1.1.1.1", "2.2.2.2"}, "SHA", "AES128"),
			},
		},
	})
}

// testAccOrganizationsSnmpSettingsResourceConfigCreate returns the configuration for creating SNMP settings
func testAccOrganizationsSnmpSettingsResourceConfigCreate() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_snmp" "test" {
		depends_on = [meraki_organization.test]
		organization_id = resource.meraki_organization.test.organization_id
		v2c_enabled     = true
		v3_enabled      = true
		v3_auth_pass    = "bjhb989*%gg"
		v3_priv_pass    = "jkjbbbj679$%"
		v3_auth_mode    = "SHA"
		v3_priv_mode    = "AES128"
		peer_ips        = ["1.1.1.1"]
	}
	`
}

// testAccOrganizationsSnmpSettingsResourceConfigUpdate returns the configuration for updating SNMP settings
func testAccOrganizationsSnmpSettingsResourceConfigUpdate() string {
	return `
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_snmp" "test" {
		depends_on = [meraki_organization.test]
		organization_id = resource.meraki_organization.test.organization_id
		v2c_enabled     = true
		v3_enabled      = true
		v3_auth_pass    = "bjhb989*%gg"
		v3_priv_pass    = "jkjbbbj679$%"
		v3_auth_mode    = "SHA"
		v3_priv_mode    = "AES128"
		peer_ips        = ["1.1.1.1", "2.2.2.2"]
	}
	`
}

// OrganizationsSnmpSettingsTestChecks returns the test check functions for verifying the SNMP settings
func OrganizationsSnmpSettingsTestChecks(peerIps []string, authMode, privMode string) resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"v2c_enabled":  "true",
		"v3_enabled":   "true",
		"v3_auth_mode": authMode,
		"v3_priv_mode": privMode,
		"peer_ips.#":   fmt.Sprintf("%d", len(peerIps)),
		"peer_ips.0":   peerIps[0],
	}

	// Add checks for additional peer IPs if more than one
	if len(peerIps) > 1 {
		expectedAttrs["peer_ips.1"] = peerIps[1]
	}

	return utils.ResourceTestCheck("meraki_organizations_snmp.test", expectedAttrs)
}
