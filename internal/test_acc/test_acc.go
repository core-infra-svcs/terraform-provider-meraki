package test_acc

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"testing"
)

func TestAccPreCheck(t *testing.T) {
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

// TestAccProtoV6ProviderFactories are used to instantiate a provider during acceptance testing.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"meraki": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// EntryPoint is the entry point for the test suite.
func EntryPoint(m *testing.M) {
	ctx := context.Background()
	fmt.Println("Starting EntryPoint")

	// Setup code here (e.g., initialize resources, set environment variables).
	Setup(ctx)

	// Run the tests
	fmt.Println("Running tests")
	exitCode := m.Run()

	// Run sweepers or other Cleanup code if tests failed.
	if exitCode != 0 {
		fmt.Println("Tests failed, running Cleanup")
		Cleanup(ctx)
	}

	// Additional Cleanup code here (e.g., close connections, remove files).
	Teardown(ctx)

	fmt.Println("Exiting EntryPoint", map[string]interface{}{"exitCode": exitCode})
	os.Exit(exitCode)
}
