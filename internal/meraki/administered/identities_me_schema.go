package administered

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var identitiesMeSchema = schema.Schema{
	MarkdownDescription: "Returns the identity of the current user",

	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			MarkdownDescription: "Username",
			Optional:            true,
			Computed:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "User email",
			Optional:            true,
			Computed:            true,
		},
		"last_used_dashboard_at": schema.StringAttribute{
			MarkdownDescription: "Last seen active on Dashboard UI",
			Optional:            true,
			Computed:            true,
			CustomType:          jsontypes.StringType,
		},
		"authentication": schema.SingleNestedAttribute{
			MarkdownDescription: "Authentication details for the user",
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"mode": schema.StringAttribute{
					MarkdownDescription: "Authentication mode",
					Computed:            true,
				},
				"api_key_created": schema.BoolAttribute{
					MarkdownDescription: "If API key is created for this user",
					Computed:            true,
				},
				"api": schema.SingleNestedAttribute{
					MarkdownDescription: "API details",
					Computed:            true,
					Attributes: map[string]schema.Attribute{
						"key": schema.SingleNestedAttribute{
							MarkdownDescription: "API key details",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"created": schema.BoolAttribute{
									MarkdownDescription: "Whether the API key is created",
									Computed:            true,
								},
							},
						},
					},
				},
				"saml_enabled": schema.BoolAttribute{
					MarkdownDescription: "If SAML authentication is enabled",
					Computed:            true,
				},
				"two_factor": schema.SingleNestedAttribute{
					MarkdownDescription: "Two-factor authentication details",
					Computed:            true,
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether two-factor authentication is enabled",
							Computed:            true,
						},
					},
				},
			},
		},
	},
}
