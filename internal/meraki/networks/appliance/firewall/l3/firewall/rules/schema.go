package rules

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rs.Schema{
		MarkdownDescription: "Manage Network Appliance L3 Firewall Rules",
		Attributes: map[string]rs.Attribute{
			"id": rs.StringAttribute{
				MarkdownDescription: "Network ID",
				Computed:            true,
				Optional:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"network_id": rs.StringAttribute{
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
			"rules": rs.ListNestedAttribute{
				Required: true,
				NestedObject: rs.NestedAttributeObject{
					Attributes: map[string]rs.Attribute{
						"comment": rs.StringAttribute{
							MarkdownDescription: "Description of the rule (optional)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_cidr": rs.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination IP address(es) (in IP or CIDR notation), fully-qualified domain names (FQDN) or 'Any'",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_port": rs.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination port(s) (integer in the range 1-65535), or 'Any'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_cidr": rs.StringAttribute{
							MarkdownDescription: "Comma-separated list of source IP address(es) (in IP or CIDR notation), or 'any' (note: FQDN not supported for source addresses)",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_port": rs.StringAttribute{
							MarkdownDescription: "Comma-separated list of source port(s) (integer in the range 1-65535), or 'Any'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"policy": rs.StringAttribute{
							MarkdownDescription: "'allow' or 'deny' traffic specified by this rule",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"protocol": rs.StringAttribute{
							MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6', 'Any', or 'any')",
							Required:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"tcp", "udp", "icmp", "icmp6", "Any", "any"}...),
							},
						},
						"syslog_enabled": rs.BoolAttribute{
							MarkdownDescription: "Log this rule to syslog (true or false, boolean value) - only applicable if a syslog has been configured (optional)",
							Optional:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ds.Schema{
		MarkdownDescription: "Get l3 firewall rules",

		Attributes: map[string]ds.Attribute{
			"id": ds.StringAttribute{
				MarkdownDescription: "Network ID",
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"network_id": ds.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"syslog_default_rule": ds.BoolAttribute{
				MarkdownDescription: "Log the special default rule (boolean value - enable only if you've configured a syslog server) (optional)",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"rules": ds.ListNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: ds.NestedAttributeObject{
					Attributes: map[string]ds.Attribute{
						"comment": ds.StringAttribute{
							MarkdownDescription: "Description of the rule (optional)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_cidr": ds.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination IP address(es) (in IP or CIDR notation), fully-qualified domain names (FQDN) or 'Any'",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_port": ds.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination port(s) (integer in the range 1-65535), or 'Any'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_cidr": ds.StringAttribute{
							MarkdownDescription: "Comma-separated list of source IP address(es) (in IP or CIDR notation), or 'any' (note: FQDN not supported for source addresses)",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_port": ds.StringAttribute{
							MarkdownDescription: "Comma-separated list of source port(s) (integer in the range 1-65535), or 'Any'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"policy": ds.StringAttribute{
							MarkdownDescription: "'allow' or 'deny' traffic specified by this rule",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"protocol": ds.StringAttribute{
							MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6', 'Any', or 'any')",
							Required:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"tcp", "udp", "icmp", "icmp6", "Any", "any"}...),
							},
						},
						"syslog_enabled": ds.BoolAttribute{
							MarkdownDescription: "Log this rule to syslog (true or false, boolean value) - only applicable if a syslog has been configured (optional)",
							Optional:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
		},
	}
}
