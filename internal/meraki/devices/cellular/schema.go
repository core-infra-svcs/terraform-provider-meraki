package cellular

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Schema returns the resource schema definition
func Schema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages the SIM and APN configurations for a cellular device.",
		Attributes: map[string]schema.Attribute{
			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial number of the device.",
				Optional:            true,
				Computed:            true,
			},
			"sim_failover": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "If true, failover to the secondary SIM is enabled.",
						Optional:            true,
						Computed:            true,
					},
					"timeout": schema.Int64Attribute{
						MarkdownDescription: "Timeout value (in seconds) for SIM failover. Defaults to 0.",
						Optional:            true,
						Computed:            true,
						Default:             utils.NewInt64Default(0),
					},
				},
			},
			"sims": SimsSchema(), // Call the helper function for sims schema
			"sim_ordering": schema.SetAttribute{
				MarkdownDescription: "Ordered list of SIM slots, prioritized for failover.",
				Computed:            true,
				ElementType:         jsontypes.StringType,
			},
		},
	}
}

// SimsSchema defines the schema for the sims attribute.
func SimsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "Configuration for device SIMs. Unspecified SIMs remain unchanged.",
		Required:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"slot": schema.StringAttribute{
					MarkdownDescription: "SIM slot being configured. Must be 'sim1' for single-SIM devices.",
					Optional:            true,
					Computed:            true,
				},
				"is_primary": schema.BoolAttribute{
					MarkdownDescription: "Indicates if this SIM is the primary SIM for boot. Must be true for single-SIM devices.",
					Optional:            true,
					Computed:            true,
				},
				"apns": ApnsSchema(), // Call the helper function for APNs schema
			},
		},
	}
}

// ApnsSchema defines the schema for the apns attribute inside sims.
func ApnsSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "APN configurations for the SIM. If empty, the default APN is used.",
		Required:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: "Name of the APN.",
					Required:            true,
				},
				"allowed_ip_types": schema.SetAttribute{
					MarkdownDescription: "Allowed IP versions for the APN (e.g., 'ipv4', 'ipv6').",
					Required:            true,
					ElementType:         jsontypes.StringType,
				},
				"authentication": AuthenticationSchema(), // Call the helper function for authentication schema
			},
		},
	}
}

// AuthenticationSchema defines the schema for the authentication attribute inside apns.
func AuthenticationSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Authentication details for the APN.",
		Attributes: map[string]schema.Attribute{
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for APN authentication.",
				Optional:            true,
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for APN authentication.",
				Optional:            true,
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of authentication. Valid values: 'chap', 'none', 'pap'.",
				Optional:            true,
				Computed:            true,
				Validators:          []validator.String{stringvalidator.OneOf("chap", "none", "pap")},
			},
		},
	}
}
