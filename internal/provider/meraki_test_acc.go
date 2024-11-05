package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"testing"
)

func TestAccPreCheck(t *testing.T) {
	// Ensure required environment variables are set before running tests
	ev := []string{"MERAKI_DASHBOARD_API_KEY", "TF_ACC_MERAKI_MX_LICENCE", "TF_ACC_MERAKI_MX_SERIAL",
		"TF_ACC_MERAKI_MS_SERIAL", "TF_ACC_MERAKI_MG_SERIAL", "TF_ACC_MERAKI_ORDER_NUMBER", "TF_ACC_MERAKI_ORGANIZATION_ID"}
	for _, v := range ev {
		if os.Getenv(v) == "" {
			t.Fatalf("%s must be set for acceptance tests", v)
		}
	}
}

// TestAccProtoV6ProviderFactories are used to instantiate a provider during acceptance testing.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"meraki": providerserver.NewProtocol6WithError(New("test")()),
}
