package licences

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

type dataSourceModel struct {
	Id             jsontypes.String      `tfsdk:"id"`
	OrganizationId jsontypes.String      `tfsdk:"organization_id"`
	PerPage        jsontypes.Int64       `tfsdk:"per_page"`
	StartingAfter  jsontypes.String      `tfsdk:"starting_after"`
	EndingBefore   jsontypes.String      `tfsdk:"ending_before"`
	DeviceSerial   jsontypes.String      `tfsdk:"device_serial"`
	NetworkId      jsontypes.String      `tfsdk:"network_id"`
	State          jsontypes.String      `tfsdk:"state"`
	List           []dataSourceModelList `tfsdk:"list"`
}

type dataSourceModelList struct {
	Id                        jsontypes.String                                                                      `tfsdk:"id"`
	LicenseType               jsontypes.String                                                                      `tfsdk:"license_type"`
	LicenseKey                jsontypes.String                                                                      `tfsdk:"license_key"`
	OrderNumber               jsontypes.String                                                                      `tfsdk:"order_number"`
	DeviceSerial              jsontypes.String                                                                      `tfsdk:"device_serial"`
	NetworkId                 jsontypes.String                                                                      `tfsdk:"network_id"`
	State                     jsontypes.String                                                                      `tfsdk:"state"`
	SeatCount                 jsontypes.Int64                                                                       `tfsdk:"seat_count"`
	TotalDurationInDays       jsontypes.Int64                                                                       `tfsdk:"total_duration_in_days"`
	DurationInDays            jsontypes.Int64                                                                       `tfsdk:"duration_in_days"`
	PermanentlyQueuedLicenses []openApiClient.GetOrganizationLicenses200ResponseInnerPermanentlyQueuedLicensesInner `tfsdk:"permanently_queued_licenses"`
	ClaimDate                 jsontypes.String                                                                      `tfsdk:"claim_date"`
	ActivationDate            jsontypes.String                                                                      `tfsdk:"activation_date"`
	ExpirationDate            jsontypes.String                                                                      `tfsdk:"expiration_date"`
	HeadLicenseId             jsontypes.String                                                                      `tfsdk:"head_license_id"`
}
