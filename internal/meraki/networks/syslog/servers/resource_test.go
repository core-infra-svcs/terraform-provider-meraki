package servers_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksSyslogServersResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			/*
				{
					ResourceName:      "meraki_networks_syslog_servers.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_syslog_servers"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_syslog_servers"),
			},

			// Update and Read Networks Syslog Servers.
			{
				Config: SyslogServersResourceConfigUpdate(),
				Check:  SyslogServersResourceConfigUpdateChecks(),
			},
		},
	})
}

func SyslogServersResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_syslog_servers" "test" {
	  depends_on = [resource.meraki_network.test]
      network_id = resource.meraki_network.test.network_id
	  servers = [{
		host = "1.2.3.67"
		port = "443"
		roles = ["URLs"]
	}] 
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_syslog_servers"),
	)
}

// SyslogServersResourceConfigUpdateChecks returns the aggregated test check functions for the syslog servers resource
func SyslogServersResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"servers.0.host":    "1.2.3.67",
		"servers.0.port":    "443",
		"servers.0.roles.0": "URLs",
	}
	return utils.ResourceTestCheck("meraki_networks_syslog_servers.test", expectedAttrs)
}
