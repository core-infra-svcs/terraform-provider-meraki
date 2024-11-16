package administered_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccAdministeredIdentitiesMeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read AdministeredIdentitiesMe
			{
				Config: testAccAdministeredIdentitiesMeDataSourceConfigCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					//testCheckTopLevelFields(),
					resource.TestCheckResourceAttrWith(
						"data.meraki_administered_identities_me.test", "last_used_dashboard_at", validateRFC3339),
					testCheckAuthenticationFields(),
				),
			},
		},
	})
}

const testAccAdministeredIdentitiesMeDataSourceConfigCreate = `
data "meraki_administered_identities_me" "test" {
}
`

// Helper function to check top-level fields
func testCheckTopLevelFields() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "name", "Miles Meraki"),
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "email", "miles@meraki.com"),
	)
}

// Helper function to check authentication fields
func testCheckAuthenticationFields() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "authentication.mode", "email"),
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "authentication.api_key_created", "true"),
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "authentication.saml_enabled", "false"),
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "authentication.two_factor.enabled", "false"),
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "authentication.api.key.created", "true"),
	)
}

// Helper function to validate RFC3339 format
func validateRFC3339(value string) error {
	re := regexp.MustCompile(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)$`)
	if !re.MatchString(value) {
		return fmt.Errorf("received timestamp does not match RFC3339 format: %s", value)
	}
	return nil
}
