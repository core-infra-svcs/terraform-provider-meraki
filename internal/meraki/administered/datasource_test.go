package administered_test

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/administered"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSource(t *testing.T) {
	// Validate schema-model consistency for the top-level DataSource schema
	t.Run("Validate Top-Level Schema", func(t *testing.T) {
		testutils.ValidateDataSourceSchemaModelConsistency(t, administered.GetDatasourceSchema.Attributes, administered.DataSourceModel{})
	})

	// Validate schema-model consistency for nested schemas
	t.Run("Validate Nested Authentication Schema", func(t *testing.T) {
		testutils.ValidateDataSourceSchemaModelConsistency(t, administered.DatasourceAuthenticationAttributes.Attributes, administered.AuthenticationModel{})
	})
	t.Run("Validate Nested API Schema", func(t *testing.T) {
		testutils.ValidateDataSourceSchemaModelConsistency(t, administered.DatasourceAPIAttributes.Attributes, administered.APIModel{})
	})
	t.Run("Validate Nested SAML Schema", func(t *testing.T) {
		testutils.ValidateDataSourceSchemaModelConsistency(t, administered.DatasourceSAMLAttributes.Attributes, administered.SAMLModel{})
	})
	t.Run("Validate Nested Two-Factor Schema", func(t *testing.T) {
		testutils.ValidateDataSourceSchemaModelConsistency(t, administered.DatasourceTwoFactorAttributes.Attributes, administered.TwoFactorModel{})
	})

	// Run Terraform acceptance test
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfig,
				Check:  testCheckAttributes(),
			},
		},
	})
}

// testCheckAttributes validates the retrieved data source attributes.
func testCheckAttributes() resource.TestCheckFunc {
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
			utils.ValidateRFC3339,
		),
	)
}

// testCheckTopLevelFields validates the top-level fields of the data source.
func testCheckTopLevelFields() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "name", "Core Infrastructure Services"),
		resource.TestCheckResourceAttr("data.meraki_administered_identities_me.test", "email", "dl-core-infra-svcs-public@starbucks.com"),
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

// Terraform configuration for the data source
const testAccDataSourceConfig = `
data "meraki_administered_identities_me" "test" {}
`
