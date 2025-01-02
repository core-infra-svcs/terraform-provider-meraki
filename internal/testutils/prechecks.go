package testutils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"

	p "github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestAccPreCheck ensures all required environment variables are set.
func TestAccPreCheck(t *testing.T) {
	requiredEnvVars := []string{
		"MERAKI_DASHBOARD_API_KEY", "TF_ACC_MERAKI_MX_LICENCE", "TF_ACC_MERAKI_MX_SERIAL",
		"TF_ACC_MERAKI_MS_SERIAL", "TF_ACC_MERAKI_MG_SERIAL", "TF_ACC_MERAKI_ORDER_NUMBER",
		"TF_ACC_MERAKI_ORGANIZATION_ID",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("%s must be set for acceptance tests", envVar)
		}
	}
}

// TestAccProtoV6ProviderFactories initializes the provider for tests.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"meraki": providerserver.NewProtocol6WithError(p.New("test")()),
}

// ResourceTestCheck validates a resource's attributes.
func ResourceTestCheck(resourceName string, expectedAttrs map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		for attr, expectedValue := range expectedAttrs {
			actualValue, exists := res.Primary.Attributes[attr]
			if !exists {
				return fmt.Errorf("attribute %q not found in resource %s", attr, resourceName)
			}
			if actualValue != expectedValue {
				return fmt.Errorf("expected %q for attribute %q but got %q", expectedValue, attr, actualValue)
			}
		}
		return nil
	}
}
