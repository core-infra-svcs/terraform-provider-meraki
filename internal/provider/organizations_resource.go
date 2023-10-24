package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &OrganizationResource{}
	_ resource.ResourceWithConfigure   = &OrganizationResource{}
	_ resource.ResourceWithImportState = &OrganizationResource{}
)

func NewOrganizationResource() resource.Resource {
	return &OrganizationResource{}
}

// OrganizationResource defines the resource implementation.
type OrganizationResource struct {
	client *openApiClient.APIClient
}

// OrganizationResourceModel describes the resource data model.
type OrganizationResourceModel struct {
	Id                     types.String     `tfsdk:"id"`
	ApiEnabled             jsontypes.Bool   `tfsdk:"api_enabled"`
	CloudRegionName        jsontypes.String `tfsdk:"cloud_region_name"`
	ManagementDetailsName  jsontypes.String `tfsdk:"management_details_name"`
	ManagementDetailsValue jsontypes.String `tfsdk:"management_details_value"`
	OrgId                  jsontypes.String `tfsdk:"organization_id"`
	LicensingModel         jsontypes.String `tfsdk:"licensing_model"`
	Name                   jsontypes.String `tfsdk:"name"`
	Url                    jsontypes.String `tfsdk:"url"`
	OrgToClone             jsontypes.String `tfsdk:"org_to_clone"`
}

func (r *OrganizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (r *OrganizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage the organizations that the user has privileges on",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"api_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enable API access",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"cloud_region_name": schema.StringAttribute{
				MarkdownDescription: "Name of region",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
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
			"licensing_model": schema.StringAttribute{
				MarkdownDescription: "Organization licensing model. Can be 'co-term', 'per-device', or 'subscription'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"co-term", "per-device", "subscription"}...),
				},
			},
			"management_details_name": schema.StringAttribute{
				MarkdownDescription: "Name of management data",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"management_details_value": schema.StringAttribute{
				MarkdownDescription: "Value of management data",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Organization name",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "Organization URL",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"org_to_clone": schema.StringAttribute{
				MarkdownDescription: "ID or org which to clone from",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
		},
	}
}

