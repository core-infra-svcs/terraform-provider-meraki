---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "meraki_organizations_saml_idps Data Source - terraform-provider-meraki"
subcategory: ""
description: |-
  Ports the SAML IdPs in your organization.
---

# meraki_organizations_saml_idps (Data Source)

Ports the SAML IdPs in your organization.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `list` (Attributes List) (see [below for nested schema](#nestedatt--list))
- `organization_id` (String) Organization ID

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--list"></a>
### Nested Schema for `list`

Optional:

- `consumer_url` (String) URL that is consuming SAML Identity Provider (IdP)
- `idp_id` (String) ID associated with the SAML Identity Provider (IdP)
- `slo_logout_url` (String) Dashboard will redirect users to this URL when they sign out.
- `x_509_cert_sha1_fingerprint` (String) Fingerprint (SHA1) of the SAML certificate provided by your Identity Provider (IdP). This will be used for encryption / validation.
