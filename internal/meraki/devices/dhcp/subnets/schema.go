package subnets

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// SchemaResource returns the schema for the appliance DHCP subnets resource.
func SchemaResource() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manage DHCP subnet configuration for an appliance.",
		Attributes:          ResourceAttributes(),
	}
}

// ResourceAttributes defines the attributes for the appliance DHCP subnet resource schema.
func ResourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"serial": schema.StringAttribute{
			MarkdownDescription: "The serial number of the appliance.",
			Required:            true,
		},
		"subnet": schema.StringAttribute{
			MarkdownDescription: "The subnet (CIDR) of the DHCP pool.",
			Required:            true,
		},
		"vlan_id": schema.Int64Attribute{
			MarkdownDescription: "The VLAN ID associated with the subnet.",
			Required:            true,
		},
		"used_count": schema.Int64Attribute{
			MarkdownDescription: "The number of IP addresses currently in use in the subnet.",
			Required:            true,
		},
		"free_count": schema.Int64Attribute{
			MarkdownDescription: "The number of IP addresses available in the subnet.",
			Required:            true,
		},
	}
}

// SchemaDataSource returns the schema for the appliance DHCP subnets data source.
func SchemaDataSource() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		MarkdownDescription: "Retrieve the DHCP subnet information for an appliance.",
		Attributes: map[string]datasourceSchema.Attribute{
			"id": datasourceSchema.StringAttribute{
				MarkdownDescription: "The ID of the data source instance.",
				Computed:            true,
			},
			"serial": datasourceSchema.StringAttribute{
				MarkdownDescription: "The serial number of the appliance.",
				Required:            true,
			},
			"resources": datasourceAttributes(),
		},
	}
}

// datasourceAttributes defines the "resources" attribute for the data source schema.
func datasourceAttributes() datasourceSchema.ListNestedAttribute {
	return datasourceSchema.ListNestedAttribute{
		MarkdownDescription: "The list of DHCP subnets.",
		Computed:            true,
		NestedObject: datasourceSchema.NestedAttributeObject{
			Attributes: ResourcesAttributeForDataSource(),
		},
	}
}

// ResourcesAttributeForDataSource defines the attributes for each DHCP subnet in the data source schema.
func ResourcesAttributeForDataSource() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"subnet": datasourceSchema.StringAttribute{
			MarkdownDescription: "The subnet (CIDR) of the DHCP pool.",
			Computed:            true,
		},
		"vlan_id": datasourceSchema.Int64Attribute{
			MarkdownDescription: "The VLAN ID associated with the subnet.",
			Computed:            true,
		},
		"used_count": datasourceSchema.Int64Attribute{
			MarkdownDescription: "The number of IP addresses currently in use in the subnet.",
			Computed:            true,
		},
		"free_count": datasourceSchema.Int64Attribute{
			MarkdownDescription: "The number of IP addresses available in the subnet.",
			Computed:            true,
		},
	}
}
