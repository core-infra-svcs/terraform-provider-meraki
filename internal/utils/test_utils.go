package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"time"
)

// GenerateTimestamp returns the current timestamp in the format YYYYMMDDHHMMSS
func GenerateTimestamp() string {
	return time.Now().Format("20060102150405")
}

// CreateResourceConfig returns a configuration string to create any resource
func CreateResourceConfig(resourceType string, attributes map[string]string) string {
	attributeStr := ""
	for key, value := range attributes {
		attributeStr += fmt.Sprintf("%s = \"%s\"\n", key, value)
	}
	return fmt.Sprintf(`
	resource "%s" "test" {
		%s
	}
	`, resourceType, attributeStr)
}

// ResourceTestCheck is a generic function to test any resource's attributes
func ResourceTestCheck(resourceName string, expectedAttrs map[string]string) resource.TestCheckFunc {
	checks := make([]resource.TestCheckFunc, 0)
	for attr, value := range expectedAttrs {
		checks = append(checks, resource.TestCheckResourceAttr(resourceName, attr, value))
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// CreateOrganizationConfig returns a configuration string to create an organization resource
func CreateOrganizationConfig(organizationName string) string {
	return CreateResourceConfig("meraki_organization", map[string]string{
		"name":        organizationName,
		"api_enabled": "true",
	})
}

// OrganizationTestChecks returns the aggregated test check functions for an Organization resource
func OrganizationTestChecks(organizationName string) resource.TestCheckFunc {
	return ResourceTestCheck("meraki_organization.test", map[string]string{
		"name":        organizationName,
		"api_enabled": "true",
	})
}

// CreateNetworkConfig returns a configuration string to create a network resource with the necessary test checks
func CreateNetworkConfig(orgName, networkName string) string {
	return fmt.Sprintf(`
	resource "meraki_organization" "test" {
		name = "%s"
		api_enabled = true
	}

	resource "meraki_network" "test" {
		organization_id = resource.meraki_organization.test.organization_id
		product_types = ["appliance", "switch", "wireless"]
		name = "%s"
		timezone = "America/Los_Angeles"
		tags = ["tag1"]
		notes = "Additional description of the network"
	}
	`, orgName, networkName)
}

// NetworkTestChecks returns the aggregated test check functions for a network resource
func NetworkTestChecks(networkName string) resource.TestCheckFunc {
	return ResourceTestCheck("meraki_network.test", map[string]string{
		"name":            networkName,
		"timezone":        "America/Los_Angeles",
		"tags.#":          "1",
		"tags.0":          "tag1",
		"product_types.#": "3",
		"product_types.0": "appliance",
		"product_types.1": "switch",
		"product_types.2": "wireless",
		"notes":           "Additional description of the network",
	})
}

// ClaimDeviceConfig returns a configuration string to claim a device by serial
func ClaimDeviceConfig(serial string) string {
	return fmt.Sprintf(`
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_claim" "test_serial" {
		organization_id = resource.meraki_organization.test.organization_id
		orders = []
		licences = []
		serials = ["%s"]
	}
	`, serial)
}

// ClaimDeviceTestChecks returns the test check functions for claiming a device by serial
func ClaimDeviceTestChecks(serial string) resource.TestCheckFunc {
	return ResourceTestCheck("meraki_organizations_claim.test_serial", map[string]string{
		"serials.#": "1",
		"serials.0": serial,
	})
}

// ClaimOrderConfig returns a configuration string to claim an order
func ClaimOrderConfig(order string) string {
	return fmt.Sprintf(`
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_claim" "test_order" {
		organization_id = resource.meraki_organization.test.organization_id
		orders = ["%s"]
		serials = []
		licences = []
	}
	`, order)
}

// ClaimOrderTestChecks returns the test check functions for claiming an order
func ClaimOrderTestChecks(order string) resource.TestCheckFunc {
	return ResourceTestCheck("meraki_organizations_claim.test_order", map[string]string{
		"orders.#": "1",
		"orders.0": order,
	})
}

// ClaimLicenseConfig returns a configuration string to claim a license
func ClaimLicenseConfig(license string) string {
	return fmt.Sprintf(`
	resource "meraki_organization" "test" {}

	resource "meraki_organizations_claim" "test_license" {
		organization_id = resource.meraki_organization.test.organization_id
		orders = []
		serials = []
		licences = [
			{
				key = "%s"
				mode = "addDevices"
			}
		]
	}
	`, license)
}

// ClaimLicenseTestChecks returns the test check functions for claiming a license
func ClaimLicenseTestChecks(license string) resource.TestCheckFunc {
	return ResourceTestCheck("meraki_organizations_claim.test_license", map[string]string{
		"licenses.#":      "1",
		"licenses.0.key":  license,
		"licenses.0.mode": "addDevices",
	})
}
