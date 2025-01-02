package policy_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkGroupPoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_group_policy"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_group_policy"),
			},

			// Create and Read Networks Group Policy
			{
				Config: NetworksGroupPolicyResourceConfigCreate(),
				Check:  NetworksGroupPolicyResourceConfigCreateChecks(),
			},

			// Read test network group policies
			{
				Config: NetworkGroupPoliciesDataSourceConfigRead(),
				Check:  NetworkGroupPoliciesDataSourceConfigReadChecks(),
			},
		},
	})
}

func NetworkGroupPoliciesDataSourceConfigRead() string {
	return fmt.Sprintf(`
	%s
data "meraki_networks_group_policies" "test" {	
	depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_group_policy"),
	)
}

// NetworkGroupPoliciesDataSourceConfigReadChecks returns the test check functions for NetworkGroupPoliciesDataSourceConfigRead
func NetworkGroupPoliciesDataSourceConfigReadChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"list.0.name":                                                                                                         "testpolicy",
		"list.0.splash_auth_settings":                                                                                         "network default",
		"list.0.bandwidth.settings":                                                                                           "network default",
		"list.0.vlan_tagging.settings":                                                                                        "network default",
		"list.0.firewall_and_traffic_shaping.settings":                                                                        "network default",
		"list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.comment":                                                     "Allow TCP traffic to subnet with HTTP servers.",
		"list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.policy":                                                      "allow",
		"list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.protocol":                                                    "tcp",
		"list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.dest_port":                                                   "443",
		"list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.dest_cidr":                                                   "192.168.1.0/24",
		"list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.dscp_tag_value":                                          "0",
		"list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.pcp_tag_value":                                           "0",
		"list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.settings":                    "custom",
		"list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_down": "100000",
		"list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_up":   "100000",
		"list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.type":                                      "host",
		"list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.value":                                     "google.com",
		"list.0.scheduling.enabled":                                                                                           "true",
		"list.0.scheduling.friday.active":                                                                                     "true",
		"list.0.scheduling.friday.from":                                                                                       "00:00",
		"list.0.scheduling.friday.to":                                                                                         "24:00",
		"list.0.scheduling.saturday.active":                                                                                   "true",
		"list.0.scheduling.saturday.from":                                                                                     "00:00",
		"list.0.scheduling.saturday.to":                                                                                       "24:00",
		"list.0.scheduling.sunday.active":                                                                                     "true",
		"list.0.scheduling.sunday.from":                                                                                       "00:00",
		"list.0.scheduling.sunday.to":                                                                                         "24:00",
		"list.0.scheduling.monday.active":                                                                                     "true",
		"list.0.scheduling.monday.from":                                                                                       "00:00",
		"list.0.scheduling.monday.to":                                                                                         "24:00",
		"list.0.scheduling.tuesday.active":                                                                                    "true",
		"list.0.scheduling.tuesday.from":                                                                                      "00:00",
		"list.0.scheduling.tuesday.to":                                                                                        "24:00",
		"list.0.scheduling.wednesday.active":                                                                                  "true",
		"list.0.scheduling.wednesday.from":                                                                                    "00:00",
		"list.0.scheduling.wednesday.to":                                                                                      "24:00",
		"list.0.scheduling.thursday.active":                                                                                   "true",
		"list.0.scheduling.thursday.from":                                                                                     "00:00",
		"list.0.scheduling.thursday.to":                                                                                       "24:00",
		"list.0.bonjour_forwarding.settings":                                                                                  "network default",
	}
	return utils.ResourceTestCheck("data.meraki_networks_group_policies.test", expectedAttrs)
}
