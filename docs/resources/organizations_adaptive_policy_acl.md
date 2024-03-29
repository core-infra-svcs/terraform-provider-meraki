---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_organizations_adaptive_policy_acl Resource - terraform-provider-meraki"
subcategory: ""
description: |-
  Manage adaptive policy ACLs in a organization
---

# meraki_organizations_adaptive_policy_acl (Resource)

Manage adaptive policy ACLs in a organization



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `acl_id` (String) ACL ID
- `created_at` (String)
- `description` (String) Description of the adaptive policy ACL
- `ip_version` (String) IP version of adaptive policy ACL. One of: 'any', 'ipv4' or 'ipv6
- `name` (String) Name of the adaptive policy ACL
- `organization_id` (String) Organization ID
- `rules` (Attributes List) An ordered array of the adaptive policy ACL rules. An empty array will clear the rules. (see [below for nested schema](#nestedatt--rules))
- `updated_at` (String)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--rules"></a>
### Nested Schema for `rules`

Optional:

- `dst_port` (String)
- `policy` (String)
- `protocol` (String)
- `src_port` (String)
