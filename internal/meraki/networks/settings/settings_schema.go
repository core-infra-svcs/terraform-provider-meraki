package networksSettings

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var settingsSchema = schema.Schema{
	MarkdownDescription: "NetworksSettings resource for updating network settings.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:   true,
			CustomType: jsontypes.StringType,
		},
		"network_id": schema.StringAttribute{
			MarkdownDescription: "Network ID",
			Required:            true,
			CustomType:          jsontypes.StringType,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 31),
			},
		},
		"local_status_page_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enables / disables access to the device status page (http://[device's LAN IP]). Optional. Can only be set if localStatusPageEnabled is set to true",
			Optional:            true,
			Computed:            true,
			CustomType:          jsontypes.BoolType,
		},
		"remote_status_page_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enables / disables access to the device status page (http://[device's LAN IP]). Optional. Can only be set if localStatusPageEnabled is set to true",
			Optional:            true,
			Computed:            true,
			CustomType:          jsontypes.BoolType,
		},
		"local_status_page": settingsLocalStatusPageSchema,
		"secure_port_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enables / disables the secure port.",
			Optional:            true,
			Computed:            true,
			CustomType:          jsontypes.BoolType,
		},
		"fips_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enables / disables FIPS on the network.",
			Optional:            true,
			Computed:            true,
			CustomType:          jsontypes.BoolType,
		},
		"named_vlans_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enables / disables Named VLANs on the Network.",
			Optional:            true,
			Computed:            true,
			CustomType:          jsontypes.BoolType,
		},
		/*
			"client_privacy_expire_data_older_than": schema.Int64Attribute{
					MarkdownDescription: "The number of days, weeks, or months in Epoch time to expire the data before",
					Optional:            true,
					Computed:            true,
					CustomType:          jsontypes.Int64Type,
				},
				"client_privacy_expire_data_before": schema.StringAttribute{
					MarkdownDescription: "The date to expire the data before",
					Optional:            true,
					Computed:            true,
					CustomType:          jsontypes.StringType,
				},
		*/
	},
}

var settingsLocalStatusPageSchema = schema.SingleNestedAttribute{
	Optional: true,
	Computed: true,
	Attributes: map[string]schema.Attribute{
		"authentication": schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					MarkdownDescription: "Enables / disables the authentication on Local Status page(s).",
					Optional:            true,
					Computed:            true,
					CustomType:          jsontypes.BoolType,
				},
				"username": schema.StringAttribute{
					MarkdownDescription: "The username used for Local Status Page(s).",
					Optional:            true,
					Computed:            true,
					CustomType:          jsontypes.StringType,
				},
				"password": schema.StringAttribute{
					MarkdownDescription: "The password used for Local Status Page(s). Set this to null to clear the password.",
					Optional:            true,
					CustomType:          jsontypes.StringType,
					Sensitive:           true,
				},
			}},
	},
}

// SettingsResourceModel describes the resource data model.
type SettingsResourceModel struct {
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
