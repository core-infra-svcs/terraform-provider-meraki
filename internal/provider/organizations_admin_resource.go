package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &OrganizationsAdminResource{}
	_ resource.ResourceWithConfigure   = &OrganizationsAdminResource{}
	_ resource.ResourceWithImportState = &OrganizationsAdminResource{}
)

func NewOrganizationsAdminResource() resource.Resource {
	return &OrganizationsAdminResource{}
}

// OrganizationsAdminResource defines the resource implementation.
type OrganizationsAdminResource struct {
	client *openApiClient.APIClient
}

// OrganizationsAdminResourceModel describes the resource data model.
type OrganizationsAdminResourceModel struct {
	Id                   jsontypes.String                         `tfsdk:"id"`
	OrgId                jsontypes.String                         `tfsdk:"organization_id" json:"organizationId"`
	AdminId              jsontypes.String                         `tfsdk:"admin_id" json:"id"`
	Name                 jsontypes.String                         `tfsdk:"name"`
	Email                jsontypes.String                         `tfsdk:"email"`
	OrgAccess            jsontypes.String                         `tfsdk:"org_access" json:"orgAccess"`
	AccountStatus        jsontypes.String                         `tfsdk:"account_status" json:"accountStatus"`
	TwoFactorAuthEnabled jsontypes.Bool                           `tfsdk:"two_factor_auth_enabled" json:"twoFactorAuthEnabled"`
	HasApiKey            jsontypes.Bool                           `tfsdk:"has_api_key" json:"hasApiKey"`
	LastActive           jsontypes.String                         `tfsdk:"last_active" json:"lastActive"`
	Tags                 []OrganizationsAdminResourceModelTag     `tfsdk:"tags" json:"tags"`
	Networks             []OrganizationsAdminResourceModelNetwork `tfsdk:"networks" json:"networks"`
	AuthenticationMethod jsontypes.String                         `tfsdk:"authentication_method" json:"authenticationMethod"`
}

type OrganizationsAdminResourceModelTag struct {
	Tag    jsontypes.String `tfsdk:"tag" json:"tag"`
	Access jsontypes.String `tfsdk:"access" json:"access"`
}

type OrganizationsAdminResourceModelNetwork struct {
	Id     jsontypes.String `tfsdk:"id" json:"id"`
	Access jsontypes.String `tfsdk:"access" json:"access"`
}

func (r *OrganizationsAdminResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_admin"
}

func (r *OrganizationsAdminResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage the dashboard administrators in this organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				CustomType: jsontypes.StringType,
				Computed:   true,
				Optional:   true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"admin_id": schema.StringAttribute{
				MarkdownDescription: "Admin ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the dashboard administrator",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the dashboard administrator. This attribute can not be updated.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"org_access": schema.StringAttribute{
				MarkdownDescription: "The privilege of the dashboard administrator on the organization. Can be one of 'full', 'read-only', 'enterprise' or 'none'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"full", "read-only", "enterprise", "none"}...),
					stringvalidator.LengthAtLeast(4),
				},
			},
			"account_status": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"two_factor_auth_enabled": schema.BoolAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"has_api_key": schema.BoolAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"last_active": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
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
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"access": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
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
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"access": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
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
				Validators: []validator.String{

					stringvalidator.OneOf([]string{"Email", "Cisco SecureX Sign-On"}...),
					stringvalidator.LengthAtLeast(5),
				},
			},
		},
	}
}

