
// Create a new Meraki organization
resource "meraki_organization" "example_organization" {
  name        = "example_organization"
  api_enabled = true
  licensing_model = "co-term"
}

// Create new Meraki organization Network
resource "meraki_network" "example_networks" {
  depends_on = [meraki_organization.example_organization]
  organization_id = resource.meraki_organization.example_organization.organization_id
  product_types = ["appliance"]
  tags = ["example_network_batch"]
  name        = "example_network"
  timezone = "America/Los_Angeles"
  notes = "Additional description of the network"
}