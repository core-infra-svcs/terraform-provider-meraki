package netflow_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksNetFlowResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_netflow"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_netflow"),
			},

			// Update and Read Networks NetFlow.
			{
				Config: NetFlowResourceConfigUpdateSettings(),
				Check:  NetFlowResourceConfigUpdateSettingsChecks(),
			},

			/*
				// Import testing
					{
						ResourceName:      "meraki_networks_netflow.test",
						ImportState:       true,
						ImportStateVerify: true,
					},

			*/
		},
	})
}

func NetFlowResourceConfigUpdateSettings() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_netflow" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  reporting_enabled = false     
      eta_enabled = false   
	  collector_ip = "1.2.3.4"
      collector_port = 443 
	  eta_dst_port = 443	  
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_netflow"),
	)
}

// NetFlowResourceConfigUpdateSettingsChecks returns the aggregated test check functions for the netflow resource
func NetFlowResourceConfigUpdateSettingsChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"reporting_enabled": "false",
		"eta_enabled":       "false",
		"collector_ip":      "1.2.3.4",
		"collector_port":    "443",
		"eta_dst_port":      "443",
	}
	return utils.ResourceTestCheck("meraki_networks_netflow.test", expectedAttrs)
}
