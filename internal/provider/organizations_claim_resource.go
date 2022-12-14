package provider

import (
	"context"
	"encoding/json"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// TODO - DON'T FORGET TO DELETE ALL "TODO" COMMENTS!

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationsClaimResource{}
var _ resource.ResourceWithImportState = &OrganizationsClaimResource{}

func NewOrganizationsClaimResource() resource.Resource {
	return &OrganizationsClaimResource{}
}

// OrganizationsClaimResource defines the resource implementation.
type OrganizationsClaimResource struct {
	client *apiclient.APIClient
}

// OrganizationsClaimResourceModel describes the resource data model.
type OrganizationsClaimResourceModel struct {
	Id       types.String                        `tfsdk:"id"`
	OrgId    types.String                        `tfsdk:"organization_id"`
	Orders   []types.String                      `tfsdk:"orders"`
	Serials  []types.String                      `tfsdk:"serials"`
	Licences []OrganizationsClaimResourceLicence `tfsdk:"licences"`
}

type OrganizationsClaimResourceLicence struct {
	Key  types.String `tfsdk:"key"`
	Mode types.String `tfsdk:"mode"`
}

func (r *OrganizationsClaimResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_claim"
}

func (r *OrganizationsClaimResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "OrganizationsClaim resource - Claim a list of devices, licenses, " +
			"and/or orders into an organization. When claiming by order, " +
			"all devices and licenses in the order will be claimed; " +
			"licenses will be added to the organization and devices will be placed in the organization's inventory.",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "Example identifier",
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"organization_id": {
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"orders": {
				Description:         "The numbers of the orders that should be claimed",
				MarkdownDescription: "The numbers of the orders that should be claimed",
				Type:                types.ListType{ElemType: types.StringType},
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"serials": {
				Description:         "The serials of the devices that should be claimed",
				MarkdownDescription: "The serials of the devices that should be claimed",
				Type:                types.ListType{ElemType: types.StringType},
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"licences": {
				Description: "The keys & modes of the list of license",
				MarkdownDescription: "Either 'renew' or 'addDevices'. 'addDevices' will increase the license limit, " +
					"while 'renew' will extend the amount of time until expiration. Defaults to 'addDevices'. " +
					"All licenses must be claimed with the same mode, and at most one renewal can be claimed at a time. " +
					"This parameter is legacy and does not apply to organizations with per-device licensing enabled.",
				Type:               types.ListType{ElemType: types.StringType},
				Required:           false,
				Optional:           true,
				Computed:           true,
				Sensitive:          false,
				Attributes:         nil,
				DeprecationMessage: "",
				Validators:         nil,
				PlanModifiers:      nil,
			},
		},
	}, nil
}

func (r *OrganizationsClaimResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *OrganizationsClaimResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsClaimResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create and Validate Payload
	claimIntoOrganization := apiclient.NewInlineObject188()

	if len(data.Orders) < 0 {
		var orders []string
		for _, order := range data.Orders {
			orders = append(orders, order.ValueString())
		}
		claimIntoOrganization.Orders = orders
	}

	if len(data.Serials) < 0 {
		var serials []string
		for _, serial := range data.Serials {
			serials = append(serials, serial.ValueString())
		}
		claimIntoOrganization.Serials = serials
	}

	if len(data.Licences) < 0 {
		var licences []apiclient.OrganizationsOrganizationIdClaimLicenses
		for _, licence := range data.Licences {
			var lic apiclient.OrganizationsOrganizationIdClaimLicenses

			// key
			lic.Key = licence.Key.ValueString()

			// mode
			mode := licence.Mode.ValueString()
			lic.Mode = &mode

			licences = append(licences, lic)
		}
		claimIntoOrganization.Licenses = licences
	}

	inlineResp, httpResp, err := r.client.OrganizationsApi.ClaimIntoOrganization(context.Background(), data.OrgId.ValueString()).ClaimIntoOrganization(*claimIntoOrganization).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}
	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")

	// orders attribute
	if orders := inlineResp["orders"]; orders != nil {
		for _, order := range orders.([]interface{}) {
			data.Orders = append(data.Orders, types.StringValue(order.(string)))
		}
	} else {
		data.Orders = nil
	}

	// serials attribute
	if serials := inlineResp["serials"]; serials != nil {
		for _, serial := range serials.([]interface{}) {
			data.Serials = append(data.Serials, types.StringValue(serial.(string)))
		}
	} else {
		data.Serials = nil
	}

	// licences attribute
	if licences := inlineResp["licences"]; licences != nil {
		for _, lic := range licences.([]interface{}) {
			var licence OrganizationsClaimResourceLicence
			_ = json.Unmarshal([]byte(lic.(string)), &licence)
			data.Licences = append(data.Licences, licence)
		}
	} else {
		data.Licences = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *OrganizationsClaimResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsClaimResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationsClaimResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsClaimResourceModel
	var stateData *OrganizationsClaimResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	// Check state for required attributes
	if len(data.OrgId.ValueString()) < 1 {
		data.OrgId = stateData.OrgId
	}

	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("OrganizationId: %s", data.OrgId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create and Validate Payload
	claimIntoOrganization := apiclient.NewInlineObject188()

	if len(data.Orders) < 0 {
		var orders []string
		for _, order := range data.Orders {
			orders = append(orders, order.ValueString())
		}
		claimIntoOrganization.Orders = orders
	}

	if len(data.Serials) < 0 {
		var serials []string
		for _, serial := range data.Serials {
			serials = append(serials, serial.ValueString())
		}
		claimIntoOrganization.Serials = serials
	}

	if len(data.Licences) < 0 {
		var licences []apiclient.OrganizationsOrganizationIdClaimLicenses
		for _, licence := range data.Licences {
			var lic apiclient.OrganizationsOrganizationIdClaimLicenses

			// key
			lic.Key = licence.Key.ValueString()

			// mode
			mode := licence.Mode.ValueString()
			lic.Mode = &mode

			licences = append(licences, lic)
		}
		claimIntoOrganization.Licenses = licences
	}

	inlineResp, httpResp, err := r.client.OrganizationsApi.ClaimIntoOrganization(context.Background(), data.OrgId.ValueString()).ClaimIntoOrganization(*claimIntoOrganization).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}
	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// save into the Terraform state.
	data.Id = types.StringValue("example-id")

	// orders attribute
	if orders := inlineResp["orders"]; orders != nil {
		for _, order := range orders.([]interface{}) {
			data.Orders = append(data.Orders, types.StringValue(order.(string)))
		}
	} else {
		data.Orders = nil
	}

	// serials attribute
	if serials := inlineResp["serials"]; serials != nil {
		for _, serial := range serials.([]interface{}) {
			data.Serials = append(data.Serials, types.StringValue(serial.(string)))
		}
	} else {
		data.Serials = nil
	}

	// licences attribute
	if licences := inlineResp["licences"]; licences != nil {
		for _, lic := range licences.([]interface{}) {
			var licence OrganizationsClaimResourceLicence
			_ = json.Unmarshal([]byte(lic.(string)), &licence)
			data.Licences = append(data.Licences, licence)
		}
	} else {
		data.Licences = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationsClaimResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsClaimResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *OrganizationsClaimResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
}
