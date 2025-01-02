package _interface

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var GetDatasourceSchema = schema.Schema{
	Description: "Retrieve the management interface settings for a device.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier for the resource.",
			Computed:    true,
		},
		"serial": schema.StringAttribute{
			Description: "Serial number of the device.",
			Required:    true,
		},
		"ddns_hostnames": schema.SingleNestedAttribute{
			Description: "Dynamic DNS hostnames for the device.",
			Optional:    true,
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"active_ddns_hostname": schema.StringAttribute{
					Description: "The active DDNS hostname.",
					Optional:    true,
					Computed:    true,
				},
				"ddns_hostname_wan1": schema.StringAttribute{
					Description: "DDNS hostname for WAN1.",
					Optional:    true,
					Computed:    true,
				},
				"ddns_hostname_wan2": schema.StringAttribute{
					Description: "DDNS hostname for WAN2.",
					Optional:    true,
					Computed:    true,
				},
			},
		},
		"wan1": schema.SingleNestedAttribute{
			Description: "WAN1 interface configuration.",
			Optional:    true,
			Computed:    true,
			Attributes:  DatasourceWanAttributes,
		},
		"wan2": schema.SingleNestedAttribute{
			Description: "WAN2 interface configuration.",
			Optional:    true,
			Computed:    true,
			Attributes:  DatasourceWanAttributes,
		},
	},
}

var DatasourceWanAttributes = map[string]schema.Attribute{
	"wan_enabled": schema.StringAttribute{
		Description: "Enable or disable the WAN interface.",
		Optional:    true,
		Computed:    true,
	},
	"using_static_ip": schema.BoolAttribute{
		Description: "Whether the WAN interface is using a static IP.",
		Optional:    true,
		Computed:    true,
	},
	"static_ip": schema.StringAttribute{
		Description: "Static IP address.",
		Optional:    true,
		Computed:    true,
	},
	"static_subnet_mask": schema.StringAttribute{
		Description: "Static subnet mask.",
		Optional:    true,
		Computed:    true,
	},
	"static_gateway_ip": schema.StringAttribute{
		Description: "Static gateway IP.",
		Optional:    true,
		Computed:    true,
	},
	"static_dns": schema.ListAttribute{
		Description: "List of static DNS IPs.",
		ElementType: types.StringType,
		Optional:    true,
		Computed:    true,
	},
	"vlan": schema.Int64Attribute{
		Description: "VLAN ID.",
		Optional:    true,
		Computed:    true,
	},
}

var DatasourceStaticDnsAttributes = schema.ListAttribute{
	Description: "A list of static DNS IP addresses.",
	ElementType: types.StringType,
	Optional:    true,
	Computed:    true,
}
