---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_networks_cellular_gateway_uplink Resource - terraform-provider-meraki"
subcategory: ""
description: |-
  Networks Cellular Gateway Uplink Updates the uplink settings for your MG network.
---

# meraki_networks_cellular_gateway_uplink (Resource)

Networks Cellular Gateway Uplink Updates the uplink settings for your MG network.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_id` (String) Network Id

### Optional

- `bandwidth_limits` (Attributes) The bandwidth settings for your MG network (see [below for nested schema](#nestedatt--bandwidth_limits))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--bandwidth_limits"></a>
### Nested Schema for `bandwidth_limits`

Optional:

- `limit_down` (Number) The maximum download limit (integer, in Kbps). null indicates no limit.
- `limit_up` (Number) The maximum upload limit (integer, in Kbps). null indicates no limit.

