
// Create a new Meraki organization
resource "meraki_organization" "example_organization" {
  name            = "example_organization"
  api_enabled     = true
  licensing_model = "co-term"
}

// terraform output -json example_organization | jq
output "meraki_example" {
  value = meraki_organization.example_organization
}

// Destroys example organization
// terraform apply -destroy -target='meraki_organization.example_organization'

// Import a pre-existing organizations by Id:
// terraform import meraki_organization.example_organization "1234567890"

// Manually destroy a single organization
//terraform state rm 'meraki_organization.example_organization'