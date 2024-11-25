package appliance_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworksApplianceStaticRoutesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_static_routes"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_networks_appliance_static_routes"),
			},

			// Create and Read Networks Appliance Static Routes.
			{
				Config: NetworksApplianceStaticRoutesResourceConfigCreate(),
				Check:  NetworksApplianceStaticRoutesResourceConfigCreateChecks(),
			},

			// Update testing
			{
				Config: NetworksApplianceStaticRoutesResourceConfigUpdate(),
				Check:  NetworksApplianceStaticRoutesResourceConfigUpdateChecks(),
			},
		},
	})
}

func NetworksApplianceStaticRoutesResourceConfigCreate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_static_routes" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id  
    name = "My route"
    subnet = "192.168.129.0/24"
    gateway_ip = "192.168.128.1"
	reserved_ip_ranges = []
	
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_static_routes"),
	)
}

// NetworksApplianceStaticRoutesResourceConfigCreateChecks returns the test check functions for NetworksApplianceStaticRoutesResourceConfigCreate
func NetworksApplianceStaticRoutesResourceConfigCreateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":       "My route",
		"subnet":     "192.168.129.0/24",
		"gateway_ip": "192.168.128.1",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_static_routes.test", expectedAttrs)
}

func NetworksApplianceStaticRoutesResourceConfigUpdate() string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_appliance_static_routes" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id    
	name = "My route"
    subnet = "192.168.129.0/24"
	fixed_ip_assignments_mac_address = "22:33:44:55:66:77"
	fixed_ip_assignments_mac_ip_address = "192.168.128.1"
	fixed_ip_assignments_mac_name = "Some client name"   
	reserved_ip_ranges = [
        {
            start = "192.168.128.1"
            end = "192.168.128.2"
            comment = "A reserved IP range"
        }
    ]
}
	
	`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_networks_appliance_static_routes"),
	)
}

// NetworksApplianceStaticRoutesResourceConfigUpdateChecks returns the test check functions for NetworksApplianceStaticRoutesResourceConfigUpdate
func NetworksApplianceStaticRoutesResourceConfigUpdateChecks() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"name":                                "My route",
		"subnet":                              "192.168.129.0/24",
		"gateway_ip":                          "192.168.128.1",
		"enable":                              "true",
		"fixed_ip_assignments_mac_address":    "22:33:44:55:66:77",
		"fixed_ip_assignments_mac_ip_address": "192.168.128.1",
		"fixed_ip_assignments_mac_name":       "Some client name",
		"reserved_ip_ranges.0.comment":        "A reserved IP range",
		"reserved_ip_ranges.0.start":          "192.168.128.1",
		"reserved_ip_ranges.0.end":            "192.168.128.2",
	}
	return utils.ResourceTestCheck("meraki_networks_appliance_static_routes.test", expectedAttrs)
}
