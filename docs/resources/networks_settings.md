---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_networks_settings Resource - terraform-provider-meraki"
subcategory: ""
description: |-
  NetworksSettings resource for updating network settings.
---

# meraki_networks_settings (Resource)

NetworksSettings resource for updating network settings.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_id` (String) Network ID

### Optional

- `fips_enabled` (Boolean) Enables / disables FIPS on the network.
- `local_status_page` (Attributes) (see [below for nested schema](#nestedatt--local_status_page))
- `local_status_page_enabled` (Boolean) Enables / disables access to the device status page (http://[device's LAN IP]). Optional. Can only be set if localStatusPageEnabled is set to true
- `named_vlans_enabled` (Boolean) Enables / disables Named VLANs on the Network.
- `remote_status_page_enabled` (Boolean) Enables / disables access to the device status page (http://[device's LAN IP]). Optional. Can only be set if localStatusPageEnabled is set to true
- `secure_port_enabled` (Boolean) Enables / disables the secure port.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--local_status_page"></a>
### Nested Schema for `local_status_page`

Optional:

- `authentication` (Attributes) (see [below for nested schema](#nestedatt--local_status_page--authentication))

<a id="nestedatt--local_status_page--authentication"></a>
### Nested Schema for `local_status_page.authentication`

Optional:

- `enabled` (Boolean) Enables / disables the authentication on Local Status page(s).
- `password` (String, Sensitive) The password used for Local Status Page(s). Set this to null to clear the password.
- `username` (String) The username used for Local Status Page(s).
