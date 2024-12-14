package _interface

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// GetResourceSchema defines the schema for the resource.
var GetResourceSchema = schema.Schema{
	MarkdownDescription: "Manage the management interface settings for a device",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the resource.",
			Computed:            true,
		},
		"serial": schema.StringAttribute{
			MarkdownDescription: "Serial number of the device.",
			Required:            true,
		},
		"ddns_hostnames": schema.SingleNestedAttribute{
			MarkdownDescription: "Dynamic DNS hostnames for the device.",
			Computed:            true,
			Attributes:          ResourceDdnsHostnamesAttributes,
		},
		"wan1": schema.SingleNestedAttribute{
			MarkdownDescription: "WAN1 interface configuration.",
			Optional:            true,
			Computed:            true,
			Attributes:          ResourceWanAttributes,
		},
		"wan2": schema.SingleNestedAttribute{
			MarkdownDescription: "WAN2 interface configuration.",
			Optional:            true,
			Computed:            true,
			Attributes:          ResourceWanAttributes,
		},
	},
}

// ResourceDdnsHostnamesAttributes defines attributes for the DDNS Hostnames in resources.
var ResourceDdnsHostnamesAttributes = map[string]schema.Attribute{
	"active_ddns_hostname": schema.StringAttribute{
		Description: "The active DDNS hostname.",
		Computed:    true,
	},
	"ddns_hostname_wan1": schema.StringAttribute{
		Description: "DDNS hostname for WAN1.",
		Computed:    true,
	},
	"ddns_hostname_wan2": schema.StringAttribute{
		Description: "DDNS hostname for WAN2.",
		Computed:    true,
	},
}

// ResourceWanAttributes defines attributes for WAN in resources.
var ResourceWanAttributes = map[string]schema.Attribute{
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
}
