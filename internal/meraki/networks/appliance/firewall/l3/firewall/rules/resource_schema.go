package rules

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var ResourceSchema = schema.Schema{
	MarkdownDescription: "Manage Network Appliance L3 Firewall Rules",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "resource import id",
			Computed:            true,
			Optional:            true,
			CustomType:          jsontypes.StringType,
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 31),
			},
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
		"syslog_default_rule": schema.BoolAttribute{
			MarkdownDescription: "Log the special default rule (boolean value - enable only if you've configured a syslog server) (optional)",
			Optional:            true,
			CustomType:          jsontypes.BoolType,
		},
		"rules": ResourceRulesSchema,
	},
}

var ResourceRulesSchema = schema.ListNestedAttribute{
	Required: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"comment": schema.StringAttribute{
				MarkdownDescription: "Description of the rule (optional)",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dest_cidr": schema.StringAttribute{
				MarkdownDescription: "Comma-separated list of destination IP address(es) (in IP or CIDR notation), fully-qualified domain names (FQDN) or 'Any'",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"dest_port": schema.StringAttribute{
				MarkdownDescription: "Comma-separated list of destination port(s) (integer in the range 1-65535), or 'Any'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"src_cidr": schema.StringAttribute{
				MarkdownDescription: "Comma-separated list of source IP address(es) (in IP or CIDR notation), or 'any' (note: FQDN not supported for source addresses)",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"src_port": schema.StringAttribute{
				MarkdownDescription: "Comma-separated list of source port(s) (integer in the range 1-65535), or 'Any'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"policy": schema.StringAttribute{
				MarkdownDescription: "'allow' or 'deny' traffic specified by this rule",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6', 'Any', or 'any')",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"tcp", "udp", "icmp", "icmp6", "Any", "any"}...),
				},
			},
			"syslog_enabled": schema.BoolAttribute{
				MarkdownDescription: "Log this rule to syslog (true or false, boolean value) - only applicable if a syslog has been configured (optional)",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
		},
	},
}
