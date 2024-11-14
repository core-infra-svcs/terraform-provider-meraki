package firewall

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
var _ resource.Resource = &NetworksApplianceFirewallL7FirewallRulesResource{}
var _ resource.ResourceWithImportState = &NetworksApplianceFirewallL7FirewallRulesResource{}

func NewNetworksApplianceFirewallL7FirewallRulesResource() resource.Resource {
	return &NetworksApplianceFirewallL7FirewallRulesResource{}
}

// NetworksApplianceFirewallL7FirewallRulesResource defines the resource implementation.
type NetworksApplianceFirewallL7FirewallRulesResource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceFirewallL7FirewallRulesResourceModel describes the resource data model.
type NetworksApplianceFirewallL7FirewallRulesResourceModel struct {
	Id        jsontypes.String                                            `tfsdk:"id"`
	NetworkId jsontypes.String                                            `tfsdk:"network_id" json:"network_id"`
	Rules     []NetworksApplianceFirewallL7FirewallRulesResourceModelRule `tfsdk:"rules" json:"rules"`
}

type NetworksApplianceFirewallL7FirewallRulesResourceModelRule struct {
	Policy jsontypes.String `tfsdk:"policy"`
	Type   jsontypes.String `tfsdk:"type"`
	Value  jsontypes.String `tfsdk:"value"`
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_firewall_l7_firewall_rules"
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Network Appliance L7 Firewall Rules",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"rules": schema.SetNestedAttribute{
				MarkdownDescription: "An ordered array of the MX L7 firewall rules",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy": schema.StringAttribute{
							MarkdownDescription: "Deny' traffic specified by this rule",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the L7 rule. One of: 'application', 'applicationCategory', 'host', 'port', 'ipRange'",
							Required:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"application", "applicationCategory", "host", "port", "ipRange"}...),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The 'value' of what you want to block. Format of 'value' varies depending on type of the rule. The application categories and application ids can be retrieved from the the 'MX L7 application categories' endpoint. The countries follow the two-letter ISO 3166-1 alpha-2 format.",
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

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data *NetworksApplianceFirewallL7FirewallRulesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceFirewallL7FirewallRules := *openApiClient.NewUpdateNetworkApplianceFirewallL7FirewallRulesRequest()
	var rules []openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner

	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner
			rule.SetPolicy(attribute.Policy.ValueString())
			rule.SetType(attribute.Type.ValueString())
			rule.SetValue(attribute.Value.ValueString())
			rules = append(rules, rule)
		}
	}

	updateNetworkApplianceFirewallL7FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceFirewallL7FirewallRulesRequest(updateNetworkApplianceFirewallL7FirewallRules).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

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

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceFirewallL7FirewallRulesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

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

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksApplianceFirewallL7FirewallRulesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceFirewallL7FirewallRules := *openApiClient.NewUpdateNetworkApplianceFirewallL7FirewallRulesRequest()
	var rules []openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner

	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner
			rule.SetPolicy(attribute.Policy.ValueString())
			rule.SetType(attribute.Type.ValueString())
			rule.SetValue(attribute.Value.ValueString())
			rules = append(rules, rule)
		}
	}

	updateNetworkApplianceFirewallL7FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceFirewallL7FirewallRulesRequest(updateNetworkApplianceFirewallL7FirewallRules).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

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

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksApplianceFirewallL7FirewallRulesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rules := []openApiClient.UpdateNetworkApplianceFirewallL7FirewallRulesRequestRulesInner{}
	updateNetworkApplianceFirewallL7FirewallRules := *openApiClient.NewUpdateNetworkApplianceFirewallL7FirewallRulesRequest()
	updateNetworkApplianceFirewallL7FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceFirewallL7FirewallRulesRequest(updateNetworkApplianceFirewallL7FirewallRules).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

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

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksApplianceFirewallL7FirewallRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
