package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"meraki": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	// Check environmental variables
	ev := []string{"MERAKI_DASHBOARD_API_KEY", "TF_ACC_MERAKI_MX_LICENCE", "TF_ACC_MERAKI_MX_SERIAL",
		"TF_ACC_MERAKI_MS_SERIAL", "TF_ACC_MERAKI_MG_SERIAL", "TF_ACC_MERAKI_ORDER_NUMBER", "TF_ACC_MERAKI_ORGANIZATION_ID"}
	for _, v := range ev {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for acceptance tests", v)
		}
	}
}
