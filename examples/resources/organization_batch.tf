// Create new Meraki organizations in bulk
resource "meraki_organization" "example_organization_batch" {
  count           = 10 # tested up to 100 organizations
  name            = "example_organization_batch_${count.index}"
  api_enabled     = true
  licensing_model = "co-term"
}

// Destroy ALL example organizations (can't touch any organization unknown to tfstate)
// terraform apply -destroy

// Manually remove from state file
//terraform state rm 'meraki_organization.{$NAME}'