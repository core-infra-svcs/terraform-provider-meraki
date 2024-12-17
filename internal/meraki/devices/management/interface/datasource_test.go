package _interface_test

import (
	"fmt"
	_interface "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/management/interface"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDevicesManagementInterfaceDataSource validates schema, resource creation, and retrieval
func TestAccDevicesManagementInterfaceDataSource(t *testing.T) {

	// Test 1: Validate schema-model consistency for the DataSource
	t.Run("Validate Schema-Model Consistency", func(t *testing.T) {
		testutils.ValidateDataSourceSchemaModelConsistency(
			t, _interface.GetDatasourceSchema.Attributes, _interface.DataSourceModel{},
		)
	})

	// Test 2: Validate creation and retrieval of management interface settings
	t.Run("Test Resource Lifecycle", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck: func() {
				testutils.TestAccPreCheck(t)
			},
			ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{

				// Step 1: Create Network and Claim Device
				{
					Config: testCreateNetworkAndDeviceConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					Check: resource.ComposeTestCheckFunc(
						utils.NetworkOrgIdTestChecks("test_acc_device_management_interface_d"),
					),
				},

				// Step 2: Retrieve Device Management Interface (Data Source)
				{
					Config: testDevicesManagementInterfaceDatasourceConfig(os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
					Check:  testDevicesManagementInterfaceDatasourceCheck(),
				},
			},
		})
	})
}

// testCreateNetworkAndDeviceConfig generates Terraform config for creating a network and claiming a device
func testCreateNetworkAndDeviceConfig(orgId, serial string) string {
	return fmt.Sprintf(`
resource "meraki_network" "test" {
    organization_id = "%s"
    name            = "test_acc_device_management_interface_d"
    type            = "appliance"
}

resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = [
      "%s"
    ]
}

resource "meraki_devices_management_interface" "test" {
    depends_on = [meraki_network.test, meraki_networks_devices_claim.test]
	serial = "%s"
	wan1 = {
		wan_enabled = "enabled"
		vlan = 2
		using_static_ip = false
		static_dns = ["1.1.1.1", "8.8.8.8"]
	}
}
`, orgId, serial, serial)
}

// testDevicesManagementInterfaceDatasourceConfig generates Terraform config for the data source
func testDevicesManagementInterfaceDatasourceConfig(serial string) string {
	return fmt.Sprintf(`
data "meraki_devices_management_interface" "test" {
    serial = "%s"
}
`, serial)
}

// testDevicesManagementInterfaceDatasourceCheck validates the expected state of the data source
func testDevicesManagementInterfaceDatasourceCheck() resource.TestCheckFunc {
	expectedAttrs := map[string]string{
		"serial":                              os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
		"wan1.wan_enabled":                    "enabled",
		"wan1.vlan":                           "2",
		"wan1.using_static_ip":                "false",
		"wan1.static_dns.#":                   "2", // Validate list length
		"wan1.static_dns.0":                   "1.1.1.1",
		"wan1.static_dns.1":                   "8.8.8.8",
		"ddns_hostnames.active_ddns_hostname": "",
	}
	return utils.ResourceTestCheck("data.meraki_devices_management_interface.test", expectedAttrs)
}
