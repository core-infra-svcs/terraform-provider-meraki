package provider

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
