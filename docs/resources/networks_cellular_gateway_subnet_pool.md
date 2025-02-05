---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_networks_cellular_gateway_subnet_pool Resource - terraform-provider-meraki"
subcategory: ""
description: |-
  Manage the subnet pool and mask configuration for MGs in the network.
---

# meraki_networks_cellular_gateway_subnet_pool (Resource)

Manage the subnet pool and mask configuration for MGs in the network.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) Network ID

### Optional

- `cidr` (String) CIDR of the pool of subnets. Each MG in this network will automatically pick a subnet from this pool.
- `deployment_mode` (String)
- `mask` (Number) Mask used for the subnet of all MGs in this network.
- `subnets` (Attributes Set) (see [below for nested schema](#nestedatt--subnets))

<a id="nestedatt--subnets"></a>
### Nested Schema for `subnets`

Optional:

- `appliance_ip` (String)
- `name` (String)
- `serial` (String)
- `subnet` (String)
