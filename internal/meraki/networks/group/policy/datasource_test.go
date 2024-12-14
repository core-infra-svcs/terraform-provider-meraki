package policy_test

import (
	"fmt"
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
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_group_policies"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_group_policies"),
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
		"name":                                  "testpolicy",
		"splash_auth_settings":                  "network default",
		"bandwidth.settings":                    "network default",
		"vlan_tagging.settings":                 "network default",
		"firewall_and_traffic_shaping.settings": "network default",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.comment":                                                     "Allow TCP traffic to subnet with HTTP servers.",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.policy":                                                      "allow",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.protocol":                                                    "tcp",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.dest_port":                                                   "443",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.dest_cidr":                                                   "192.168.1.0/24",
		"firewall_and_traffic_shaping.traffic_shaping_rules.0.dscp_tag_value":                                          "0",
		"firewall_and_traffic_shaping.traffic_shaping_rules.0.pcp_tag_value":                                           "0",
		"firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.settings":                    "custom",
		"firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_down": "100000",
		"firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_up":   "100000",
		"firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.type":                                      "host",
		"firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.value":                                     "google.com",
		"scheduling.enabled":          "true",
		"scheduling.friday.active":    "true",
		"scheduling.friday.from":      "00:00",
		"scheduling.friday.to":        "24:00",
		"scheduling.saturday.active":  "true",
		"scheduling.saturday.from":    "00:00",
		"scheduling.saturday.to":      "24:00",
		"scheduling.sunday.active":    "true",
		"scheduling.sunday.from":      "00:00",
		"scheduling.sunday.to":        "24:00",
		"scheduling.monday.active":    "true",
		"scheduling.monday.from":      "00:00",
		"scheduling.monday.to":        "24:00",
		"scheduling.tuesday.active":   "true",
		"scheduling.tuesday.from":     "00:00",
		"scheduling.tuesday.to":       "24:00",
		"scheduling.wednesday.active": "true",
		"scheduling.wednesday.from":   "00:00",
		"scheduling.wednesday.to":     "24:00",
		"scheduling.thursday.active":  "true",
		"scheduling.thursday.from":    "00:00",
		"scheduling.thursday.to":      "24:00",
		"bonjour_forwarding.settings": "network default",
	}
	return utils.ResourceTestCheck("data.meraki_networks_group_policies.test", expectedAttrs)
}
