package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"strings"

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

var (
	_ resource.Resource                = &NetworksWirelessSsidsFirewallL7FirewallRulesResource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &NetworksWirelessSsidsFirewallL7FirewallRulesResource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &NetworksWirelessSsidsFirewallL7FirewallRulesResource{} // Interface for resources with import state functionality
)

// The NewNetworksWirelessSsidsFirewallL7FirewallRulesResource function is a constructor for the resource.
func NewNetworksWirelessSsidsFirewallL7FirewallRulesResource() resource.Resource {
	return &NetworksWirelessSsidsFirewallL7FirewallRulesResource{}
}

// NetworksWirelessSsidsFirewallL7FirewallRulesResource struct defines the structure for this resource.
type NetworksWirelessSsidsFirewallL7FirewallRulesResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The NetworksWirelessSsidsFirewallL7FirewallRulesResourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type NetworksWirelessSsidsFirewallL7FirewallRulesResourceModel struct {
	Id        jsontypes.String                                                `tfsdk:"id"`
	NetworkId jsontypes.String                                                `tfsdk:"network_id" json:"network_id"`
	Number    jsontypes.String                                                `tfsdk:"number"`
	Rules     []NetworksWirelessSsidsFirewallL7FirewallRulesResourceModelRule `tfsdk:"rules" json:"rules"`
}

type NetworksWirelessSsidsFirewallL7FirewallRulesResourceModelRule struct {
	Policy jsontypes.String `tfsdk:"policy"`
	Type   jsontypes.String `tfsdk:"type"`
	Value  jsontypes.String `tfsdk:"value"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *NetworksWirelessSsidsFirewallL7FirewallRulesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	// The TypeName attribute is important as it provides the user-friendly name for the resource/data source.
	// This is the name users will use to reference the resource/data source and it's also used in the acceptance tests.
	resp.TypeName = req.ProviderTypeName + "_networks_wireless_ssids_firewall_l7_firewall_rules"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *NetworksWirelessSsidsFirewallL7FirewallRulesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksWirelessSsidsFirewallL7FirewallRules updates Networks Wireless Ssids Firewall L7FirewallRules",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
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
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"number": schema.StringAttribute{
				MarkdownDescription: "SsIds SsidNumber",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"rules": schema.SetNestedAttribute{
				MarkdownDescription: "An array of L7 firewall rules for this SSID. Rules will get applied in the same order user has specified in request. Empty array will clear the L7 firewall rule configuration.",
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
							MarkdownDescription: "The value of what needs to get blocked. Format of the value varies depending on type of the firewall rule selected.",
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

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
func (r *NetworksWirelessSsidsFirewallL7FirewallRulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// The provider must be properly configured before it can be used.
	if req.ProviderData == nil {
		return
	}

	// Here we expect the provider data to be of type *openApiClient.APIClient.
	client, ok := req.ProviderData.(*openApiClient.APIClient)

	// This is a fatal error and the provider cannot proceed without it.
	// If you see this error, it means there is an issue with the provider setup.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	// This allows the resource to use the configured provider for any API calls it needs to make.
	r.client = client
}

// Create method is responsible for creating a new resource.
// It takes a CreateRequest containing the planned state of the new resource and returns a CreateResponse
// with the final state of the new resource or an error.
func (r *NetworksWirelessSsidsFirewallL7FirewallRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksWirelessSsidsFirewallL7FirewallRulesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkWirelessSsidFirewallL7FirewallRules := *openApiClient.NewUpdateNetworkWirelessSsidFirewallL7FirewallRulesRequest()
	var rules []openApiClient.UpdateNetworkWirelessSsidFirewallL7FirewallRulesRequestRulesInner
	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateNetworkWirelessSsidFirewallL7FirewallRulesRequestRulesInner
			rule.SetPolicy(attribute.Policy.ValueString())
			rule.SetType(attribute.Type.ValueString())
			rule.SetValue(attribute.Value.ValueString())
			rules = append(rules, rule)
		}
	}
	updateNetworkWirelessSsidFirewallL7FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.FirewallApi.UpdateNetworkWirelessSsidFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidFirewallL7FirewallRulesRequest(updateNetworkWirelessSsidFirewallL7FirewallRules).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
// It takes a ReadRequest and returns a ReadResponse with the current state of the resource or an error.
func (r *NetworksWirelessSsidsFirewallL7FirewallRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksWirelessSsidsFirewallL7FirewallRulesResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.FirewallApi.GetNetworkWirelessSsidFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *NetworksWirelessSsidsFirewallL7FirewallRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksWirelessSsidsFirewallL7FirewallRulesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkWirelessSsidFirewallL7FirewallRules := *openApiClient.NewUpdateNetworkWirelessSsidFirewallL7FirewallRulesRequest()
	var rules []openApiClient.UpdateNetworkWirelessSsidFirewallL7FirewallRulesRequestRulesInner
	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateNetworkWirelessSsidFirewallL7FirewallRulesRequestRulesInner
			rule.SetPolicy(attribute.Policy.ValueString())
			rule.SetType(attribute.Type.ValueString())
			rule.SetValue(attribute.Value.ValueString())
			rules = append(rules, rule)
		}
	}
	updateNetworkWirelessSsidFirewallL7FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.FirewallApi.UpdateNetworkWirelessSsidFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidFirewallL7FirewallRulesRequest(updateNetworkWirelessSsidFirewallL7FirewallRules).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
// It uses a DeleteRequest and responds with a DeleteResponse which contains the updated state of the resource or an error.
func (r *NetworksWirelessSsidsFirewallL7FirewallRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksWirelessSsidsFirewallL7FirewallRulesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there was an error reading the state, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkWirelessSsidFirewallL7FirewallRules := *openApiClient.NewUpdateNetworkWirelessSsidFirewallL7FirewallRulesRequest()
	var rules []openApiClient.UpdateNetworkWirelessSsidFirewallL7FirewallRulesRequestRulesInner
	if len(data.Rules) > 0 {
		for _, attribute := range data.Rules {
			var rule openApiClient.UpdateNetworkWirelessSsidFirewallL7FirewallRulesRequestRulesInner
			rule.SetPolicy(attribute.Policy.ValueString())
			rule.SetType(attribute.Type.ValueString())
			rule.SetValue(attribute.Value.ValueString())
			rules = append(rules, rule)
		}
	}
	updateNetworkWirelessSsidFirewallL7FirewallRules.SetRules(rules)

	_, httpResp, err := r.client.FirewallApi.UpdateNetworkWirelessSsidFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidFirewallL7FirewallRulesRequest(updateNetworkWirelessSsidFirewallL7FirewallRules).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// TODO: The resource has been deleted, so remove it from the state.
	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
// The function expects an ImportStateRequest and responds with an ImportStateResponse which contains
// the new state of the resource or an error.
func (r *NetworksWirelessSsidsFirewallL7FirewallRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, serial number. Got: %q", req.ID),
		)
		return
	}

	// Set the attributes required for making a Read API call in the state.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), idParts[1])...)

	// If there were any errors setting the attributes, return early.
	if resp.Diagnostics.HasError() {
		return
	}
}
