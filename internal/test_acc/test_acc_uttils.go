package test_acc

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"os"
)

// Setup is a placeholder for Setup code.
func Setup(ctx context.Context) {
	fmt.Println("Setup code running...")
}

// Cleanup is a placeholder for Cleanup code executed on test failures.
func Cleanup(ctx context.Context) {
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

// Teardown is a placeholder for Teardown code.
func Teardown(ctx context.Context) {
	fmt.Println("Teardown code running...")
}
