terraform {
  required_providers {
    meraki = {
      source = "core-infra-svcs/meraki"
    }
  }
}


resource "meraki_devices" "test" {
  serial = "ABCD-1234"
}

resource "meraki_devices" "test" {
  serial = "%s"
  name = "My AP"
  tags = ["sfo", "ca"]
}