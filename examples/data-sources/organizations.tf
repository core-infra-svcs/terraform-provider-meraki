// list of Meraki organizations
data "meraki_organizations" "list" {
}

// terraform output -json organizations | jq
output "organizations" {
  value = data.meraki_organizations.list
}
