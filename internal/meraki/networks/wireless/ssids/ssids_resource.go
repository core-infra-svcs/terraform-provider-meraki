package ssids

import (
	"bytes"
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksWirelessSsidsResource{}
	_ resource.ResourceWithConfigure   = &NetworksWirelessSsidsResource{}
	_ resource.ResourceWithImportState = &NetworksWirelessSsidsResource{}
)

func NewNetworksWirelessSsidsResource() resource.Resource {
	return &NetworksWirelessSsidsResource{
		typeName: "meraki_networks_wireless_ssids",
	}
}

// NetworksWirelessSsidsResource defines the resource implementation.
type NetworksWirelessSsidsResource struct {
	client        *openApiClient.APIClient
	typeName      string
	encryptionKey string
}

// updateNetworksWirelessSsidsResourceState updates the resource state with the provided api data.
func updateNetworksWirelessSsidsResourceState(ctx context.Context, plan *NetworksWirelessSsidResourceModel, state *NetworksWirelessSsidResourceModel, data *openApiClient.GetNetworkWirelessSsids200ResponseInner, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics

	rawResp, err := utils.ExtractResponseToMap(httpResp)
	if err != nil {
		diags.AddError("Failed to Unmarshal HttpResp", err.Error())
		return diags
	}

	// NetworkId
	if state.NetworkId.IsNull() || state.NetworkId.IsUnknown() {
		state.NetworkId = plan.NetworkId
	}

	if state.NetworkId.IsNull() || state.NetworkId.IsUnknown() {
		state.NetworkId = types.StringNull()
	}

	// Number
	if state.Number.IsNull() || state.Number.IsUnknown() {
		number := int64(*data.Number)
		state.Number = types.Int64Value(number)
	}

	// Import ID

	if state.Id.IsNull() && state.Id.IsUnknown() || !state.NetworkId.IsNull() || !state.NetworkId.IsUnknown() {
		id := state.NetworkId.ValueString() + "," + strconv.FormatInt(state.Number.ValueInt64(), 10)
		state.Id = types.StringValue(id)
	}

	if state.Id.IsNull() && state.Id.IsUnknown() {
		state.Id = types.StringNull()
	}

	// Name
	if state.Name.IsNull() || state.Name.IsUnknown() {
		state.Name, diags = utils.ExtractStringAttr(rawResp, "name")
		if diags.HasError() {
			diags.AddError("Name Attribute", state.Name.ValueString())
			return diags
		}

	}

	// Enabled
	if state.Enabled.IsNull() || state.Enabled.IsUnknown() {
		state.Enabled, diags = utils.ExtractBoolAttr(rawResp, "enabled")
		if diags.HasError() {
			diags.AddError("Enabled Attribute", "")
			return diags
		}
	}

	// SplashPage
	if state.SplashPage.IsNull() || state.SplashPage.IsUnknown() {
		state.SplashPage, diags = utils.ExtractStringAttr(rawResp, "splashPage")

		if diags.HasError() {
			tflog.Error(ctx, "SplashPage Attribute")
			return diags
		}

	}

	// SsidAdminAccessible
	if state.SsidAdminAccessible.IsNull() || state.SsidAdminAccessible.IsUnknown() {

		state.SsidAdminAccessible, diags = utils.ExtractBoolAttr(rawResp, "ssidAdminAccessible")
		if diags.HasError() {
			diags.AddError("SSIDAdminAccessible Attribute", "")
			return diags
		}
	}

	// LocalAuth
	if state.LocalAuth.IsNull() || state.LocalAuth.IsUnknown() {

		state.LocalAuth, diags = utils.ExtractBoolAttr(rawResp, "localAuth")
		if diags.HasError() {
			diags.AddError("LocalAuth Attribute", "")
			return diags
		}
	}

	// AuthMode
	if state.AuthMode.IsNull() || state.AuthMode.IsUnknown() {

		state.AuthMode, diags = utils.ExtractStringAttr(rawResp, "authMode")
		if diags.HasError() {
			diags.AddError("AuthMode Attribute", "")
			return diags
		}

	}

	// EncryptionMode
	if state.EncryptionMode.IsNull() || state.EncryptionMode.IsUnknown() {

		state.EncryptionMode, diags = utils.ExtractStringAttr(rawResp, "encryptionMode")
		if diags.HasError() {
			diags.AddError("Encryption Attribute", "")
			return diags
		}

	}

	// WPAEncryptionMode
	if state.WPAEncryptionMode.IsNull() || state.WPAEncryptionMode.IsUnknown() {

		state.WPAEncryptionMode, diags = utils.ExtractStringAttr(rawResp, "wpaEncryptionMode")
		if diags.HasError() {
			diags.AddError("WPAEncryptionMode Attribute", "")
			return diags
		}
	}

	// RadiusAuthenticationNASID
	if state.RadiusAuthenticationNASID.IsNull() || state.RadiusAuthenticationNASID.IsUnknown() {
		state.RadiusAuthenticationNASID, diags = utils.ExtractStringAttr(rawResp, "radiusAuthenticationNasId")
		if diags.HasError() {
			diags.AddError("radiusAuthenticationNasId Attribute", "")
			return diags
		}
	}

	// RadiusServers
	if state.RadiusServers.IsNull() || state.RadiusServers.IsUnknown() {

		radiusServers, diags := NetworksWirelessSsidStateRadiusServers(ctx, *plan, rawResp)
		if diags.HasError() {
			diags.AddError("Radius Servers Attribute", "")
			return diags
		}
		state.RadiusServers = radiusServers

	}

	// RadiusAccountingServers
	if state.RadiusAccountingServers.IsNull() || state.RadiusAccountingServers.IsUnknown() {

		radiusAccountingServers, diags := NetworksWirelessSsidStateRadiusAccountingServers(ctx, *plan, rawResp)
		if diags.HasError() {
			diags.AddError("Radius Accounting Servers Attribute", "")
			return diags
		}
		state.RadiusAccountingServers = radiusAccountingServers
	}

	// RadiusAccountingEnabled
	if state.RadiusAccountingEnabled.IsNull() || state.RadiusAccountingEnabled.IsUnknown() {
		state.RadiusAccountingEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusAccountingEnabled")
		if diags.HasError() {
			diags.AddError("Radius Accounting Enabled Attribute", "")
			return diags
		}
	}

	// RadiusEnabled
	if state.RadiusEnabled.IsNull() || state.RadiusEnabled.IsUnknown() {

		state.RadiusEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusEnabled")
		if diags.HasError() {
			diags.AddError("Radius Enabled Attribute", "")
			return diags
		}
	}

	// RadiusAttributeForGroupPolicies
	if state.RadiusAttributeForGroupPolicies.IsNull() || state.RadiusAttributeForGroupPolicies.IsUnknown() {

		state.RadiusAttributeForGroupPolicies, diags = utils.ExtractStringAttr(rawResp, "radiusAttributeForGroupPolicies")
		if diags.HasError() {
			diags.AddError("Radius Attribute for Group Policies Attribute", "")
			return diags
		}

	}

	// RadiusFailOverPolicy
	if state.RadiusFailOverPolicy.IsNull() || state.RadiusFailOverPolicy.IsUnknown() {

		state.RadiusFailOverPolicy, diags = utils.ExtractStringAttr(rawResp, "radiusFailoverPolicy")
		if diags.HasError() {
			diags.AddError("Radius Failover Attribute", "")
			return diags
		}

		//// If the API returns null or unknown, use the value from the plan
		//if state.RadiusFailOverPolicy.IsNull() || state.RadiusFailOverPolicy.IsUnknown() {
		//	state.RadiusFailOverPolicy = plan.RadiusFailOverPolicy
		//}

	}

	// RadiusLoadBalancingPolicy
	if state.RadiusLoadBalancingPolicy.IsNull() || state.RadiusLoadBalancingPolicy.IsUnknown() {

		state.RadiusLoadBalancingPolicy, diags = utils.ExtractStringAttr(rawResp, "radiusLoadBalancingPolicy")
		if diags.HasError() {
			diags.AddError("Radius load balancing policyAttribute", "")
			return diags
		}

		//// If the API returns null or unknown, use the value from the plan
		//if state.RadiusLoadBalancingPolicy.IsNull() || state.RadiusLoadBalancingPolicy.IsUnknown() {
		//	state.RadiusLoadBalancingPolicy = plan.RadiusLoadBalancingPolicy
		//}

	}

	// IPAssignmentMode
	if state.IPAssignmentMode.IsNull() || state.IPAssignmentMode.IsUnknown() {

		state.IPAssignmentMode, diags = utils.ExtractStringAttr(rawResp, "ipAssignmentMode")
		if diags.HasError() {
			diags.AddError("IP Assignment mode Attribute", "")
			return diags
		}

	}

	// AdminSplashUrl
	if state.AdminSplashUrl.IsNull() || state.AdminSplashUrl.IsUnknown() {

		state.AdminSplashUrl, diags = networksWirelessSsidAdminSplashUrl(data)
		if diags.HasError() {
			diags.AddError("AdminSplashUrl Attribute", "")
			return diags
		}

	}

	// SplashTimeout
	if state.SplashTimeout.IsNull() || state.SplashTimeout.IsUnknown() {

		state.SplashTimeout, diags = utils.ExtractStringAttr(rawResp, "splashTimeout")
		if diags.HasError() {
			diags.AddError("SplashTimeout Attribute", "")
			return diags
		}

	}

	// WalledGardenEnabled
	if state.WalledGardenEnabled.IsNull() || state.WalledGardenEnabled.IsUnknown() {

		state.WalledGardenEnabled, diags = utils.ExtractBoolAttr(rawResp, "walledGardenEnabled")
		if diags.HasError() {
			diags.AddError("WalledGardanEnabled Attribute", "")
			return diags
		}

	}

	// WalledGardenRanges
	if state.WalledGardenRanges.IsNull() || state.WalledGardenRanges.IsUnknown() {

		state.WalledGardenRanges, diags = utils.ExtractListStringAttr(rawResp, "walledGardenRanges")
		if diags.HasError() {
			diags.AddError("Walled Garden Ranges Attribute", "")
			return diags
		}
	}

	// MinBitRate
	if state.MinBitRate.IsNull() || state.MinBitRate.IsUnknown() {

		// Attempt to extract the value as an int64 directly
		var minBitrateInt types.Int64
		minBitrateInt, diags := utils.ExtractInt64FromFloat(rawResp, "minBitrate")
		fmt.Printf("Type of minBitrateInt: %T\n", minBitrateInt)
		if diags.HasError() {
			diags.AddError("Min Bit Rate Attribute", "")
			return diags
		}

		// Directly assign the extracted int64 value to the state
		state.MinBitRate = minBitrateInt
	}

	// BandSelection
	if state.BandSelection.IsNull() || state.BandSelection.IsUnknown() {

		state.BandSelection, diags = utils.ExtractStringAttr(rawResp, "bandSelection")
		if diags.HasError() {
			diags.AddError("Band Selection Attribute", "")
			return diags
		}

	}

	// PerClientBandwidthLimitUp
	if state.PerClientBandwidthLimitUp.IsNull() || state.PerClientBandwidthLimitUp.IsUnknown() {

		state.PerClientBandwidthLimitUp, diags = utils.ExtractInt64FromFloat(rawResp, "perClientBandwidthLimitUp")
		if diags.HasError() {
			diags.AddError("Per client Bandwidth limit up Attribute", "")
			return diags
		}

	}

	// PerClientBandwidthLimitDown
	if state.PerClientBandwidthLimitDown.IsNull() || state.PerClientBandwidthLimitDown.IsUnknown() {

		state.PerClientBandwidthLimitDown, diags = utils.ExtractInt64FromFloat(rawResp, "perClientBandwidthLimitDown")
		if diags.HasError() {
			diags.AddError("Per client Bandwidth limit down Attribute", "")
			return diags
		}

	}

	// Visible
	if state.Visible.IsNull() || state.Visible.IsUnknown() {

		state.Visible, diags = utils.ExtractBoolAttr(rawResp, "visible")
		if diags.HasError() {
			diags.AddError("Visible Attribute", "")
			return diags
		}

	}

	// AvailableOnAllAps
	if state.AvailableOnAllAps.IsNull() || state.AvailableOnAllAps.IsUnknown() {

		state.AvailableOnAllAps, diags = utils.ExtractBoolAttr(rawResp, "availableOnAllAps")
		if diags.HasError() {
			diags.AddError("AvailableOnAllAPs Attribute", "")
			return diags
		}

	}

	// AvailabilityTags
	if state.AvailabilityTags.IsNull() || state.AvailabilityTags.IsUnknown() {

		state.AvailabilityTags, diags = utils.ExtractListStringAttr(rawResp, "availabilityTags")
		if diags.HasError() {
			diags.AddError("AvailabilityTags Attribute", "")
			return diags
		}

	}

	// PerSsidBandwidthLimitUp
	if state.PerSsidBandwidthLimitUp.IsNull() || state.PerSsidBandwidthLimitUp.IsUnknown() {

		state.PerSsidBandwidthLimitUp, diags = utils.ExtractInt64FromFloat(rawResp, "perSsidBandwidthLimitUp")
		if diags.HasError() {
			diags.AddError("perSsidBandwidthLimitUp Attribute", "")
			return diags
		}

	}

	// PerSsidBandwidthLimitDown
	if state.PerSsidBandwidthLimitDown.IsNull() || state.PerSsidBandwidthLimitDown.IsUnknown() {

		state.PerSsidBandwidthLimitDown, diags = utils.ExtractInt64FromFloat(rawResp, "perSsidBandwidthLimitDown")
		if diags.HasError() {
			diags.AddError("perSsidBandwidthLimitDown Attribute", "")
			return diags
		}

	}

	// MandatoryDhcpEnabled
	if state.MandatoryDhcpEnabled.IsNull() || state.MandatoryDhcpEnabled.IsUnknown() {

		state.MandatoryDhcpEnabled, diags = utils.ExtractBoolAttr(rawResp, "mandatoryDhcpEnabled")
		if diags.HasError() {
			diags.AddError("mandatoryDhcpEnabled Attribute", "")
			return diags
		}

	}

	// Active Directory
	if state.ActiveDirectory.IsNull() || state.ActiveDirectory.IsUnknown() {
		state.ActiveDirectory, diags = NetworksWirelessSsidStateActiveDirectory(rawResp)
		if diags.HasError() {
			diags.AddError("Active Directory Attribute", "")
			return diags
		}
	}

	// Ensure the PSK value from the state is preserved if the API does not return it
	if state.PSK.IsNull() || state.PSK.IsUnknown() {
		state.PSK, diags = utils.ExtractStringAttr(rawResp, "psk")
		if diags.HasError() {
			diags.AddError("PSK Attribute", "Error extracting PSK attribute")
			return diags
		}

		// If the API returns null or unknown, use the value from the plan
		if state.PSK.IsNull() || state.PSK.IsUnknown() {
			state.PSK = plan.PSK
		}
	} else {
		// Ensure the state value matches the planned value
		state.PSK = plan.PSK
	}

	// EnterpriseAdminAccess
	if state.EnterpriseAdminAccess.IsNull() || state.EnterpriseAdminAccess.IsUnknown() {
		state.EnterpriseAdminAccess, diags = utils.ExtractStringAttr(rawResp, "enterpriseAdminAccess")
		if diags.HasError() {
			diags.AddError("enterpriseAdminAccess Attribute", "")
			return diags
		}
	}

	// Dot11w
	if state.Dot11w.IsNull() || state.Dot11w.IsUnknown() {

		state.Dot11w, diags = NetworksWirelessSsidStateDot11w(rawResp)
		if diags.HasError() {
			diags.AddError("Dot11w Attribute", "")
			return diags
		}
	}

	// Dot11r
	if state.Dot11r.IsNull() || state.Dot11r.IsUnknown() {

		state.Dot11r, diags = NetworksWirelessSsidStateDot11r(rawResp)
		if diags.HasError() {
			diags.AddError("Dot11r Attribute", "")
			return diags
		}
	}

	// SplashGuestSponsorDomains
	if state.SplashGuestSponsorDomains.IsNull() || state.SplashGuestSponsorDomains.IsUnknown() {
		state.SplashGuestSponsorDomains, diags = utils.ExtractListStringAttr(rawResp, "splashGuestSponsorDomains")
		if diags.HasError() {
			diags.AddError("splashGuestSponsorDomains Attribute", "")
			return diags
		}
	}

	// OAuth
	if state.OAuth.IsNull() || state.OAuth.IsUnknown() {

		state.OAuth, diags = NetworksWirelessSsidStateOauth(rawResp)
		if diags.HasError() {
			diags.AddError("oAuth Attribute", "")
			return diags
		}
	}

	// LocalRadius
	if state.LocalRadius.IsNull() || state.LocalRadius.IsUnknown() {

		state.LocalRadius, diags = NetworksWirelessSsidStateLocalRadius(rawResp)
		if diags.HasError() {
			diags.AddError("LocalRadius Attribute", "")
			return diags
		}
	}

	// LDAP
	if state.LDAP.IsNull() || state.LDAP.IsUnknown() {
		state.LDAP, diags = NetworksWirelessSsidStateLdap(rawResp)
		if diags.HasError() {
			diags.AddError("LDAP Attribute", "")
			return diags
		}
	}

	// RadiusProxyEnabled
	if state.RadiusProxyEnabled.IsNull() || state.RadiusProxyEnabled.IsUnknown() {
		state.RadiusProxyEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusProxyEnabled")
		if diags.HasError() {
			diags.AddError("radiusProxyEnabled Attribute", "")
			return diags
		}
	}

	// RadiusTestingEnabled
	if state.RadiusTestingEnabled.IsNull() || state.RadiusTestingEnabled.IsUnknown() {
		state.RadiusTestingEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusTestingEnabled")
		if diags.HasError() {
			diags.AddError("radiusTestingEnabled Attribute", "")
			return diags
		}
	}

	// RadiusCalledStationID
	if state.RadiusCalledStationID.IsNull() || state.RadiusCalledStationID.IsUnknown() {
		state.RadiusCalledStationID, diags = utils.ExtractStringAttr(rawResp, "radiusCalledStationId")
		if diags.HasError() {
			diags.AddError("radiusCalledStationId Attribute", "")
			return diags
		}
	}

	// RadiusServerTimeout
	if state.RadiusServerTimeout.IsNull() || state.RadiusServerTimeout.IsUnknown() {
		state.RadiusServerTimeout, diags = utils.ExtractInt64Attr(rawResp, "radiusServerTimeout")
		if diags.HasError() {
			diags.AddError("radiusServerTimeout Attribute", "")
			return diags
		}
	}

	// RadiusServerAttemptsLimit
	if state.RadiusServerAttemptsLimit.IsNull() || state.RadiusServerAttemptsLimit.IsUnknown() {
		state.RadiusServerAttemptsLimit, diags = utils.ExtractInt64Attr(rawResp, "radiusServerAttemptsLimit")
		if diags.HasError() {
			diags.AddError("radiusServerAttemptsLimit Attribute", "")
			return diags
		}
	}

	// RadiusFallbackEnabled
	if state.RadiusFallbackEnabled.IsNull() || state.RadiusFallbackEnabled.IsUnknown() {
		state.RadiusFallbackEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusFallbackEnabled")
		if diags.HasError() {
			diags.AddError("radiusFallbackEnabled Attribute", "")
			return diags
		}
	}

	// RadiusCoaEnabled
	if state.RadiusCoaEnabled.IsNull() || state.RadiusCoaEnabled.IsUnknown() {
		state.RadiusCoaEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusCoaEnabled")
		if diags.HasError() {
			diags.AddError("radiusCoaEnabled Attribute", "")
			return diags
		}
	}

	// RadiusAccountingInterimInterval
	if state.RadiusAccountingInterimInterval.IsNull() || state.RadiusAccountingInterimInterval.IsUnknown() {
		state.RadiusAccountingInterimInterval, diags = utils.ExtractInt64Attr(rawResp, "radiusAccountingInterimInterval")
		if diags.HasError() {
			diags.AddError("radiusAccountingInterimInterval Attribute", "")
			return diags
		}
	}

	// UseVlanTagging
	if state.UseVlanTagging.IsNull() || state.UseVlanTagging.IsUnknown() {
		state.UseVlanTagging, diags = utils.ExtractBoolAttr(rawResp, "useVlanTagging")
		if diags.HasError() {
			diags.AddError("useVlanTagging Attribute", "")
			return diags
		}
	}

	// ConcentratorNetworkID
	if state.ConcentratorNetworkID.IsNull() || state.ConcentratorNetworkID.IsUnknown() {
		state.ConcentratorNetworkID, diags = utils.ExtractStringAttr(rawResp, "concentratorNetworkId")
		if diags.HasError() {
			diags.AddError("concentratorNetworkId Attribute", "")
			return diags
		}
	}

	// SecondaryConcentratorNetworkID
	if state.SecondaryConcentratorNetworkID.IsNull() || state.SecondaryConcentratorNetworkID.IsUnknown() {
		state.SecondaryConcentratorNetworkID, diags = utils.ExtractStringAttr(rawResp, "secondaryConcentratorNetworkId")
		if diags.HasError() {
			diags.AddError("secondaryConcentratorNetworkId Attribute", "")
			return diags
		}
	}

	// DisassociateClientsOnVpnFailOver
	if state.DisassociateClientsOnVpnFailOver.IsNull() || state.DisassociateClientsOnVpnFailOver.IsUnknown() {
		state.DisassociateClientsOnVpnFailOver, diags = utils.ExtractBoolAttr(rawResp, "disassociateClientsOnVpnFailOver")
		if diags.HasError() {
			diags.AddError("disassociateClientsOnVpnFailOver Attribute", "")
			return diags
		}
	}

	// VlanID
	if state.VlanID.IsNull() || state.VlanID.IsUnknown() {
		state.VlanID, diags = utils.ExtractInt64Attr(rawResp, "vlanId")
		if diags.HasError() {
			diags.AddError("vlanId Attribute", "")
			return diags
		}
	}

	// DefaultVlanID
	if state.DefaultVlanID.IsNull() || state.DefaultVlanID.IsUnknown() {
		state.DefaultVlanID, diags = utils.ExtractInt64Attr(rawResp, "defaultVlanId")
		if diags.HasError() {
			diags.AddError("defaultVlanId Attribute", "")
			return diags
		}
	}

	// ApTagsAndVlanIDs
	if state.ApTagsAndVlanIDs.IsNull() || state.ApTagsAndVlanIDs.IsUnknown() {
		state.ApTagsAndVlanIDs, diags = NetworksWirelessSsidStateApTagsAndVlanIds(rawResp)
		if diags.HasError() {
			diags.AddError("ApTagsAndVlanIDs Attribute", "")
			return diags
		}
	}

	// GRE
	if state.GRE.IsNull() || state.GRE.IsUnknown() {

		state.GRE, diags = NetworksWirelessSsidStateGre(rawResp)
		if diags.HasError() {
			diags.AddError("Gre Attribute", "")
			return diags
		}
	}

	// RadiusOverride
	if state.RadiusOverride.IsNull() || state.RadiusOverride.IsUnknown() {
		state.RadiusOverride, diags = utils.ExtractBoolAttr(rawResp, "radiusOverride")
		if diags.HasError() {
			diags.AddError("radiusOverride Attribute", "")
			return diags
		}
	}

	// RadiusGuestVlanEnabled
	if state.RadiusGuestVlanEnabled.IsNull() || state.RadiusGuestVlanEnabled.IsUnknown() {
		state.RadiusGuestVlanEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusGuestVlanEnabled")
		if diags.HasError() {
			diags.AddError("radiusGuestVlanEnabled Attribute", "")
			return diags
		}
	}

	// RadiusGuestVlanId
	if state.RadiusGuestVlanId.IsNull() || state.RadiusGuestVlanId.IsUnknown() {
		state.RadiusGuestVlanId, diags = utils.ExtractInt64Attr(rawResp, "radiusGuestVlanId")
		if diags.HasError() {
			diags.AddError("radiusGuestVlanId Attribute", "")
			return diags
		}
	}

	// LanIsolationEnabled
	if state.LanIsolationEnabled.IsNull() || state.LanIsolationEnabled.IsUnknown() {
		state.LanIsolationEnabled, diags = utils.ExtractBoolAttr(rawResp, "lanIsolationEnabled")
		if diags.HasError() {
			diags.AddError("lanIsolationEnabled Attribute", "")
			return diags
		}
	}

	// AdultContentFilteringEnabled
	if state.AdultContentFilteringEnabled.IsNull() || state.AdultContentFilteringEnabled.IsUnknown() {
		state.AdultContentFilteringEnabled, diags = utils.ExtractBoolAttr(rawResp, "adultContentFilteringEnabled")
		if diags.HasError() {
			diags.AddError("AdultContentFilteringEnabled Attribute", "")
			return diags
		}
	}

	// DnsRewrite
	if state.DnsRewrite.IsNull() || state.DnsRewrite.IsUnknown() {

		state.DnsRewrite, diags = NetworksWirelessSsidStateDnsRewrite(rawResp)
		if diags.HasError() {
			diags.AddError("DnsRewrite Attribute", "")
			return diags
		}

	}

	// SpeedBurst
	if state.SpeedBurst.IsNull() || state.SpeedBurst.IsUnknown() {
		state.SpeedBurst, diags = NetworksWirelessSsidStateSpeedBurst(rawResp)
		if diags.HasError() {
			diags.AddError("SpeedBurst Attribute", "")
			return diags
		}

	}

	// NamedVlans
	if state.NamedVlans.IsNull() || state.NamedVlans.IsUnknown() {
		state.NamedVlans, diags = NetworksWirelessSsidStateNamedVlans(rawResp)
		if diags.HasError() {
			diags.AddError("NamedVlans Attribute", "")
			return diags
		}
	}

	return diags
}

