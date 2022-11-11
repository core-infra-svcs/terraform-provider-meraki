package provider

import (
	"context"
	"encoding/json"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &DevicesResource{}
var _ resource.ResourceWithImportState = &DevicesResource{}

func NewDevicesResource() resource.Resource {
	return &DevicesResource{}
}

// DevicesResource defines the resource implementation.
type DevicesResource struct {
	client *apiclient.APIClient
}

// DevicesResourceModel describes the resource data model.
type DevicesResourceModel struct {
	Id                  types.String  `tfsdk:"id"`
	Serial              types.String  `tfsdk:"serial"`
	Name                types.String  `tfsdk:"name"`
	Mac                 types.String  `tfsdk:"mac"`
	Model               types.String  `tfsdk:"model"`
	Tags                []string      `tfsdk:"tags"`
	LanIp               types.String  `tfsdk:"lan_ip"`
	Firmware            types.String  `tfsdk:"firmware"`
	Lat                 types.Float64 `tfsdk:"lat"`
	Lng                 types.Float64 `tfsdk:"lng"`
	Address             types.String  `tfsdk:"address"`
	Notes               types.String  `tfsdk:"notes"`
	MoveMapMarker       types.Bool    `tfsdk:"move_map_marker"`
	FloorPlanId         types.String  `tfsdk:"floor_plan_id"`
	NetworkId           types.String  `tfsdk:"network_id"`
	BeaconIdParamsUuid  types.String  `tfsdk:"beacon_id_params_uuid"`
	BeaconIdParamsMajor types.Int64   `tfsdk:"beacon_id_params_major"`
	BeaconIdParamsMinor types.Int64   `tfsdk:"beacon_id_params_minor"`
	SwitchProfileId     types.String  `tfsdk:"switch_profile_id"`
}

func (r *DevicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

func (r *DevicesResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Devices resource - Update the attributes of a device",
		Attributes: map[string]tfsdk.Attribute{

			"id": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            false,
			},
			"serial": {
				Description:         "The devices serial number",
				MarkdownDescription: "The devices serial number",
				Type:                types.StringType,
				Required:            true,
				Optional:            false,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"name": {
				Description:         "The name of a device",
				MarkdownDescription: "The name of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"mac": {
				Description:         "The mac address of a device",
				MarkdownDescription: "The mac address of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"model": {
				Description:         "The model of a device",
				MarkdownDescription: "The model of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"tags": {
				Description:         "The list of tags of a device",
				MarkdownDescription: "The list of tags of a device",
				Type:                types.ListType{ElemType: types.StringType},
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"lan_ip": {
				Description:         "The ipv4 lan ip of a device",
				MarkdownDescription: "The  ipv4 lan ip of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"firmware": {
				Description:         "The firmware version of a device",
				MarkdownDescription: "The firmware version of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"lat": {
				Description:         "The latitude of a device",
				MarkdownDescription: "The latitude of a device",
				Type:                types.Float64Type,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"lng": {
				Description:         "The longitude of a device",
				MarkdownDescription: "The longitude of a device",
				Type:                types.Float64Type,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"address": {
				Description:         "The address of a device",
				MarkdownDescription: "The address of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"notes": {
				Description:         "The notes for the device. String. Limited to 255 characters.",
				MarkdownDescription: "The notes for the device. String. Limited to 255 characters.",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"move_map_marker": {
				Description:         "Whether or not to set the latitude and longitude of a device based on the new address. Only applies when lat and lng are not specified.",
				MarkdownDescription: "Whether or not to set the latitude and longitude of a device based on the new address. Only applies when lat and lng are not specified.",
				Type:                types.BoolType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"floor_plan_id": {
				Description:         "The floor plan to associate to this device. null disassociates the device from the floor plan.",
				MarkdownDescription: "The floor plan to associate to this device. null disassociates the device from the floor plan.",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"network_id": {
				Description:         "The networkId of a device",
				MarkdownDescription: "The networkId of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"beacon_id_params_uuid": {
				Description:         "The beacon id params uuid of a device",
				MarkdownDescription: "The name of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"beacon_id_params_major": {
				Description:         "The beacon id params major of a device",
				MarkdownDescription: "The beacon id params major of a device",
				Type:                types.Int64Type,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"beacon_id_params_minor": {
				Description:         "The beacon id params minor of a device",
				MarkdownDescription: "The beacon id params minor of a device",
				Type:                types.Int64Type,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"switch_profile_id": {
				Description:         "The ID of a switch profile to bind to the device (for available switch profiles, see the 'Switch Profiles' endpoint). Use null to unbind the switch device from the current profile. For a device to be bindable to a switch profile, it must (1) be a switch, and (2) belong to a network that is bound to a configuration template.",
				MarkdownDescription: "The ID of a switch profile to bind to the device (for available switch profiles, see the 'Switch Profiles' endpoint). Use null to unbind the switch device from the current profile. For a device to be bindable to a switch profile, it must (1) be a switch, and (2) belong to a network that is bound to a configuration template.",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
		},
	}, nil
}

func (r *DevicesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DevicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := apiclient.NewInlineObject()
	payload.SetName(data.Name.ValueString())
	payload.SetAddress(data.Address.ValueString())
	payload.SetLat(float32(data.Lat.ValueFloat64()))
	payload.SetLng(float32(data.Lng.ValueFloat64()))
	payload.SetFloorPlanId(data.FloorPlanId.ValueString())
	payload.SetMoveMapMarker(data.MoveMapMarker.ValueBool())
	payload.SetNotes(data.Notes.ValueString())
	payload.SetSwitchProfileId(data.SwitchProfileId.ValueString())
	payload.SetTags(data.Tags)

	response, d, err := r.client.DevicesApi.UpdateDevice(context.Background(),
		data.Serial.Value).UpdateDevice(*payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Create Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response --",
			fmt.Sprintf("%v\n", d),
		)
		return
	}

	responseData, _ := json.Marshal(response)
	var results apiclient.InlineObject
	json.Unmarshal(responseData, &results)

	data.Id = types.String{Value: "example-id"}
	data.Name = types.String{Value: results.GetName()}
	data.Notes = types.String{Value: results.GetNotes()}
	data.FloorPlanId = types.String{Value: results.GetFloorPlanId()}
	data.Address = types.String{Value: results.GetAddress()}
	data.Tags = results.GetTags()
	data.Lat = types.Float64{Value: float64(results.GetLat())}
	data.Lng = types.Float64{Value: float64(results.GetLng())}
	data.SwitchProfileId = types.String{Value: results.GetSwitchProfileId()}

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DevicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, d, err := r.client.DevicesApi.GetDevice(context.Background(), data.Serial.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Read Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response --",
			fmt.Sprintf("%v\n", d),
		)
		return
	}

	responseData, _ := json.Marshal(response)
	var results apiclient.InlineObject
	json.Unmarshal(responseData, &results)

	data.Id = types.String{Value: "example-id"}
	data.Name = types.String{Value: results.GetName()}
	data.Notes = types.String{Value: results.GetNotes()}
	data.FloorPlanId = types.String{Value: results.GetFloorPlanId()}
	data.Address = types.String{Value: results.GetAddress()}
	data.Tags = results.GetTags()
	data.Lat = types.Float64{Value: float64(results.GetLat())}
	data.Lng = types.Float64{Value: float64(results.GetLng())}
	data.SwitchProfileId = types.String{Value: results.GetSwitchProfileId()}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DevicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DevicesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := apiclient.NewInlineObject()
	payload.SetName(data.Name.ValueString())
	payload.SetAddress(data.Address.ValueString())
	payload.SetLat(float32(data.Lat.ValueFloat64()))
	payload.SetLng(float32(data.Lng.ValueFloat64()))
	payload.SetFloorPlanId(data.FloorPlanId.ValueString())
	payload.SetMoveMapMarker(data.MoveMapMarker.ValueBool())
	payload.SetNotes(data.Notes.ValueString())
	payload.SetSwitchProfileId(data.SwitchProfileId.ValueString())
	payload.SetTags(data.Tags)

	response, d, err := r.client.DevicesApi.UpdateDevice(context.Background(),
		data.Serial.Value).UpdateDevice(*payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Update Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response --",
			fmt.Sprintf("%v\n", d),
		)
		return
	}

	responseData, _ := json.Marshal(response)
	var results apiclient.InlineObject
	json.Unmarshal(responseData, &results)

	data.Id = types.String{Value: "example-id"}
	data.Name = types.String{Value: results.GetName()}
	data.Notes = types.String{Value: results.GetNotes()}
	data.FloorPlanId = types.String{Value: results.GetFloorPlanId()}
	data.Address = types.String{Value: results.GetAddress()}
	data.Tags = results.GetTags()
	data.Lat = types.Float64{Value: float64(results.GetLat())}
	data.Lng = types.Float64{Value: float64(results.GetLng())}
	data.SwitchProfileId = types.String{Value: results.GetSwitchProfileId()}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DevicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := apiclient.NewInlineObject()
	payload.SetName("")
	payload.SetAddress("")
	payload.SetLat(0)
	payload.SetLng(0)
	payload.SetFloorPlanId("")
	payload.SetMoveMapMarker(false)
	payload.SetNotes("")
	payload.SetSwitchProfileId("")
	payload.SetTags([]string{})

	response, d, err := r.client.DevicesApi.UpdateDevice(context.Background(),
		data.Serial.Value).UpdateDevice(*payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"-- Delete Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response --",
			fmt.Sprintf("%v\n", d),
		)
		return
	}

	responseData, _ := json.Marshal(response)
	var results apiclient.InlineObject
	json.Unmarshal(responseData, &results)

	if d.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"-- Delete Error --",
			fmt.Sprintf("%v\n", err.Error()),
		)
		resp.Diagnostics.AddError(
			"-- Response --",
			fmt.Sprintf("%v\n", d),
		)
		resp.Diagnostics.AddError(
			"-- Results --",
			fmt.Sprintf("%v\n", results),
		)
		return
	}

	resp.State.RemoveResource(ctx)

}

func (r *DevicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
