package _interface_test

import (
	"fmt"
	_interface "github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/devices/management/interface"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/testutils"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResource(t *testing.T) {
	t.Run("Validate Schema-Model Consistency", func(t *testing.T) {
		validateResourceSchemaModelConsistency(t)
	})

	t.Run("Test Resource Lifecycle", func(t *testing.T) {
		testResourceLifecycle(t)
	})
}

// Validate schema-model consistency for the data source
func validateResourceSchemaModelConsistency(t *testing.T) {
	testutils.ValidateResourceSchemaModelConsistency(
		t, _interface.GetResourceSchema.Attributes, _interface.ResourceModel{},
	)
}

// Test the full resource lifecycle for the data source
func testResourceLifecycle(t *testing.T) {
	mxSerial, msSerial, mrSerial, orgId := getResourceTestEnvironmentVariables(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createNetworkConfigResource(orgId),
				Check:  networkCheckResource("test_acc_device_management_interface"),
			},
			{
				Config: claimDevicesConfigResource(mxSerial, msSerial, mrSerial),
				Check:  claimDevicesCheckResource(),
			},
			{
				Config: ResourceConfig(mxSerial),
				Check:  ResourceCheck(),
			},
		},
	})
}

// Retrieve required environment variables for tests
func getResourceTestEnvironmentVariables(t *testing.T) (string, string, string, string) {
	mxSerial := os.Getenv("TF_ACC_MERAKI_MX_SERIAL")
	msSerial := os.Getenv("TF_ACC_MERAKI_MS_SERIAL")
	mrSerial := os.Getenv("TF_ACC_MERAKI_MR_SERIAL")
	orgId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")

	if mxSerial == "" || msSerial == "" || mrSerial == "" || orgId == "" {
		t.Fatal("Environment variables TF_ACC_MERAKI_MX_SERIAL, TF_ACC_MERAKI_MS_SERIAL, TF_ACC_MERAKI_MR_SERIAL, and TF_ACC_MERAKI_ORGANIZATION_ID must be set for acceptance tests")
	}

	return mxSerial, msSerial, mrSerial, orgId
}

// Create network configuration
func createNetworkConfigResource(orgId string) string {
	return utils.CreateNetworkOrgIdConfig(orgId, "test_acc_device_management_interface")
}

// Validate network creation
func networkCheckResource(networkName string) resource.TestCheckFunc {
	return utils.NetworkOrgIdTestChecks(networkName)
}

// Claim devices configuration
func claimDevicesConfigResource(mxSerial, msSerial, mrSerial string) string {
	return fmt.Sprintf(`
%s
resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials = ["%s", "%s", "%s"]
}
`, createNetworkConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")), mxSerial, msSerial, mrSerial)
}

// Validate device claiming
func claimDevicesCheckResource() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources["meraki_networks_devices_claim.test"]
		if !ok {
			return fmt.Errorf("resource not found in state")
		}

		serialCountStr := rs.Primary.Attributes["serials.#"]
		if serialCountStr == "" {
			return fmt.Errorf("serial count not found in state")
		}

		serialCount, err := strconv.Atoi(serialCountStr)
		if err != nil {
			return fmt.Errorf("invalid serial count: %s", serialCountStr)
		}

		actualSerials := make(map[string]bool)
		for i := 0; i < serialCount; i++ {
			serialKey := fmt.Sprintf("serials.%d", i)
			actualSerial := rs.Primary.Attributes[serialKey]
			if actualSerial == "" {
				return fmt.Errorf("serial %d not found in state", i)
			}
			actualSerials[actualSerial] = true
		}

		expectedSerials := []string{
			os.Getenv("TF_ACC_MERAKI_MX_SERIAL"),
			os.Getenv("TF_ACC_MERAKI_MS_SERIAL"),
			os.Getenv("TF_ACC_MERAKI_MR_SERIAL"),
		}

		for _, expected := range expectedSerials {
			if !actualSerials[expected] {
				return fmt.Errorf("expected serial %s not found in actual serials: %+v", expected, actualSerials)
			}
		}

		return nil
	}
}

// Data source configuration
func ResourceConfig(serial string) string {
	return fmt.Sprintf(`
%s

resource "meraki_networks_devices_claim" "test" {
    depends_on = [meraki_network.test]
    network_id = meraki_network.test.network_id
    serials    = ["%s"]
}

resource "meraki_devices_management_interface" "test" {
    depends_on = [meraki_network.test, meraki_networks_devices_claim.test]
    serial = "%s"
    wan1 = {
        wan_enabled = "enabled"
        vlan = 2
        using_static_ip = false
    }
}
`, createNetworkConfig(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")), serial, serial)
}

// Validate resource attributes
func ResourceCheck() resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "serial", os.Getenv("TF_ACC_MERAKI_MX_SERIAL")),
		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.wan_enabled", "enabled"),
		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.vlan", "2"),
		resource.TestCheckResourceAttr("meraki_devices_management_interface.test", "wan1.using_static_ip", "false"),
	)
}
