package networksSettings

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resourceModel describes the resource data model.
type resourceModel struct {
	Id                      jsontypes.String `tfsdk:"id"`
	NetworkId               jsontypes.String `tfsdk:"network_id" json:"network_id"`
	LocalStatusPageEnabled  jsontypes.Bool   `tfsdk:"local_status_page_enabled" json:"localStatusPageEnabled"`
	RemoteStatusPageEnabled jsontypes.Bool   `tfsdk:"remote_status_page_enabled" json:"remoteStatusPageEnabled"`
	LocalStatusPage         types.Object     `tfsdk:"local_status_page" json:"localStatusPage"`
	SecurePortEnabled       jsontypes.Bool   `tfsdk:"secure_port_enabled" json:"securePort"`
	FipsEnabled             jsontypes.Bool   `tfsdk:"fips_enabled" json:"fipsEnabled"`
	NamedVlansEnabled       jsontypes.Bool   `tfsdk:"named_vlans_enabled" json:"namedVlansEnabled"`
	//ClientPrivacyExpireDataOlderThan      jsontypes.Int64                              `tfsdk:"client_privacy_expire_data_older_than"`
	//ClientPrivacyExpireDataBefore         jsontypes.String                             `tfsdk:"client_privacy_expire_data_before"`
}
