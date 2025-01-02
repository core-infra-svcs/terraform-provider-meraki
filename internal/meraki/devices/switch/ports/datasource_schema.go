package ports

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var portsDataSourceSchema = schema.Schema{

	MarkdownDescription: "List the switch ports for a switch",

	// The Attributes map describes the fields of the data source.
	Attributes: map[string]schema.Attribute{

		// Every data source must have an ID attribute. This is computed by the framework.
		"id": schema.StringAttribute{
			Computed:   true,
			CustomType: types.StringType,
		},
		"serial": schema.StringAttribute{
			MarkdownDescription: "A list of serial numbers. The returned devices will be filtered to only include these serials.",
			CustomType:          types.StringType,
			Required:            true,
		},
		"list": schema.SetNestedAttribute{
			MarkdownDescription: "List of switch ports",
			Optional:            true,
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"port_id": schema.StringAttribute{
						MarkdownDescription: "The identifier of the switch port.",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the switch port.",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"tags": schema.SetAttribute{
						MarkdownDescription: "The list of tags of the switch port.",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The status of the switch port.",
						CustomType:          types.BoolType,
						Optional:            true,
						Computed:            true,
					},
					"poe_enabled": schema.BoolAttribute{
						MarkdownDescription: "The PoE status of the switch port.",
						CustomType:          types.BoolType,
						Optional:            true,
						Computed:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "The type of the switch port ('trunk' or 'access').",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"vlan": schema.Int64Attribute{
						MarkdownDescription: "The VLAN of the switch port. A null value will clear the value set for trunk ports.",
						CustomType:          types.Int64Type,
						Optional:            true,
						Computed:            true,
					},
					"voice_vlan": schema.Int64Attribute{
						MarkdownDescription: "The voice VLAN of the switch port. Only applicable to access ports.",
						CustomType:          types.Int64Type,
						Optional:            true,
						Computed:            true,
					},
					"allowed_vlans": schema.StringAttribute{
						MarkdownDescription: "The VLANs allowed on the switch port. Only applicable to trunk ports.",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"isolation_enabled": schema.BoolAttribute{
						MarkdownDescription: "The isolation status of the switch port.",
						CustomType:          types.BoolType,
						Optional:            true,
						Computed:            true,
					},
					"rstp_enabled": schema.BoolAttribute{
						MarkdownDescription: "The rapid spanning tree protocol status.",
						CustomType:          types.BoolType,
						Optional:            true,
						Computed:            true,
					},
					"stp_guard": schema.StringAttribute{
						MarkdownDescription: "The state of the STP guard ('disabled', 'root guard', 'bpdu guard' or 'loop guard').",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"access_policy_type": schema.StringAttribute{
						MarkdownDescription: "The type of the access policy of the switch port. Only applicable to access ports. Can be one of 'Open', 'Custom access policy', 'MAC allow list' or 'Sticky MAC allow list'.",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"access_policy_number": schema.Int64Attribute{
						MarkdownDescription: "The number of a custom access policy to configure on the switch port. Only applicable when 'accessPolicyType' is 'Custom access policy'.",
						CustomType:          types.Int64Type,
						Optional:            true,
						Computed:            true,
					},
					"link_negotiation": schema.StringAttribute{
						MarkdownDescription: "The link speed for the switch port.",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"port_schedule_id": schema.StringAttribute{
						MarkdownDescription: "The ID of the port schedule. A value of null will clear the port schedule.",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"udld": schema.StringAttribute{
						MarkdownDescription: "The action to take when Unidirectional Link is detected (Alert only, Enforce). Default configuration is Alert only.",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"sticky_mac_white_list_limit": schema.Int64Attribute{
						MarkdownDescription: "The maximum number of MAC addresses for sticky MAC allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
						CustomType:          types.Int64Type,
						Optional:            true,
						Computed:            true,
					},
					"storm_control_enabled": schema.BoolAttribute{
						MarkdownDescription: "The storm control status of the switch port.",
						CustomType:          types.BoolType,
						Optional:            true,
						Computed:            true,
					},
					"mac_white_list": schema.SetAttribute{
						MarkdownDescription: "Only devices with MAC addresses specified in this list will have access to this port. Up to 20 MAC addresses can be defined. Only applicable when 'accessPolicyType' is 'MAC allow list'.",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"sticky_mac_white_list": schema.SetAttribute{
						MarkdownDescription: "The initial list of MAC addresses for sticky Mac allow list. Only applicable when 'accessPolicyType' is 'Sticky MAC allow list'.",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"adaptive_policy_group_id": schema.StringAttribute{
						MarkdownDescription: "The adaptive policy group ID that will be used to tag traffic through this switch port. This ID must pre-exist during the configuration, else needs to be created using adaptivePolicy/groups API. Cannot be applied to a port on a switch bound to profile.",
						CustomType:          types.StringType,
						Optional:            true,
						Computed:            true,
					},
					"peer_sgt_capable": schema.BoolAttribute{
						MarkdownDescription: "If true, Peer SGT is enabled for traffic through this switch port. Applicable to trunk port only, not access port. Cannot be applied to a port on a switch bound to profile.",
						CustomType:          types.BoolType,
						Optional:            true,
						Computed:            true,
					},
					"flexible_stacking_enabled": schema.BoolAttribute{
						MarkdownDescription: "For supported switches (e.g. MS420/MS425), whether or not the port has flexible stacking enabled.",
						CustomType:          types.BoolType,
						Optional:            true,
						Computed:            true,
					},
					"dai_trusted": schema.BoolAttribute{
						MarkdownDescription: "If true, ARP packets for this port will be considered trusted, and Dynamic ARP Inspection will allow the traffic.",
						CustomType:          types.BoolType,
						Optional:            true,
						Computed:            true,
					},
					"profile": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "When enabled, override this port's configuration with a port profile.",
								Optional:            true,
								Computed:            true,
								CustomType:          types.BoolType,
							},
							"id": schema.StringAttribute{
								MarkdownDescription: "When enabled, the ID of the port profile used to override the port's configuration.",
								Optional:            true,
								Computed:            true,
								CustomType:          types.StringType,
							},
							"iname": schema.StringAttribute{
								MarkdownDescription: "When enabled, the IName of the profile.",
								Optional:            true,
								Computed:            true,
								CustomType:          types.StringType,
							},
						},
					},
				},
			},
		},
	},
}
