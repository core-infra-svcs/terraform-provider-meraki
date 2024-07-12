package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"os"
	"testing"
)

// TestMain is the entry point for the test suite.
func TestMain(m *testing.M) {
	ctx := context.Background()
	fmt.Println("Starting TestMain")

	// Setup code here (e.g., initialize resources, set environment variables).
	setup(ctx)

	// Run the tests
	fmt.Println("Running tests")
	exitCode := m.Run()

	// Run sweepers or other cleanup code if tests failed.
	if exitCode != 0 {
		fmt.Println("Tests failed, running cleanup")
		cleanup(ctx)
	}

	// Additional cleanup code here (e.g., close connections, remove files).
	teardown(ctx)

	fmt.Println("Exiting TestMain", map[string]interface{}{"exitCode": exitCode})
	os.Exit(exitCode)
}

// setup is a placeholder for setup code.
func setup(ctx context.Context) {
	fmt.Println("Setup code running...")
}

// cleanup is a placeholder for cleanup code executed on test failures.
func cleanup(ctx context.Context) {
	fmt.Println("Cleanup code running...")

	organizationId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	if organizationId == "" {
		fmt.Println("TF_ACC_MERAKI_ORGANIZATION_ID must be set for sweeper to run")
		os.Exit(1)
	}

	client, clientErr := provider.SweeperHTTPClient()
	if clientErr != nil {
		fmt.Println("Error getting HTTP client", map[string]interface{}{
			"error": clientErr,
		})
	}

	// Set default retry and wait limit for provider client
	client.GetConfig().MaximumRetries = 3
	client.GetConfig().Retry4xxErrorWaitTime = 5

	// Sweep a Specified Static Organization
	fmt.Println("Running terraform sweepers due to test failures...")
	err := provider.SweepMerakiOrganization(ctx, client, organizationId)
	if err != nil {
		fmt.Println("Error running organization sweeper", map[string]interface{}{
			"error": err,
		})
	} else {
		fmt.Println("Organization sweeper ran successfully")
	}

	// Targeted "test_acc" Organizations Sweeper
	err = provider.SweepMerakiOrganizations(ctx, client)
	if err != nil {
		fmt.Println("Error running organizations sweeper", map[string]interface{}{
			"error": err,
		})
	} else {
		fmt.Println("Organizations sweeper ran successfully")
	}

}

// teardown is a placeholder for teardown code.
func teardown(ctx context.Context) {
	fmt.Println("Teardown code running...")
}
