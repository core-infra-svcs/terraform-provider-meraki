package subnets

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// GetDataSourceSchema returns the schema for the appliance DHCP subnets data source.
var GetDataSourceSchema = schema.Schema{
	MarkdownDescription: "Retrieve the DHCP subnet information for an appliance.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The ID of the data source instance.",
			Computed:            true,
		},
		"serial": schema.StringAttribute{
			MarkdownDescription: "The serial number of the appliance.",
			Required:            true,
		},
		"resources": DatasourceDataAttributes,
	},
}

// DatasourceDataAttributes defines the "resources" attribute for the data source schema.
var DatasourceDataAttributes = schema.ListNestedAttribute{
	MarkdownDescription: "The list of DHCP subnets.",
	Computed:            true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the data source instance.",
				Computed:            true,
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "The subnet (CIDR) of the DHCP pool.",
				Computed:            true,
			},
			"vlan_id": schema.Int64Attribute{
				MarkdownDescription: "The VLAN ID associated with the subnet.",
				Computed:            true,
			},
			"used_count": schema.Int64Attribute{
				MarkdownDescription: "The number of IP addresses currently in use in the subnet.",
				Computed:            true,
			},
			"free_count": schema.Int64Attribute{
				MarkdownDescription: "The number of IP addresses available in the subnet.",
				Computed:            true,
			},
		},
	},
}
