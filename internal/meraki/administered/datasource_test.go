package administered_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAdministeredIdentitiesMeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAdministeredIdentitiesMeDataSourceConfig,
				Check:  testCheckAdministeredIdentitiesMe(),
			},
		},
	})
}

// testCheckAdministeredIdentitiesMe validates the retrieved data source attributes.
func testCheckAdministeredIdentitiesMe() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		// Check the `id` field
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "id", "identities_me"),
		// Validate top-level fields
		testCheckTopLevelFields(),
		// Validate authentication fields
		testCheckAuthenticationFields(),
		// Validate `last_used_dashboard_at` format
		resource.TestCheckResourceAttrWith(
			"data.meraki_administered_identities_me.test",
			"last_used_dashboard_at",
			validateRFC3339,
		),
	)
}

// testCheckTopLevelFields validates the top-level fields of the data source.
func testCheckTopLevelFields() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
	//resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "name", "Miles Meraki"),
	//resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "email", "miles@meraki.com"),
	)
}

// testCheckAuthenticationFields validates the authentication-related fields.
func testCheckAuthenticationFields() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "authentication.mode", "email"),
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "authentication.saml.enabled", "false"),
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "authentication.two_factor.enabled", "false"),
		resource.TestCheckResourceAttr(
			"data.meraki_administered_identities_me.test", "authentication.api.key.created", "true"),
	)
}

// validateRFC3339 checks if a string is a valid RFC3339 date-time format.
func validateRFC3339(value string) error {
	ctx := context.Background()
	if err := utils.ValidateRFC3339(ctx, value); err != nil {
		return fmt.Errorf("value validation failed: %v", err)
	}
	return nil
}

// Terraform configuration for the data source
const testAccAdministeredIdentitiesMeDataSourceConfig = `
data "meraki_administered_identities_me" "test" {}
`
