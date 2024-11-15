package networks_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_settings"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_settings"),
			},

			// Update and Read Network Settings.
			{
				Config: NetworkSettingsResourceConfigUpdate(),
				Check:  NetworkSettingsResourceConfigUpdateChecks(),
			},

			/*
				{
						ResourceName:      "meraki_networks_settings.test",
						ImportState:       true,
						ImportStateVerify: true,
					},
			*/
		},
	})
}

func NetworkSettingsResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_settings" "test" {
	depends_on = [resource.meraki_network.test]
  	network_id = resource.meraki_network.test.network_id
	local_status_page_enabled = false
	remote_status_page_enabled = false
	local_status_page = { 
		authentication = { 
			enabled = true
			username = "admin"
			password = "testpassword"
	}
	}
  	secure_port_enabled = false
	named_vlans_enabled = true
}
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_settings"),
	)
}

// NetworkSettingsResourceConfigUpdateChecks returns the aggregated test check functions for the settings resource
func NetworkSettingsResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"local_status_page_enabled":                 "false",
		"remote_status_page_enabled":                "false",
		"local_status_page.authentication.enabled":  "true",
		"local_status_page.authentication.username": "admin",
		"local_status_page.authentication.password": "testpassword",
		"secure_port_enabled":                       "false",
	}
	return utils.ResourceTestCheck("meraki_networks_settings.test", expectedAttrs)
}
