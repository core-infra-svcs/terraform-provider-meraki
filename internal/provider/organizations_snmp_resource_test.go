package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationsSnmpSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create Organization
			{
				Config: testAccOrganizationsSnmpSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_organizations_snmp"),
					resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),
				),
			},

			// Update and Read Org Snmp Settings.
			{
				Config: testAccOrganizationsSnmpSettingsResourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "v2c_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "v3_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "v3_auth_mode", "SHA"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "v3_priv_mode", "AES128"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "peer_ips.#", "1"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "peer_ips.0", "1.1.1.1"),
				),
			},

			// Update and Read Org Snmp Settings.
			{
				Config: testAccOrganizationsSnmpSettingsResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "v2c_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "v3_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "v3_auth_mode", "SHA"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "v3_priv_mode", "AES128"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "peer_ips.#", "2"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "peer_ips.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("meraki_organizations_snmp.test", "peer_ips.1", "2.2.2.2"),
				),
			},
		},
	})
}

const testAccOrganizationsSnmpSettingsResourceConfig = `
resource "meraki_organization" "test" {
	name = "test_acc_organizations_snmp"
	api_enabled = true
}
`

const testAccOrganizationsSnmpSettingsResourceConfigCreate = `
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

const testAccOrganizationsSnmpSettingsResourceConfigUpdate = `
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
