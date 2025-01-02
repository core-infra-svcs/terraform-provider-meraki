package admins

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *openApiClient.APIClient
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_admins"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Ports the dashboard administrators in this organization",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				CustomType:          types.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"list": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Admin ID",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the dashboard administrator",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email of the dashboard administrator. This attribute can not be updated.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"org_access": schema.StringAttribute{
							MarkdownDescription: "The privilege of the dashboard administrator on the organization. Can be one of 'full', 'read-only', 'enterprise' or 'none'",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"account_status": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"two_factor_auth_enabled": schema.BoolAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.BoolType,
						},
						"has_api_key": schema.BoolAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.BoolType,
						},
						"last_active": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"tags": schema.SetNestedAttribute{
							Description: "The list of tags that the dashboard administrator has privileges on",
							Optional:    true,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"tag": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"access": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
						},
						"networks": schema.SetNestedAttribute{
							Description: "The list of networks that the dashboard administrator has privileges on",
							Optional:    true,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"access": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
						},
						"authentication_method": schema.StringAttribute{
							MarkdownDescription: "The method of authentication the user will use to sign in to the Meraki dashboard. Can be one of 'Email' or 'Cisco SecureX Sign-On'. The default is Email authentication",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceModel
	var diags diag.Diagnostics

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := d.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(d.client.GetConfig().Retry4xxErrorWaitTime)

	// usage of CustomHttpRequestRetry with a slice of strongly typed structs
	apiCallSlice := func() ([]openApiClient.GetOrganizationAdmins200ResponseInner, *http.Response, error, diag.Diagnostics) {
		inline, httpResp, err := d.client.AdminsApi.GetOrganizationAdmins(context.Background(), data.OrganizationId.ValueString()).Execute()
		return inline, httpResp, err, diags
	}

	resultSlice, httpResp, errSlice, tfDiag := utils.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCallSlice)
	if errSlice != nil {

		if tfDiag.HasError() {
			fmt.Printf(": %s\n", tfDiag.Errors())
		}
		fmt.Printf("Error creating group policy: %s\n", errSlice)
		if httpResp != nil {
			var responseBody string
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			fmt.Printf("Failed to create resource. HTTP Status Code: %d, Response Body: %s\n", httpResp.StatusCode, responseBody)
		}
		return
	}

	// Type assert apiResp to the expected []GetOrganizationAdmins200ResponseInner type
	_, ok := any(resultSlice).([]openApiClient.GetOrganizationAdmins200ResponseInner)
	if !ok {
		fmt.Println("Failed to assert API response type to []GetOrganizationAdmins200ResponseInner. Please ensure the API response structure matches the expected type.")
		return
	}
	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err := json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = types.StringValue(data.OrganizationId.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
