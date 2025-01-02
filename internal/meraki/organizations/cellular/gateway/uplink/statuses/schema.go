package statuses

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{

		MarkdownDescription: "OrganizationsCellularGatewayUplinkStatuses",

		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{

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
			"serials": schema.SetAttribute{
				MarkdownDescription: "A list of serial numbers. The returned devices will be filtered to only include these serials.",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
			"network_ids": schema.SetAttribute{
				MarkdownDescription: "A list of network IDs. The returned devices will be filtered to only include these networks.",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
			"iccids": schema.SetAttribute{
				MarkdownDescription: "A list of ICCIDs. The returned devices will be filtered to only include these ICCIDs.",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
			"list": schema.SetNestedAttribute{
				MarkdownDescription: "Ports of organization acls",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"network_id": schema.StringAttribute{
							MarkdownDescription: "ID of the network.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"serial": schema.StringAttribute{
							MarkdownDescription: "Serial number of the device.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"model": schema.StringAttribute{
							MarkdownDescription: "Device model.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"last_reported_at": schema.StringAttribute{
							MarkdownDescription: "Last reported time for the device.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"uplinks": schema.ListNestedAttribute{
							MarkdownDescription: "Uplinks info.",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"interface": schema.StringAttribute{
										MarkdownDescription: "Uplink interface.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"status": schema.StringAttribute{
										MarkdownDescription: "Uplink status.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"ip": schema.StringAttribute{
										MarkdownDescription: "Uplink ip.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"provider": schema.StringAttribute{
										MarkdownDescription: "Network Provider.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"public_ip": schema.StringAttribute{
										MarkdownDescription: "Public IP.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"model": schema.StringAttribute{
										MarkdownDescription: "Uplink model.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"signal_stat": schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{
											"rsrp": schema.StringAttribute{
												MarkdownDescription: "Reference Signal Received Power.",
												CustomType:          jsontypes.StringType,
												Optional:            true,
												Computed:            true,
											},
											"rsrq": schema.StringAttribute{
												MarkdownDescription: "Reference Signal Received Quality.",
												CustomType:          jsontypes.StringType,
												Optional:            true,
												Computed:            true,
											},
										},
									},
									"connection_type": schema.StringAttribute{
										MarkdownDescription: "Connection Type.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"apn": schema.StringAttribute{
										MarkdownDescription: "Access Point Name.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"gateway": schema.StringAttribute{
										MarkdownDescription: "Gateway IP.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"dns1": schema.StringAttribute{
										MarkdownDescription: "Primary DNS IP.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"dns2": schema.StringAttribute{
										MarkdownDescription: "Secondary DNS IP.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"signal_type": schema.StringAttribute{
										MarkdownDescription: "Signal Type.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"iccid": schema.StringAttribute{
										MarkdownDescription: "Integrated Circuit Card Identification SsidNumber.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
								},
							},
						},
					}}},
		},
	}
}
