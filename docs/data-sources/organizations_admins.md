---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_organizations_admins Data Source - terraform-provider-meraki"
subcategory: ""
description: |-
  OrganizationsAdmins data source - get all list of  admins in an organization
---

# meraki_organizations_admins (Data Source)

OrganizationsAdmins data source - get all list of  admins in an organization



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) org id

### Optional

- `list` (Attributes List) List of organization admins (see [below for nested schema](#nestedatt--list))

<a id="nestedatt--list"></a>
### Nested Schema for `list`

Optional:

- `account_status` (String) Account Status
- `authentication_method` (String) Authentication method
- `email` (String) Email of the dashboard administrator
- `has_api_key` (Boolean) Api key exists or not
- `id` (String) id of the organization
- `last_active` (String) Last Time Active
- `name` (String) name of the dashboard administrator
- `orgaccess` (String) Organization Access
- `two_factor_auth_enabled` (Boolean) Two Factor Auth Enabled or Not

Read-Only:

- `networks` (Attributes List) list of networks that the dashboard administrator has privileges on. (see [below for nested schema](#nestedatt--list--networks))
- `tags` (Attributes List) list of tags that the dashboard administrator has privileges on. (see [below for nested schema](#nestedatt--list--tags))

<a id="nestedatt--list--networks"></a>
### Nested Schema for `list.networks`

Optional:

- `access` (String) network access
- `id` (String) network id


<a id="nestedatt--list--tags"></a>
### Nested Schema for `list.tags`

Optional:

- `access` (String) access
- `tag` (String) tag

