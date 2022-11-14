package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizationResource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizationResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "testOrg1"),
					resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organization.test", "id", "example-id"),
					// save postgres Id to variable
					resource.TestCheckResourceAttrWith("meraki_organization.test", "organization_id",
						func(value string) error {
							err := os.Setenv("TF_MERAKI_DASHBOARD_ORGANIZATION_ID", value)
							if err != nil {
								return fmt.Errorf(fmt.Sprintf("Unable to add id to Env Var. Value: %s", value))
							}
							return nil
						}),
				),
			},

			// Update testing
			{
				PreConfig: func() { testOrgIdExistsPreCheck(t) },
				Config:    testAccOrganizationResourceConfigUpdate, //testAccOrganizationResourceUpdate(os.Getenv("TF_MERAKI_DASHBOARD_ORGANIZATION_ID")),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_organization.test", "organization_id", os.Getenv("TF_MERAKI_DASHBOARD_ORGANIZATION_ID")),
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "testOrg2"),
					resource.TestCheckResourceAttr("meraki_organization.test", "api_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_organization.test", "management_details_name", "MSP ID"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizationResourceConfig = `
resource "meraki_organization" "test" {
	name = "testOrg1"
	api_enabled = true
}
`

const testAccOrganizationResourceConfigUpdate = `
resource "meraki_organization" "test" {
	name = "testOrg2"
	api_enabled = true
	management_details_name = "MSP ID"
	management_details_value = "123456"
}
`

func testAccOrganizationResourceUpdate(postgresId string) string {
	return fmt.Sprintf(`
resource "meraki_organization" "test" {
	organization_id = "%s"
	name = "testOrg2"
	api_enabled = true
	management_details_name = "MSP ID"
	management_details_value = "123456"

}
`, postgresId)
}

func testOrgIdExistsPreCheck(t *testing.T) {
	if v := os.Getenv("TF_MERAKI_DASHBOARD_ORGANIZATION_ID"); v == "" {
		t.Error(fmt.Sprintf("Unable to read id to Env Var: %s", os.Getenv("TF_MERAKI_DASHBOARD_ORGANIZATION_ID")))
		t.Fatal("TF_MERAKI_DASHBOARD_ORGANIZATION_ID must be set for acceptance tests")
	}
}
