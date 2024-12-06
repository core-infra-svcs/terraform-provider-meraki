package policy_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccNetworksGroupPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
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

			// Import test
			{
				ResourceName:            "meraki_networks_group_policy.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},

			// Update Networks Group Policy
			{
				Config: NetworksGroupPolicyResourceConfigUpdate(),
				Check:  NetworksGroupPolicyResourceConfigUpdateChecks(),
			},
		},
	})

}

func NetworksGroupPolicyResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_group_policy" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    name = "testpolicy"
    splash_auth_settings = "network default"
    scheduling = {
        enabled = true
        monday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        tuesday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        wednesday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        thursday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        friday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        saturday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        sunday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
    }

    bandwidth = {
        settings = "network default"
    }

    bonjour_forwarding = {
        settings = "network default"
        rules = []
    }
    firewall_and_traffic_shaping = {
        settings = "network default"
        l3_firewall_rules = [{
            comment = "Allow TCP traffic to subnet with HTTP servers."
            policy = "allow"
            protocol = "tcp"
            dest_port = "443"
            dest_cidr = "192.168.1.0/24"
        }]
        l7_firewall_rules = []
       traffic_shaping_rules = [{
            dscp_tag_value = 0
            pcp_tag_value = 0
            per_client_bandwidth_limits = {
                settings = "custom"
                bandwidth_limits = {
                    limit_up = 100000
                    limit_down = 100000
                }
            }
            definitions = [{
                type = "host"
                value = "google.com"
            }]
        }]
    }

    vlan_tagging = {
        settings = "network default"
		vlan_id = null
    }
    content_filtering = {
        allowed_url_patterns = {}
        blocked_url_categories = {}
        blocked_url_patterns = {}
    }
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_group_policy"),
	)
}

// NetworksGroupPolicyResourceConfigCreateChecks returns the test check functions for NetworksGroupPolicyResourceConfigCreate
func NetworksGroupPolicyResourceConfigCreateChecks() resource.TestCheckFunc {
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
	return utils.ResourceTestCheck("meraki_networks_group_policy.test", expectedAttrs)
}

func NetworksGroupPolicyResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_group_policy" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    name = "testpolicy"
    splash_auth_settings = "network default"
    scheduling = {
        enabled = true
		saturday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        friday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
		thursday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        wednesday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        tuesday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
        monday = {
            active = true
            from = "00:00"
            to = "24:00"
        }
		sunday = {
			active = true
			from = "00:00"
			to = "24:00"
		}
    }

    bandwidth = {
        settings = "network default"
    }

    bonjour_forwarding = {
        settings = "network default"
        rules = []
    }
    firewall_and_traffic_shaping = {
        settings = "network default"
        l3_firewall_rules = [{
            comment = "Allow TCP traffic to subnet with HTTP servers."
            policy = "allow"
            protocol = "tcp"
            dest_port = "443"
            dest_cidr = "192.168.1.0/24"
        }]
        l7_firewall_rules = [{
            policy = "deny"
            type = "host"
            value = "google.com"
        }]
       traffic_shaping_rules = [{
            dscp_tag_value = 0
            pcp_tag_value = 0
            per_client_bandwidth_limits = {
                settings = "custom"
                bandwidth_limits = {
                    limit_up = 100000
                    limit_down = 100000
                }
            }
            definitions = [{
                type = "host"
                value = "google.com"
            }]
        }]
    }

    vlan_tagging = {
        settings = "network default"
		vlan_id = null
    }
    content_filtering = {
        allowed_url_patterns = {}
        blocked_url_categories = {}
        blocked_url_patterns = {}
    }
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_group_policy"),
	)
}

// NetworksGroupPolicyResourceConfigUpdateChecks returns the test check functions for NetworksGroupPolicyResourceConfigUpdate
func NetworksGroupPolicyResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":                                  "testpolicy",
		"splash_auth_settings":                  "network default",
		"bandwidth.settings":                    "network default",
		"vlan_tagging.settings":                 "network default",
		"firewall_and_traffic_shaping.settings": "network default",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.comment":   "Allow TCP traffic to subnet with HTTP servers.",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.policy":    "allow",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.protocol":  "tcp",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.dest_port": "443",
		"firewall_and_traffic_shaping.l3_firewall_rules.0.dest_cidr": "192.168.1.0/24",

		"firewall_and_traffic_shaping.l7_firewall_rules.0.policy": "deny",
		"firewall_and_traffic_shaping.l7_firewall_rules.0.type":   "host",
		"firewall_and_traffic_shaping.l7_firewall_rules.0.value":  "google.com",

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
	return utils.ResourceTestCheck("meraki_networks_group_policy.test", expectedAttrs)
}
