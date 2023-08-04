package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkGroupPoliciesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworkGroupPoliciesDataSourceConfigCreateOrganizations,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_network_group_policies"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworkGroupPoliciesDataSourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network"),
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

			// Create and Read Networks Group Policy.
			{
				Config: testAccNetworkGroupPoliciesDataSourceConfigCreateNetworksGroupPolicy,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "name", "testpolicy"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "splash_auth_settings", "network default"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bandwidth.settings", "custom"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bandwidth.bandwidth_limits.limit_up", "100000"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bandwidth.bandwidth_limits.limit_down", "100000"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "vlan_tagging.settings", "custom"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "vlan_tagging.vlan_id", "1"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.settings", "network default"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.dest_port", "443"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l3_firewall_rules.0.dest_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.0.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.0.type", "host"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.0.value", "google.com"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.1.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.1.type", "ipRange"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.1.value", "10.11.12.00/24"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.2.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.2.type", "ipRange"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.2.value", "10.11.12.00/24:5555"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.3.policy", "deny"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.3.type", "port"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.l7_firewall_rules.3.value", "23"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.dscp_tag_value", "0"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.pcp_tag_value", "0"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.settings", "custom"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_down", "100000"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_up", "100000"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.type", "host"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.value", "google.com"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.friday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.friday.from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.friday.to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.saturday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.saturday.from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.saturday.to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.sunday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.sunday.from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.sunday.to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.monday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.monday.from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.monday.to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.tuesday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.tuesday.from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.tuesday.to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.wednesday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.wednesday.from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.wednesday.to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.thursday.active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.thursday.from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling.thursday.to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding.settings", "custom"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding.rules.0.description", "A simple bonjour rule"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding.rules.0.vlan_id", "1"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding.rules.0.services.0", "All Services"),
				),
			},

			// Read test network group policies
			{
				Config: testAccNetworkGroupPoliciesDataSourceConfigRead,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.#", "1"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.name", "testpolicy"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.splash_auth_settings", "network default"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.bandwidth.settings", "custom"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.bandwidth.bandwidth_limits.limit_up", "100000"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.bandwidth.bandwidth_limits.limit_down", "100000"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.vlan_tagging.settings", "custom"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.vlan_tagging.vlan_id", "1"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.settings", "network default"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.comment", "Allow TCP traffic to subnet with HTTP servers."),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.policy", "allow"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.dest_port", "443"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l3_firewall_rules.0.dest_cidr", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.0.policy", "deny"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.0.type", "host"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.0.value", "google.com"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.1.policy", "deny"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.1.type", "ipRange"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.1.value", "10.11.12.00/24"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.2.policy", "deny"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.2.type", "ipRange"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.2.value", "10.11.12.00/24:5555"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.3.policy", "deny"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.3.type", "port"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.l7_firewall_rules.3.value", "23"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.dscp_tag_value", "0"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.pcp_tag_value", "0"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.settings", "custom"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_down", "100000"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits.bandwidth_limits.limit_up", "100000"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.type", "host"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.value", "google.com"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.enabled", "true"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.friday.active", "true"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.friday.from", "09:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.friday.to", "17:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.saturday.active", "true"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.saturday.from", "09:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.saturday.to", "17:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.sunday.active", "true"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.sunday.from", "09:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.sunday.to", "17:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.monday.active", "true"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.monday.from", "09:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.monday.to", "17:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.tuesday.active", "true"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.tuesday.from", "09:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.tuesday.to", "17:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.wednesday.active", "true"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.wednesday.from", "09:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.wednesday.to", "17:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.thursday.active", "true"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.thursday.from", "09:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.scheduling.thursday.to", "17:00"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.bonjour_forwarding.settings", "custom"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.bonjour_forwarding.rules.0.description", "A simple bonjour rule"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.bonjour_forwarding.rules.0.vlan_id", "1"),
					resource.TestCheckResourceAttr("data.meraki_network_group_policies.test", "list.0.bonjour_forwarding.rules.0.services.0", "All Services"),
				),
			},
		},
	})
}

