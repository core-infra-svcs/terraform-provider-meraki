package _switch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
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
var _ resource.Resource = &NetworksSwitchMtuResource{}
var _ resource.ResourceWithImportState = &NetworksSwitchMtuResource{}

func NewNetworksSwitchMtuResource() resource.Resource {
	return &NetworksSwitchMtuResource{}
}

// NetworksSwitchMtuResource defines the resource implementation.
type NetworksSwitchMtuResource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchMtuResourceModel describes the resource data model.
type NetworksSwitchMtuResourceModel struct {
	Id             jsontypes.String                         `tfsdk:"id"`
	NetworkId      jsontypes.String                         `tfsdk:"network_id" json:"network_id"`
	DefaultMtuSize jsontypes.Int64                          `tfsdk:"default_mtu_size" json:"defaultMtuSize"`
	Overrides      []NetworksSwitchMtuResourceModelOverride `tfsdk:"overrides" json:"overrides"`
}

type NetworksSwitchMtuResourceModelOverride struct {
	Switches       []string        `tfsdk:"switches" json:"switches"`
	SwitchProfiles []string        `tfsdk:"switch_profiles" json:"switchProfiles"`
	MtuSize        jsontypes.Int64 `tfsdk:"mtu_size" json:"mtuSize"`
}

func (r *NetworksSwitchMtuResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_mtu"
}

func (r *NetworksSwitchMtuResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Networks switch mtu resource for updating networks switch mtu.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"default_mtu_size": schema.Int64Attribute{
				MarkdownDescription: "MTU size for the entire network. Default value is 9578.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"overrides": schema.SetNestedAttribute{
				Description: "Override MTU size for individual switches or switch profiles. An empty array will clear overrides.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"switches": schema.SetAttribute{
							MarkdownDescription: "Ports of switch serials. Applicable only for switch network.",
							CustomType:          jsontypes.SetType[jsontypes.String](),
							Optional:            true,
							Computed:            true,
						},
						"switch_profiles": schema.SetAttribute{
							MarkdownDescription: "Ports of switch profile IDs. Applicable only for template network.",
							CustomType:          jsontypes.SetType[jsontypes.String](),
							Optional:            true,
							Computed:            true,
						},
						"mtu_size": schema.Int64Attribute{
							MarkdownDescription: "MTU size for the switches or switch profiles..",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (r *NetworksSwitchMtuResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksSwitchMtuResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSwitchMtuResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSwitchMtu := *openApiClient.NewUpdateNetworkSwitchMtuRequest()
	if !data.DefaultMtuSize.IsUnknown() {
		updateNetworkSwitchMtu.SetDefaultMtuSize(int32(data.DefaultMtuSize.ValueInt64()))
	}
	var overrides []openApiClient.GetNetworkSwitchMtu200ResponseOverridesInner
	if len(data.Overrides) > 0 {
		for _, attribute := range data.Overrides {
			var override openApiClient.GetNetworkSwitchMtu200ResponseOverridesInner
			if !attribute.MtuSize.IsUnknown() {
				override.SetMtuSize(int32(attribute.MtuSize.ValueInt64()))
			}
			if len(attribute.Switches) > 0 {
				override.SetSwitches(attribute.Switches)
			}
			if len(attribute.SwitchProfiles) > 0 {
				override.SetSwitchProfiles(attribute.SwitchProfiles)
			}
			overrides = append(overrides, override)
			updateNetworkSwitchMtu.SetOverrides(overrides)
		}
	}

	_, httpResp, err := r.client.MtuApi.UpdateNetworkSwitchMtu(ctx, data.NetworkId.ValueString()).UpdateNetworkSwitchMtuRequest(updateNetworkSwitchMtu).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksSwitchMtuResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSwitchMtuResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.MtuApi.GetNetworkSwitchMtu(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
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
	} else {
		resp.Diagnostics.Append()
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksSwitchMtuResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksSwitchMtuResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSwitchMtu := *openApiClient.NewUpdateNetworkSwitchMtuRequest()
	if !data.DefaultMtuSize.IsUnknown() {
		updateNetworkSwitchMtu.SetDefaultMtuSize(int32(data.DefaultMtuSize.ValueInt64()))
	}
	var overrides []openApiClient.GetNetworkSwitchMtu200ResponseOverridesInner
	if len(data.Overrides) > 0 {
		for _, attribute := range data.Overrides {
			var override openApiClient.GetNetworkSwitchMtu200ResponseOverridesInner
			if !attribute.MtuSize.IsUnknown() {
				override.SetMtuSize(int32(attribute.MtuSize.ValueInt64()))
			}
			if len(attribute.Switches) > 0 {
				override.SetSwitches(attribute.Switches)
			}
			if len(attribute.SwitchProfiles) > 0 {
				override.SetSwitchProfiles(attribute.SwitchProfiles)
			}
			overrides = append(overrides, override)
			updateNetworkSwitchMtu.SetOverrides(overrides)
		}
	}

	_, httpResp, err := r.client.MtuApi.UpdateNetworkSwitchMtu(ctx, data.NetworkId.ValueString()).UpdateNetworkSwitchMtuRequest(updateNetworkSwitchMtu).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksSwitchMtuResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksSwitchMtuResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSwitchMtu := *openApiClient.NewUpdateNetworkSwitchMtuRequest()
	if !data.DefaultMtuSize.IsUnknown() {
		updateNetworkSwitchMtu.SetDefaultMtuSize(int32(data.DefaultMtuSize.ValueInt64()))
	}
	var overrides []openApiClient.GetNetworkSwitchMtu200ResponseOverridesInner
	if len(data.Overrides) > 0 {
		for _, attribute := range data.Overrides {
			var override openApiClient.GetNetworkSwitchMtu200ResponseOverridesInner
			if !attribute.MtuSize.IsUnknown() {
				override.SetMtuSize(int32(attribute.MtuSize.ValueInt64()))
			}
			if len(attribute.Switches) > 0 {
				override.SetSwitches(attribute.Switches)
			}
			if len(attribute.SwitchProfiles) > 0 {
				override.SetSwitchProfiles(attribute.SwitchProfiles)
			}
			overrides = append(overrides, override)
			updateNetworkSwitchMtu.SetOverrides(overrides)
		}
	}

	_, httpResp, err := r.client.MtuApi.UpdateNetworkSwitchMtu(ctx, data.NetworkId.ValueString()).UpdateNetworkSwitchMtuRequest(updateNetworkSwitchMtu).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksSwitchMtuResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
