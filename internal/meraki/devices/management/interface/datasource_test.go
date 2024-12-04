package _interface_test

import (
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDevicesManagementInterfaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			provider.TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{

			// Create and Read Network
			{
				Config: utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"),
				Check:  utils.NetworkOrgIdTestChecks("test_acc_device_management_interface"),
			},

			// Claim Appliance To Network
			{
				Config: DevicesManagementInterfaceDatasource(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
				Check:  DevicesManagementInterfaceDatasourceCheck(),
			},
		},
	})
}

func DevicesManagementInterfaceDatasource(serial string) string {
	return fmt.Sprintf(`
	%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [resource.meraki_network.test]
    network_id = resource.meraki_network.test.network_id
    serials = [
      "%s"
  ]
}

resource "meraki_devices_management_interface" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test]
	serial = "%s"
	wan1 = {
		wan_enabled = "enabled"
		vlan = 2
		using_static_ip = false
	}
}

data "meraki_devices_management_interface" "test" {
    depends_on = [resource.meraki_network.test, resource.meraki_networks_devices_claim.test, resource.meraki_devices_management_interface.test]
	serial = "%s"
}
`,
		utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"),
		serial, serial, serial,
	)
}

// DevicesManagementInterfaceDatasourceCheck returns the test check functions for DevicesManagementInterfaceDatasource
func DevicesManagementInterfaceDatasourceCheck() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"serial":               os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
		"wan1.wan_enabled":     "enabled",
		"wan1.vlan":            "2",
		"wan1.using_static_ip": "false",
	}
	return utils.ResourceTestCheck("meraki_devices_management_interface.test", expectedAttrs)
}
