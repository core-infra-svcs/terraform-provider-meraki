package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
	"time"
)

func TestAccNetworksGroupPolicyResource(t *testing.T) {
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	policies := 2 // Number of policies to create, correlates to the amount of retries set for MaximumRetries
	uniqueId := fmt.Sprintf("%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Network
			{
				Config: testAccNetworksGroupPolicyResourceConfigCreateNetwork(orgId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_networks_group_policy"),
					resource.TestCheckResourceAttr("meraki_network.test", "timezone", "America/Los_Angeles"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("meraki_network.test", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.#", "3"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.0", "appliance"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.1", "switch"),
					resource.TestCheckResourceAttr("meraki_network.test", "product_types.2", "wireless"),
					resource.TestCheckResourceAttr("meraki_network.test", "notes", "Additional description of the network"),
				),
			},

			// Create and Read Networks Group Policy
			{
				Config: testAccNetworksGroupPolicyResourceConfigCreateNetworksGroupPolicy,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "name", "testpolicy"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "splash_auth_settings", "network default"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bandwidth.settings", "network default"),
					resource.TestCheckNoResourceAttr("meraki_networks_group_policy.test", "bandwidth.bandwidth_limits.limit_up"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "vlan_tagging.settings", "network default"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.settings", "network default"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.dest_port", "443"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.dest_cidr", "192.168.1.0/24"),
					//resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.0.policy", "deny"),
					//resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.0.type", "host"),
					//resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.0.value", "google.com"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.dscp_tag_value", "0"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.pcp_tag_value", "0"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.settings", "custom"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_down", "100000"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_up", "100000"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.type", "host"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.value", "google.com"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.friday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.friday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.friday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.saturday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.saturday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.saturday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.sunday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.sunday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.sunday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.monday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.monday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.monday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.tuesday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.tuesday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.tuesday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.wednesday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.wednesday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.wednesday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.thursday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.thursday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.thursday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding.settings", "network default"),
				),
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
				Config: testAccNetworksGroupPolicyResourceConfigUpdateNetworksGroupPolicy,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "name", "testpolicy"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bandwidth.settings", "network default"),
					resource.TestCheckNoResourceAttr("meraki_networks_group_policy.test", "bandwidth.bandwidth_limits.limit_up"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "vlan_tagging.settings", "network default"),
					//resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "vlan_tagging.vlan_id", "null"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.settings", "network default"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.dest_port", "443"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.dest_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.0.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.0.type", "host"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.0.value", "google.com"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.dscp_tag_value", "0"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.pcp_tag_value", "0"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.settings", "custom"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_down", "100000"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_up", "100000"),
					//resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.type", "host"),
					//resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.value", "test.com"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.friday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.friday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.friday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.saturday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.saturday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.saturday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.sunday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.sunday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.sunday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.monday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.monday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.monday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.tuesday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.tuesday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.tuesday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.wednesday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.wednesday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.wednesday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.thursday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.thursday.from", "00:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.thursday.to", "24:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding.settings", "network default"),
					//resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding.rules.description", "update simple bonjour rule"),
					//resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding.rules.vlan_id", "2"),
					//resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding.rules.services.0", "AirPlay"),
				),
			},

			// Test the creation of multiple group policies
			{
				Config: testAccNetworksGroupPolicyResourceConfigMultiplePolicies(orgId, policies, uniqueId),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						var checks []resource.TestCheckFunc
						// Dynamically generate checks for each group policy
						for i := 1; i <= policies; i++ {
							resourceName := fmt.Sprintf("meraki_networks_group_policy.test%d", i)
							checks = append(checks,
								resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("test_acc_policy_%d_%s", i, uniqueId)),
								resource.TestCheckResourceAttr(resourceName, "splash_auth_settings", "network default"),
								resource.TestCheckResourceAttr(resourceName, "bandwidth.settings", "network default"),
								resource.TestCheckNoResourceAttr(resourceName, "bandwidth.bandwidth_limits.limit_up"),
								resource.TestCheckResourceAttr(resourceName, "vlan_tagging.settings", "network default"),
								resource.TestCheckResourceAttr(resourceName, "bonjour_forwarding.settings", "network default"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.settings", "network default"),

								// Not working
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.l3_firewall_rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.l3_firewall_rules.0.policy", "allow"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.l3_firewall_rules.0.protocol", "tcp"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.l3_firewall_rules.0.dest_port", "443"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.l3_firewall_rules.0.dest_cidr", "192.168.1.0/24"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.l7_firewall_rules.0.policy", "deny"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.l7_firewall_rules.0.type", "host"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.l7_firewall_rules.0.value", "google.com"),

								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.traffic_shaping_rules.0.dscp_tag_value", "0"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.traffic_shaping_rules.0.pcp_tag_value", "0"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.settings", "custom"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_down", "100000"),
								resource.TestCheckResourceAttr(resourceName, "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_up", "100000"),

								resource.TestCheckResourceAttr(resourceName, "scheduling.enabled", "true"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.friday.active", "true"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.friday.from", "00:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.friday.to", "24:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.saturday.active", "true"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.saturday.from", "00:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.saturday.to", "24:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.sunday.active", "true"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.sunday.from", "00:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.sunday.to", "24:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.monday.active", "true"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.monday.from", "00:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.monday.to", "24:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.tuesday.active", "true"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.tuesday.from", "00:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.tuesday.to", "24:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.wednesday.active", "true"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.wednesday.from", "00:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.wednesday.to", "24:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.thursday.active", "true"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.thursday.from", "00:00"),
								resource.TestCheckResourceAttr(resourceName, "scheduling.thursday.to", "24:00"),
							)
						}
						return resource.ComposeAggregateTestCheckFunc(checks...)(s)
					},
				),
			},
		},
	})

}

func testAccNetworksGroupPolicyResourceConfigCreateNetwork(orgId string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = "%s"
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_networks_group_policy"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}
`, orgId)
}

const testAccNetworksGroupPolicyResourceConfigCreateNetworksGroupPolicy = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}

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
`

const testAccNetworksGroupPolicyResourceConfigUpdateNetworksGroupPolicy = `
resource "meraki_network" "test" {
    product_types = ["appliance", "switch", "wireless"]
}
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
`

func testAccNetworksGroupPolicyResourceConfigMultiplePolicies(orgId string, policies int, uniqueId string) string {
	config := fmt.Sprintf(`
resource "meraki_network" "test" {
	organization_id = "%s"
    product_types = ["appliance", "switch", "wireless"]
    tags = ["tag1"]
    name = "test_acc_networks_group_policy"
    timezone = "America/Los_Angeles"
    notes = "Additional description of the network"
}

`, orgId)

	for i := 1; i <= policies; i++ {
		config += fmt.Sprintf(`
resource "meraki_networks_group_policy" "test%d" {
  depends_on = [meraki_network.test]
  network_id = meraki_network.test.network_id
  name   = "test_acc_policy_%d_%s"
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
`, i, i, uniqueId)
	}
	return config
}
