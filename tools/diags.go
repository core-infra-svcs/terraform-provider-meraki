package tools

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

// LogCurrentState LogCurrentStat Custom test check function to log the current state
func LogCurrentState(t *testing.T) resource.TestCheckFunc {
	/*
		// Example test step
		{
			Config: testAccDevicesClaimResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
			Check: resource.ComposeAggregateTestCheckFunc(
				tools.LogCurrentState(t), // Add this line to log the state
				resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_claim_device"),
			),
		},
	*/

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			rsType := fmt.Sprintf("[DEBUG] Resource: %s", rs.Type)
			fmt.Println(rsType)
			for key, value := range rs.Primary.Attributes {
				rsAttr := fmt.Sprintf("[DEBUG] %s: %s", key, value)
				fmt.Println(rsAttr)
			}
		}
		return nil
	}
}