const testAccNetworkGroupPoliciesDataSourceConfigCreateOrganizations = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_network_group_policies"
	api_enabled = true
}
`

const testAccNetworkGroupPoliciesDataSourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "test_acc_network"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

const testAccNetworkGroupPoliciesDataSourceConfigCreateNetworksGroupPolicy = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_group_policy" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
    network_id = resource.meraki_network.test.network_id
    name = "testpolicy"
    splash_auth_settings = "network default"
    scheduling = {
    enabled = true
    friday = {
    active = true
    from = "09:00"
    to = "17:00"
    },
    saturday = {
    active = true
    from = "09:00"
    to = "17:00"
    },
    sunday = { 
    active = true
    from = "09:00"
    to = "17:00"
    },
    monday = { 
    active = true
    from = "09:00"
    to = "17:00"
    },
    tuesday = {
    active = true
    from = "09:00"
    to = "17:00"
    },
    wednesday = { 
    active = true
    from = "09:00"
    to = "17:00"
    },
    thursday = { 
    active = true
    from = "09:00"
    to = "17:00"
    }
    }
	bandwidth = {
		settings = "custom"
		bandwidth_limits = {			
          limit_up = 100000
          limit_down = 100000
	    }
	}
    bonjour_forwarding = { 
        settings = "custom"
        rules = [
            {
        description = "A simple bonjour rule"
        vlan_id = "1"
        services = [ "All Services" ]
        }
     ] 
    }
    firewall_and_traffic_shaping = {
        settings = "network default"
        l3_firewall_rules = [{
            comment =  "Allow TCP traffic to subnet with HTTP servers."
            policy = "allow"
            protocol = "tcp"
            dest_port = "443"
            dest_cidr = "192.168.1.0/24"
        }]
        l7_firewall_rules = [{
            policy =  "deny"
            type = "host"
            value = "google.com"
        },
        {
            policy = "deny"
            type = "port"
            value = "23"
        },
        {
            policy = "deny"
            type = "ipRange"
            value = "10.11.12.00/24"
        },
        {
            policy = "deny",
            type = "ipRange"
            value = "10.11.12.00/24:5555"
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
                value =  "google.com"
            },
            {
                type =  "port"
                value =  "9090"
            },
            {
                type = "ipRange",
                value = "192.1.0.0"
            },
            {
                type = "ipRange"
                value = "192.1.0.0/16"
            },
            {
                type =  "ipRange"
                value = "10.1.0.0/16:80"
            },
            {
                type = "localNet"
                value = "192.168.0.0/16"
            }]
        }]
    }
    vlan_tagging = {
    settings = "custom"
    vlan_id = 1
    }
    content_filtering = {
        allowed_url_patterns = {
            patterns = []
        }
        blocked_url_categories = {
            categories = []
        }
        blocked_url_patterns = {
            patterns = []
        }
    }
    
       

  
}
`

const testAccNetworkGroupPoliciesDataSourceConfigRead = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	product_types = ["appliance", "switch", "wireless"]
}
resource "meraki_networks_group_policy" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_organization.test]
    network_id = resource.meraki_network.test.network_id
    name = "testpolicy"
    splash_auth_settings = "network default"
    scheduling = {
    enabled = true
    friday = {
    active = true
    from = "09:00"
    to = "17:00"
    },
    saturday = {
    active = true
    from = "09:00"
    to = "17:00"
    },
    sunday = { 
    active = true
    from = "09:00"
    to = "17:00"
    },
    monday = { 
    active = true
    from = "09:00"
    to = "17:00"
    },
    tuesday = {
    active = true
    from = "09:00"
    to = "17:00"
    },
    wednesday = { 
    active = true
    from = "09:00"
    to = "17:00"
    },
    thursday = { 
    active = true
    from = "09:00"
    to = "17:00"
    }
    }
	bandwidth = {
		settings = "custom"
		bandwidth_limits = {			
          limit_up = 100000
          limit_down = 100000
	    }
	}
    bonjour_forwarding = { 
        settings = "custom"
        rules = [
            {
        description = "A simple bonjour rule"
        vlan_id = "1"
        services = [ "All Services" ]
        }
     ] 
    }
    firewall_and_traffic_shaping = {
        settings = "network default"
        l3_firewall_rules = [{
            comment =  "Allow TCP traffic to subnet with HTTP servers."
            policy = "allow"
            protocol = "tcp"
            dest_port = "443"
            dest_cidr = "192.168.1.0/24"
        }]
        l7_firewall_rules = [{
            policy =  "deny"
            type = "host"
            value = "google.com"
        },
        {
            policy = "deny"
            type = "port"
            value = "23"
        },
        {
            policy = "deny"
            type = "ipRange"
            value = "10.11.12.00/24"
        },
        {
            policy = "deny",
            type = "ipRange"
            value = "10.11.12.00/24:5555"
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
                value =  "google.com"
            },
            {
                type =  "port"
                value =  "9090"
            },
            {
                type = "ipRange",
                value = "192.1.0.0"
            },
            {
                type = "ipRange"
                value = "192.1.0.0/16"
            },
            {
                type =  "ipRange"
                value = "10.1.0.0/16:80"
            },
            {
                type = "localNet"
                value = "192.168.0.0/16"
            }]
        }]
    }
    vlan_tagging = {
    settings = "custom"
    vlan_id = 1
    }
    content_filtering = {
        allowed_url_patterns = {
            patterns = []
        }
        blocked_url_categories = {
            categories = []
        }
        blocked_url_patterns = {
            patterns = []
        }
    }
}

data "meraki_network_group_policies" "test" {	
	depends_on = [resource.meraki_networks_group_policy.test, resource.meraki_network.test, resource.meraki_organization.test]
    network_id = resource.meraki_network.test.network_id
}
`
