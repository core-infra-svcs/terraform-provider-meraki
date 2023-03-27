package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksGroupPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create test Organization
			{
				Config: testAccNetworksGroupPolicyResourceConfigCreateOrganization,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_organization.test", "name", "test_acc_meraki_organizations_networks_group_policy"),
				),
			},

			// TODO - ImportState testing - This only works when hard-coded networkId.
			/*
				{
					ResourceName:      "meraki_networks_group_policy.test",
					ImportState:       true,
					ImportStateVerify: false,
					ImportStateId:     "657525545596096508",
				},
			*/

			// Create and Read Network.
			{
				Config: testAccNetworksGroupPolicyResourceConfigCreateNetwork,
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
				Config: testAccNetworksGroupPolicyResourceConfigCreateNetworksGroupPolicy,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "name", "testpolicy"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "splash_auth_settings", "network default"),
				),
			},

			// Update testing
			{
				Config: testAccNetworksGroupPolicyResourceConfigUpdateNetworksGroupPolicy,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "name", "testpolicy"),
					resource.TestCheckResourceAttr("meraki_networks_group_policy.test", "splash_auth_settings", "network default"),
				),
			},
		},
	})
}

const testAccNetworksGroupPolicyResourceConfigCreateOrganization = `
 resource "meraki_organization" "test" {
 	name = "test_acc_meraki_organizations_networks_group_policy"
 	api_enabled = true
 }
 `
const testAccNetworksGroupPolicyResourceConfigCreateNetwork = `
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

const testAccNetworksGroupPolicyResourceConfigCreateNetworksGroupPolicy = `
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
	  bandwidth = {
		settings = "network default" 
		bandwidth_limits = {
			limit_down = 567
			limit_up = 879
		}
	}
	bonjour_forwarding = {
        settings = "custom"
        rules = [
            {
                description = "A simple bonjour rule"
                vlan_id = 1
                services = [ "All Services" ]
            }
        ]
    }
	content_filtering = {
        allowed_url_patterns = {}
		blocked_url_categories = {}
		blocked_url_patterns = {}
    }

	scheduling =  {
        enabled = true
        monday = {
            active = true
            from = "9:00"
            to = "17:00"
        },
        tuesday = {
            active = true
            from = "9:00"
            to = "17:00"
           
        }
        wednesday = {
			active = true
            from = "9:00"
            to = "17:00"
        }
        thursday = {
            active = true
            from = "9:00"
            to = "17:00"
        }
        friday = {
			active = true
            from = "9:00"
            to = "17:00"
        }
        saturday = {
            active = true
            from = "9:00"
            to = "17:00"
        }
        sunday = {
            active = true
            from = "9:00"
            to = "17:00"
        }
    }

	vlan_tagging = {
        settings = "custom"
        vlan_id =  "1"
    }

	firewall_and_traffic_shaping =  {
        settings = "custom"
        traffic_shaping_rules = [
            {
                definitions = [
                    {
                        type = "host"
                        value = "google.com"
                    },
                    {
                        type = "port"
                        value ="9090"
                    },
                    {
                        type = "ipRange"
                        value = "192.1.0.0"
                    },
                    {
                        type = "ipRange"
                        value = "192.1.0.0/16"
                    },
                    {
                        type = "ipRange"
                        value = "10.1.0.0/16:80"
                    },
                    {
                        type = "localNet"
                        value = "192.168.0.0/16"
                    }
                ]
                per_client_bandwidth_limits = {
                    settings = "custom"
                    bandwidth_limits = {
                        limit_up = 1000000
                        limit_down = 1000000
                    }
                }
                dscp_tag_value = 0
                pcp_tag_value = 0
            }
        ]
        l3_firewall_rules = [
            {
                comment = "Allow TCP traffic to subnet with HTTP servers."
                policy = "allow"
                protocol = "tcp"
                dest_port =  "443"
                dest_cidr = "192.168.1.0"
            }
        ]
        l7_firewall_rules = [
            {
                policy = "deny"
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
                value = "10.11.12.00"
            },
            {
                policy = "deny"
                type = "ipRange"
                value = "10.11.12.00"
            }
        ]
    }

	 
	  

}
`

const testAccNetworksGroupPolicyResourceConfigUpdateNetworksGroupPolicy = `
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
	  bandwidth = {
		settings = "network default" 
		bandwidth_limits = {
			limit_down = 567
			limit_up = 879
		}
	}
	bonjour_forwarding = {
        settings = "custom"
        rules = [
            {
                description = "A simple bonjour rule"
                vlan_id = 1
                services = [ "All Services" ]
            }
        ]
    }
	content_filtering = {
        allowed_url_patterns = {}
		blocked_url_categories = {}
		blocked_url_patterns = {}
    }

	scheduling =  {
        enabled = true
        monday = {
            active = true
            from = "9:00"
            to = "17:00"
        },
        tuesday = {
            active = true
            from = "9:00"
            to = "17:00"
           
        }
        wednesday = {
			active = true
            from = "9:00"
            to = "17:00"
        }
        thursday = {
            active = true
            from = "9:00"
            to = "17:00"
        }
        friday = {
			active = true
            from = "9:00"
            to = "17:00"
        }
        saturday = {
            active = true
            from = "9:00"
            to = "17:00"
        }
        sunday = {
            active = true
            from = "9:00"
            to = "17:00"
        }
    }

	vlan_tagging = {
        settings = "custom"
        vlan_id =  "1"
    }

	firewall_and_traffic_shaping =  {
        settings = "custom"
        traffic_shaping_rules = [
            {
                definitions = [
                    {
                        type = "host"
                        value = "google.com"
                    },
                    {
                        type = "port"
                        value ="9090"
                    },
                    {
                        type = "ipRange"
                        value = "192.1.0.0"
                    },
                    {
                        type = "ipRange"
                        value = "192.1.0.0/16"
                    },
                    {
                        type = "ipRange"
                        value = "10.1.0.0/16:80"
                    },
                    {
                        type = "localNet"
                        value = "192.168.0.0/16"
                    }
                ]
                per_client_bandwidth_limits = {
                    settings = "custom"
                    bandwidth_limits = {
                        limit_up = 1000000
                        limit_down = 1000000
                    }
                }
                dscp_tag_value = 0
                pcp_tag_value = 0
            }
        ]
        l3_firewall_rules = [
            {
                comment = "Allow TCP traffic to subnet with HTTP servers."
                policy = "allow"
                protocol = "tcp"
                dest_port =  "443"
                dest_cidr = "192.168.1.0"
            }
        ]
        l7_firewall_rules = [
            {
                policy = "deny"
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
                value = "10.11.12.00"
            },
            {
                policy = "deny"
                type = "ipRange"
                value = "10.11.12.00"
            }
        ]
    }

	 
	  

}
`
