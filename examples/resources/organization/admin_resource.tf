terraform {
  required_providers {
    meraki = {
      source = "core-infra-svcs/meraki"
    }
  }
}

resource "meraki_organizations_admin" "testAdmin" {
  id        = "784752235069308981"
  name      = "testAdmin123456"
  email     = "kirankumar600270092456016661289@gmail.com"
  orgaccess = "read-only"
  tags = [
    {
      tag    = "east"
      access = "monitor-only"
    }
  ]



}


// terraform output -json testOrg1 | jq
output "terraform1" {
  value = meraki_organizations_admin.testAdmin
}


// Destroy an org after creation & remove from state
// terraform show
// terraform apply -destroy -target='meraki_organization.testOrg1'

