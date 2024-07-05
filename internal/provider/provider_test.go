package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"testing"
)

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

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"meraki": providerserver.NewProtocol6WithError(New("test")()),
}

// TestMain is the entry point for the test suite.
func TestMain(m *testing.M) {
	ctx := context.Background()
	fmt.Println(ctx, "Starting TestMain")

	// Setup code here (e.g., initialize resources, set environment variables).
	setup(ctx)

	// Run the tests
	fmt.Println(ctx, "Running tests")
	exitCode := m.Run()

	// Run sweepers or other cleanup code if tests failed.
	if exitCode != 0 {
		fmt.Println(ctx, "Tests failed, running cleanup")
		cleanup(ctx)
	}

	// Additional cleanup code here (e.g., close connections, remove files).
	teardown(ctx)

	fmt.Println(ctx, "Exiting TestMain", map[string]interface{}{"exitCode": exitCode})
	os.Exit(exitCode)
}

// setup is a placeholder for setup code.
func setup(ctx context.Context) {
	fmt.Println(ctx, "Setup code running...")
}

// cleanup is a placeholder for cleanup code executed on test failures.
func cleanup(ctx context.Context) {
	fmt.Println(ctx, "Cleanup code running...")

	organizationId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	if organizationId == "" {
		fmt.Println(ctx, "TF_ACC_MERAKI_ORGANIZATION_ID must be set for sweeper to run")
		os.Exit(1)
	}

	client, clientErr := SweeperHTTPClient()
	if clientErr != nil {
		fmt.Println(ctx, "Error getting HTTP client", map[string]interface{}{
			"error": clientErr,
		})
	}

	// Sweep a Specified Static Organization
	fmt.Println(ctx, "Running terraform sweepers due to test failures...")
	err := sweepMerakiOrganization(ctx, client, organizationId)
	if err != nil {
		fmt.Println(ctx, "Error running organization sweeper", map[string]interface{}{
			"error": err,
		})
	} else {
		fmt.Println(ctx, "Organization sweeper ran successfully")
	}

	// Targeted "test_acc" Organizations Sweeper
	err = sweepMerakiOrganizations(ctx, client)
	if err != nil {
		fmt.Println(ctx, "Error running organizations sweeper", map[string]interface{}{
			"error": err,
		})
	} else {
		fmt.Println(ctx, "Organizations sweeper ran successfully")
	}

}

// teardown is a placeholder for teardown code.
func teardown(ctx context.Context) {
	fmt.Println(ctx, "Teardown code running...")
}
