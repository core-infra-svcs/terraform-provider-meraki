terraform {
  required_providers {
    meraki = {
      source = "core-infra-svcs/meraki"
    }
  }
}

provider "meraki" {
  # example configuration here
  apikey = var.MERAKI_DASHBOARD_API_KEY
  path   = var.path
  host   = var.host
}


// Get List of Organizations
data "meraki_organizations" "list" {
}


// terraform output -json organizations | jq
output "organizations" {
  value = data.meraki_organizations.list
}

/*
// import any existing organizations by Id:
terraform import meraki_organization.testOrg1 "1234567890"

// Manually remove from state file
terraform state rm 'meraki_organization.testOrg1'

*/

// Create a new Meraki Organization. api_enabled and Id are required fields when modifying an organization.
resource "meraki_organization" "testOrg1" {
  id = "762234236932456967"
  name = "testOrg1"
  api_enabled = true
}


// terraform output -json testOrg1 | jq
output "terraform1" {
  value = meraki_organization.testOrg1
}


// Destroy an org after creation & remove from state
// terraform show
// terraform apply -destroy -target='meraki_organization.testOrg1'
