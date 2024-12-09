package _interface

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// GetDatasourceSchema provides the schema for the management interface data source.
func GetDatasourceSchema() schema.Schema {
	return schema.Schema{
		Description: "Retrieve the management interface settings for a device.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource.",
				Computed:    true,
			},
			"serial": schema.StringAttribute{
				Description: "Serial number of the device.",
				Optional:    true,
				Computed:    true,
			},
			"ddns_hostnames": schema.SingleNestedAttribute{
				Description: "Dynamic DNS hostnames for the device.",
				Computed:    true,
				Attributes:  DatasourceDdnsHostnamesAttributes,
			},
			"wan1": schema.SingleNestedAttribute{
				Description: "WAN1 interface configuration.",
				Optional:    true,
				Attributes:  DatasourceWanAttributes,
			},
			"wan2": schema.SingleNestedAttribute{
				Description: "WAN2 interface configuration.",
				Optional:    true,
				Attributes:  DatasourceWanAttributes,
			},
		},
	}
}

// DatasourceDdnsHostnamesAttributes defines attributes for the DDNS Hostnames in data sources.
var DatasourceDdnsHostnamesAttributes = map[string]schema.Attribute{
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

// DatasourceWanAttributes defines attributes for WAN1 in data sources.
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
}
