package device

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Schema provides a way to define the structure of the resource data.
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage network Devices resource. This only works for devices associated with a network.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "The devices serial number",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(14, 14),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of a device",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				Description: "Network tags",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"lat": schema.Float64Attribute{
				MarkdownDescription: "The latitude of a device",
				Optional:            true,
				Computed:            true,
			},
			"lng": schema.Float64Attribute{
				MarkdownDescription: "The longitude of a device",
				Optional:            true,
				Computed:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The address of a device",
				Optional:            true,
				Computed:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Notes for the network",
				Optional:            true,
				Computed:            true,
			},
			"details": schema.ListNestedAttribute{
				Description: "Network tags",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of a device",
						Optional:            true,
						Computed:            true,
					},
					"value": schema.StringAttribute{
						MarkdownDescription: "The value of a device",
						Optional:            true,
						Computed:            true,
					},
				}},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"move_map_marker": schema.BoolAttribute{
				MarkdownDescription: "Whether or not to set the latitude and longitude of a device based on the new address. Only applies when lat and lng are not specified.",
				Optional:            true,
				Computed:            true,
			},
			"switch_profile_id": schema.StringAttribute{
				MarkdownDescription: "The ID of a switch profile to bind to the device (for available switch profiles, see the 'Switch Profiles' endpoint). Use null to unbind the switch device from the current profile. For a device to be bindable to a switch profile, it must (1) be a switch, and (2) belong to a network that is bound to a configuration template.",
				Optional:            true,
				Computed:            true,
			},
			"floor_plan_id": schema.StringAttribute{
				MarkdownDescription: "The floor plan to associate to this device. null disassociates the device from the floor plan.",
				Optional:            true,
				Computed:            true,
			},
			"mac": schema.StringAttribute{
				MarkdownDescription: "The mac address of a device",
				Optional:            true,
				Computed:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The model of a device",
				Optional:            true,
				Computed:            true,
			},
			"lan_ip": schema.StringAttribute{
				MarkdownDescription: "The ipv4 lan ip of a device",
				Optional:            true,
				Computed:            true,
			},
			"firmware": schema.StringAttribute{
				MarkdownDescription: "The firmware version of a device",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The url for the network associated with the device.",
				Optional:            true,
				Computed:            true,
			},
			"beacon_id_params": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Computed: true,
					},
					"major": schema.Int64Attribute{
						Computed: true,
					},
					"minor": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
		},
	}
}
