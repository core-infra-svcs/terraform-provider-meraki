package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworkApplianceStaticRoutesResource{}
	_ resource.ResourceWithConfigure   = &NetworkApplianceStaticRoutesResource{}
	_ resource.ResourceWithImportState = &NetworkApplianceStaticRoutesResource{}
)

func NewNetworkApplianceStaticRoutesResource() resource.Resource {
	return &NetworkApplianceStaticRoutesResource{}
}

// NetworkApplianceStaticRoutesResource defines the resource implementation.
type NetworkApplianceStaticRoutesResource struct {
	client *openApiClient.APIClient
}

// NetworkApplianceStaticRoutesResourceModel describes the resource data model.
type NetworkApplianceStaticRoutesResourceModel struct {
	Id                             types.String      `tfsdk:"id"`
	NetworkId                      jsontypes.String  `tfsdk:"network_id" json:"networkId"`
	StaticRoutId                   jsontypes.String  `tfsdk:"static_route_id" json:"id"`
	Ipversion                      jsontypes.Int64   `tfsdk:"ip_version" json:"ipVersion"`
	Enable                         jsontypes.Bool    `tfsdk:"enable" json:"enable"`
	Name                           jsontypes.String  `tfsdk:"name" json:"name"`
	GatewayIp                      jsontypes.String  `tfsdk:"gateway_ip" json:"gatewayIp"`
	Subnet                         jsontypes.String  `tfsdk:"subnet" json:"subnet"`
	GatewayVlanId                  jsontypes.String  `tfsdk:"gateway_vlan_id" json:"gatewayVlanId"`
	FixedIpAssignmentsMacAddress   jsontypes.String  `tfsdk:"fixed_ip_assignments_mac_address"`
	FixedIpAssignmentsMacIpAddress jsontypes.String  `tfsdk:"fixed_ip_assignments_mac_ip_address"`
	FixedIpAssignmentsMacName      jsontypes.String  `tfsdk:"fixed_ip_assignments_mac_name"`
	ReservedIpRanges               []ReservedIpRange `tfsdk:"reserved_ip_ranges" json:"reserved_ip_ranges"`
}

type MacData struct {
	Ip   string `json:"ip"`
	Name string `json:"name"`
}

type ReservedIpRange struct {
	Comment jsontypes.String `tfsdk:"comment" json:"comment"`
	End     jsontypes.String `tfsdk:"end" json:"end"`
	Start   jsontypes.String `tfsdk:"start" json:"start"`
}

func (r *NetworkApplianceStaticRoutesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_static_routes"
}

func (r *NetworkApplianceStaticRoutesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage static route for an MX or teleworker network.",

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
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"static_route_id": schema.StringAttribute{
				MarkdownDescription: "Static Route ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"ip_version": schema.Int64Attribute{
				MarkdownDescription: "Ip Version",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"gateway_ip": schema.StringAttribute{
				MarkdownDescription: "The gateway IP (next hop) of the static route",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: "The enabled state of the static route",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the new static route",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "The subnet of the static route",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"gateway_vlan_id": schema.StringAttribute{
				MarkdownDescription: "The gateway IP (next hop) VLAN ID of the static route",
				Computed:            true,
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"fixed_ip_assignments_mac_address": schema.StringAttribute{
				MarkdownDescription: "The DHCP fixed IP assignments on the static route. MAC address",
				Computed:            true,
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"fixed_ip_assignments_mac_ip_address": schema.StringAttribute{
				MarkdownDescription: "The DHCP fixed IP assignments on the static route. MAC IP address",
				Computed:            true,
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"fixed_ip_assignments_mac_name": schema.StringAttribute{
				MarkdownDescription: "The DHCP fixed IP assignments on the static route. MAC Name",
				Computed:            true,
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"reserved_ip_ranges": schema.SetNestedAttribute{
				Description: "The DHCP reserved IP ranges on the static route",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"comment": schema.StringAttribute{
							MarkdownDescription: "A text comment for the reserved range",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"end": schema.StringAttribute{
							MarkdownDescription: "The last IP in the reserved range",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"start": schema.StringAttribute{
							MarkdownDescription: "The first IP in the reserved range",
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

func (r *NetworkApplianceStaticRoutesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworkApplianceStaticRoutesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworkApplianceStaticRoutesResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createNetworkApplianceStaticRoutes := *openApiClient.NewInlineObject48(data.Name.ValueString(), data.Subnet.ValueString(), data.GatewayIp.ValueString())
	if !data.GatewayVlanId.IsUnknown() {
		createNetworkApplianceStaticRoutes.SetGatewayVlanId(data.GatewayVlanId.ValueString())
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceStaticRoute(context.Background(), data.NetworkId.ValueString()).CreateNetworkApplianceStaticRoute(createNetworkApplianceStaticRoutes).Execute()
	if err != nil {

		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)

	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	if reservedIpRangesResponse := inlineResp["reservedIpRanges"]; reservedIpRangesResponse != nil {
		var reservedIpRanges []ReservedIpRange
		jsonData, _ := json.Marshal(reservedIpRangesResponse)
		json.Unmarshal(jsonData, &reservedIpRanges)
		for _, attribute := range reservedIpRanges {
			var reservedIpRange ReservedIpRange
			reservedIpRange.Comment = attribute.Comment
			reservedIpRange.End = attribute.End
			reservedIpRange.Start = attribute.Start
			data.ReservedIpRanges = append(data.ReservedIpRanges, reservedIpRange)
		}
	} else {
		data.ReservedIpRanges = nil
	}

	if fixedIpAssignmentsResponse := inlineResp["fixedIpAssignments"]; fixedIpAssignmentsResponse != nil {
		if macresponse := fixedIpAssignmentsResponse.(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()]; macresponse != nil {
			data.FixedIpAssignmentsMacAddress = jsontypes.StringValue(inlineResp["fixedIpAssignments"].(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()].(string))
			var macData MacData
			jsonData, _ := json.Marshal(fixedIpAssignmentsResponse.(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()])
			json.Unmarshal(jsonData, &macData)
			data.FixedIpAssignmentsMacIpAddress = jsontypes.StringValue(macData.Ip)
			data.FixedIpAssignmentsMacName = jsontypes.StringValue(macData.Name)
		} else {
			data.FixedIpAssignmentsMacIpAddress = jsontypes.StringNull()
			data.FixedIpAssignmentsMacAddress = jsontypes.StringNull()
			data.FixedIpAssignmentsMacName = jsontypes.StringNull()
		}
	} else {
		data.FixedIpAssignmentsMacIpAddress = jsontypes.StringNull()
		data.FixedIpAssignmentsMacAddress = jsontypes.StringNull()
		data.FixedIpAssignmentsMacName = jsontypes.StringNull()
	}

	data.Id = types.StringValue("example-id")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *NetworkApplianceStaticRoutesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworkApplianceStaticRoutesResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceStaticRoute(ctx, data.NetworkId.ValueString(), data.StaticRoutId.ValueString()).Execute()
	if err != nil {

		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)

	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	if reservedIpRangesResponse := inlineResp["reservedIpRanges"]; reservedIpRangesResponse != nil {
		var reservedIpRanges []ReservedIpRange
		jsonData, _ := json.Marshal(reservedIpRangesResponse)
		json.Unmarshal(jsonData, &reservedIpRanges)
		for _, attribute := range reservedIpRanges {
			var reservedIpRange ReservedIpRange
			reservedIpRange.Comment = attribute.Comment
			reservedIpRange.End = attribute.End
			reservedIpRange.Start = attribute.Start
			data.ReservedIpRanges = append(data.ReservedIpRanges, reservedIpRange)
		}
	} else {
		data.ReservedIpRanges = nil
	}

	if fixedIpAssignmentsResponse := inlineResp["fixedIpAssignments"]; fixedIpAssignmentsResponse != nil {
		if macresponse := fixedIpAssignmentsResponse.(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()]; macresponse != nil {
			data.FixedIpAssignmentsMacAddress = jsontypes.StringValue(inlineResp["fixedIpAssignments"].(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()].(string))
			var macData MacData
			jsonData, _ := json.Marshal(fixedIpAssignmentsResponse.(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()])
			json.Unmarshal(jsonData, &macData)
			data.FixedIpAssignmentsMacIpAddress = jsontypes.StringValue(macData.Ip)
			data.FixedIpAssignmentsMacName = jsontypes.StringValue(macData.Name)
		} else {
			data.FixedIpAssignmentsMacIpAddress = jsontypes.StringNull()
			data.FixedIpAssignmentsMacAddress = jsontypes.StringNull()
			data.FixedIpAssignmentsMacName = jsontypes.StringNull()
		}
	} else {
		data.FixedIpAssignmentsMacIpAddress = jsontypes.StringNull()
		data.FixedIpAssignmentsMacAddress = jsontypes.StringNull()
		data.FixedIpAssignmentsMacName = jsontypes.StringNull()
	}

	data.Id = types.StringValue("example-id")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworkApplianceStaticRoutesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworkApplianceStaticRoutesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	updateNetworkApplianceStaticRoutes := *openApiClient.NewInlineObject49()
	if !data.GatewayVlanId.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetGatewayVlanId(data.GatewayVlanId.ValueString())
	}
	if !data.Name.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetName(data.Name.ValueString())
	}
	if !data.Subnet.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetSubnet(data.Subnet.ValueString())
	}
	if !data.GatewayIp.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetGatewayIp(data.GatewayIp.ValueString())
	}
	if !data.GatewayVlanId.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetGatewayVlanId(data.GatewayVlanId.ValueString())
	}
	if !data.Enable.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetEnabled(data.Enable.ValueBool())
	}

	if len(data.ReservedIpRanges) > 0 {
		var networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRanges []openApiClient.NetworksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRanges
		for _, attribute := range data.ReservedIpRanges {
			var networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRangesData openApiClient.NetworksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRanges
			networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRangesData.SetComment(attribute.Comment.ValueString())
			networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRangesData.SetComment(attribute.End.ValueString())
			networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRangesData.SetComment(attribute.Start.ValueString())
			networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRanges = append(networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRanges, networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRangesData)
		}
		updateNetworkApplianceStaticRoutes.SetReservedIpRanges(networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRanges)
	}

	if !data.FixedIpAssignmentsMacAddress.IsUnknown() {

		if !data.FixedIpAssignmentsMacIpAddress.IsUnknown() || !data.FixedIpAssignmentsMacName.IsUnknown() {

			fixedIpAssignmentsMapData := map[string]interface{}{
				data.FixedIpAssignmentsMacAddress.ValueString(): map[string]interface{}{
					"ip":   data.FixedIpAssignmentsMacIpAddress.ValueString(),
					"name": data.FixedIpAssignmentsMacName.ValueString(),
				},
			}

			updateNetworkApplianceStaticRoutes.SetFixedIpAssignments(fixedIpAssignmentsMapData)
		} else {
			resp.Diagnostics.AddError(
				"Ip address or Name Missing",
				"Add Ip address and name to terraform config",
			)
			return
		}
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceStaticRoute(context.Background(), data.NetworkId.ValueString(), data.StaticRoutId.ValueString()).UpdateNetworkApplianceStaticRoute(updateNetworkApplianceStaticRoutes).Execute()
	if err != nil {

		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)

	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	if reservedIpRangesResponse := inlineResp["reservedIpRanges"]; reservedIpRangesResponse != nil {
		var reservedIpRanges []ReservedIpRange
		jsonData, _ := json.Marshal(inlineResp["reservedIpRanges"])
		json.Unmarshal(jsonData, &reservedIpRanges)
		for _, attribute := range reservedIpRanges {
			var reservedIpRange ReservedIpRange
			reservedIpRange.Comment = attribute.Comment
			reservedIpRange.End = attribute.End
			reservedIpRange.Start = attribute.Start
			data.ReservedIpRanges = append(data.ReservedIpRanges, reservedIpRange)
		}
	} else {
		data.ReservedIpRanges = nil
	}

	if fixedIpAssignmentsResponse := inlineResp["fixedIpAssignments"]; fixedIpAssignmentsResponse != nil {
		if macresponse := inlineResp["fixedIpAssignments"].(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()]; macresponse != nil {
			data.FixedIpAssignmentsMacAddress = jsontypes.StringValue(inlineResp["fixedIpAssignments"].(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()].(string))
			var macData MacData
			jsonData, _ := json.Marshal(inlineResp["fixedIpAssignments"].(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()])
			json.Unmarshal(jsonData, &macData)
			data.FixedIpAssignmentsMacIpAddress = jsontypes.StringValue(macData.Ip)
			data.FixedIpAssignmentsMacName = jsontypes.StringValue(macData.Name)
		} else {
			data.FixedIpAssignmentsMacIpAddress = jsontypes.StringNull()
			data.FixedIpAssignmentsMacAddress = jsontypes.StringNull()
			data.FixedIpAssignmentsMacName = jsontypes.StringNull()
		}
	} else {
		data.FixedIpAssignmentsMacIpAddress = jsontypes.StringNull()
		data.FixedIpAssignmentsMacAddress = jsontypes.StringNull()
		data.FixedIpAssignmentsMacName = jsontypes.StringNull()
	}

	data.Id = types.StringValue("example-id")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworkApplianceStaticRoutesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworkApplianceStaticRoutesResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	httpResp, err := r.client.ApplianceApi.DeleteNetworkApplianceStaticRoute(ctx, data.NetworkId.ValueString(), data.StaticRoutId.ValueString()).Execute()
	if err != nil {

		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)

	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworkApplianceStaticRoutesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, static_route_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("static_route_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
