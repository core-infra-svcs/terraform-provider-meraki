package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworkResource{}
	_ resource.ResourceWithConfigure   = &NetworkResource{}
	_ resource.ResourceWithImportState = &NetworkResource{}
)

func NewNetworkResource() resource.Resource {
	return &NetworkResource{}
}

// NetworkResource defines the resource implementation.
type NetworkResource struct {
	client *openApiClient.APIClient
}

type Tag string

func (t *Tag) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*t = Tag(strings.Trim(s, `"`))
	return nil
}

// NetworkResourceModel describes the resource data model.
type NetworkResourceModel struct {
	Id                      types.String                    `tfsdk:"id"`
	NetworkId               jsontypes.String                `tfsdk:"network_id" json:"id"`
	OrganizationId          jsontypes.String                `tfsdk:"organization_id" json:"organizationId"`
	Name                    jsontypes.String                `tfsdk:"name"`
	ProductTypes            jsontypes.Set[jsontypes.String] `tfsdk:"product_types" json:"productTypes"`
	Timezone                jsontypes.String                `tfsdk:"timezone" json:"timeZone"`
	Tags                    []Tag                           `tfsdk:"tags"`
	EnrollmentString        jsontypes.String                `tfsdk:"enrollment_string" json:"enrollmentString"`
	Url                     jsontypes.String                `tfsdk:"url"`
	Notes                   jsontypes.String                `tfsdk:"notes"`
	IsBoundToConfigTemplate jsontypes.Bool                  `tfsdk:"is_bound_to_config_template" json:"IsBoundToConfigTemplate"`
	ConfigTemplateId        jsontypes.String                `tfsdk:"config_template_id" json:"configTemplateId"`
	CopyFromNetworkId       jsontypes.String                `tfsdk:"copy_from_network_id" json:"copyFromNetworkId"`
	AutoBind                types.Bool                      `tfsdk:"auto_bind" json:"autoBind"`
}

func (r *NetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network"
}

