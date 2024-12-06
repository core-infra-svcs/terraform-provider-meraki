package ports

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rs.Schema{
		MarkdownDescription: "NetworksAppliancePorts resource for updating Network Appliance Firewall L3 Firewall Rules.",
		Attributes: map[string]rs.Attribute{
			"id": rs.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
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
			"port_id": rs.StringAttribute{
				MarkdownDescription: "Port ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"access_policy": rs.StringAttribute{
				MarkdownDescription: "The name of the policy. Only applicable to Access ports.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"allowed_vlans": rs.StringAttribute{
				MarkdownDescription: "Comma-delimited list of the VLAN ID's allowed on the port, or 'all' to permit all VLAN's on the port.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"drop_untagged_traffic": rs.BoolAttribute{
				MarkdownDescription: "Whether the trunk port can drop all untagged traffic.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"enabled": rs.BoolAttribute{
				Description:         "The status of the port",
				MarkdownDescription: "The status of the port",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"number": rs.Int64Attribute{
				MarkdownDescription: "SsidNumber of the port",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"type": rs.StringAttribute{
				MarkdownDescription: "The type of the port: 'access' or 'trunk'.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"access", "trunk"}...),
					stringvalidator.LengthAtLeast(4),
				},
				CustomType: jsontypes.StringType,
			},
			"vlan": rs.Int64Attribute{
				MarkdownDescription: "Native VLAN when the port is in Trunk mode. Access VLAN when the port is in Access mode.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
		},
	}
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ds.Schema{
		MarkdownDescription: "Get appliance ports",

		Attributes: map[string]ds.Attribute{

			"id": ds.StringAttribute{
				Description:         "Example identifier",
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": ds.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"list": ds.SetNestedAttribute{
				MarkdownDescription: "Ports of Network Appliance Ports",
				Optional:            true,
				Computed:            true,
				NestedObject: ds.NestedAttributeObject{
					Attributes: map[string]ds.Attribute{
						"access_policy": ds.StringAttribute{
							MarkdownDescription: "The name of the policy. Only applicable to Access ports.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"allowed_vlans": ds.StringAttribute{
							MarkdownDescription: "Comma-delimited list of the VLAN ID's allowed on the port, or 'all' to permit all VLAN's on the port.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"drop_untagged_traffic": ds.BoolAttribute{
							MarkdownDescription: "Whether the trunk port can drop all untagged traffic.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"enabled": ds.BoolAttribute{
							Description:         "The status of the port",
							MarkdownDescription: "The status of the port",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
						"number": ds.Int64Attribute{
							MarkdownDescription: "SsidNumber of the port",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"type": ds.StringAttribute{
							MarkdownDescription: "The type of the port: 'access' or 'trunk'.",
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"access", "trunk"}...),
								stringvalidator.LengthAtLeast(4),
							},
							CustomType: jsontypes.StringType,
						},
						"vlan": ds.Int64Attribute{
							MarkdownDescription: "Native VLAN when the port is in Trunk mode. Access VLAN when the port is in Access mode.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
					},
				},
			},
		},
	}
}
