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
	  bandwidth_settings = "network default" 
	  bandwidth_limit_down = 567
	  bandwidth_limit_up = 879
      bonjour_forwarding_rules = [
            {
                description = "A simple bonjour rule"
                vlan_id = 1
                services = [ "All Services" ]
            }
        ]
      bonjour_forwarding_settings = "network default"  
      vlan_tagging_settings = "custom" 
      vlan_tagging_vlan_id =  "1"
      scheduling_enabled = true 
      scheduling_friday_active = true  
      scheduling_friday_from = "09:00"
      scheduling_friday_to = "17:00"
     content_filtering_allowed_url_patterns_settings = "network default"
	 url_patterns = []
	 content_filtering_blocked_url_categories_settings = "network default"
	 categories = []
	 content_filtering_blocked_url_patterns_settings = "network default"
	 blocked_url_patterns = []
      
       
        
    
       
    
    
		
	

	 
	  

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
    bandwidth_settings = "network default" 
    bandwidth_limit_down = 567
    bandwidth_limit_up = 879
    bonjour_forwarding_rules = [
          {
              description = "A simple bonjour rule"
              vlan_id = 1
              services = [ "All Services" ]
          }
      ]
    bonjour_forwarding_settings = "network default"  
    vlan_tagging_settings = "custom" 
    vlan_tagging_vlan_id =  "1"
    scheduling_enabled = true 
    scheduling_friday_active = true  
    scheduling_friday_from = "09:00"
    scheduling_friday_to = "17:00"
   content_filtering_allowed_url_patterns_settings = "network default"
   url_patterns = []
   content_filtering_blocked_url_categories_settings = "network default"
   categories = []
   content_filtering_blocked_url_patterns_settings = "network default"
   blocked_url_patterns = []

	 
	  

}
`
