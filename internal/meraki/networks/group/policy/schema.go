package policy

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ScheduleDaySchema = schema.SingleNestedAttribute{
	Optional: true,
	Computed: true,
	Attributes: map[string]schema.Attribute{
		"active": schema.BoolAttribute{
			Required: true,
		},
		"from": schema.StringAttribute{
			Required: true,
		},
		"to": schema.StringAttribute{
			Required: true,
		},
	},
}

var SchedulingSchema = schema.SingleNestedAttribute{
	Optional:    true,
	Computed:    true,
	Description: "The scheduling settings of the group policy.",
	Attributes: map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			Required: true,
		},
		"sunday":    ScheduleDaySchema,
		"monday":    ScheduleDaySchema,
		"tuesday":   ScheduleDaySchema,
		"wednesday": ScheduleDaySchema,
		"thursday":  ScheduleDaySchema,
		"friday":    ScheduleDaySchema,
		"saturday":  ScheduleDaySchema,
	},
}

var TrafficShapingRulesSchema = schema.ListNestedAttribute{
	Optional: true,
	Computed: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"dscp_tag_value": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"pcp_tag_value": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"per_client_bandwidth_limits": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						Required: true,
					},
					"bandwidth_limits": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"limit_up": schema.Int64Attribute{
								Optional: true,
								Computed: true,
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"limit_down": schema.Int64Attribute{
								Optional: true,
								Computed: true,
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"definitions": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	},
}

var FirewallAndTrafficShapingSchema = schema.SingleNestedAttribute{
	Optional:    true,
	Computed:    true,
	Description: "The firewall and traffic shaping settings of the group policy.",
	Attributes: map[string]schema.Attribute{
		"settings": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"l3_firewall_rules": schema.ListNestedAttribute{
			Optional: true,
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"comment": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"policy": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"protocol": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"dest_port": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"dest_cidr": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"l7_firewall_rules": schema.ListNestedAttribute{
			Optional: true,
			Computed: true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"policy": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"type": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"value": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
		},
		"traffic_shaping_rules": TrafficShapingRulesSchema,
	},
}

var UrlPatternsSchema = schema.SingleNestedAttribute{
	Optional: true,
	Computed: true,
	Attributes: map[string]schema.Attribute{
		"patterns": schema.ListAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"settings": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Default:  utils.NewStringDefault("network default"),
			Validators: []validator.String{
				stringvalidator.OneOf("network default", "append", "override"),
			},
		},
	},
}

var ContentFilteringSchema = schema.SingleNestedAttribute{
	Optional:    true,
	Computed:    true,
	Description: "The content filtering settings of the group policy.",
	Attributes: map[string]schema.Attribute{
		"allowed_url_patterns": UrlPatternsSchema,
		"blocked_url_patterns": UrlPatternsSchema,
		"blocked_url_categories": schema.SingleNestedAttribute{
			Optional: true,
			Computed: true,
			Attributes: map[string]schema.Attribute{
				"categories": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Computed:    true,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
				},
				"settings": schema.StringAttribute{
					Optional: true,
					Computed: true,
					Default:  utils.NewStringDefault("network default"),
					Validators: []validator.String{
						stringvalidator.OneOf("network default", "append", "override"),
					},
				},
			},
		},
	},
}

var ResourceSchema = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		Computed:    true,
		Description: "The unique identifier for the resource, generated by the Meraki API.",
	},
	"group_policy_id": schema.StringAttribute{
		Computed:    true,
		Description: "The unique identifier for the group policy.",
	},
	"network_id": schema.StringAttribute{
		Required:    true,
		Description: "The network Id where the group policy is applied.",
		Validators: []validator.String{
			stringvalidator.LengthBetween(1, 31),
		},
	},
	"name": schema.StringAttribute{
		Required:    true,
		Description: "The name of the group policy.",
	},
	"scheduling": SchedulingSchema,
	"bandwidth": schema.SingleNestedAttribute{
		Optional:    true,
		Computed:    true,
		Description: "The bandwidth settings of the group policy.",
		Attributes: map[string]schema.Attribute{
			"settings": schema.StringAttribute{
				Required: true,
			},
			"bandwidth_limits": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"limit_up": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "The upload bandwidth limit. Can be null.",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"limit_down": schema.Int64Attribute{
						Optional:    true,
						Computed:    true,
						Description: "The download bandwidth limit. Can be null.",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
	},
	"firewall_and_traffic_shaping": FirewallAndTrafficShapingSchema,
	"content_filtering":            ContentFilteringSchema,
	"splash_auth_settings": schema.StringAttribute{
		Optional:    true,
		Computed:    true,
		Description: "The splash authentication settings of the group policy.",
	},
	"vlan_tagging": schema.SingleNestedAttribute{
		Optional:    true,
		Computed:    true,
		Description: "The VLAN tagging settings of the group policy.",
		Attributes: map[string]schema.Attribute{
			"settings": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"vlan_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	},
	"bonjour_forwarding": schema.SingleNestedAttribute{
		Optional:    true,
		Computed:    true,
		Description: "The bonjour forwarding settings of the group policy.",
		Attributes: map[string]schema.Attribute{
			"settings": schema.StringAttribute{
				Required: true,
			},
			"rules": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"vlan_id": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"services": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	},
}
