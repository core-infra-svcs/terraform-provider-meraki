---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_networks_wireless_ssids_firewall_l3_firewall_rules Resource - terraform-provider-meraki"
subcategory: ""
description: |-
  NetworksWirelessSsidsFirewallL3FirewallRules for Updating Networks Wireless Ssids Firewall L3FirewallRules
---

# meraki_networks_wireless_ssids_firewall_l3_firewall_rules (Resource)

NetworksWirelessSsidsFirewallL3FirewallRules for Updating Networks Wireless Ssids Firewall L3FirewallRules



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_id` (String) Network ID
- `number` (String) SsIds Number

### Optional

- `allow_lan_access` (Boolean) Allow wireless client access to local LAN (boolean value - true allows access and false denies access) (optional)
- `rules` (Attributes Set) An ordered array of the firewall rules for this SSID (not including the local LAN access rule or the default rule) (see [below for nested schema](#nestedatt--rules))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Required:

- `dest_cidr` (String) Comma-separated list of destination IP address(es) (in IP or CIDR notation), fully-qualified domain names (FQDN) or 'Any'
- `policy` (String) 'allow' or 'deny' traffic specified by this rule
- `protocol` (String) The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6' or 'Any')

Optional:

- `comment` (String) Description of the rule (optional)
- `dest_port` (String) Comma-separated list of destination port(s) (integer in the range 1-65535), or 'Any'