func (r *NetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage the networks that the user has privileges on in an organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
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
			"name": schema.StringAttribute{
				MarkdownDescription: "Network name",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"product_types": schema.SetAttribute{
				//ElementType: types.StringType,
				CustomType: jsontypes.SetType[jsontypes.String](),
				Required:   true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf([]string{"appliance", "switch", "wireless", "systemsManager", "camera", "cellularGateway", "sensor", "cloudGateway"}...), //
						stringvalidator.LengthAtLeast(5),
					),
				},
			},
			"timezone": schema.StringAttribute{
				MarkdownDescription: "Timezone of the network",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.SetAttribute{
				Description: "Network tags",
				ElementType: jsontypes.StringType,
				CustomType:  jsontypes.SetType[jsontypes.String](),
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"enrollment_string": schema.StringAttribute{
				MarkdownDescription: "A unique identifier which can be used for device enrollment or easy access through the Meraki SM Registration page or the Self Service Portal. Once enabled, a network enrollment strings can be changed but they cannot be deleted.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.NoneOf([]string{";", ":", "@", "=", "&", "$", "!", "‘", "“", ",", "?", ".", "(", ")", "{", "}", "[", "]", "\\", "*", "+", "/", "#", "<", ">", "|", "^", "%"}...),
					stringvalidator.LengthBetween(4, 50),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL to the network Dashboard UI",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Notes for the network",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_bound_to_config_template": schema.BoolAttribute{
				MarkdownDescription: "If the network is bound to a config template",
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"config_template_id": schema.StringAttribute{
				MarkdownDescription: "Config Template Id",
				Optional:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"copy_from_network_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the network to copy configuration from. Other provided parameters will override the copied configuration, except type which must match this network's type exactly.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_bind": schema.BoolAttribute{
				MarkdownDescription: "Optional boolean indicating whether the network's switches should automatically bind to profiles of the same model. Defaults to false if left unspecified. This option only affects switch networks and switch templates. Auto-bind is not valid unless the switch template has at least one profile and has at most one profile per switch model.",
				Optional:            true,
			},
		},
	}
}

func (r *NetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworkResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	createOrganizationNetwork := openApiClient.NewCreateOrganizationNetworkRequest(data.Name.ValueString(), nil)
	if !data.Timezone.IsUnknown() {
		createOrganizationNetwork.SetTimeZone(data.Timezone.ValueString())
	}

	// ProductTypes
	var productTypes []string
	for _, product := range data.ProductTypes.Elements() {
		pt := fmt.Sprint(strings.Trim(product.String(), "\""))
		productTypes = append(productTypes, pt)
	}
	createOrganizationNetwork.SetProductTypes(productTypes)

	if resp.Diagnostics.HasError() {
		return
	}

	createOrganizationNetwork.SetProductTypes(productTypes)

	// Tags
	if len(data.Tags) > 0 {
		var tags []string
		for _, attribute := range data.Tags {
			tags = append(tags, string(attribute))
		}
		createOrganizationNetwork.SetTags(tags)
	}

	// Notes
	if !data.Notes.IsUnknown() {
		createOrganizationNetwork.SetNotes(data.Notes.ValueString())
	}

	// CopyFromNetworkId
	if !data.CopyFromNetworkId.IsUnknown() {
		createOrganizationNetwork.SetCopyFromNetworkId(data.CopyFromNetworkId.ValueString())
	}

	// Initialize provider client and make API call
	_, httpResp, err := r.client.OrganizationsApi.CreateOrganizationNetwork(ctx, data.OrganizationId.ValueString()).CreateOrganizationNetworkRequest(*createOrganizationNetwork).Execute()
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
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// save inlineResp data into Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	// Note: Id is properly returned here, shouldnt be example
	data.Id = data.NetworkId.StringValue

	if data.CopyFromNetworkId.IsUnknown() {
		data.CopyFromNetworkId = jsontypes.StringNull()
	}

	// Validate Whether we need to bind a template
	if !data.ConfigTemplateId.IsUnknown() {
		createNetworkBindRequest := openApiClient.NewBindNetworkRequest(data.ConfigTemplateId.ValueString())

		if !data.AutoBind.IsUnknown() {
			createNetworkBindRequest.SetAutoBind(data.AutoBind.ValueBool())
		}

		_, bindHttpResp, bindErr := r.client.NetworksApi.BindNetwork(ctx, data.NetworkId.ValueString()).BindNetworkRequest(*createNetworkBindRequest).Execute()

		if bindErr != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				tools.HttpDiagnostics(bindHttpResp),
			)
			return
		}

		// Check for API success response code
		if bindHttpResp.StatusCode != 200 {
			resp.Diagnostics.AddError(
				"Unexpected HTTP Response Status Code",
				fmt.Sprintf("%v", bindHttpResp.StatusCode),
			)
			return
		}
		var bindRes *NetworkResourceModel

		req.Plan.Get(ctx, &bindRes)

		if err = json.NewDecoder(bindHttpResp.Body).Decode(bindRes); err != nil {
			resp.Diagnostics.AddError(
				"JSON Decode issue",
				fmt.Sprintf("%v", bindHttpResp.StatusCode),
			)

			return
		}

		data.ConfigTemplateId = bindRes.ConfigTemplateId
		data.IsBoundToConfigTemplate = bindRes.IsBoundToConfigTemplate
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *NetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworkResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Initialize provider client and make API call
	_, httpResp, err := r.client.NetworksApi.GetNetwork(context.Background(), data.NetworkId.ValueString()).Execute()
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
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {

		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		return
	}
	// save inlineResp data into Terraform state.
	data.Id = data.NetworkId.StringValue

	if data.CopyFromNetworkId.IsUnknown() {
		data.CopyFromNetworkId = jsontypes.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state, plan *NetworkResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	updateNetwork := openApiClient.NewUpdateNetworkRequest()
	updateNetwork.SetName(data.Name.ValueString())
	updateNetwork.SetTimeZone(data.Timezone.ValueString())

	// Tags
	if len(data.Tags) > 0 {
		var tags []string
		for _, attribute := range data.Tags {
			tags = append(tags, string(attribute))
		}
		updateNetwork.SetTags(tags)
	}

	// Enrollment String
	if !data.EnrollmentString.IsUnknown() {
		updateNetwork.SetEnrollmentString(data.EnrollmentString.ValueString())
	}

	// Notes
	updateNetwork.SetNotes(data.Notes.ValueString())

	// Initialize provider client and make API call
	_, httpResp, err := r.client.NetworksApi.UpdateNetwork(context.Background(),
		data.NetworkId.ValueString()).UpdateNetworkRequest(*updateNetwork).Execute()
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
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// save inlineResp data into Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	// If a new template is being added, unbind
	if (plan.ConfigTemplateId.IsUnknown() && len(state.ConfigTemplateId.ValueString()) > 0) ||
		state.ConfigTemplateId.StringValue != plan.ConfigTemplateId.StringValue && len(state.ConfigTemplateId.ValueString()) > 0 {

		_, httpResp, err := r.client.NetworksApi.UnbindNetwork(context.Background(), plan.NetworkId.ValueString()).Execute()

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
			return
		}
		var bindRes *NetworkResourceModel

		req.Plan.Get(ctx, &bindRes)

		if err = json.NewDecoder(httpResp.Body).Decode(bindRes); err != nil {
			resp.Diagnostics.AddError(
				"JSON Decode issue",
				fmt.Sprintf("%v", httpResp.StatusCode),
			)

			return
		}

		data = bindRes
	}

	// If new template, or swappng, bind
	if (state.ConfigTemplateId.IsUnknown() && len(plan.ConfigTemplateId.ValueString()) > 0) ||
		state.ConfigTemplateId.StringValue != plan.ConfigTemplateId.StringValue && len(plan.ConfigTemplateId.ValueString()) > 0 {
		tflog.Debug(ctx, "binding")

		createNetworkBindRequest := openApiClient.NewBindNetworkRequest(plan.ConfigTemplateId.ValueString())

		if !plan.AutoBind.IsUnknown() {
			createNetworkBindRequest.SetAutoBind(plan.AutoBind.ValueBool())
		}

		_, bindHttpResp, bindErr := r.client.NetworksApi.BindNetwork(ctx, plan.NetworkId.ValueString()).BindNetworkRequest(*createNetworkBindRequest).Execute()

		if bindErr != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				tools.HttpDiagnostics(bindHttpResp),
			)
			return
		}

		// Check for API success response code
		if bindHttpResp.StatusCode != 200 {
			resp.Diagnostics.AddError(
				"Unexpected HTTP Response Status Code",
				fmt.Sprintf("%v", bindHttpResp.StatusCode),
			)
			return
		}
		var bindRes *NetworkResourceModel

		req.Plan.Get(ctx, &bindRes)

		if err = json.NewDecoder(bindHttpResp.Body).Decode(bindRes); err != nil {
			resp.Diagnostics.AddError(
				"JSON Decode issue",
				fmt.Sprintf("%v", bindHttpResp.StatusCode),
			)

			return
		}

		data = bindRes
	}

	data.Id = data.NetworkId.StringValue

	if data.CopyFromNetworkId.IsUnknown() {
		data.CopyFromNetworkId = jsontypes.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworkResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// HTTP DELETE METHOD does not leverage the retry-after header and throws 400 errors.
	retries := 100 // This is the only way to ensure a scaled result.
	wait := 1
	var deletedFromMerakiPortal bool
	deletedFromMerakiPortal = false

	for retries > 0 {

		// Initialize provider client and make API call
		httpResp, err := r.client.NetworksApi.DeleteNetwork(context.Background(), data.NetworkId.ValueString()).Execute()

		if httpResp.StatusCode == 204 {
			// check for HTTP errors
			if err != nil {
				resp.Diagnostics.AddError(
					"HTTP Client Failure",
					tools.HttpDiagnostics(httpResp),
				)
				return
			}

			// Check for errors after diagnostics collected
			if resp.Diagnostics.HasError() {
				return
			}

			deletedFromMerakiPortal = true

			// escape loop
			break

		} else {

			// decrement retry counter
			retries -= 1

			// exponential wait
			time.Sleep(time.Duration(wait) * time.Second)
			wait += 1
		}
	}

	// Check if deleted from Meraki Portal was successful
	if deletedFromMerakiPortal {
		resp.State.RemoveResource(ctx)

		// Write logs using the tflog package
		tflog.Trace(ctx, "removed resource")
	} else {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			"Failed to delete resource",
		)
		return
	}

}

func (r *NetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, network_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