func (r *OrganizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var inlineResp *openApiClient.GetOrganizations200ResponseInner
	var httpResp *http.Response
	var err error

	if !data.OrgToClone.IsNull() {
		cloneOrganization := *openApiClient.NewCloneOrganizationRequest(data.Name.ValueString())

		// Initialize provider client and make API call
		inlineResp, httpResp, err = r.client.OrganizationsApi.CloneOrganization(context.Background(), data.OrgToClone.ValueString()).CloneOrganizationRequest(cloneOrganization).Execute()
	} else {
		// Create HTTP request body
		createOrganization := *openApiClient.NewCreateOrganizationRequest(data.Name.ValueString())

		// Set management details
		var name = data.ManagementDetailsName.ValueString()
		var value = data.ManagementDetailsValue.ValueString()
		detail := openApiClient.GetOrganizations200ResponseInnerManagementDetailsInner{
			Name:  &name,
			Value: &value,
		}
		details := []openApiClient.GetOrganizations200ResponseInnerManagementDetailsInner{detail}
		organizationsManagement := openApiClient.GetOrganizations200ResponseInnerManagement{Details: details}
		createOrganization.SetManagement(organizationsManagement)

		// Initialize provider client and make API call
		inlineResp, httpResp, err = r.client.OrganizationsApi.CreateOrganization(context.Background()).CreateOrganizationRequest(createOrganization).Execute()
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}
	// Check for API success response code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")
	data.OrgId = jsontypes.StringValue(inlineResp.GetId())
	data.Name = jsontypes.StringValue(inlineResp.GetName())
	data.CloudRegionName = jsontypes.StringValue(inlineResp.Cloud.Region.GetName())
	data.Url = jsontypes.StringValue(inlineResp.GetUrl())
	data.ApiEnabled = jsontypes.BoolValue(inlineResp.Api.GetEnabled())
	data.LicensingModel = jsontypes.StringValue(inlineResp.Licensing.GetModel())

	// Management Details Response
	if len(inlineResp.Management.Details) > 0 {
		responseDetails := inlineResp.Management.GetDetails()

		// name attribute
		if managementDetailName := responseDetails[0].GetName(); responseDetails[0].HasName() {
			data.ManagementDetailsName = jsontypes.StringValue(managementDetailName)
		} else {
			data.ManagementDetailsName = jsontypes.StringNull()
		}

		// Value attribute
		if managementDetailValue := responseDetails[0].GetValue(); responseDetails[0].HasValue() {
			data.ManagementDetailsValue = jsontypes.StringValue(managementDetailValue)
		} else {
			data.ManagementDetailsValue = jsontypes.StringNull()
		}

	} else {
		data.ManagementDetailsName = jsontypes.StringNull()
		data.ManagementDetailsValue = jsontypes.StringNull()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *OrganizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.OrganizationsApi.GetOrganization(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
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
		return
	}

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")
	data.OrgId = jsontypes.StringValue(inlineResp.GetId())
	data.Name = jsontypes.StringValue(inlineResp.GetName())
	data.CloudRegionName = jsontypes.StringValue(inlineResp.Cloud.Region.GetName())
	data.Url = jsontypes.StringValue(inlineResp.GetUrl())
	data.ApiEnabled = jsontypes.BoolValue(inlineResp.Api.GetEnabled())
	data.LicensingModel = jsontypes.StringValue(inlineResp.Licensing.GetModel())

	// Management Details Response
	if len(inlineResp.Management.Details) > 0 {
		responseDetails := inlineResp.Management.GetDetails()

		// name attribute
		if managementDetailName := responseDetails[0].GetName(); responseDetails[0].HasName() {
			data.ManagementDetailsName = jsontypes.StringValue(managementDetailName)
		} else {
			data.ManagementDetailsName = jsontypes.StringNull()
		}

		// Value attribute
		if managementDetailValue := responseDetails[0].GetValue(); responseDetails[0].HasValue() {
			data.ManagementDetailsValue = jsontypes.StringValue(managementDetailValue)
		} else {
			data.ManagementDetailsValue = jsontypes.StringNull()
		}
	} else {
		data.ManagementDetailsName = jsontypes.StringNull()
		data.ManagementDetailsValue = jsontypes.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Create HTTP request body
	updateOrganization := openApiClient.NewUpdateOrganizationRequest()
	updateOrganization.SetName(data.Name.ValueString())

	// Set enabled attribute
	var enabled = data.ApiEnabled.ValueBool()
	Api := openApiClient.UpdateOrganizationRequestApi{Enabled: &enabled}
	updateOrganization.SetApi(Api)

	// Set management details
	var name = data.ManagementDetailsName.ValueString()
	var value = data.ManagementDetailsValue.ValueString()
	detail := openApiClient.GetOrganizations200ResponseInnerManagementDetailsInner{
		Name:  &name,
		Value: &value,
	}
	details := []openApiClient.GetOrganizations200ResponseInnerManagementDetailsInner{detail}
	organizationsManagement := openApiClient.GetOrganizations200ResponseInnerManagement{Details: details}
	updateOrganization.SetManagement(organizationsManagement)

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.OrganizationsApi.UpdateOrganization(context.Background(),
		data.OrgId.ValueString()).UpdateOrganizationRequest(*updateOrganization).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// save inlineResp data into Terraform state
	data.Id = types.StringValue("example-id")
	data.OrgId = jsontypes.StringValue(inlineResp.GetId())
	data.Name = jsontypes.StringValue(inlineResp.GetName())
	data.CloudRegionName = jsontypes.StringValue(inlineResp.Cloud.Region.GetName())
	data.Url = jsontypes.StringValue(inlineResp.GetUrl())
	data.ApiEnabled = jsontypes.BoolValue(inlineResp.Api.GetEnabled())
	data.LicensingModel = jsontypes.StringValue(inlineResp.Licensing.GetModel())

	// Management Details Response
	if len(inlineResp.Management.Details) > 0 {
		responseDetails := inlineResp.Management.GetDetails()

		// name attribute
		if managementDetailName := responseDetails[0].GetName(); responseDetails[0].HasName() {
			data.ManagementDetailsName = jsontypes.StringValue(managementDetailName)
		} else {
			data.ManagementDetailsName = jsontypes.StringNull()
		}

		// Value attribute
		if managementDetailValue := responseDetails[0].GetValue(); responseDetails[0].HasValue() {
			data.ManagementDetailsValue = jsontypes.StringValue(managementDetailValue)
		} else {
			data.ManagementDetailsValue = jsontypes.StringNull()
		}
	} else {
		data.ManagementDetailsName = jsontypes.StringNull()
		data.ManagementDetailsValue = jsontypes.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Initialize provider client and make API call
	httpResp, err := r.client.OrganizationsApi.DeleteOrganization(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

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

	// Remove from state
	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *OrganizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), req.ID)...)
}
