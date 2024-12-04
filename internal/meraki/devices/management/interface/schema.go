package _interface

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// resourceSchema defines the schema for the management interface resource.
var resourceSchema = schema.Schema{
	MarkdownDescription: "Manage the management interface settings for a device",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the resource.",
			Computed:            true,
			CustomType:          types.StringType,
		},
		"serial": schema.StringAttribute{
			MarkdownDescription: "Serial number of the device.",
			Required:            true,
			CustomType:          types.StringType,
		},
		"ddns_hostnames": schema.SingleNestedAttribute{
			MarkdownDescription: "Dynamic DNS hostnames for the device.",
			Computed:            true,
			Attributes:          ddnsHostnameAttributes(),
		},
		"wan1": schema.SingleNestedAttribute{
			MarkdownDescription: "WAN1 interface configuration.",
			Optional:            true,
			Computed:            true,
			Attributes:          wanAttributes,
		},
		"wan2": schema.SingleNestedAttribute{
			MarkdownDescription: "WAN2 interface configuration.",
			Optional:            true,
			Computed:            true,
			Attributes:          wanAttributes,
		},
	},
}

// dataSourceSchema defines the schema for the management interface data source.
var dataSourceSchema = datasourceSchema.Schema{
	MarkdownDescription: "Retrieve the management interface settings for a device.",
	Attributes: map[string]datasourceSchema.Attribute{
		"id": schema.StringAttribute{
			Computed:   true,
			CustomType: types.StringType,
		},
		"serial": schema.StringAttribute{
			MarkdownDescription: "Serial number",
			Optional:            true,
			Computed:            true,
			CustomType:          types.StringType,
		},
		"ddns_hostnames": schema.SingleNestedAttribute{
			MarkdownDescription: "Dynamic DNS hostnames for the device.",
			Computed:            true,
			Attributes:          ddnsHostnameAttributes(),
		},
		"wan1": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: wanAttributes,
		},
		"wan2": schema.SingleNestedAttribute{
			Optional:   true,
			Attributes: wanAttributes,
		},
	},
}

// ddnsHostnameAttributes defines the attributes for DDNS hostnames.
func ddnsHostnameAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"active_ddns_hostname": schema.StringAttribute{
			MarkdownDescription: "The active DDNS hostname.",
			Computed:            true,
		},
		"ddns_hostname_wan1": schema.StringAttribute{
			MarkdownDescription: "DDNS hostname for WAN1.",
			Computed:            true,
		},
		"ddns_hostname_wan2": schema.StringAttribute{
			MarkdownDescription: "DDNS hostname for WAN2.",
			Computed:            true,
		},
	}
}

// wanAttributes defines the shared attributes for WAN1 and WAN2 configurations.
var wanAttributes = map[string]schema.Attribute{
	"wan_enabled": schema.StringAttribute{
		MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled' for MX devices. Leave value null for MR and MS devices.",
		Optional:            true,
		Computed:            true,
		CustomType:          types.StringType,
	},
	"using_static_ip": schema.BoolAttribute{
		MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
		Optional:            true,
		Computed:            true,
		CustomType:          types.BoolType,
	},
	"vlan": schema.Int64Attribute{
		MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
		Optional:            true,
		Computed:            true,
		CustomType:          types.Int64Type,
	},
	"static_ip": schema.StringAttribute{
		MarkdownDescription: "The IP the device should use on the WAN.",
		Optional:            true,
		Computed:            true,
		CustomType:          types.StringType,
	},
	"static_subnet_mask": schema.StringAttribute{
		MarkdownDescription: "The subnet mask for the WAN.",
		Optional:            true,
		Computed:            true,
		CustomType:          types.StringType,
	},
	"static_gateway_ip": schema.StringAttribute{
		MarkdownDescription: "The IP of the gateway on the WAN.",
		Optional:            true,
		Computed:            true,
		CustomType:          types.StringType,
	},
	"static_dns": schema.ListAttribute{
		MarkdownDescription: "Up to two DNS IPs.",
		Optional:            true,
		Computed:            true,
		ElementType:         types.StringType,
	},
}
