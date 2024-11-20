package networks_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksTrafficAnalysisResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_traffic_analysis.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_traffic_analysis"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_traffic_analysis"),
			},

			// Update and Read Networks Traffic Analysis.
			{
				Config: NetworksTrafficAnalysisResourceConfigUpdate(),
				Check:  NetworksTrafficAnalysisResourceConfigUpdateChecks(),
			},
		},
	})
}

func NetworksTrafficAnalysisResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_traffic_analysis" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  mode = "basic"
	  custom_pie_chart_items = [
		{
			"name": "Item from hostname",
			"type": "host",
			"value": "example.com"
		}
	  ]
	 
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "nw name"),
	)
}

// NetworksTrafficAnalysisResourceConfigUpdateChecks returns the aggregated test check functions for the traffic analysis resource
func NetworksTrafficAnalysisResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"mode":                           "basic",
		"custom_pie_chart_items.0.name":  "Item from hostname",
		"custom_pie_chart_items.0.type":  "host",
		"custom_pie_chart_items.0.value": "example.com",
	}
	return utils.ResourceTestCheck("meraki_networks_traffic_analysis.test", expectedAttrs)
}
