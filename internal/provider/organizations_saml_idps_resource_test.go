package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrganizations{Organizationid}SamlIdps{Idpid}Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccOrganizations{Organizationid}SamlIdps{Idpid}ResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(

				    // TODO - Check return data matches expected result
					// TODO - Example: resource.TestCheckResourceAttr("meraki_organizations_{organization_id}_saml_idps_{idp_id}.test", "name", "testOrg1"),
				),
			},

			// TODO - Update+Read Test

			// Delete testing automatically occurs in TestCase
		},
	})
}

const testAccOrganizations{Organizationid}SamlIdps{Idpid}ResourceConfig = `
resource "meraki_organizations_{organization_id}_saml_idps_{idp_id}" "test" {
    // TODO - Add configuration
}
`