func (r *OrganizationsAdminResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *OrganizationsAdminResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsAdminResourceModel
	var diags diag.Diagnostics

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Creating and Validating Payload for Creating Administrator
	createOrganizationAdmin := *openApiClient.NewCreateOrganizationAdminRequest(
		data.Email.ValueString(),
		data.Name.ValueString(),
		data.OrgAccess.ValueString())

	// Tags
	if len(data.Tags) > 0 {
		var tags []openApiClient.CreateOrganizationAdminRequestTagsInner
		for _, attribute := range data.Tags {
			var tag openApiClient.CreateOrganizationAdminRequestTagsInner
			tag.Tag = attribute.Tag.ValueString()
			tag.Access = attribute.Access.ValueString()
			tags = append(tags, tag)
		}
		createOrganizationAdmin.SetTags(tags)
	}

	// Networks
	if len(data.Networks) > 0 {
		var networks []openApiClient.CreateOrganizationAdminRequestNetworksInner
		for _, attribute := range data.Networks {
			var network openApiClient.CreateOrganizationAdminRequestNetworksInner
			network.Id = attribute.Id.ValueString()
			network.Access = attribute.Access.ValueString()
			networks = append(networks, network)
		}
		createOrganizationAdmin.SetNetworks(networks)
	}

	if !data.AuthenticationMethod.IsNull() {
		createOrganizationAdmin.SetAuthenticationMethod(data.AuthenticationMethod.ValueString())
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetOrganizationAdmins200ResponseInner, *http.Response, error, diag.Diagnostics) {
		inline, httpResp, err := r.client.AdminsApi.CreateOrganizationAdmin(context.Background(), data.OrgId.ValueString()).CreateOrganizationAdminRequest(createOrganizationAdmin).Execute()
		return inline, httpResp, err, diags
	}

	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	_, httpResp, err, tfDiags := tools.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group policy",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating admin resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	if tfDiags.HasError() {

	}
	// Check for API success response code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue(fmt.Sprintf("%s,%s", data.OrgId.ValueString(), data.AdminId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *OrganizationsAdminResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsAdminResourceModel
	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// Usage of CustomHttpRequestRetry with a slice of strongly typed structs
	apiCallSlice := func() ([]openApiClient.GetOrganizationAdmins200ResponseInner, *http.Response, error) {
		inline, httpResp, err := r.client.AdminsApi.GetOrganizationAdmins(context.Background(), data.OrgId.ValueString()).Execute()
		return inline, httpResp, err
	}

	// Directly use the type returned by the function
	resultSlice, httpRespSlice, errSlice := tools.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCallSlice)
	if errSlice != nil {
		resp.Diagnostics.AddError(
			"Error reading admins",
			fmt.Sprintf("Error retrieving admins from the API: %s", errSlice),
		)
		if httpRespSlice != nil {
			var responseBody string
			if httpRespSlice.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpRespSlice.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			resp.Diagnostics.AddError(
				"HTTP Response",
				fmt.Sprintf("Failed to read admins. HTTP Status Code: %d, Response Body: %s", httpRespSlice.StatusCode, responseBody),
			)
		}
		return
	}

	// Check for API success inlineResp code
	if httpRespSlice.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpRespSlice.StatusCode),
		)
		return
	}

	// Iterate through the resultSlice directly since it is already of the expected type
	var foundAdmin *OrganizationsAdminResourceModel
	for _, admin := range resultSlice {
		if admin.GetId() == data.AdminId.ValueString() {

			// tags
			var tags []OrganizationsAdminResourceModelTag
			for _, t := range admin.GetTags() {
				var tag OrganizationsAdminResourceModelTag
				tag.Tag = jsontypes.StringValue(t.GetTag())
				tag.Access = jsontypes.StringValue(t.GetAccess())

				tags = append(tags, tag)
			}

			// networks
			var networks []OrganizationsAdminResourceModelNetwork
			for _, n := range admin.GetNetworks() {
				var network OrganizationsAdminResourceModelNetwork
				network.Id = jsontypes.StringValue(n.GetId())
				network.Access = jsontypes.StringValue(n.GetAccess())
				networks = append(networks, network)
			}

			foundAdmin = &OrganizationsAdminResourceModel{
				Id:                   jsontypes.StringValue(fmt.Sprintf("%s,%s", data.OrgId.ValueString(), data.AdminId.ValueString())),
				OrgId:                data.OrgId,
				AdminId:              jsontypes.StringValue(admin.GetId()),
				Name:                 jsontypes.StringValue(admin.GetName()),
				Email:                jsontypes.StringValue(admin.GetEmail()),
				OrgAccess:            jsontypes.StringValue(admin.GetOrgAccess()),
				AccountStatus:        jsontypes.StringValue(admin.GetAccountStatus()),
				TwoFactorAuthEnabled: jsontypes.BoolValue(admin.GetTwoFactorAuthEnabled()),
				HasApiKey:            jsontypes.BoolValue(admin.GetHasApiKey()),
				LastActive:           jsontypes.StringValue(admin.GetLastActive().String()),
				Tags:                 tags,
				Networks:             networks,
				AuthenticationMethod: jsontypes.StringValue(admin.GetAuthenticationMethod()),
			}

			break
		}
	}

	if foundAdmin == nil {
		resp.Diagnostics.AddError(
			"Admin not found",
			fmt.Sprintf("No admin found with ID: %s", data.AdminId.ValueString()),
		)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, foundAdmin)...)
}

