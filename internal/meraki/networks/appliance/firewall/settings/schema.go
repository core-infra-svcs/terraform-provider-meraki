package settings

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
		MarkdownDescription: "Manage network appliance firewall settings.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"spoofing_protection": schema.SingleNestedAttribute{
				MarkdownDescription: "Spoofing protection settings",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"ip_source_guard": schema.SingleNestedAttribute{
						MarkdownDescription: "IP source address spoofing settings",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								MarkdownDescription: "Mode of protection.",
								Required:            true,
								CustomType:          jsontypes.StringType,
							},
						},
					},
				},
			},
		},
	}
}
