package subnets

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

// GetResourceSchema returns the schema for the appliance DHCP subnets resource.
func GetResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manage DHCP subnet configuration for an appliance.",
		Attributes: map[string]schema.Attribute{
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
		},
	}
}
