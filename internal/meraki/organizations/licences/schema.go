package licences

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Schema provides a way to define the structure of the data source data.
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{

		// It should provide a clear and concise description of the data source.
		MarkdownDescription: "Ports Organization Licenses",

		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{

			// Every data source must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization Id",
				CustomType:          jsontypes.StringType,
				Required:            true,
			},
			"per_page": schema.Int64Attribute{
				MarkdownDescription: "The number of entries per page returned. Acceptable range is 3 - 1000. Default is 1000.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"starting_after": schema.StringAttribute{
				MarkdownDescription: "A token used by the server to indicate the start of the page. Often this is a timestamp or an ID but it is not limited to those. This parameter should not be defined by client applications. The link for the first, last, prev, or next page in the HTTP Link header should define it.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"ending_before": schema.StringAttribute{
				MarkdownDescription: "A token used by the server to indicate the end of the page. Often this is a timestamp or an ID but it is not limited to those. This parameter should not be defined by client applications. The link for the first, last, prev, or next page in the HTTP Link header should define it.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"device_serial": schema.StringAttribute{
				MarkdownDescription: "Filter the licenses to those assigned to a particular device. Returned in the same order that they are queued to the device",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Filter the licenses to those assigned in a particular network",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "Filter the licenses to those in a particular state. Can be one of 'active', 'expired', 'expiring', 'unused', 'unusedActive' or 'recentlyQueued'",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"list": schema.SetNestedAttribute{
				MarkdownDescription: "Ports of organization acls",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"license_id": schema.StringAttribute{
							MarkdownDescription: "License ID",
							CustomType:          jsontypes.StringType,
							Required:            true,
						},
						"device_serial": schema.StringAttribute{
							MarkdownDescription: "The serial number of the device to assign this license to. Set this to null to unassign the license. If a different license is already active on the device, this parameter will control queueing/dequeuing this license.",
							CustomType:          jsontypes.StringType,
							Required:            true,
						},
						"license_type": schema.StringAttribute{
							MarkdownDescription: "License Type.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"license_key": schema.StringAttribute{
							MarkdownDescription: "License Key.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"order_number": schema.StringAttribute{
							MarkdownDescription: "Order SsidNumber.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"network_id": schema.StringAttribute{
							MarkdownDescription: "ID of the network the license is assigned to.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"state": schema.StringAttribute{
							MarkdownDescription: "The state of the license. All queued licenses have a status of `recentlyQueued`.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"claim_date": schema.StringAttribute{
							MarkdownDescription: "The date the license was claimed into the organization.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"activation_date": schema.StringAttribute{
							MarkdownDescription: "The date the license started burning.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"expiration_date": schema.StringAttribute{
							MarkdownDescription: "The date the license will expire.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"head_license_id": schema.StringAttribute{
							MarkdownDescription: "The id of the head license this license is queued behind. If there is no head license, it returns nil.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"seat_count": schema.Int64Attribute{
							MarkdownDescription: "The number of seats of the license. Only applicable to SM licenses.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"total_duration_in_days": schema.Int64Attribute{
							MarkdownDescription: "The duration of the license plus all permanently queued licenses associated with it.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"duration_in_days": schema.Int64Attribute{
							MarkdownDescription: "The duration of the individual license.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"permanently_queued_licenses": schema.SingleNestedAttribute{
							MarkdownDescription: "DEPRECATED Ports of permanently queued licenses attached to the license. Instead, use /organizations/{organizationId}/licenses?deviceSerial= to retrieved queued licenses for a given device.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Permanently queued license ID.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"license_type": schema.StringAttribute{
									MarkdownDescription: "License type.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"license_key": schema.StringAttribute{
									MarkdownDescription: "License key.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"order_number": schema.StringAttribute{
									MarkdownDescription: "Order number.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"duration_in_days": schema.Int64Attribute{
									MarkdownDescription: "The duration of the individual license.",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
							},
						},
					},
				},
			},
		},
	}
}
