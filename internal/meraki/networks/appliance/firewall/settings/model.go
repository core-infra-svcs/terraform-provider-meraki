package settings

import "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"

// resourceModel describes the resource data model.
type resourceModel struct {
	Id                 jsontypes.String                `tfsdk:"id"`
	NetworkId          jsontypes.String                `tfsdk:"network_id" json:"network_id"`
	SpoofingProtection resourceModelSpoofingProtection `tfsdk:"spoofing_protection" json:"spoofingProtection"`
}

type resourceModelSpoofingProtection struct {
	IpSourceGuard resourceModelIpSourceGuard `tfsdk:"ip_source_guard" json:"ipSourceGuard"`
}

type resourceModelIpSourceGuard struct {
	Mode jsontypes.String `tfsdk:"mode" json:"mode"`
}
