terraform {
  required_providers {
    meraki = {
      source = "core-infra-svcs/meraki"
    }
  }
}

provider "meraki" {
  # example configuration here
  api_key  = var.MERAKI_DASHBOARD_API_KEY
  base_url = var.MERAKI_DASHBOARD_API_URL

}