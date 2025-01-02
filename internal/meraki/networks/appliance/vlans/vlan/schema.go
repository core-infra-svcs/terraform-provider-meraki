package vlan

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ds "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rs "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rs.Schema{

		MarkdownDescription: "Manage the VLANs for an MX network",
		Attributes: map[string]rs.Attribute{
			"id": rs.StringAttribute{
				Computed: true,
			},
			"vlan_id": rs.Int64Attribute{
				Required: true,
			},
			"network_id": rs.StringAttribute{
				MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"interface_id": rs.StringAttribute{
				MarkdownDescription: "The Interface ID",
				Optional:            true,
				Computed:            true,
			},
			"name": rs.StringAttribute{
				MarkdownDescription: "The name of the new VLAN",
				Optional:            true,
				Computed:            true,
			},
			"subnet": rs.StringAttribute{
				MarkdownDescription: "The subnet of the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"appliance_ip": rs.StringAttribute{
				MarkdownDescription: "The local IP of the appliance on the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"group_policy_id": rs.StringAttribute{
				MarkdownDescription: " desired group policy to apply to the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"vpn_nat_subnet": rs.StringAttribute{
				MarkdownDescription: "The translated VPN subnet if VPN and VPN subnet translation are enabled on the VLAN",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_handling": rs.StringAttribute{
				MarkdownDescription: "The appliance's handling of DHCP requests on this VLAN. One of: 'Run a DHCP server', 'Relay DHCP to another server' or 'Do not respond to DHCP requests'",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_relay_server_ips": rs.ListAttribute{
				ElementType: types.StringType,
				Description: "An array of DHCP relay server IPs to which DHCP packets would get relayed for this VLAN",
				Optional:    true,
				Computed:    true,
			},
			"dhcp_lease_time": rs.StringAttribute{
				MarkdownDescription: "The term of DHCP leases if the appliance is running a DHCP server on this VLAN. One of: '30 minutes', '1 hour', '4 hours', '12 hours', '1 day' or '1 week'",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_boot_options_enabled": rs.BoolAttribute{
				MarkdownDescription: "Use DHCP boot options specified in other properties",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_boot_next_server": rs.StringAttribute{
				MarkdownDescription: "DHCP boot option to direct boot clients to the server to load the boot file from",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_boot_filename": rs.StringAttribute{
				MarkdownDescription: "DHCP boot option for boot filename ",
				Optional:            true,
				Computed:            true,
			},
			"fixed_ip_assignments": rs.MapNestedAttribute{
				Description: "The DHCP fixed IP assignments on the VLAN. This should be an object that contains mappings from MAC addresses to objects that themselves each contain \"ip\" and \"name\" string fields. See the sample request/response for more details",
				Optional:    true,
				Computed:    true,
				NestedObject: rs.NestedAttributeObject{
					Attributes: map[string]rs.Attribute{
						"ip": rs.StringAttribute{
							MarkdownDescription: "Enable IPv6 on VLAN.",
							Required:            true,
						},
						"name": rs.StringAttribute{
							MarkdownDescription: "Enable IPv6 on VLAN.",
							Optional:            true,
						},
					},
				},
			},
			"reserved_ip_ranges": rs.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The DHCP reserved IP ranges on the VLAN",
				NestedObject: rs.NestedAttributeObject{
					Attributes: map[string]rs.Attribute{
						"start": rs.StringAttribute{
							MarkdownDescription: "The first IP in the reserved range",
							Optional:            true,
							Computed:            true,
						},
						"end": rs.StringAttribute{
							MarkdownDescription: "The last IP in the reserved range",
							Optional:            true,
							Computed:            true,
						},
						"comment": rs.StringAttribute{
							MarkdownDescription: "A text comment for the reserved range",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"dns_nameservers": rs.StringAttribute{
				MarkdownDescription: "The DNS nameservers used for DHCP responses, either \"upstream_dns\", \"google_dns\", \"opendns\", or a newline seperated string of IP addresses or domain names",
				Optional:            true,
				Computed:            true,
			},
			"dhcp_options": rs.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The list of DHCP options that will be included in DHCP responses. Each object in the list should have \"code\", \"type\", and \"value\" properties.",
				NestedObject: rs.NestedAttributeObject{
					Attributes: map[string]rs.Attribute{
						"code": rs.StringAttribute{
							MarkdownDescription: "The code for the DHCP option. This should be an integer between 2 and 254.",
							Optional:            true,
							Computed:            true,
						},
						"type": rs.StringAttribute{
							MarkdownDescription: "The type for the DHCP option. One of: 'text', 'ip', 'hex' or 'integer'",
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("text", "ip", "hex", "integer"),
							},
						},
						"value": rs.StringAttribute{
							MarkdownDescription: "The value for the DHCP option",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"template_vlan_type": rs.StringAttribute{
				MarkdownDescription: "Type of subnetting of the VLAN. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("same", "unique"),
				},
			},
			"cidr": rs.StringAttribute{
				MarkdownDescription: "CIDR of the pool of subnets. Applicable only for template network. Each network bound to the template will automatically pick a subnet from this pool to build its own VLAN.",
				Optional:            true,
				Computed:            true,
			},
			"mask": rs.Int64Attribute{
				MarkdownDescription: "Mask used for the subnet of all bound to the template networks. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
			},
			"ipv6": rs.SingleNestedAttribute{
				Description: "IPv6 configuration on the VLAN",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]rs.Attribute{
					"enabled": rs.BoolAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						Computed:            true,
					},
					"prefix_assignments": rs.ListNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Prefix assignments on the VLAN",
						NestedObject: rs.NestedAttributeObject{
							Attributes: map[string]rs.Attribute{
								"autonomous": rs.BoolAttribute{
									MarkdownDescription: "Auto assign a /64 prefix from the origin to the VLAN",
									Optional:            true,
									Computed:            true,
								},
								"static_prefix": rs.StringAttribute{
									MarkdownDescription: "Manual configuration of a /64 prefix on the VLAN",
									Optional:            true,
									Computed:            true,
								},
								"static_appliance_ip6": rs.StringAttribute{
									MarkdownDescription: "Manual configuration of the IPv6 Appliance IP",
									Optional:            true,
									Computed:            true,
								},
								"origin": rs.SingleNestedAttribute{
									MarkdownDescription: "The origin of the prefix",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]rs.Attribute{
										"type": rs.StringAttribute{
											MarkdownDescription: "Type of the origin",
											Optional:            true,
											Computed:            true,
											Validators: []validator.String{
												stringvalidator.OneOf("independent", "internet"),
											},
										},
										"interfaces": rs.SetAttribute{
											ElementType: types.StringType,
											Description: "Interfaces associated with the prefix",
											Optional:    true,
											Computed:    true,
										},
									},
								},
							}},
					},
				},
			},
			"mandatory_dhcp": rs.SingleNestedAttribute{
				Description: "Mandatory DHCP will enforce that clients connecting to this VLAN must use the IP address assigned by the DHCP server. Clients who use a static IP address won't be able to associate. Only available on firmware versions 17.0 and above",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]rs.Attribute{
					"enabled": rs.BoolAttribute{
						MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
		},
	}
}

func (r *Datasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = ds.Schema{
		MarkdownDescription: ".",
		Attributes: map[string]ds.Attribute{

			"id": ds.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"network_id": ds.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"list": ds.ListNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: ds.NestedAttributeObject{Attributes: map[string]ds.Attribute{
					"id": ds.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"vlan_id": ds.Int64Attribute{
						Computed: true,
						Optional: true,
					},
					"network_id": ds.StringAttribute{
						MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
						Computed:            true,
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthBetween(8, 31),
						},
					},
					"interface_id": ds.StringAttribute{
						MarkdownDescription: "The Interface ID",
						Optional:            true,
						Computed:            true,
					},
					"name": ds.StringAttribute{
						MarkdownDescription: "The name of the new VLAN",
						Optional:            true,
						Computed:            true,
					},
					"subnet": ds.StringAttribute{
						MarkdownDescription: "The subnet of the VLAN",
						Optional:            true,
						Computed:            true,
					},
					"appliance_ip": ds.StringAttribute{
						MarkdownDescription: "The local IP of the appliance on the VLAN",
						Optional:            true,
						Computed:            true,
					},
					"group_policy_id": ds.StringAttribute{
						MarkdownDescription: " desired group policy to apply to the VLAN",
						Optional:            true,
						Computed:            true,
					},
					"vpn_nat_subnet": ds.StringAttribute{
						MarkdownDescription: "The translated VPN subnet if VPN and VPN subnet translation are enabled on the VLAN",
						Optional:            true,
						Computed:            true,
					},
					"dhcp_handling": ds.StringAttribute{
						MarkdownDescription: "The appliance's handling of DHCP requests on this VLAN. One of: 'Run a DHCP server', 'Relay DHCP to another server' or 'Do not respond to DHCP requests'",
						Optional:            true,
						Computed:            true,
					},
					"dhcp_relay_server_ips": ds.ListAttribute{
						ElementType: types.StringType,
						Description: "An array of DHCP relay server IPs to which DHCP packets would get relayed for this VLAN",
						Optional:    true,
						Computed:    true,
					},
					"dhcp_lease_time": ds.StringAttribute{
						MarkdownDescription: "The term of DHCP leases if the appliance is running a DHCP server on this VLAN. One of: '30 minutes', '1 hour', '4 hours', '12 hours', '1 day' or '1 week'",
						Optional:            true,
						Computed:            true,
					},
					"dhcp_boot_options_enabled": ds.BoolAttribute{
						MarkdownDescription: "Use DHCP boot options specified in other properties",
						Optional:            true,
						Computed:            true,
					},
					"dhcp_boot_next_server": ds.StringAttribute{
						MarkdownDescription: "DHCP boot option to direct boot clients to the server to load the boot file from",
						Optional:            true,
						Computed:            true,
					},
					"dhcp_boot_filename": ds.StringAttribute{
						MarkdownDescription: "DHCP boot option for boot filename ",
						Optional:            true,
						Computed:            true,
					},
					"fixed_ip_assignments": ds.MapNestedAttribute{
						Description: "The DHCP fixed IP assignments on the VLAN. This should be an object that contains mappings from MAC addresses to objects that themselves each contain \"ip\" and \"name\" string fields. See the sample request/response for more details",
						Optional:    true,
						Computed:    true,
						NestedObject: ds.NestedAttributeObject{
							Attributes: map[string]ds.Attribute{
								"ip": ds.StringAttribute{
									MarkdownDescription: "Enable IPv6 on VLAN.",
									Required:            true,
								},
								"name": ds.StringAttribute{
									MarkdownDescription: "Enable IPv6 on VLAN.",
									Optional:            true,
								},
							},
						},
					},
					"reserved_ip_ranges": ds.ListNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The DHCP reserved IP ranges on the VLAN",
						NestedObject: ds.NestedAttributeObject{
							Attributes: map[string]ds.Attribute{
								"start": ds.StringAttribute{
									MarkdownDescription: "The first IP in the reserved range",
									Optional:            true,
									Computed:            true,
								},
								"end": ds.StringAttribute{
									MarkdownDescription: "The last IP in the reserved range",
									Optional:            true,
									Computed:            true,
								},
								"comment": ds.StringAttribute{
									MarkdownDescription: "A text comment for the reserved range",
									Optional:            true,
									Computed:            true,
								},
							},
						},
					},
					"dns_nameservers": ds.StringAttribute{
						MarkdownDescription: "The DNS nameservers used for DHCP responses, either \"upstream_dns\", \"google_dns\", \"opendns\", or a newline seperated string of IP addresses or domain names",
						Optional:            true,
						Computed:            true,
					},
					"dhcp_options": ds.ListNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "The list of DHCP options that will be included in DHCP responses. Each object in the list should have \"code\", \"type\", and \"value\" properties.",
						NestedObject: ds.NestedAttributeObject{
							Attributes: map[string]ds.Attribute{
								"code": ds.StringAttribute{
									MarkdownDescription: "The code for the DHCP option. This should be an integer between 2 and 254.",
									Optional:            true,
									Computed:            true,
								},
								"type": ds.StringAttribute{
									MarkdownDescription: "The type for the DHCP option. One of: 'text', 'ip', 'hex' or 'integer'",
									Optional:            true,
									Computed:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("text", "ip", "hex", "integer"),
									},
								},
								"value": ds.StringAttribute{
									MarkdownDescription: "The value for the DHCP option",
									Optional:            true,
									Computed:            true,
								},
							},
						},
					},
					"template_vlan_type": ds.StringAttribute{
						MarkdownDescription: "Type of subnetting of the VLAN. Applicable only for template network.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("same", "unique"),
						},
					},
					"cidr": ds.StringAttribute{
						MarkdownDescription: "CIDR of the pool of subnets. Applicable only for template network. Each network bound to the template will automatically pick a subnet from this pool to build its own VLAN.",
						Optional:            true,
						Computed:            true,
					},
					"mask": ds.Int64Attribute{
						MarkdownDescription: "Mask used for the subnet of all bound to the template networks. Applicable only for template network.",
						Optional:            true,
						Computed:            true,
					},
					"ipv6": ds.SingleNestedAttribute{
						Description: "IPv6 configuration on the VLAN",
						Optional:    true,
						Computed:    true,
						Attributes: map[string]ds.Attribute{
							"enabled": ds.BoolAttribute{
								MarkdownDescription: "Enable IPv6 on VLAN.",
								Optional:            true,
								Computed:            true,
							},
							"prefix_assignments": ds.ListNestedAttribute{
								Optional:    true,
								Computed:    true,
								Description: "Prefix assignments on the VLAN",
								NestedObject: ds.NestedAttributeObject{
									Attributes: map[string]ds.Attribute{
										"autonomous": ds.BoolAttribute{
											MarkdownDescription: "Auto assign a /64 prefix from the origin to the VLAN",
											Optional:            true,
											Computed:            true,
										},
										"static_prefix": ds.StringAttribute{
											MarkdownDescription: "Manual configuration of a /64 prefix on the VLAN",
											Optional:            true,
											Computed:            true,
										},
										"static_appliance_ip6": ds.StringAttribute{
											MarkdownDescription: "Manual configuration of the IPv6 Appliance IP",
											Optional:            true,
											Computed:            true,
										},
										"origin": ds.SingleNestedAttribute{
											MarkdownDescription: "The origin of the prefix",
											Optional:            true,
											Computed:            true,
											Attributes: map[string]ds.Attribute{
												"type": ds.StringAttribute{
													MarkdownDescription: "Type of the origin",
													Optional:            true,
													Computed:            true,
													Validators: []validator.String{
														stringvalidator.OneOf("independent", "internet"),
													},
												},
												"interfaces": ds.SetAttribute{
													ElementType: types.StringType,
													Description: "Interfaces associated with the prefix",
													Optional:    true,
													Computed:    true,
												},
											},
										},
									}},
							},
						},
					},
					"mandatory_dhcp": ds.SingleNestedAttribute{
						Description: "Mandatory DHCP will enforce that clients connecting to this VLAN must use the IP address assigned by the DHCP server. Clients who use a static IP address won't be able to associate. Only available on firmware versions 17.0 and above",
						Optional:    true,
						Computed:    true,
						Attributes: map[string]ds.Attribute{
							"enabled": ds.BoolAttribute{
								MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
								Optional:            true,
								Computed:            true,
							},
						},
					},
				}}},
		},
	}
}
