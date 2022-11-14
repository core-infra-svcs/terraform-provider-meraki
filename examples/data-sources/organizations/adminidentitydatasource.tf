
// Get identity
data "meraki_administered_identities_me" "test" {

}

// terraform output -json organizations | jq
output "admindetails" {
  value = data.meraki_administered_identities_me.test
}