func (r *OrganizationsAdminResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsAdminResourceModel
	var diags diag.Diagnostics
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Creating and Validating Payload for Creating Administrator
	updateOrganizationAdmin := *openApiClient.NewUpdateOrganizationAdminRequest()
	updateOrganizationAdmin.SetName(data.Name.ValueString())
	updateOrganizationAdmin.SetOrgAccess(data.OrgAccess.ValueString())

	// Tags
	if len(data.Tags) > 0 {
		var tags []openApiClient.CreateOrganizationAdminRequestTagsInner
		for _, attribute := range data.Tags {
			var tag openApiClient.CreateOrganizationAdminRequestTagsInner
			tag.Tag = attribute.Tag.ValueString()
			tag.Access = attribute.Access.ValueString()
			tags = append(tags, tag)
		}
		updateOrganizationAdmin.SetTags(tags)
	}

	// Networks
	if len(data.Networks) > 0 {
		var networks []openApiClient.CreateOrganizationAdminRequestNetworksInner
		for _, attribute := range data.Networks {
			var network openApiClient.CreateOrganizationAdminRequestNetworksInner
			network.Id = attribute.Id.ValueString()
			network.Access = attribute.Access.ValueString()
			networks = append(networks, network)
		}
		updateOrganizationAdmin.SetNetworks(networks)
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetOrganizationAdmins200ResponseInner, *http.Response, error, diag.Diagnostics) {
		inline, httpResp, err := r.client.AdminsApi.UpdateOrganizationAdmin(context.Background(), data.OrgId.ValueString(), data.AdminId.ValueString()).UpdateOrganizationAdminRequest(updateOrganizationAdmin).Execute()
		return inline, httpResp, err, diags
	}

	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	_, httpResp, err, tfDiags := tools.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating group policy",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if tfDiags.HasError() {
		}

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to update resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error updating admin resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue(fmt.Sprintf("%s,%s", data.OrgId.ValueString(), data.AdminId.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsAdminResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsAdminResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (map[string]interface{}, *http.Response, error) {
		httpResp, err := r.client.AdminsApi.DeleteOrganizationAdmin(context.Background(), data.OrgId.ValueString(), data.AdminId.ValueString()).Execute()

		return nil, httpResp, err
	}

	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	_, httpResp, err := tools.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting admin",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		// Check for API success response code
		if httpResp.StatusCode != 204 {
			resp.Diagnostics.AddError(
				"Unexpected HTTP Response Status Code",
				fmt.Sprintf("%v", httpResp.StatusCode),
			)
		}

		// Check for errors after diagnostics collected
		if resp.Diagnostics.HasError() {
			return
		}

		resp.State.RemoveResource(ctx)

	}
}

func (r *OrganizationsAdminResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, admin_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("admin_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
