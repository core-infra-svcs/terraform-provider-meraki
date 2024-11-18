package administered

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// identitiesMeSchema defines the schema for the current user's identity.
var identitiesMeSchema = schema.Schema{
	MarkdownDescription: "Returns the identity of the current user.",

	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Unique identifier for this data source. Always set to 'identities_me'.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the user.",
			Computed:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "The email of the user.",
			Computed:            true,
		},
		"last_used_dashboard_at": schema.StringAttribute{
			MarkdownDescription: "The last time the user was active on the Dashboard UI.",
			Computed:            true,
		},
		"authentication": identitiesMeAuthenticationSchema,
	},
}

// identitiesMeAuthenticationSchema defines the schema for authentication details.
var identitiesMeAuthenticationSchema = schema.SingleNestedAttribute{
	MarkdownDescription: "Authentication details for the user.",
	Computed:            true,
	Attributes: map[string]schema.Attribute{
		"mode": schema.StringAttribute{
			MarkdownDescription: "The authentication mode.",
			Computed:            true,
		},
		"api":        identitiesMeAPISchema,
		"saml":       identitiesMeSAMLSchema,
		"two_factor": identitiesMeTwoFactorSchema,
	},
}

// identitiesMeAPISchema defines the schema for API details.
var identitiesMeAPISchema = schema.SingleNestedAttribute{
	MarkdownDescription: "API details for the user.",
	Computed:            true,
	Attributes: map[string]schema.Attribute{
		"key": schema.SingleNestedAttribute{
			MarkdownDescription: "API key details.",
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"created": schema.BoolAttribute{
					MarkdownDescription: "Whether the API key is created.",
					Computed:            true,
				},
			},
		},
	},
}

// identitiesMeSAMLSchema defines the schema for SAML authentication.
var identitiesMeSAMLSchema = schema.SingleNestedAttribute{
	MarkdownDescription: "Details about SAML authentication.",
	Computed:            true,
	Attributes: map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Whether SAML authentication is enabled.",
			Computed:            true,
		},
	},
}

// identitiesMeTwoFactorSchema defines the schema for two-factor authentication.
var identitiesMeTwoFactorSchema = schema.SingleNestedAttribute{
	MarkdownDescription: "Details about two-factor authentication.",
	Computed:            true,
	Attributes: map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Whether two-factor authentication is enabled.",
			Computed:            true,
		},
	},
}
