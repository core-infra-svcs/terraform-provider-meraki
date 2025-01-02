package device

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var GetResourceSchema = schema.Schema{
	MarkdownDescription: "Manage Meraki devices resource. This resource allows updating device attributes for devices associated with a network.",

	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"network_id": schema.StringAttribute{
			MarkdownDescription: "Network ID to which the device belongs.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"serial": schema.StringAttribute{
			MarkdownDescription: "The serial number of the device.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(14, 14),
			},
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the device.",
			Optional:            true,
			Computed:            true,
		},
		"tags": schema.ListAttribute{
			MarkdownDescription: "Tags associated with the device.",
			Optional:            true,
			Computed:            true,
			ElementType:         types.StringType,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"lat": schema.Float64Attribute{
			MarkdownDescription: "The latitude of the device.",
			Optional:            true,
			Computed:            true,
		},
		"lng": schema.Float64Attribute{
			MarkdownDescription: "The longitude of the device.",
			Optional:            true,
			Computed:            true,
		},
		"address": schema.StringAttribute{
			MarkdownDescription: "The physical address of the device.",
			Optional:            true,
			Computed:            true,
		},
		"notes": schema.StringAttribute{
			MarkdownDescription: "Notes about the device, limited to 255 characters.",
			Optional:            true,
			Computed:            true,
		},
		"details": DetailsSchema,
		"move_map_marker": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to set latitude and longitude based on the address. Ignored if `lat` and `lng` are provided.",
			Computed:            true,
			Default:             utils.NewBoolDefault(false),
		},
		"switch_profile_id": schema.StringAttribute{
			MarkdownDescription: "ID of the switch profile to bind to the device. Use `null` to unbind.",
			Optional:            true,
			Computed:            true,
		},
		"floor_plan_id": schema.StringAttribute{
			MarkdownDescription: "Floor plan ID associated with the device. Use `null` to disassociate.",
			Optional:            true,
			Computed:            true,
		},
		"mac": schema.StringAttribute{
			MarkdownDescription: "MAC address of the device.",
			Computed:            true,
		},
		"model": schema.StringAttribute{
			MarkdownDescription: "Model of the device.",
			Computed:            true,
		},
		"lan_ip": schema.StringAttribute{
			MarkdownDescription: "LAN IP address of the device.",
			Computed:            true,
		},
		"firmware": schema.StringAttribute{
			MarkdownDescription: "Firmware version of the device.",
			Computed:            true,
		},
		"url": schema.StringAttribute{
			MarkdownDescription: "URL of the network associated with the device.",
			Computed:            true,
		},
		"beacon_id_params": BeaconIdParamsSchema,
	},
}

var DetailsSchema = schema.ListNestedAttribute{
	MarkdownDescription: "Additional details about the device.",
	Computed:            true,
	Optional:            true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the detail.",
				Required:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value of the detail.",
				Required:            true,
			},
		},
	},
}

var BeaconIdParamsSchema = schema.SingleNestedAttribute{
	MarkdownDescription: "Beacon ID parameters of the device.",
	Computed:            true,
	Optional:            true,
	Attributes: map[string]schema.Attribute{
		"uuid": schema.StringAttribute{
			MarkdownDescription: "UUID of the beacon identifier.",
			Computed:            true,
			Optional:            true,
		},
		"major": schema.Int64Attribute{
			MarkdownDescription: "Major number of the beacon identifier.",
			Computed:            true,
			Optional:            true,
		},
		"minor": schema.Int64Attribute{
			MarkdownDescription: "Minor number of the beacon identifier.",
			Computed:            true,
			Optional:            true,
		},
		"beacon_id": schema.StringAttribute{
			MarkdownDescription: "Beacon ID (alternative to UUID).",
			Computed:            true,
			Optional:            true,
		},
		"proximity": schema.StringAttribute{
			MarkdownDescription: "Proximity of the beacon.",
			Computed:            true,
			Optional:            true,
		},
	},
}
