package bandwidth

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
		MarkdownDescription: "Manage Network Appliance Traffic Shaping UplinkBandWidth",
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
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"bandwidth_limit_cellular_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'cellular' uplink. The maximum upload limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"bandwidth_limit_cellular_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'cellular' uplink. The maximum download limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"bandwidth_limit_wan2_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'wan2' uplink. The maximum upload limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"bandwidth_limit_wan2_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'wan2' uplink. The maximum download limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"bandwidth_limit_wan1_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'wan1' uplink. The maximum upload limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"bandwidth_limit_wan1_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'wan1' uplink. The maximum download limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
		},
	}
}
