package _interface

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

// ddnsHostnameSchema creates the schema for the DDNS hostnames.
func ddnsHostnameSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Dynamic DNS hostnames for the device.",
		Computed:            true,
		Attributes:          ddnsHostnameAttributes(),
	}
}

// resourceAttributes defines the attributes for the management interface resource schema.
var resourceAttributes = map[string]schema.Attribute{
	"id": schema.StringAttribute{
		MarkdownDescription: "The unique identifier for the resource.",
		Computed:            true,
	},
	"serial": schema.StringAttribute{
		MarkdownDescription: "Serial number of the device.",
		Required:            true,
	},
	"ddns_hostnames": ddnsHostnameSchema(),
	"wan1": schema.SingleNestedAttribute{
		MarkdownDescription: "WAN1 interface configuration.",
		Computed:            true,
		Attributes:          wanAttributes,
	},
	"wan2": schema.SingleNestedAttribute{
		MarkdownDescription: "WAN2 interface configuration.",
		Computed:            true,
		Attributes:          wanAttributes,
	},
}

// resourceSchema defines the schema for the management interface resource.
var resourceSchema = schema.Schema{
	MarkdownDescription: "Manage the management interface settings for a device.",
	Attributes:          resourceAttributes,
}

// dataSourceSchema defines the schema for the management interface data source.
var dataSourceSchema = datasourceSchema.Schema{
	MarkdownDescription: "Retrieve the management interface settings for a device.",
	Attributes:          utils.ConvertResourceSchemaToDataSourceSchema(resourceAttributes),
}

// wanAttributes defines the shared attributes for WAN1 and WAN2 configurations.
var wanAttributes = map[string]schema.Attribute{
	"wan_enabled": schema.StringAttribute{
		MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
		Computed:            true,
	},
	"using_static_ip": schema.BoolAttribute{
		MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
		Computed:            true,
	},
	"vlan": schema.Int64Attribute{
		MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether `using_static_ip` is true or false.",
		Computed:            true,
	},
	"static_ip": schema.StringAttribute{
		MarkdownDescription: "The IP the device should use on the WAN.",
		Computed:            true,
	},
	"static_subnet_mask": schema.StringAttribute{
		MarkdownDescription: "The subnet mask for the WAN.",
		Computed:            true,
	},
	"static_gateway_ip": schema.StringAttribute{
		MarkdownDescription: "The IP of the gateway on the WAN.",
		Computed:            true,
	},
	"static_dns": schema.ListAttribute{
		MarkdownDescription: "Up to two DNS IPs.",
		Computed:            true,
		ElementType:         types.StringType,
	},
}
