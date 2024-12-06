package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"strings"

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
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewNetworkApplianceStaticRoutesResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	client *openApiClient.APIClient
}

// resourceModel describes the resource data model.
type resourceModel struct {
	Id                             types.String                   `tfsdk:"id"`
	NetworkId                      jsontypes.String               `tfsdk:"network_id" json:"networkId"`
	StaticRoutId                   jsontypes.String               `tfsdk:"static_route_id" json:"id"`
	Enable                         jsontypes.Bool                 `tfsdk:"enable" json:"enabled"`
	Name                           jsontypes.String               `tfsdk:"name" json:"name"`
	GatewayIp                      jsontypes.String               `tfsdk:"gateway_ip" json:"gatewayIp"`
	Subnet                         jsontypes.String               `tfsdk:"subnet" json:"subnet"`
	FixedIpAssignmentsMacAddress   jsontypes.String               `tfsdk:"fixed_ip_assignments_mac_address"`
	FixedIpAssignmentsMacIpAddress jsontypes.String               `tfsdk:"fixed_ip_assignments_mac_ip_address"`
	FixedIpAssignmentsMacName      jsontypes.String               `tfsdk:"fixed_ip_assignments_mac_name"`
	ReservedIpRanges               []ReservedIpRangeResourceModel `tfsdk:"reserved_ip_ranges" json:"reservedIpRanges"`
}

type resourceModelMacData struct {
	Ip   string `json:"ip"`
	Name string `json:"name"`
}

type ReservedIpRangeResourceModel struct {
	Comment jsontypes.String `tfsdk:"comment" json:"comment"`
	End     jsontypes.String `tfsdk:"end" json:"end"`
	Start   jsontypes.String `tfsdk:"start" json:"start"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_static_routes"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gateway_ip": schema.StringAttribute{
				MarkdownDescription: "The gateway IP (next hop) of the static route",
				Computed:            true,
				Optional:            true,
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
				Optional:    true,
				Computed:    true,
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

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *resourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createNetworkApplianceStaticRoutes := *openApiClient.NewCreateNetworkApplianceStaticRouteRequest(data.Name.ValueString(), data.Subnet.ValueString(), data.GatewayIp.ValueString())

	inlineResp, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceStaticRoute(context.Background(), data.NetworkId.ValueString()).CreateNetworkApplianceStaticRouteRequest(createNetworkApplianceStaticRoutes).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
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

	if staticRouteId := inlineResp["id"]; staticRouteId != nil {
		data.StaticRoutId = jsontypes.StringValue(staticRouteId.(string))
	}

	if reservedIpRangesResponse := inlineResp["reservedIpRanges"]; reservedIpRangesResponse != nil {
		var reservedIpRanges []ReservedIpRangeResourceModel
		jsonData, _ := json.Marshal(reservedIpRangesResponse)
		json.Unmarshal(jsonData, &reservedIpRanges)
		data.ReservedIpRanges = make([]ReservedIpRangeResourceModel, 0)
		for _, attribute := range reservedIpRanges {
			var reservedIpRange ReservedIpRangeResourceModel
			reservedIpRange.Comment = attribute.Comment
			reservedIpRange.End = attribute.End
			reservedIpRange.Start = attribute.Start
			data.ReservedIpRanges = append(data.ReservedIpRanges, reservedIpRange)
		}

	}

	if fixedIpAssignmentsResponse := inlineResp["fixedIpAssignments"]; fixedIpAssignmentsResponse != nil {
		if macresponse := fixedIpAssignmentsResponse.(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()]; macresponse != nil {
			var macData resourceModelMacData
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

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceStaticRoute(ctx, data.NetworkId.ValueString(), data.StaticRoutId.ValueString()).Execute()
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
		var reservedIpRanges []ReservedIpRangeResourceModel
		jsonData, _ := json.Marshal(reservedIpRangesResponse)
		json.Unmarshal(jsonData, &reservedIpRanges)
		data.ReservedIpRanges = make([]ReservedIpRangeResourceModel, 0)
		for _, attribute := range reservedIpRanges {
			var reservedIpRange ReservedIpRangeResourceModel
			reservedIpRange.Comment = attribute.Comment
			reservedIpRange.End = attribute.End
			reservedIpRange.Start = attribute.Start
			data.ReservedIpRanges = append(data.ReservedIpRanges, reservedIpRange)
		}
	}

	if fixedIpAssignmentsResponse := inlineResp["fixedIpAssignments"]; fixedIpAssignmentsResponse != nil {
		if macresponse := fixedIpAssignmentsResponse.(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()]; macresponse != nil {
			var macData resourceModelMacData
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

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	updateNetworkApplianceStaticRoutes := *openApiClient.NewUpdateNetworkApplianceStaticRouteRequest()

	if !data.Name.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetName(data.Name.ValueString())
	}
	if !data.Subnet.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetSubnet(data.Subnet.ValueString())
	}
	if !data.GatewayIp.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetGatewayIp(data.GatewayIp.ValueString())
	}
	if !data.Enable.IsUnknown() {
		updateNetworkApplianceStaticRoutes.SetEnabled(data.Enable.ValueBool())
	}

	if len(data.ReservedIpRanges) > 0 {
		var networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRanges []openApiClient.UpdateNetworkApplianceStaticRouteRequestReservedIpRangesInner
		for _, attribute := range data.ReservedIpRanges {
			var networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRangesData openApiClient.UpdateNetworkApplianceStaticRouteRequestReservedIpRangesInner
			networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRangesData.SetComment(attribute.Comment.ValueString())
			networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRangesData.SetEnd(attribute.End.ValueString())
			networksNetworkIdApplianceStaticRoutesStaticRouteIdReservedIpRangesData.SetStart(attribute.Start.ValueString())
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

	inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceStaticRoute(context.Background(), data.NetworkId.ValueString(), data.StaticRoutId.ValueString()).UpdateNetworkApplianceStaticRouteRequest(updateNetworkApplianceStaticRoutes).Execute()
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
		var reservedIpRanges []ReservedIpRangeResourceModel
		jsonData, _ := json.Marshal(reservedIpRangesResponse)
		json.Unmarshal(jsonData, &reservedIpRanges)
		data.ReservedIpRanges = make([]ReservedIpRangeResourceModel, 0)
		for _, attribute := range reservedIpRanges {
			var reservedIpRange ReservedIpRangeResourceModel
			reservedIpRange.Comment = attribute.Comment
			reservedIpRange.End = attribute.End
			reservedIpRange.Start = attribute.Start
			data.ReservedIpRanges = append(data.ReservedIpRanges, reservedIpRange)
		}
	}

	if fixedIpAssignmentsResponse := inlineResp["fixedIpAssignments"]; fixedIpAssignmentsResponse != nil {
		if macresponse := inlineResp["fixedIpAssignments"].(map[string]interface{})[data.FixedIpAssignmentsMacAddress.ValueString()]; macresponse != nil {
			var macData resourceModelMacData
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

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *resourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	httpResp, err := r.client.ApplianceApi.DeleteNetworkApplianceStaticRoute(ctx, data.NetworkId.ValueString(), data.StaticRoutId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
