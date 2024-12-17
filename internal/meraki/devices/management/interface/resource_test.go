package _interface_test

import (
	"fmt"
	"os"
	"testing"

	_interface "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/management/interface"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccResource validates the resource lifecycle and schema-model consistency
func TestAccResource(t *testing.T) {
	// Step 1: Validate schema-model consistency
	t.Run("Validate Schema-Model Consistency", func(t *testing.T) {
		testutils.ValidateResourceSchemaModelConsistency(
			t, _interface.GetResourceSchema.Attributes, _interface.ResourceModel{},
		)
	})

	// Step 2: Test resource lifecycle with multiple steps
	t.Run("Test Resource Lifecycle", func(t *testing.T) {
		// Test environment variables
		mxSerial := os.Getenv("TF_ACC_MERAKI_MX_SERIAL")
		msSerial := os.Getenv("TF_ACC_MERAKI_MS_SERIAL")
		mrSerial := os.Getenv("TF_ACC_MERAKI_MR_SERIAL")
		orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")

		if mxSerial == "" || msSerial == "" || mrSerial == "" || orgId == "" {
			t.Fatal("Environment variables TF_ACC_MERAKI_MX_SERIAL, TF_ACC_MERAKI_MS_SERIAL, TF_ACC_MERAKI_MR_SERIAL, and TF_ACC_MERAKI_ORGANIZATION_ID must be set for acceptance tests")
		}

		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testutils.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Step 1: Create Network
				{
					Config: utils.CreateNetworkOrgIdConfig(orgId, "test_acc_device_management_interface"),
					Check:  utils.NetworkOrgIdTestChecks("test_acc_device_management_interface"),
				},

				// Step 2: Claim devices to Network
				{
					Config: DevicesManagementInterfaceResourceConfigCreate(mxSerial, msSerial, mrSerial),
					Check:  DevicesManagementInterfaceResourceConfigCreateCheck(),
				},

				// Step 3: Update MX interface
				{
					Config: DevicesManagementInterfaceResourceConfigUpdate(mxSerial),
					Check:  DevicesManagementInterfaceResourceConfigUpdateCheck(),
				},

				// Step 4: Configure MS interface
				{
					Config: DevicesManagementInterfaceResourceConfigCreateMS(msSerial),
					Check:  DevicesManagementInterfaceResourceConfigCreateMSCheck(),
				},

				// Step 5: Configure MR interface
				{
					Config: DevicesManagementInterfaceResourceConfigCreateMR(mrSerial),
					Check:  DevicesManagementInterfaceResourceConfigCreateMRCheck(),
				},
			},
		})
	})
}

// DevicesManagementInterfaceResourceConfigCreate claims devices and creates the MX interface
func DevicesManagementInterfaceResourceConfigCreate(mxSerial, msSerial, mrSerial string) string {
	return fmt.Sprintf(`
%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = ["%s", "%s", "%s"]
}

resource "meraki_devices_management_interface" "mx" {
    depends_on = [meraki_networks_devices_claim.test]
    serial = "%s"
    wan1 = {
        wan_enabled = "disabled"
        vlan = 2
        using_static_ip = false
    }
}
`, utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"),
		mxSerial, msSerial, mrSerial, mxSerial,
	)
}

func DevicesManagementInterfaceResourceConfigCreateCheck() resource.TestCheckFunc {
	return utils.ResourceTestCheck("meraki_devices_management_interface.mx", map[string]string{
		"serial":               os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
		"wan1.wan_enabled":     "disabled",
		"wan1.vlan":            "2",
		"wan1.using_static_ip": "false",
	})
}

// DevicesManagementInterfaceResourceConfigUpdate updates the MX interface
func DevicesManagementInterfaceResourceConfigUpdate(serial string) string {
	return fmt.Sprintf(`
%s
resource "meraki_devices_management_interface" "mx" {
    serial = "%s"
    wan1 = {
        wan_enabled = "enabled"
        vlan = 20
        using_static_ip = false
    }
}
`, utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"), serial)
}

func DevicesManagementInterfaceResourceConfigUpdateCheck() resource.TestCheckFunc {
	return utils.ResourceTestCheck("meraki_devices_management_interface.mx", map[string]string{
		"serial":               os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
		"wan1.wan_enabled":     "enabled",
		"wan1.vlan":            "20",
		"wan1.using_static_ip": "false",
	})
}

// DevicesManagementInterfaceResourceConfigCreateMS configures the MS interface
func DevicesManagementInterfaceResourceConfigCreateMS(serial string) string {
	return fmt.Sprintf(`
%s
resource "meraki_devices_management_interface" "ms" {
    serial = "%s"
    wan1 = {
        vlan = 2
        using_static_ip = false
    }
}
`, utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"), serial)
}

func DevicesManagementInterfaceResourceConfigCreateMSCheck() resource.TestCheckFunc {
	return utils.ResourceTestCheck("meraki_devices_management_interface.ms", map[string]string{
		"serial":               os.Getenv("TF_ACC_MERAKI_MS_SERIAL"),
		"wan1.vlan":            "2",
		"wan1.using_static_ip": "false",
	})
}

// DevicesManagementInterfaceResourceConfigCreateMR configures the MR interface
func DevicesManagementInterfaceResourceConfigCreateMR(serial string) string {
	return fmt.Sprintf(`
%s
resource "meraki_devices_management_interface" "mr" {
    serial = "%s"
    wan1 = {
        wan_enabled = "not configured"
        vlan = 2
        using_static_ip = false
    }
}
`, utils.CreateNetworkOrgIdConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID"), "test_acc_device_management_interface"), serial)
}

func DevicesManagementInterfaceResourceConfigCreateMRCheck() resource.TestCheckFunc {
	return utils.ResourceTestCheck("meraki_devices_management_interface.mr", map[string]string{
		"serial":               os.Getenv("TF_ACC_MERAKI_MR_SERIAL"),
		"wan1.wan_enabled":     "not configured",
		"wan1.vlan":            "2",
		"wan1.using_static_ip": "false",
	})
}