func updateNetworksWirelessSsidsResourcePayload(ctx context.Context, plan *NetworksWirelessSsidResourceModel) (openApiClient.UpdateNetworkWirelessSsidRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	var payload openApiClient.UpdateNetworkWirelessSsidRequest

	// Set simple attributes with checks
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		payload.SetName(plan.Name.ValueString())
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		payload.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.AuthMode.IsNull() && !plan.AuthMode.IsUnknown() {
		payload.SetAuthMode(plan.AuthMode.ValueString())
	}

	if !plan.EncryptionMode.IsNull() && !plan.EncryptionMode.IsUnknown() {
		payload.SetEncryptionMode(plan.EncryptionMode.ValueString())
	}

	if !plan.PSK.IsNull() && !plan.PSK.IsUnknown() {
		payload.SetPsk(plan.PSK.ValueString())
	}

	if !plan.WPAEncryptionMode.IsNull() && !plan.WPAEncryptionMode.IsUnknown() {
		payload.SetWpaEncryptionMode(plan.WPAEncryptionMode.ValueString())
	}

	if !plan.SplashPage.IsNull() && !plan.SplashPage.IsUnknown() {
		payload.SetSplashPage(plan.SplashPage.ValueString())
	}

	if !plan.RadiusProxyEnabled.IsNull() && !plan.RadiusProxyEnabled.IsUnknown() {
		payload.SetRadiusProxyEnabled(plan.RadiusProxyEnabled.ValueBool())
	}

	if !plan.RadiusTestingEnabled.IsNull() && !plan.RadiusTestingEnabled.IsUnknown() {
		payload.SetRadiusTestingEnabled(plan.RadiusTestingEnabled.ValueBool())
	}

	if !plan.RadiusCalledStationID.IsNull() && !plan.RadiusCalledStationID.IsUnknown() {
		payload.SetRadiusCalledStationId(plan.RadiusCalledStationID.ValueString())
	}

	if !plan.RadiusAuthenticationNASID.IsNull() && !plan.RadiusAuthenticationNASID.IsUnknown() {
		payload.SetRadiusAuthenticationNasId(plan.RadiusAuthenticationNASID.ValueString())
	}

	if !plan.RadiusFallbackEnabled.IsNull() && !plan.RadiusFallbackEnabled.IsUnknown() {
		payload.SetRadiusFallbackEnabled(plan.RadiusFallbackEnabled.ValueBool())
	}

	if !plan.RadiusCoaEnabled.IsNull() && !plan.RadiusCoaEnabled.IsUnknown() {
		payload.SetRadiusCoaEnabled(plan.RadiusCoaEnabled.ValueBool())
	}

	if !plan.RadiusFailOverPolicy.IsNull() && !plan.RadiusFailOverPolicy.IsUnknown() {
		payload.SetRadiusFailoverPolicy(plan.RadiusFailOverPolicy.ValueString())
	}

	if !plan.RadiusLoadBalancingPolicy.IsNull() && !plan.RadiusLoadBalancingPolicy.IsUnknown() {
		payload.SetRadiusLoadBalancingPolicy(plan.RadiusLoadBalancingPolicy.ValueString())
	}

	if !plan.RadiusAccountingEnabled.IsNull() && !plan.RadiusAccountingEnabled.IsUnknown() {
		payload.SetRadiusAccountingEnabled(plan.RadiusAccountingEnabled.ValueBool())
	}

	if !plan.IPAssignmentMode.IsNull() && !plan.IPAssignmentMode.IsUnknown() {
		payload.SetIpAssignmentMode(plan.IPAssignmentMode.ValueString())
	}

	if !plan.UseVlanTagging.IsNull() && !plan.UseVlanTagging.IsUnknown() {
		payload.SetUseVlanTagging(plan.UseVlanTagging.ValueBool())
	}

	if !plan.ConcentratorNetworkID.IsNull() && !plan.ConcentratorNetworkID.IsUnknown() {
		payload.SetConcentratorNetworkId(plan.ConcentratorNetworkID.ValueString())
	}

	if !plan.SecondaryConcentratorNetworkID.IsNull() && !plan.SecondaryConcentratorNetworkID.IsUnknown() {
		payload.SetSecondaryConcentratorNetworkId(plan.SecondaryConcentratorNetworkID.ValueString())
	}

	if !plan.DisassociateClientsOnVpnFailOver.IsNull() && !plan.DisassociateClientsOnVpnFailOver.IsUnknown() {
		payload.SetDisassociateClientsOnVpnFailover(plan.DisassociateClientsOnVpnFailOver.ValueBool())
	}

	if !plan.WalledGardenEnabled.IsNull() && !plan.WalledGardenEnabled.IsUnknown() {
		payload.SetWalledGardenEnabled(plan.WalledGardenEnabled.ValueBool())
	}

	if !plan.RadiusOverride.IsNull() && !plan.RadiusOverride.IsUnknown() {
		payload.SetRadiusOverride(plan.RadiusOverride.ValueBool())
	}

	if !plan.RadiusGuestVlanEnabled.IsNull() && !plan.RadiusGuestVlanEnabled.IsUnknown() {
		payload.SetRadiusGuestVlanEnabled(plan.RadiusGuestVlanEnabled.ValueBool())
	}

	if !plan.RadiusGuestVlanId.IsNull() && !plan.RadiusGuestVlanId.IsUnknown() {
		payload.SetRadiusGuestVlanId(int32(plan.RadiusGuestVlanId.ValueInt64()))
	}

	if !plan.BandSelection.IsNull() && !plan.BandSelection.IsUnknown() {
		payload.SetBandSelection(plan.BandSelection.ValueString())
	}

	if !plan.LanIsolationEnabled.IsNull() && !plan.LanIsolationEnabled.IsUnknown() {
		payload.SetLanIsolationEnabled(plan.LanIsolationEnabled.ValueBool())
	}

	if !plan.Visible.IsNull() && !plan.Visible.IsUnknown() {
		payload.SetVisible(plan.Visible.ValueBool())
	}

	if !plan.AvailableOnAllAps.IsNull() && !plan.AvailableOnAllAps.IsUnknown() {
		payload.SetAvailableOnAllAps(plan.AvailableOnAllAps.ValueBool())
	}

	if !plan.MandatoryDhcpEnabled.IsNull() && !plan.MandatoryDhcpEnabled.IsUnknown() {
		payload.SetMandatoryDhcpEnabled(plan.MandatoryDhcpEnabled.ValueBool())
	}

	if !plan.AdultContentFilteringEnabled.IsNull() && !plan.AdultContentFilteringEnabled.IsUnknown() {
		payload.SetAdultContentFilteringEnabled(plan.AdultContentFilteringEnabled.ValueBool())
	}

	if !plan.RadiusServerTimeout.IsNull() && !plan.RadiusServerTimeout.IsUnknown() {
		radiusServerTimeout, err := utils.Int32Pointer(plan.RadiusServerTimeout.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.SetRadiusServerTimeout(*radiusServerTimeout)
	}

	if !plan.RadiusServerAttemptsLimit.IsNull() && !plan.RadiusServerAttemptsLimit.IsUnknown() {
		radiusServerAttemptsLimit, err := utils.Int32Pointer(plan.RadiusServerAttemptsLimit.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.RadiusServerAttemptsLimit = radiusServerAttemptsLimit
	}

	if !plan.RadiusAccountingInterimInterval.IsNull() && !plan.RadiusAccountingInterimInterval.IsUnknown() {
		radiusAccountingInterimInterval, err := utils.Int32Pointer(plan.RadiusAccountingInterimInterval.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.RadiusAccountingInterimInterval = radiusAccountingInterimInterval
	}

	if !plan.VlanID.IsNull() && !plan.VlanID.IsUnknown() {
		vlanId, err := utils.Int32Pointer(plan.VlanID.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.VlanId = vlanId
	}

	if !plan.DefaultVlanID.IsNull() && !plan.DefaultVlanID.IsUnknown() {
		defaultVlanId, err := utils.Int32Pointer(plan.DefaultVlanID.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.DefaultVlanId = defaultVlanId
	}

	if !plan.MinBitRate.IsNull() && !plan.MinBitRate.IsUnknown() {
		minBitRate, err := utils.Float32Pointer(float64(plan.MinBitRate.ValueInt64()))
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.MinBitrate = minBitRate
	}

	if !plan.PerClientBandwidthLimitUp.IsNull() && !plan.PerClientBandwidthLimitUp.IsUnknown() {
		perClientBandwidthLimitUp, err := utils.Int32Pointer(plan.PerClientBandwidthLimitUp.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.PerClientBandwidthLimitUp = perClientBandwidthLimitUp
	}

	if !plan.PerClientBandwidthLimitDown.IsNull() && !plan.PerClientBandwidthLimitDown.IsUnknown() {
		perClientBandwidthLimitDown, err := utils.Int32Pointer(plan.PerClientBandwidthLimitDown.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.PerClientBandwidthLimitDown = perClientBandwidthLimitDown
	}

	if !plan.PerSsidBandwidthLimitUp.IsNull() && !plan.PerSsidBandwidthLimitUp.IsUnknown() {
		perSsidBandwidthLimitUp, err := utils.Int32Pointer(plan.PerSsidBandwidthLimitUp.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.PerSsidBandwidthLimitUp = perSsidBandwidthLimitUp
	}

	if !plan.PerSsidBandwidthLimitDown.IsNull() && !plan.PerSsidBandwidthLimitDown.IsUnknown() {
		perSsidBandwidthLimitDown, err := utils.Int32Pointer(plan.PerSsidBandwidthLimitDown.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.PerSsidBandwidthLimitDown = perSsidBandwidthLimitDown
	}

	if !plan.WalledGardenRanges.IsNull() && !plan.WalledGardenRanges.IsUnknown() {
		walledGardenRanges, err := utils.ExtractStringsFromList(plan.WalledGardenRanges)
		if err.HasError() {
			diags.Append(err...)
		}
		payload.WalledGardenRanges = walledGardenRanges
	}

	availabilityTags, err := utils.ExtractStringsFromList(plan.AvailabilityTags)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.AvailabilityTags = availabilityTags

	splashGuestSponsorDomains, err := utils.ExtractStringsFromList(plan.SplashGuestSponsorDomains)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.SplashGuestSponsorDomains = splashGuestSponsorDomains

	// Check enterprise admin access
	if !plan.EnterpriseAdminAccess.IsNull() && !plan.EnterpriseAdminAccess.IsUnknown() {
		payload.SetEnterpriseAdminAccess(plan.EnterpriseAdminAccess.ValueString())
	}

	// Expand complex attributes
	dot11w, err := NetworksWirelessSsidPayloadDot11w(plan.Dot11w)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Dot11w = dot11w

	dot11r, err := NetworksWirelessSsidPayloadDot11r(plan.Dot11r)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Dot11r = dot11r

	oauth, err := NetworksWirelessSsidPayloadOauth(plan.OAuth)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Oauth = oauth

	localRadius, err := NetworksWirelessSsidPayloadLocalRadius(plan.LocalRadius)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.LocalRadius = &localRadius

	ldap, err := NetworksWirelessSsidPayloadLdap(plan.LDAP)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Ldap = ldap

	activeDirectory, err := NetworksWirelessSsidPayloadActiveDirectory(plan.ActiveDirectory)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.ActiveDirectory = activeDirectory

	radiusServers, err := NetworksWirelessSsidPayloadRadiusServers(ctx, plan.RadiusServers)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.RadiusServers = radiusServers

	radiusAccountingServers, err := NetworksWirelessSsidPayloadRadiusAccountingServers(ctx, plan.RadiusAccountingServers)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.RadiusAccountingServers = radiusAccountingServers

	apTagsAndVlanIds, err := NetworksWirelessSsidPayloadApTagsAndVlanIds(plan.ApTagsAndVlanIDs)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.ApTagsAndVlanIds = apTagsAndVlanIds

	gre, err := NetworksWirelessSsidPayloadGre(plan.GRE)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Gre = gre

	dnsRewrite, err := NetworksWirelessSsidPayloadDnsRewrite(plan.DnsRewrite)
	payload.DnsRewrite = dnsRewrite
	if err.HasError() {
		diags.Append(err...)
	}

	speedBurst, err := NetworksWirelessSsidPayloadSpeedBurst(plan.SpeedBurst)
	payload.SpeedBurst = speedBurst
	if err.HasError() {
		diags.Append(err...)
	}

	namedVlans, err := NetworksWirelessSsidPayloadNamedVlans(plan.NamedVlans)
	payload.NamedVlans = namedVlans
	if err.HasError() {
		diags.Append(err...)
	}

	return payload, diags

}

func (r *NetworksWirelessSsidsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.typeName

}

func (r *NetworksWirelessSsidsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource, generated by the Meraki API.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"active_directory": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for Active Directory. Only valid if splashPage is 'Password-protected with Active Directory'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"credentials": schema.SingleNestedAttribute{
						MarkdownDescription: `(Optional) The credentials of the user account to be used by the AP to bind to your Active Directory server. The Active Directory account should have permissions on all your Active Directory servers. Only valid if the splashPage is 'Password-protected with Active Directory'.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"login_name": schema.StringAttribute{
								MarkdownDescription: `The login name of the Active Directory account.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"password": schema.StringAttribute{
								MarkdownDescription: `The password to the Active Directory user account.`,
								Sensitive:           true,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"servers": schema.ListNestedAttribute{
						MarkdownDescription: `The Active Directory servers to be used for authentication.`,
						Computed:            true,
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{

								"host": schema.StringAttribute{
									MarkdownDescription: `IP address (or FQDN) of your Active Directory server.`,
									Computed:            true,
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: `(Optional) UDP port the Active Directory server listens on. By default, uses port 3268.`,
									Computed:            true,
									Optional:            true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
				},
			},
			"admin_splash_url": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"adult_content_filtering_enabled": schema.BoolAttribute{
				MarkdownDescription: `Boolean indicating whether or not adult content will be blocked`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ap_tags_and_vlan_ids": schema.ListNestedAttribute{
				MarkdownDescription: `The list of tags and VLAN IDs used for VLAN tagging. This param is only valid when the ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{

						"tags": schema.ListAttribute{
							MarkdownDescription: `Array of AP tags`,
							Computed:            true,
							Optional:            true,

							ElementType: types.StringType,
						},
						"vlan_id": schema.Int64Attribute{
							MarkdownDescription: `Numerical identifier that is assigned to the VLAN`,
							Computed:            true,
							Optional:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"auth_mode": schema.StringAttribute{
				MarkdownDescription: `The association control method for the SSID`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"8021x-google",
						"8021x-localradius",
						"8021x-meraki",
						"8021x-nac",
						"8021x-radius",
						"ipsk-with-nac",
						"ipsk-with-radius",
						"ipsk-without-radius",
						"open",
						"open-enhanced",
						"open-with-nac",
						"open-with-radius",
						"psk",
					),
				},
			},
			"availability_tags": schema.ListAttribute{
				MarkdownDescription: `List of tags for this SSID. If availableOnAllAps is false, then the SSID is only broadcast by APs with tags matching any of the tags in this list`,
				Optional:            true,
				Computed:            true,

				ElementType: types.StringType,
			},
			"available_on_all_aps": schema.BoolAttribute{
				MarkdownDescription: `Whether all APs broadcast the SSID or if it's restricted to APs matching any availability tags`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"band_selection": schema.StringAttribute{
				MarkdownDescription: `The client-serving radio frequencies of this SSID in the default indoor RF profile`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"5 GHz band only",
						"Dual band operation",
						"Dual band operation with Band Steering",
					),
				},
			},
			"concentrator_network_id": schema.StringAttribute{
				MarkdownDescription: `The concentrator to use when the ipAssignmentMode is 'Layer 3 roaming with a concentrator' or 'VPN'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_vlan_id": schema.Int64Attribute{
				MarkdownDescription: `The default VLAN ID used for 'all other APs'. This param is only valid when the ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"disassociate_clients_on_vpn_fail_over": schema.BoolAttribute{
				MarkdownDescription: `Disassociate clients when 'VPN' concentrator failover occurs in order to trigger clients to re-associate and generate new DHCP requests. This param is only valid if ipAssignmentMode is 'VPN'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dns_rewrite": schema.SingleNestedAttribute{
				MarkdownDescription: `DNS servers rewrite settings`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"dns_custom_name_servers": schema.ListAttribute{
						MarkdownDescription: `User specified DNS servers (up to two servers)`,
						Computed:            true,
						Optional:            true,
						ElementType:         types.StringType,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: `Boolean indicating whether or not DNS server rewrite is enabled. If disabled, upstream DNS will be used`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"dot11r": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for 802.11r`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"adaptive": schema.BoolAttribute{
						MarkdownDescription: `(Optional) Whether 802.11r is adaptive or not.`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: `Whether 802.11r is enabled or not.`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"dot11w": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for Protected Management Frames (802.11w).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"enabled": schema.BoolAttribute{
						MarkdownDescription: `Whether 802.11w is enabled or not.`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"required": schema.BoolAttribute{
						MarkdownDescription: `(Optional) Whether 802.11w is required or not.`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether or not the SSID is enabled`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption_mode": schema.StringAttribute{
				MarkdownDescription: `The psk encryption mode for the SSID`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enterprise_admin_access": schema.StringAttribute{
				MarkdownDescription: `Whether or not an SSID is accessible by 'enterprise' administrators ('access disabled' or 'access enabled')`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"access disabled",
						"access enabled",
					),
				},
			},
			"gre": schema.SingleNestedAttribute{
				MarkdownDescription: `Ethernet over GRE settings`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"concentrator": schema.SingleNestedAttribute{
						MarkdownDescription: `The EoGRE concentrator's settings`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"host": schema.StringAttribute{
								MarkdownDescription: `The EoGRE concentrator's IP or FQDN. This param is required when ipAssignmentMode is 'Ethernet over GRE'.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"key": schema.Int64Attribute{
						MarkdownDescription: `Optional numerical identifier that will add the GRE key field to the GRE header. Used to identify an individual traffic flow within a tunnel.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"ip_assignment_mode": schema.StringAttribute{
				MarkdownDescription: `The client IP assignment mode`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Bridge mode",
						"Ethernet over GRE",
						"Layer 3 roaming",
						"Layer 3 roaming with a concentrator",
						"NAT mode",
						"VPN",
					),
				},
			},
			"lan_isolation_enabled": schema.BoolAttribute{
				MarkdownDescription: `Boolean indicating whether Layer 2 LAN isolation should be enabled or disabled. Only configurable when ipAssignmentMode is 'Bridge mode'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for LDAP. Only valid if splashPage is 'Password-protected with LDAP'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"base_distinguished_name": schema.StringAttribute{
						MarkdownDescription: `The base distinguished name of users on the LDAP server.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"credentials": schema.SingleNestedAttribute{
						MarkdownDescription: `(Optional) The credentials of the user account to be used by the AP to bind to your LDAP server. The LDAP account should have permissions on all your LDAP servers.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"distinguished_name": schema.StringAttribute{
								MarkdownDescription: `The distinguished name of the LDAP user account (example: cn=user,dc=meraki,dc=com).`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"password": schema.StringAttribute{
								MarkdownDescription: `The password of the LDAP user account.`,
								Sensitive:           true,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"server_ca_certificate": schema.SingleNestedAttribute{
						MarkdownDescription: `The CA certificate used to sign the LDAP server's key.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"contents": schema.StringAttribute{
								MarkdownDescription: `The contents of the CA certificate. Must be in PEM or DER format.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"servers": schema.ListNestedAttribute{
						MarkdownDescription: `The LDAP servers to be used for authentication.`,
						Computed:            true,
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{

								"host": schema.StringAttribute{
									MarkdownDescription: `IP address (or FQDN) of your LDAP server.`,
									Computed:            true,
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: `UDP port the LDAP server listens on.`,
									Computed:            true,
									Optional:            true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
				},
			},
			"local_auth": schema.BoolAttribute{
				MarkdownDescription: `Extended local auth flag for Enterprise NAC`,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"local_radius": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for Local Authentication, a built-in RADIUS server on the access point. Only valid if authMode is '8021x-localradius'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"cache_timeout": schema.Int64Attribute{
						MarkdownDescription: `The duration (in seconds) for which LDAP and OCSP lookups are cached.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"certificate_authentication": schema.SingleNestedAttribute{
						MarkdownDescription: `The current setting for certificate verification.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"client_root_ca_certificate": schema.SingleNestedAttribute{
								MarkdownDescription: `The Client CA Certificate used to sign the client certificate.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{

									"contents": schema.StringAttribute{
										MarkdownDescription: `The contents of the Client CA Certificate. Must be in PEM or DER format.`,
										Computed:            true,
										Optional:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
							"enabled": schema.BoolAttribute{
								MarkdownDescription: `Whether or not to use EAP-TLS certificate-based authentication to validate wireless clients.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"ocsp_responder_url": schema.StringAttribute{
								MarkdownDescription: `(Optional) The URL of the OCSP responder to verify client certificate status.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"use_ldap": schema.BoolAttribute{
								MarkdownDescription: `Whether or not to verify the certificate with LDAP.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"use_ocsp": schema.BoolAttribute{
								MarkdownDescription: `Whether or not to verify the certificate with OCSP.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"password_authentication": schema.SingleNestedAttribute{
						MarkdownDescription: `The current setting for password-based authentication.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"enabled": schema.BoolAttribute{
								MarkdownDescription: `Whether or not to use EAP-TTLS/PAP or PEAP-GTC password-based authentication via LDAP lookup.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"mandatory_dhcp_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether clients connecting to this SSID must use the IP address assigned by the DHCP server`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"min_bitrate": schema.Int64Attribute{
				MarkdownDescription: `The minimum bitrate in Mbps of this SSID in the default indoor RF profile`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: `The name of the SSID`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"named_vlans": schema.SingleNestedAttribute{
				MarkdownDescription: `Named VLAN settings.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"radius": schema.SingleNestedAttribute{
						MarkdownDescription: `RADIUS settings. This param is only valid when authMode is 'open-with-radius' and ipAssignmentMode is not 'NAT mode'.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"guest_vlan": schema.SingleNestedAttribute{
								MarkdownDescription: `guest vlan settings`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{

									"enabled": schema.BoolAttribute{
										MarkdownDescription: `Whether or not RADIUS guest named VLAN is enabled.`,
										Computed:            true,
										Optional:            true,
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
									"name": schema.StringAttribute{
										MarkdownDescription: `RADIUS guest VLAN name.`,
										Computed:            true,
										Optional:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
						},
					},
					"tagging": schema.SingleNestedAttribute{
						MarkdownDescription: `VLAN tagging settings. This param is only valid when ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"by_ap_tags": schema.ListNestedAttribute{
								MarkdownDescription: `The list of AP tags and VLAN names used for named VLAN tagging. If an AP has a tag matching one in the list, then traffic on this SSID will be directed to use the VLAN name associated to the tag.`,
								Computed:            true,
								Optional:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"tags": schema.ListAttribute{
											MarkdownDescription: `List of AP tags.`,
											Computed:            true,
											Optional:            true,
											ElementType:         types.StringType,
										},
										"vlan_name": schema.StringAttribute{
											MarkdownDescription: `VLAN name that will be used to tag traffic.`,
											Computed:            true,
											Optional:            true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
									},
								},
							},
							"default_vlan_name": schema.StringAttribute{
								MarkdownDescription: `The default VLAN name used to tag traffic in the absence of a matching AP tag.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"enabled": schema.BoolAttribute{
								MarkdownDescription: `Whether or not traffic should be directed to use specific VLAN names.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: `networkId path parameter. Network ID`,
				Required:            true,
			},
			"number": schema.Int64Attribute{
				MarkdownDescription: `Unique identifier of the SSID`,
				Required:            true,
				//            Differents_types: `   parameter: schema.TypeString, item: schema.TypeInt`,
			},
			"oauth": schema.SingleNestedAttribute{
				MarkdownDescription: `The OAuth settings of this SSID. Only valid if splashPage is 'Google OAuth'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"allowed_domains": schema.ListAttribute{
						MarkdownDescription: `(Optional) The list of domains allowed access to the network.`,
						Computed:            true,
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"per_client_bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: `The download bandwidth limit in Kbps. (0 represents no limit.)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"per_client_bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: `The upload bandwidth limit in Kbps. (0 represents no limit.)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"per_ssid_bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: `The total download bandwidth limit in Kbps (0 represents no limit)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"per_ssid_bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: `The total upload bandwidth limit in Kbps (0 represents no limit)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"psk": schema.StringAttribute{
				MarkdownDescription: `The passkey for the SSID. This param is only valid if the authMode is 'psk'`,
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					utils.NewSensitivePlanModifier(r.encryptionKey),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_accounting_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether or not RADIUS accounting is enabled`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_accounting_interim_interval": schema.Int64Attribute{
				MarkdownDescription: `The interval (in seconds) in which accounting information is updated and sent to the RADIUS accounting server.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_accounting_servers": schema.ListNestedAttribute{
				MarkdownDescription: `List of RADIUS accounting 802.1X servers to be used for authentication`,
				Optional:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						/*
							"server_id": schema.StringAttribute{  // not in api spec and changes all the time
									MarkdownDescription: `ServerId of your RADIUS server`,
									Computed:            true,
									PlanModifiers: []planmodifier.String{
										utils.NewSuppressDiffServerIDModifier(),
									},
								},
						*/
						"ca_certificate": schema.StringAttribute{
							MarkdownDescription: `Certificate used for authorization for the RADSEC Server`,
							Optional:            true,
							Computed:            true,
							Sensitive:           true,
							PlanModifiers: []planmodifier.String{
								utils.NewSensitivePlanModifier(r.encryptionKey),
							},
						},
						"host": schema.StringAttribute{
							MarkdownDescription: `IP address (or FQDN) to which the APs will send RADIUS accounting messages`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: `Port on the RADIUS server that is listening for accounting messages`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"open_roaming_certificate_id": schema.Int64Attribute{
							MarkdownDescription: `The ID of the Open roaming Certificate attached to radius server`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"rad_sec_enabled": schema.BoolAttribute{
							MarkdownDescription: `Use RADSEC (TLS over TCP) to connect to this RADIUS accounting server. Requires radiusProxyEnabled.`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: `Shared key used to authenticate messages between the APs and RADIUS server`,
							Optional:            true,
							Computed:            true,
							Sensitive:           true,
							PlanModifiers: []planmodifier.String{
								utils.NewSensitivePlanModifier(r.encryptionKey),
							},
						},
					},
				},
			},
			"radius_attribute_for_group_policies": schema.StringAttribute{
				MarkdownDescription: `RADIUS attribute used to look up group policies`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Airespace-ACL-Name",
						"Aruba-User-Role",
						"Filter-ServerId",
						"Filter-Id",
						"Reply-Message",
					),
				},
			},
			"radius_authentication_nas_id": schema.StringAttribute{
				MarkdownDescription: `The template of the NAS identifier to be used for RADIUS authentication (ex. $NODE_MAC$:$VAP_NUM$).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_called_station_id": schema.StringAttribute{
				MarkdownDescription: `The template of the called station identifier to be used for RADIUS (ex. $NODE_MAC$:$VAP_NUM$).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_coa_enabled": schema.BoolAttribute{
				MarkdownDescription: `If true, Meraki devices will act as a RADIUS Dynamic Authorization Server and will respond to RADIUS Change-of-Authorization and Disconnect messages sent by the RADIUS server.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether RADIUS authentication is enabled`,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_fail_over_policy": schema.StringAttribute{
				MarkdownDescription: `Policy which determines how authentication requests should be handled in the event that all of the configured RADIUS servers are unreachable`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Allow access",
						"Deny access",
					),
				},
			},
			"radius_fallback_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether or not higher priority RADIUS servers should be retried after 60 seconds.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_guest_vlan_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether or not RADIUS Guest VLAN is enabled. This param is only valid if the authMode is 'open-with-radius' and addressing mode is not set to 'isolated' or 'nat' mode`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_guest_vlan_id": schema.Int64Attribute{
				MarkdownDescription: `VLAN ID of the RADIUS Guest VLAN. This param is only valid if the authMode is 'open-with-radius' and addressing mode is not set to 'isolated' or 'nat' mode`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_load_balancing_policy": schema.StringAttribute{
				MarkdownDescription: `Policy which determines which RADIUS server will be contacted first in an authentication attempt, and the ordering of any necessary retry attempts`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Round robin",
						"Strict priority order",
					),
				},
			},
			"radius_override": schema.BoolAttribute{
				MarkdownDescription: `If true, the RADIUS response can override VLAN tag. This is not valid when ipAssignmentMode is 'NAT mode'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_proxy_enabled": schema.BoolAttribute{
				MarkdownDescription: `If true, Meraki devices will proxy RADIUS messages through the Meraki cloud to the configured RADIUS auth and accounting servers.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_server_attempts_limit": schema.Int64Attribute{
				MarkdownDescription: `The maximum number of transmit attempts after which a RADIUS server is failed over (must be between 1-5).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_server_timeout": schema.Int64Attribute{
				MarkdownDescription: `The amount of time for which a RADIUS client waits for a reply from the RADIUS server (must be between 1-10 seconds).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_servers": schema.ListNestedAttribute{
				MarkdownDescription: `The RADIUS 802.1X servers to be used for authentication. This param is only valid if the authMode is 'open-with-radius', '8021x-radius' or 'ipsk-with-radius'`,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{

						"ca_certificate": schema.StringAttribute{
							MarkdownDescription: `Certificate used for authorization for the RADSEC Server`,
							Optional:            true,
							Computed:            true,
							Sensitive:           true,
							PlanModifiers: []planmodifier.String{
								utils.NewSensitivePlanModifier(r.encryptionKey),
							},
						},
						/*
							"server_id": schema.StringAttribute{  // not in api spec and changes all the time
									MarkdownDescription: `ServerId of your RADIUS server`,
									Computed:            true,
									PlanModifiers: []planmodifier.String{
										utils.NewSuppressDiffServerIDModifier(),
									},
								},
						*/
						"host": schema.StringAttribute{
							MarkdownDescription: `IP address of your RADIUS server`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"open_roaming_certificate_id": schema.Int64Attribute{
							MarkdownDescription: `The ID of the Openroaming Certificate attached to radius server.`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: `UDP port the RADIUS server listens on for Access-requests`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"rad_sec_enabled": schema.BoolAttribute{
							MarkdownDescription: `Use RADSEC (TLS over TCP) to connect to this RADIUS server. Requires radiusProxyEnabled.`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: `RADIUS client shared secret`,
							Optional:            true,
							Computed:            true,
							Sensitive:           true,
							PlanModifiers: []planmodifier.String{
								utils.NewSensitivePlanModifier(r.encryptionKey),
							},
						},
					},
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_testing_enabled": schema.BoolAttribute{
				MarkdownDescription: `If true, Meraki devices will periodically send Access-Request messages to configured RADIUS servers using identity 'meraki_8021x_test' to ensure that the RADIUS servers are reachable.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"secondary_concentrator_network_id": schema.StringAttribute{
				MarkdownDescription: `The secondary concentrator to use when the ipAssignmentMode is 'VPN'. If configured, the APs will switch to using this concentrator if the primary concentrator is unreachable. This param is optional. ('disabled' represents no secondary concentrator.)`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"speed_burst": schema.SingleNestedAttribute{
				MarkdownDescription: `The SpeedBurst setting for this SSID'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: `Boolean indicating whether or not to allow users to temporarily exceed the bandwidth limit for short periods while still keeping them under the bandwidth limit over time.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"splash_guest_sponsor_domains": schema.ListAttribute{
				MarkdownDescription: `Array of valid sponsor email domains for sponsored guest splash type.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},

				ElementType: types.StringType,
			},
			"splash_page": schema.StringAttribute{
				MarkdownDescription: `The type of splash page for the SSID`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Billing",
						"Cisco ISE",
						"Click-through splash page",
						"Facebook Wi-Fi",
						"Google Apps domain",
						"Google OAuth",
						"None",
						"Password-protected with Active Directory",
						"Password-protected with LDAP",
						"Password-protected with Meraki RADIUS",
						"Password-protected with custom RADIUS",
						"SMS authentication",
						"Sponsored guest",
						"Systems Manager Sentry",
					),
				},
			},
			"splash_timeout": schema.StringAttribute{
				MarkdownDescription: `Splash page timeout`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"ssid_admin_accessible": schema.BoolAttribute{
				MarkdownDescription: `SSID Administrator access status`,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"use_vlan_tagging": schema.BoolAttribute{
				MarkdownDescription: `Whether or not traffic should be directed to use specific VLANs. This param is only valid if the ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"visible": schema.BoolAttribute{
				MarkdownDescription: `Whether the SSID is advertised or hidden by the AP`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"vlan_id": schema.Int64Attribute{
				MarkdownDescription: `The VLAN ID used for VLAN tagging. This param is only valid when the ipAssignmentMode is 'Layer 3 roaming with a concentrator' or 'VPN'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"walled_garden_enabled": schema.BoolAttribute{
				MarkdownDescription: `Allow users to access a configurable list of IP ranges prior to sign-on`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"walled_garden_ranges": schema.ListAttribute{
				MarkdownDescription: `Domain names and IP address ranges available in Walled Garden mode`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},

				ElementType: types.StringType,
			},
			"wpa_encryption_mode": schema.StringAttribute{
				MarkdownDescription: `The types of WPA encryption`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"WPA1 and WPA2",
						"WPA1 only",
						"WPA2 only",
						"WPA3 192-bit Security",
						"WPA3 Transition Mode",
						"WPA3 only",
					),
				},
			},
		},
	}
}

func (r *NetworksWirelessSsidsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Retrieve the encryption key and client from the provider configuration
	client, ok := req.ProviderData.(*openApiClient.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Client Type",
			fmt.Sprintf("Expected *openApiClient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client

	// Since we are passing only the client directly for resources, we need to handle the encryption key separately
	// Retrieve the encryption key from the context
	encryptionKey, ok := ctx.Value("encryption_key").(string)
	if ok {
		r.encryptionKey = encryptionKey
	} else {
		r.encryptionKey = ""
	}

}

// Create creates the resource and sets the initial Terraform state.
func (r *NetworksWirelessSsidsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state, plan NetworksWirelessSsidResourceModel

	// Read the Terraform configuration into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request payload
	payload, payloadDiags := updateNetworksWirelessSsidsResourcePayload(ctx, &plan)
	if payloadDiags.HasError() {
		tflog.Error(ctx, "Failed to create resource payload", map[string]interface{}{
			"error": payloadDiags,
		})
		resp.Diagnostics.AddError(
			"Error creating ssid payload",
			fmt.Sprintf("Unexpected error: %s", payloadDiags),
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
		return inline, respHttp, err
	})

	// Capture the response body for logging
	var responseBody string
	if httpResp != nil && httpResp.Body != nil {
		bodyBytes, readErr := io.ReadAll(httpResp.Body)
		if readErr == nil {
			responseBody = string(bodyBytes)
		}
		// Reset the response body so it can be read again later if necessary
		httpResp.Body = io.NopCloser(io.NopCloser(bytes.NewBuffer(bodyBytes)))
	}

	// Check if the error matches a specific condition
	if err != nil {
		// Terminate early if specific error condition is met
		if strings.Contains(responseBody, "Open Roaming certificate 0 not found") {
			tflog.Error(ctx, "Terminating early due to specific error condition", map[string]interface{}{
				"error":        err.Error(),
				"responseBody": responseBody,
			})
			resp.Diagnostics.AddError(
				"HTTP Call Failed",
				fmt.Sprintf("Details: %s", responseBody),
			)
			return
		}

		// Check for the specific unmarshalling error
		if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
			tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
		} else {
			tflog.Error(ctx, "HTTP Call Failed", map[string]interface{}{
				"error":        err.Error(),
				"responseBody": responseBody,
			})
			resp.Diagnostics.AddError(
				"HTTP Call Failed",
				fmt.Sprintf("Details: %s", err.Error()),
			)
		}
		return
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%d", httpResp.StatusCode),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, &state, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *NetworksWirelessSsidsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state, plan NetworksWirelessSsidResourceModel

	// Read Terraform prior state into the model
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.GetNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).Execute()
		if err != nil {
			// Check for specific error
			if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
				err = nil
			}

			// Check for specific error
			if strings.Contains(err.Error(), "Open Roaming certificate 0 not found") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: Open Roaming certificate 0 not found")
				err = nil
			}
		}
		return inline, respHttp, err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
		return
	}

	// Ensure inlineResp and httpResp are not nil before using them
	if inlineResp == nil {
		resp.Diagnostics.AddError(
			"Received nil response",
			"Expected a valid response but received nil",
		)
		return
	}

	if httpResp == nil {
		resp.Diagnostics.AddError(
			"Received nil HTTP response",
			"Expected a valid HTTP response but received nil",
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		responseBody, _ := utils.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			utils.NewHttpDiagnostics(httpResp, responseBody),
		)
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, &state, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworksWirelessSsidsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan NetworksWirelessSsidResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request payload
	payload, payloadDiags := updateNetworksWirelessSsidsResourcePayload(ctx, &plan)
	if payloadDiags.HasError() {
		tflog.Error(ctx, "Failed to create resource payload", map[string]interface{}{
			"error": payloadDiags,
		})
		resp.Diagnostics.AddError(
			"Error creating ssid payload",
			fmt.Sprintf("Unexpected error: %s", payloadDiags),
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
		if err != nil {
			// Check for specific error
			if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
				err = nil
			}

			// Check for specific error
			if strings.Contains(err.Error(), "Open Roaming certificate 0 not found") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: Open Roaming certificate 0 not found")
				err = nil
			}

		}
		return inline, respHttp, err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Update Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := utils.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			utils.NewHttpDiagnostics(httpResp, responseBody),
		)
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, &state, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *NetworksWirelessSsidsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *NetworksWirelessSsidResourceModel
	// Read Terraform plan state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkWirelessSsidRequest()
	payload.SetEnabled(false)
	payload.SetName("")
	payload.SetAuthMode("open")
	payload.SetVlanId(1)

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	_, httpResp, err := utils.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), state.NetworkId.ValueString(), fmt.Sprint(state.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
		if err != nil {
			// Check for specific error
			if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
				err = nil
			}

			// Check for specific error
			if strings.Contains(err.Error(), "Open Roaming certificate 0 not found") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: Open Roaming certificate 0 not found")
				err = nil
			}

		}
		return inline, respHttp, err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := utils.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			utils.NewHttpDiagnostics(httpResp, responseBody),
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksWirelessSsidsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, number. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	i, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		panic(err)
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("number"), i)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
