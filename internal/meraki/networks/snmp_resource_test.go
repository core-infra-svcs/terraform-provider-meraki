package networks_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSnmpSettingsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_snmp_settings"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_snmp_settings"),
			},

			// Create
			{
				Config: NetworkSnmpSettingsResourceCreate(),
				Check:  NetworkSnmpSettingsResourceCreateChecks(),
			},

			{
				Config: NetworkSnmpSettingsResourceUpdate(),
				Check:  NetworkSnmpSettingsResourceUpdateChecks(),
			},
			{
				ResourceName:      "meraki_networks_snmp.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources["meraki_networks_snmp.test"]
					if !ok {
						return "", fmt.Errorf("not found: %s", "meraki_networks_snmp.test")
					}
					return rs.Primary.ID, nil
				},
			},
		},
	})
}

func NetworkSnmpSettingsResourceCreate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_snmp" "test" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	access = "community"
	community_string = "public"
	users = []
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_snmp_settings"),
	)
}

// NetworkSnmpSettingsResourceCreateChecks returns the aggregated test check functions for the SNMP Settings resource
func NetworkSnmpSettingsResourceCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"access":           "community",
		"community_string": "public",
	}
	return utils.ResourceTestCheck("meraki_networks_snmp.test", expectedAttrs)
}

func NetworkSnmpSettingsResourceUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_snmp" "test" {
	depends_on = [resource.meraki_network.test]
	network_id = resource.meraki_network.test.network_id
	access = "users"
	users = [{
		username = "snmp_user"
		passphrase = "snmp_passphrase"
	}]
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_snmp_settings"),
	)
}

// NetworkSnmpSettingsResourceUpdateChecks returns the aggregated test check functions for the SNMP Settings resource
func NetworkSnmpSettingsResourceUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"access":             "users",
		"users.0.username":   "snmp_user",
		"users.0.passphrase": "snmp_passphrase",
	}
	return utils.ResourceTestCheck("meraki_networks_snmp.test", expectedAttrs)
}
