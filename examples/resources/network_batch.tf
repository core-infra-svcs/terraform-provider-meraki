
// Create a new Meraki organization
resource "meraki_organization" "example_organization" {
  name        = "example_organization"
  api_enabled = true
  licensing_model = "co-term"
}

// Create new Meraki organization Networks in bulk
resource "meraki_network" "example_networks_batch" {
  count = 100 # tested up to 100 networks
  depends_on = [meraki_organization.example_organization]
  organization_id = resource.meraki_organization.example_organization.organization_id
  product_types = ["appliance"]
  tags = ["example_network_batch"]
  name        = "example_network_batch_${count.index}"
  timezone = "America/Los_Angeles"
  notes = "Additional description of the network"
}