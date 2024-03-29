---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_networks_switch_settings Resource - terraform-provider-meraki"
subcategory: ""
description: |-
  NetworksSwitchSettings resource for updating network switch settings.
---

# meraki_networks_switch_settings (Resource)

NetworksSwitchSettings resource for updating network switch settings.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `network_id` (String) Network Id
- `power_exceptions` (Attributes List) Exceptions on a per switch basis to &quot;useCombinedPower&quot; (see [below for nested schema](#nestedatt--power_exceptions))

### Optional

- `use_combined_power` (Boolean) The use combined Power as the default behavior of secondary power supplies on supported devices.
- `vlan` (Number) Management VLAN

### Read-Only

- `id` (String) Example identifier

<a id="nestedatt--power_exceptions"></a>
### Nested Schema for `power_exceptions`

Optional:

- `power_type` (String) Per switch exception (combined, redundant, useNetworkSetting)
- `serial` (String) Serial number of the switch
