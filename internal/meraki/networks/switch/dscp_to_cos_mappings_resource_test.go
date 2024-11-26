package _switch_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccNetworksSwitchDscpToCosMappingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_switch_dscp_to_cos_mappings"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_switch_dscp_to_cos_mappings"),
			},

			// Create and Read Test
			{
				Config: NetworksSwitchDscpToCosMappingsResourceConfigCreate(),
				Check:  NetworksSwitchDscpToCosMappingsResourceConfigCreateChecks(),
			},

			// Update and Read Test
			{
				Config: NetworksSwitchDscpToCosMappingsResourceConfigUpdate(),
				Check:  NetworksSwitchDscpToCosMappingsResourceConfigUpdateChecks(),
			},
		},
	})
}

func NetworksSwitchDscpToCosMappingsResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_switch_dscp_to_cos_mappings" "test" {
  depends_on                = [meraki_network.test]
  network_id                = resource.meraki_network.test.network_id
  mappings = [
	{
		dscp = 1
		cos = 1
	}
  ]
}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_switch_dscp_to_cos_mappings"),
	)
}

// NetworksSwitchDscpToCosMappingsResourceConfigCreateChecks returns the test check functions for NetworksSwitchDscpToCosMappingsResourceConfigCreate
func NetworksSwitchDscpToCosMappingsResourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"mappings.0.dscp": "1",
		"mappings.0.cos":  "1",
	}
	return utils.ResourceTestCheck("meraki_networks_switch_dscp_to_cos_mappings.test", expectedAttrs)
}

func NetworksSwitchDscpToCosMappingsResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_switch_dscp_to_cos_mappings" "test" {
  depends_on                = [meraki_network.test]
  network_id                = resource.meraki_network.test.network_id
  mappings = [
	{
		dscp = 63
		cos = 5
	}
  ]
}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_switch_dscp_to_cos_mappings"),
	)
}

// NetworksSwitchDscpToCosMappingsResourceConfigUpdateChecks returns the test check functions for NetworksSwitchDscpToCosMappingsResourceConfigUpdate
func NetworksSwitchDscpToCosMappingsResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"mappings.0.dscp": "63",
		"mappings.0.cos":  "5",
	}
	return utils.ResourceTestCheck("meraki_networks_switch_dscp_to_cos_mappings.test", expectedAttrs)
}
