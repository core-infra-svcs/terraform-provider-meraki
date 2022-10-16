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