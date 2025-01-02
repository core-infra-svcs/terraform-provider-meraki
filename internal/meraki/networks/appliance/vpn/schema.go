package vpn

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage networks appliance vpn site to site vpn. Only valid for MX networks in NAT mode.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "The site-to-site VPN mode.",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"hubs": schema.ListNestedAttribute{
				Description: "The list of VPN hubs, in order of preference.",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hub_id": schema.StringAttribute{
							MarkdownDescription: "The network ID of the hub",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"use_default_route": schema.BoolAttribute{
							MarkdownDescription: "Indicates whether default route traffic should be sent to this hub.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
			"subnets": schema.ListNestedAttribute{
				Description: "The list of subnets and their VPN presence.",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"local_subnet": schema.StringAttribute{
							MarkdownDescription: "The CIDR notation subnet used within the VPN",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"use_vpn": schema.BoolAttribute{
							MarkdownDescription: "Indicates the presence of the subnet in the VPN.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
		},
	}
}
