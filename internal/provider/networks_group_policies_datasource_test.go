package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkGroupPolicysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworkGroupPolicysDataSourceConfigCreateOrganizations,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_network_group_policys"),
				),
			},

			// Create and Read Network.
			{
				Config: testAccNetworkGroupPolicysDataSourceConfigCreateNetwork,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_network.test", "name", "Main Office"),
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
				Config: testAccNetworkGroupPolicysDataSourceConfigCreateNetworksGroupPolicy,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "name", "testpolicy"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "splash_auth_settings", "network default"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bandwidth_settings", "custom"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bandwidth_limit_up", "100000"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bandwidth_limit_down", "100000"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "vlan_tagging_settings", "custom"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "vlan_tagging_vlan_id", "1"),
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
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.per_client_bandwidth_limits_settings", "network default"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.type", "host"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "firewall_and_traffic_shaping.traffic_shaping_rules.0.definitions.0.value", "google.com"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_enabled", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_friday_active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_friday_from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_friday_to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_saturday_active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_saturday_from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_saturday_to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_sunday_active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_sunday_from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_sunday_to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_monday_active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_monday_from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_monday_to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_tuesday_active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_tuesday_from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_tuesday_to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_wednesday_active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_wednesday_from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_wednesday_to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_thursday_active", "true"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_thursday_from", "09:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "scheduling_thursday_to", "17:00"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding_rules.0.description", "A simple bonjour rule"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding_rules.0.vlan_id", "1"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "bonjour_forwarding_rules.0.services.0", "All Services"),
				),
			},

			// Read test group policies
			{
				Config: testAccNetworkGroupPolicysDataSourceConfigRead,
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

const testAccNetworkGroupPolicysDataSourceConfigCreateOrganizations = `
resource "meraki_organization" "test" {
	name = "test_acc_meraki_network_group_policys"
	api_enabled = true
}
`

const testAccNetworkGroupPolicysDataSourceConfigCreateNetwork = `
resource "meraki_organization" "test" {}
resource "meraki_network" "test" {
	depends_on = [resource.meraki_organization.test]
	organization_id = resource.meraki_organization.test.organization_id
	product_types = ["appliance", "switch", "wireless"]
	tags = ["tag1"]
	name = "Main Office"
	timezone = "America/Los_Angeles"
	notes = "Additional description of the network"
}
`

const testAccNetworkGroupPolicysDataSourceConfigCreateNetworksGroupPolicy = `
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
    scheduling_enabled = true
    scheduling_friday_active = true
    scheduling_friday_from = "09:00"
    scheduling_friday_to = "17:00"
    scheduling_saturday_active = true
    scheduling_saturday_from = "09:00"
    scheduling_saturday_to = "17:00"
    scheduling_sunday_active = true
    scheduling_sunday_from = "09:00"
    scheduling_sunday_to = "17:00"
    scheduling_monday_active = true
    scheduling_monday_from = "09:00"
    scheduling_monday_to = "17:00"
    scheduling_tuesday_active = true
    scheduling_tuesday_from = "09:00"
    scheduling_tuesday_to = "17:00"
    scheduling_wednesday_active = true
    scheduling_wednesday_from = "09:00"
    scheduling_wednesday_to = "17:00"
    scheduling_thursday_active = true
    scheduling_thursday_from = "09:00"
    scheduling_thursday_to = "17:00"
    bandwidth_settings = "custom"
    bandwidth_limit_up = 100000
    bandwidth_limit_down = 100000
    bonjour_forwarding_rules = [{
        description = "A simple bonjour rule"
        vlan_id = "1"
        services = [ "All Services" ]
    }]
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
            per_client_bandwidth_limits_settings = "network default"
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
    vlan_tagging_settings = "custom"
    vlan_tagging_vlan_id = 1
    content_filtering_allow_url_patterns = []
    content_filtering_blocked_url_categories = []
    content_filtering_blocked_url_patterns = []
}
`

const testAccNetworkGroupPolicysDataSourceConfigRead = `
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
    scheduling_enabled = true
    scheduling_friday_active = true
    scheduling_friday_from = "09:00"
    scheduling_friday_to = "17:00"
    scheduling_saturday_active = true
    scheduling_saturday_from = "09:00"
    scheduling_saturday_to = "17:00"
    scheduling_sunday_active = true
    scheduling_sunday_from = "09:00"
    scheduling_sunday_to = "17:00"
    scheduling_monday_active = true
    scheduling_monday_from = "09:00"
    scheduling_monday_to = "17:00"
    scheduling_tuesday_active = true
    scheduling_tuesday_from = "09:00"
    scheduling_tuesday_to = "17:00"
    scheduling_wednesday_active = true
    scheduling_wednesday_from = "09:00"
    scheduling_wednesday_to = "17:00"
    scheduling_thursday_active = true
    scheduling_thursday_from = "09:00"
    scheduling_thursday_to = "17:00"
    bandwidth_settings = "custom"
    bandwidth_limit_up = 100000
    bandwidth_limit_down = 100000
    bonjour_forwarding_rules = [{
        description = "A simple bonjour rule"
        vlan_id = "1"
        services = [ "All Services" ]
    }]
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
            per_client_bandwidth_limits_settings = "network default"
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
    vlan_tagging_settings = "custom"
    vlan_tagging_vlan_id = 1
    content_filtering_allow_url_patterns = []
    content_filtering_blocked_url_categories = []
    content_filtering_blocked_url_patterns = []
}

data "meraki_network_group_policys" "test" {
	
	depends_on = [resource.meraki_networks_group_policy.test, resource.meraki_network.test, resource.meraki_organization.test]
    network_id = resource.meraki_network.test.network_id
}
`
