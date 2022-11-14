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

					/*
						// save postgres Id to variable
							resource.TestCheckResourceAttrWith("meraki_organization.test", "organization_id",
								func(value string) error {
									fmt.Println(fmt.Sprintf("ENV VAR: %s", os.Getenv("TF_MERAKI_DASHBOARD_ORG_ID")))
									fmt.Println(fmt.Sprintf("Value: %s", value))
									return nil
								}),
					*/
				),
			},

			// Update testing
			{
				// PreConfig: func() { testOrgIdExistsPreCheck(t) },
				Config: testAccOrganizationResourceUpdate(testAccOrganizationResourceConfigUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("meraki_organization.test", "organization_id", os.Getenv("TF_MERAKI_DASHBOARD_ORG_ID")),
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
	if v := os.Getenv("TF_MERAKI_DASHBOARD_ORG_ID"); v == "" {
		t.Fatal("TF_MERAKI_DASHBOARD_ORG_ID must be set for acceptance tests")
	}
}
