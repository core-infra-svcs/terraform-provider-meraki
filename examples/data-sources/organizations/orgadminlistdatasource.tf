// Get admin list
data "meraki_organizations_admins" "list" {
  id = "784752235069308980"
}

// terraform output -json organizations | jq
output "orgadminlist" {
  value = data.meraki_organizations_admins.list
}
