package ssids

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

func NetworksWirelessSsidPayloadDot11w(input types.Object) (*openApiClient.UpdateNetworkApplianceSsidRequestDot11w, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var dot11wObject Dot11w

	err := input.As(context.Background(), &dot11wObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	return &openApiClient.UpdateNetworkApplianceSsidRequestDot11w{
		Enabled:  dot11wObject.Enabled.ValueBoolPointer(),
		Required: dot11wObject.Required.ValueBoolPointer(),
	}, diags
}

func NetworksWirelessSsidPayloadDot11r(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestDot11r, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var dot11rObject Dot11r

	err := input.As(context.Background(), &dot11rObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestDot11r{
		Enabled:  dot11rObject.Enabled.ValueBoolPointer(),
		Adaptive: dot11rObject.Adaptive.ValueBoolPointer(),
	}, diags
}

func NetworksWirelessSsidPayloadOauth(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestOauth, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var oauthObject OAuth

	err := input.As(context.Background(), &oauthObject, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty: true,
	})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	var allowedDomains []string
	allowedDomainsList := oauthObject.AllowedDomains.Elements()

	for _, domain := range allowedDomainsList {
		allowedDomains = append(allowedDomains, domain.String())
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestOauth{
		AllowedDomains: allowedDomains,
	}, diags
}

func NetworksWirelessSsidPayloadLocalRadius(input types.Object) (openApiClient.UpdateNetworkWirelessSsidRequestLocalRadius, diag.Diagnostics) {
	var result openApiClient.UpdateNetworkWirelessSsidRequestLocalRadius
	if input.IsNull() || input.IsUnknown() {
		return result, nil
	}

	var diags diag.Diagnostics

	// Unmarshalling input to LocalRadius struct
	var localRadius LocalRadius
	err := input.As(context.Background(), &localRadius, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, diag.NewErrorDiagnostic("Error converting LocalRadius", fmt.Sprintf("%s", err.Errors())))
		return result, diags
	}

	// CacheTimeout
	if !localRadius.CacheTimeout.IsUnknown() && !localRadius.CacheTimeout.IsNull() {
		cacheTimeout := int32(localRadius.CacheTimeout.ValueInt64())
		result.SetCacheTimeout(cacheTimeout)
	}

	// PasswordAuthentication
	if !localRadius.PasswordAuthentication.IsUnknown() && !localRadius.PasswordAuthentication.IsNull() {
		var passwordAuth PasswordAuthentication
		err := localRadius.PasswordAuthentication.As(context.Background(), &passwordAuth, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting PasswordAuthentication", fmt.Sprintf("%s", err.Errors())))
		} else {
			var passwordAuthentication openApiClient.UpdateNetworkWirelessSsidRequestLocalRadiusPasswordAuthentication
			passwordAuthentication.SetEnabled(passwordAuth.Enabled.ValueBool())
			result.SetPasswordAuthentication(passwordAuthentication)
		}
	}

	// CertificateAuthentication
	if !localRadius.CertificateAuthentication.IsNull() && !localRadius.CertificateAuthentication.IsUnknown() {
		var certificateAuthentication CertificateAuthentication
		err := localRadius.CertificateAuthentication.As(context.Background(), &certificateAuthentication, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags = append(diags, err.Errors()...)
		}
		var clientRootCaCertificate openApiClient.UpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthenticationClientRootCaCertificate
		if !certificateAuthentication.ClientRootCaCertificate.IsNull() {
			var clientRootCaCert CaCertificate
			err := certificateAuthentication.ClientRootCaCertificate.As(context.Background(), &clientRootCaCert, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags = append(diags, err.Errors()...)
			}
			clientRootCaCertificate.SetContents(clientRootCaCert.Contents.ValueString())
		}
		var certAuth openApiClient.UpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthentication
		certAuth.SetEnabled(certificateAuthentication.Enabled.ValueBool())
		certAuth.SetUseLdap(certificateAuthentication.UseLdap.ValueBool())
		certAuth.SetUseOcsp(certificateAuthentication.UseOcsp.ValueBool())
		certAuth.SetOcspResponderUrl(certificateAuthentication.OcspResponderUrl.ValueString())
		certAuth.SetClientRootCaCertificate(clientRootCaCertificate)
		result.SetCertificateAuthentication(certAuth)
	}

	return result, diags
}

func NetworksWirelessSsidPayloadActiveDirectory(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectory, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Unmarshalling input to ActiveDirectory struct
	var activeDirectoryObject ActiveDirectory
	err := input.As(context.Background(), &activeDirectoryObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	// Processing servers
	var activeDirectoryServers []ActiveDirectoryServer
	err = activeDirectoryObject.Servers.ElementsAs(context.Background(), &activeDirectoryServers, true)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic("Error converting ActiveDirectory Servers", fmt.Sprintf("%s", err.Errors())))
	}

	var activeDirectoryServersArray []openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner
	for _, svr := range activeDirectoryServers {
		var activeDirectoryServer openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner
		activeDirectoryServer.SetHost(svr.Host.ValueString())
		activeDirectoryServer.SetPort(int32(svr.Port.ValueInt64()))
		activeDirectoryServersArray = append(activeDirectoryServersArray, activeDirectoryServer)
	}

	// Processing credentials
	var credentialsObject AdCredentials
	err = activeDirectoryObject.Credentials.As(context.Background(), &credentialsObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectory{
		Servers: activeDirectoryServersArray,
		Credentials: &openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryCredentials{
			LogonName: credentialsObject.LoginName.ValueStringPointer(),
			Password:  credentialsObject.Password.ValueStringPointer(),
		},
	}, diags
}

func NetworksWirelessSsidPayloadRadiusServers(ctx context.Context, input types.List) ([]openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var servers []openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner
	var radiusServers []RadiusServer

	err := input.ElementsAs(ctx, &radiusServers, true)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	// Retrieve the encryption key from the context
	encryptionKey, ok := ctx.Value("encryption_key").(string)
	if !ok {
		// If encryption key is not available, log a warning and proceed without decryption
		tflog.Warn(ctx, "The encryption key is not available in the context, proceeding without decryption")
	}

	for _, server := range radiusServers {

		var serverPayload openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner

		// Host
		serverPayload.SetHost(server.Host.ValueString())

		// Port
		port, err := utils.Int32Pointer(server.Port.ValueInt64())
		if err != nil {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting Port", fmt.Sprintf("%s", err.Errors())))
		}
		serverPayload.SetPort(*port)

		// Secret
		if encryptionKey != "" {
			decryptedSecret, err := utils.Decrypt(encryptionKey, server.Secret.ValueString())
			if err != nil {
				diags = append(diags, diag.NewErrorDiagnostic("Error Decrypting Secret", err.Error()))
			} else {
				serverPayload.SetSecret(decryptedSecret)
			}
		} else {
			serverPayload.SetSecret(server.Secret.ValueString())
		}

		// RadSecEnabled
		serverPayload.SetRadsecEnabled(server.RadSecEnabled.ValueBool())

		// OpenRoamingCertificateId
		openRoamingCertificateId, err := utils.Int32Pointer(server.OpenRoamingCertificateID.ValueInt64())
		if err != nil {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting OpenRoamingCertificateID", fmt.Sprintf("%s", err.Errors())))
		}

		if *openRoamingCertificateId == 0 {
			openRoamingCertificateId = nil
		}

		if openRoamingCertificateId != nil {
			serverPayload.SetOpenRoamingCertificateId(*openRoamingCertificateId)
		}

		// CaCertificate
		serverPayload.SetCaCertificate(server.CaCertificate.ValueString())

		servers = append(servers, serverPayload)
	}

	return servers, diags
}

func NetworksWirelessSsidPayloadRadiusAccountingServers(ctx context.Context, input types.List) ([]openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner, diag.Diagnostics) {
	var diags diag.Diagnostics
	var servers []openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner
	var radiusServers []RadiusServer

	err := input.ElementsAs(ctx, &radiusServers, true)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	// Retrieve the encryption key from the context
	encryptionKey, ok := ctx.Value("encryption_key").(string)
	if !ok {
		// If encryption key is not available, log a warning and proceed without decryption
		tflog.Warn(ctx, "The encryption key is not available in the context, proceeding without decryption")
	}

	for _, server := range radiusServers {

		var serverPayload openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner

		// Host
		serverPayload.SetHost(server.Host.ValueString())

		// Port
		port, err := utils.Int32Pointer(server.Port.ValueInt64())
		if err != nil {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting Port", fmt.Sprintf("%s", err.Errors())))
		}
		serverPayload.SetPort(*port)

		// Secret
		if encryptionKey != "" {
			decryptedSecret, err := utils.Decrypt(encryptionKey, server.Secret.ValueString())
			if err != nil {
				diags = append(diags, diag.NewErrorDiagnostic("Error Decrypting Secret", err.Error()))
			} else {
				serverPayload.SetSecret(decryptedSecret)
			}
		} else {
			serverPayload.SetSecret(server.Secret.ValueString())
		}

		// RadSecEnabled
		serverPayload.SetRadsecEnabled(server.RadSecEnabled.ValueBool())

		// CaCertificate
		serverPayload.SetCaCertificate(server.CaCertificate.ValueString())

		servers = append(servers, serverPayload)
	}

	return servers, diags
}

func NetworksWirelessSsidPayloadApTagsAndVlanIds(input types.List) ([]openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var tagsAndVlans []openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner
	var tagsAndVlansList []ApTagsAndVlanID

	err := input.ElementsAs(context.Background(), &tagsAndVlansList, true)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	for _, tagAndVlan := range tagsAndVlansList {
		vlanId, err := utils.Int32Pointer(tagAndVlan.VlanId.ValueInt64())
		if err != nil {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting VlanID", fmt.Sprintf("%s", err.Errors())))
		}

		var tags []string
		for _, tag := range tagAndVlan.Tags.Elements() {
			tags = append(tags, tag.String())
		}

		tagsAndVlans = append(tagsAndVlans, openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner{
			Tags:   tags,
			VlanId: vlanId,
		})
	}

	return tagsAndVlans, diags
}

func NetworksWirelessSsidPayloadGre(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestGre, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics

	var gre GRE

	err := input.As(context.Background(), &gre, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	var concentrator GreConcentrator

	err = gre.Concentrator.As(context.Background(), &concentrator, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	key := int32(gre.Key.ValueInt64())

	return &openApiClient.UpdateNetworkWirelessSsidRequestGre{
		Key: &key,
		Concentrator: &openApiClient.UpdateNetworkWirelessSsidRequestGreConcentrator{
			Host: concentrator.Host.ValueString(),
		},
	}, diags
}

func NetworksWirelessSsidPayloadDnsRewrite(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestDnsRewrite, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var dnsRewriteObject DnsRewrite

	err := input.As(context.Background(), &dnsRewriteObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	var dnsCustomNameservers []string
	dnsCustomNameserversList := dnsRewriteObject.DnsCustomNameservers.Elements()

	for _, dns := range dnsCustomNameserversList {
		dnsCustomNameservers = append(dnsCustomNameservers, dns.String())
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestDnsRewrite{
		Enabled:              dnsRewriteObject.Enabled.ValueBoolPointer(),
		DnsCustomNameservers: dnsCustomNameservers,
	}, diags
}

func NetworksWirelessSsidPayloadSpeedBurst(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestSpeedBurst, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var speedBurst SpeedBurst

	err := input.As(context.Background(), &speedBurst, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestSpeedBurst{
		Enabled: speedBurst.Enabled.ValueBoolPointer(),
	}, diags
}

func NetworksWirelessSsidPayloadNamedVlans(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestNamedVlans, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var namedVlansObject NamedVlans

	err := input.As(context.Background(), &namedVlansObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	// Tagging
	var tagging Tagging
	err = namedVlansObject.Tagging.As(context.Background(), &tagging, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	// Radius
	var radius Radius
	err = namedVlansObject.Radius.As(context.Background(), &radius, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	byApTags, err := NetworksWirelessSsidPayloadByApTags(tagging.ByApTags)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	// guestVlan
	guestVlan, err := NetworksWirelessSsidPayloadRadiusGuestVlan(radius.GuestVlan)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestNamedVlans{
		Tagging: &openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTagging{
			Enabled:         tagging.Enabled.ValueBoolPointer(),
			DefaultVlanName: tagging.DefaultVlanName.ValueStringPointer(),
			ByApTags:        byApTags,
		},
		Radius: &openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansRadius{
			GuestVlan: &guestVlan,
		},
	}, diags
}

func NetworksWirelessSsidPayloadLdap(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestLdap, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Unmarshalling input to LDAP struct
	var ldapObject LDAP
	err := input.As(context.Background(), &ldapObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	// Processing servers
	var servers []LdapServer
	err = ldapObject.Servers.ElementsAs(context.Background(), &servers, true)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic("Error converting Servers", fmt.Sprintf("%s", err.Errors())))
	}

	var serversArray []openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner
	for _, svr := range servers {
		var server openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner
		server.SetHost(svr.Host.ValueString())
		server.SetPort(int32(svr.Port.ValueInt64()))
		serversArray = append(serversArray, server)
	}

	// Processing credentials
	var creds LdapCredentials
	err = ldapObject.Credentials.As(context.Background(), &creds, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	// Processing server CA certificate
	var serverCaCertificate CaCertificate
	err = ldapObject.ServerCaCertificate.As(context.Background(), &serverCaCertificate, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestLdap{
		Servers: serversArray,
		Credentials: &openApiClient.UpdateNetworkWirelessSsidRequestLdapCredentials{
			DistinguishedName: creds.DistinguishedName.ValueStringPointer(),
			Password:          creds.Password.ValueStringPointer(),
		},
		BaseDistinguishedName: ldapObject.BaseDistinguishedName.ValueStringPointer(),
		ServerCaCertificate: &openApiClient.UpdateNetworkWirelessSsidRequestLdapServerCaCertificate{
			Contents: serverCaCertificate.Contents.ValueStringPointer(),
		},
	}, diags
}

func NetworksWirelessSsidPayloadByApTags(input types.List) ([]openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var byApTags []openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner

	var byApTagsList []ByApTag
	err := input.ElementsAs(context.Background(), &byApTagsList, false)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	for _, byApTag := range byApTagsList {
		var tags []string
		for _, tag := range byApTag.Tags.Elements() {
			tags = append(tags, tag.String())
		}

		byApTags = append(byApTags, openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner{
			Tags:     tags,
			VlanName: byApTag.VlanName.ValueStringPointer(),
		})
	}
	return byApTags, diags
}

func NetworksWirelessSsidPayloadRadiusGuestVlan(input types.Object) (openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansRadiusGuestVlan, diag.Diagnostics) {
	var diags diag.Diagnostics
	var guestVlans openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansRadiusGuestVlan

	var data RadiusGuestVlan

	err := input.As(context.Background(), data, basetypes.ObjectAsOptions{})
	if err.HasError() {
		return guestVlans, err
	}

	// enabled
	guestVlans.SetEnabled(data.Enabled.ValueBool())

	// name
	guestVlans.SetName(data.Name.ValueString())

	return guestVlans, diags
}

// update terraform state funcs //

func networksWirelessSsidAdminSplashUrl(data *openApiClient.GetNetworkWirelessSsids200ResponseInner) (basetypes.StringValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	result := types.StringNull()
	adminSplashUrl, ok := data.GetAdminSplashUrlOk()
	if ok {
		result = types.StringValue(*adminSplashUrl)
		return result, diags
	}

	return result, diags
}

func NetworksWirelessSsidStateRadiusServers(ctx context.Context, plan NetworksWirelessSsidResourceModel, httpResp map[string]interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var radiusServers []RadiusServer

	radiusServerAttr := map[string]attr.Type{
		"host": types.StringType,
		// "server_id":                   types.StringType,  // not in api spec and changes all the time
		"port":                        types.Int64Type,
		"secret":                      types.StringType,
		"rad_sec_enabled":             types.BoolType,
		"open_roaming_certificate_id": types.Int64Type,
		"ca_certificate":              types.StringType,
	}

	// Retrieve the encryption key from the context
	encryptionKey, ok := ctx.Value("encryption_key").(string)
	if !ok {
		// If encryption key is not available, log a warning and proceed without encryption
		tflog.Warn(ctx, "The encryption key is not available in the context, proceeding without encryption")
	}

	// Process the response from the API
	if radiusServersResp, ok := httpResp["radiusServers"].([]interface{}); ok {
		for _, rsr := range radiusServersResp {
			if rs, ok := rsr.(map[string]interface{}); ok {
				var radiusServerResp RadiusServer

				// Extract attributes from the response
				radiusServerResp.Host, _ = utils.ExtractStringAttr(rs, "host")

				portFloat, _ := utils.ExtractFloat64Attr(rs, "port")
				if !portFloat.IsNull() && !portFloat.IsUnknown() {
					radiusServerResp.Port = types.Int64Value(int64(portFloat.ValueFloat64()))
				} else {
					radiusServerResp.Port = types.Int64Null()
				}

				// radiusServerResp.ServerId, _ = utils.ExtractStringAttr(rs, "id")  // not in api spec and changes all the time

				radiusServerResp.OpenRoamingCertificateID, _ = utils.ExtractInt32Attr(rs, "openRoamingCertificateId")
				radiusServerResp.CaCertificate, _ = utils.ExtractStringAttr(rs, "caCertificate")
				radiusServerResp.RadSecEnabled, _ = utils.ExtractBoolAttr(rs, "radsecEnabled")

				// Secret not returned by API, will be set from the plan
				radiusServerResp.Secret = types.StringNull()

				radiusServers = append(radiusServers, radiusServerResp)
			}
		}
	}

	var newRadiusServers []attr.Value

	var radiusServersPlan []RadiusServer
	err := plan.RadiusServers.ElementsAs(ctx, &radiusServersPlan, true)
	if err.HasError() {
		diags.Append(err...)
	}

	// Process the plan to extract secret then handle encryption
	for i, radiusServerPlan := range radiusServersPlan {

		if i < len(radiusServers) {

			if !radiusServerPlan.Secret.IsNull() && !radiusServerPlan.Secret.IsUnknown() {
				// Extract the secret from the plan
				radiusServers[i].Secret = radiusServerPlan.Secret
			} else {
				radiusServers[i].Secret = types.StringNull()
			}

			// Encrypt the secret if the encryption key is available
			if encryptionKey != "" {
				encryptedSecret, err := utils.Encrypt(encryptionKey, radiusServerPlan.Secret.ValueString())
				if err != nil {
					diags.Append(diag.NewErrorDiagnostic("Error Encrypting Secret", err.Error()))
				} else {
					radiusServers[i].Secret = types.StringValue(encryptedSecret)
				}
			}

			// Encrypt the ca_certificate if the encryption key is available
			if encryptionKey != "" {
				encryptedCaCertificate, err := utils.Encrypt(encryptionKey, radiusServers[i].CaCertificate.ValueString())
				if err != nil {
					diags.Append(diag.NewErrorDiagnostic("Error Encrypting CA Certificate", err.Error()))
				} else {
					radiusServers[i].CaCertificate = types.StringValue(encryptedCaCertificate)
				}
			}

			// Convert the RadiusServer object to a types.ObjectValue
			radiusServerObject, radiusServerObjectErr := types.ObjectValueFrom(ctx, radiusServerAttr, radiusServers[i])
			if radiusServerObjectErr.HasError() {
				diags.Append(radiusServerObjectErr...)
				continue
			}

			newRadiusServers = append(newRadiusServers, radiusServerObject)
		}
	}

	// Return a populated or empty list instead of a null value
	if newRadiusServers != nil && len(newRadiusServers) == 0 {
		newRadiusServersList, err := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: radiusServerAttr}, []attr.Value{})
		if err.HasError() {
			diags.Append(err...)
		}
		return newRadiusServersList, diags
	}

	newRadiusServersList, err := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: radiusServerAttr}, newRadiusServers)
	if err.HasError() {
		diags.Append(err...)
	}

	return newRadiusServersList, diags
}

func NetworksWirelessSsidStateRadiusAccountingServers(ctx context.Context, plan NetworksWirelessSsidResourceModel, httpResp map[string]interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var radiusServers []RadiusServer

	radiusServerAttr := map[string]attr.Type{
		"host": types.StringType,
		// "server_id":                   types.StringType,  // not in api spec and changes all the time
		"port":                        types.Int64Type,
		"secret":                      types.StringType,
		"rad_sec_enabled":             types.BoolType,
		"open_roaming_certificate_id": types.Int64Type,
		"ca_certificate":              types.StringType,
	}

	// Retrieve the encryption key from the context
	encryptionKey, ok := ctx.Value("encryption_key").(string)
	if !ok {
		// If encryption key is not available, log a warning and proceed without encryption
		tflog.Warn(ctx, "The encryption key is not available in the context, proceeding without encryption")
	}

	// Process the response from the API
	if radiusServersResp, ok := httpResp["radiusAccountingServers"].([]interface{}); ok {
		for _, rsr := range radiusServersResp {
			if rs, ok := rsr.(map[string]interface{}); ok {
				var radiusServerResp RadiusServer

				// Extract attributes from the response
				radiusServerResp.Host, _ = utils.ExtractStringAttr(rs, "host")

				portFloat, _ := utils.ExtractFloat64Attr(rs, "port")
				if !portFloat.IsNull() && !portFloat.IsUnknown() {
					radiusServerResp.Port = types.Int64Value(int64(portFloat.ValueFloat64()))
				} else {
					radiusServerResp.Port = types.Int64Null()
				}

				// radiusServerResp.ServerId, _ = utils.ExtractStringAttr(rs, "id")  // not in api spec and changes all the time

				radiusServerResp.OpenRoamingCertificateID, _ = utils.ExtractInt32Attr(rs, "openRoamingCertificateId")
				radiusServerResp.CaCertificate, _ = utils.ExtractStringAttr(rs, "caCertificate")
				radiusServerResp.RadSecEnabled, _ = utils.ExtractBoolAttr(rs, "radsecEnabled")

				// Secret not returned by API, will be set from the plan
				radiusServerResp.Secret = types.StringNull()

				radiusServers = append(radiusServers, radiusServerResp)
			}
		}
	}

	var newRadiusServers []attr.Value

	var radiusServersPlan []RadiusServer
	err := plan.RadiusAccountingServers.ElementsAs(ctx, &radiusServersPlan, true)
	if err.HasError() {
		diags.Append(err...)
	}

	// Process the plan to extract secret then handle encryption
	for i, radiusServerPlan := range radiusServersPlan {

		if i < len(radiusServers) {

			if !radiusServerPlan.Secret.IsNull() && !radiusServerPlan.Secret.IsUnknown() {
				// Extract the secret from the plan
				radiusServers[i].Secret = radiusServerPlan.Secret
			} else {
				radiusServers[i].Secret = types.StringNull()
			}

			// Encrypt the secret if the encryption key is available
			if encryptionKey != "" {
				encryptedSecret, err := utils.Encrypt(encryptionKey, radiusServerPlan.Secret.ValueString())
				if err != nil {
					diags.Append(diag.NewErrorDiagnostic("Error Encrypting Secret", err.Error()))
				} else {
					radiusServers[i].Secret = types.StringValue(encryptedSecret)
				}
			}

			// Encrypt the ca_certificate if the encryption key is available
			if encryptionKey != "" {
				encryptedCaCertificate, err := utils.Encrypt(encryptionKey, radiusServers[i].CaCertificate.ValueString())
				if err != nil {
					diags.Append(diag.NewErrorDiagnostic("Error Encrypting CA Certificate", err.Error()))
				} else {
					radiusServers[i].CaCertificate = types.StringValue(encryptedCaCertificate)
				}
			}

			// Convert the RadiusServer object to a types.ObjectValue
			radiusServerObject, radiusServerObjectErr := types.ObjectValueFrom(ctx, radiusServerAttr, radiusServers[i])
			if radiusServerObjectErr.HasError() {
				diags.Append(radiusServerObjectErr...)
				continue
			}

			newRadiusServers = append(newRadiusServers, radiusServerObject)
		}
	}

	// Return a populated or empty list instead of a null value
	if newRadiusServers != nil && len(newRadiusServers) == 0 {
		newRadiusServersList, err := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: radiusServerAttr}, []attr.Value{})
		if err.HasError() {
			diags.Append(err...)
		}
		return newRadiusServersList, diags
	}

	newRadiusServersList, err := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: radiusServerAttr}, newRadiusServers)
	if err.HasError() {
		diags.Append(err...)
	}

	return newRadiusServersList, diags
}

func NetworksWirelessSsidStateDot11w(rawResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var dot11w Dot11w

	dot11wAttrs := map[string]attr.Type{
		"enabled":  types.BoolType,
		"required": types.BoolType,
	}

	if d, ok := rawResp["dot11w"].(map[string]interface{}); ok {

		// enabled
		enabled, err := utils.ExtractBoolAttr(d, "enabled")
		if diags.HasError() {
			diags.Append(err...)
		}
		dot11w.Enabled = enabled

		// required
		required, err := utils.ExtractBoolAttr(d, "required")
		if diags.HasError() {
			diags.Append(err...)
		}
		dot11w.Required = required

	} else {
		Dot11wNull := types.ObjectNull(dot11wAttrs)
		return Dot11wNull, diags
	}

	dot11wObj, err := types.ObjectValueFrom(context.Background(), dot11wAttrs, dot11w)
	if err.HasError() {
		diags.Append(err...)
	}

	return dot11wObj, diags
}

func NetworksWirelessSsidStateDot11r(rawResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var dot11r Dot11r

	dot11rAttrs := map[string]attr.Type{
		"enabled":  types.BoolType,
		"adaptive": types.BoolType,
	}

	if d, ok := rawResp["dot11r"].(map[string]interface{}); ok {

		// enabled
		enabled, err := utils.ExtractBoolAttr(d, "enabled")
		if diags.HasError() {
			diags.Append(err...)
		}
		dot11r.Enabled = enabled

		// adaptive
		adaptive, err := utils.ExtractBoolAttr(d, "adaptive")
		if diags.HasError() {
			diags.Append(err...)
		}
		dot11r.Adaptive = adaptive

	} else {
		dot11rNull := types.ObjectNull(dot11rAttrs)
		return dot11rNull, diags
	}

	outputObj, err := types.ObjectValueFrom(context.Background(), dot11rAttrs, dot11r)
	if err.HasError() {
		diags.Append(err...)
	}

	return outputObj, diags
}

func NetworksWirelessSsidStateOauth(rawResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var oauth OAuth

	oauthAttrs := map[string]attr.Type{
		"allowed_domains": types.ListType{ElemType: types.StringType},
	}

	//oauth
	if oa, ok := rawResp["oauth"].(map[string]interface{}); ok {

		// allowed domains
		if ad, ok := oa["allowed_domains"].([]string); ok {
			var allowedDomains []types.String
			for _, domain := range ad {
				allowedDomains = append(allowedDomains, types.StringValue(domain))
			}

			// returns a populated or empty list instead of a null value
			if allowedDomains != nil {
				allowedDomainsObj, err := types.ListValueFrom(context.Background(), types.StringType, allowedDomains)
				if err.HasError() {
					diags.Append(err...)
				}
				oauth.AllowedDomains = allowedDomainsObj
			} else {
				allowedDomainsObj, err := types.ListValueFrom(context.Background(), types.StringType, []attr.Value{})
				if err.HasError() {
					diags.Append(err...)
				}
				oauth.AllowedDomains = allowedDomainsObj
			}

		} else {
			allowedDomainsObjNull := types.ListNull(types.StringType)
			oauth.AllowedDomains = allowedDomainsObjNull
		}

	} else {
		oauthObjNull := types.ObjectNull(oauthAttrs)
		return oauthObjNull, diags
	}

	oauthObj, err := types.ObjectValueFrom(context.Background(), oauthAttrs, oauth)
	if err.HasError() {
		diags.Append(err...)
	}

	return oauthObj, diags
}

func NetworksWirelessSsidStateLocalRadius(rawResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	passwordAuthenticationAttrs := map[string]attr.Type{
		"enabled": types.BoolType,
	}

	contentsAttr := map[string]attr.Type{
		"contents": types.StringType,
	}

	certificateAuthenticationAttr := map[string]attr.Type{
		"enabled":                    types.BoolType,
		"use_ldap":                   types.BoolType,
		"use_ocsp":                   types.BoolType,
		"ocsp_responder_url":         types.StringType,
		"client_root_ca_certificate": types.ObjectType{AttrTypes: contentsAttr},
	}

	localRadiusAttrs := map[string]attr.Type{
		"cache_timeout":              types.Int64Type,
		"password_authentication":    types.ObjectType{AttrTypes: passwordAuthenticationAttrs},
		"certificate_authentication": types.ObjectType{AttrTypes: certificateAuthenticationAttr},
	}

	var localRadius LocalRadius

	// cacheTimeout
	cacheTimeout, err := utils.ExtractInt64Attr(rawResp, "cacheTimeOut")
	if diags.HasError() {
		diags.Append(err...)
	}
	localRadius.CacheTimeout = cacheTimeout

	// Password Authentication
	if pa, ok := rawResp["passwordAuthentication"].(map[string]interface{}); ok {
		var passwordAuth PasswordAuthentication

		// enabled
		enabled, err := utils.ExtractBoolAttr(pa, "enabled")
		if diags.HasError() {
			diags.Append(err...)
		}
		passwordAuth.Enabled = enabled

	} else {
		passwordAuthenticationObjNull := types.ObjectNull(passwordAuthenticationAttrs)
		localRadius.PasswordAuthentication = passwordAuthenticationObjNull
	}

	// certificateAuthentication
	if ca, ok := rawResp["certificateAuthentication"].(map[string]interface{}); ok {
		var certificateAuthentication CertificateAuthentication

		//   Enabled
		if _, ok := ca["enabled"].(types.Bool); ok {

			caEnabled, err := utils.ExtractBoolAttr(ca, "enabled")
			if diags.HasError() {
				diags.Append(err...)
			}

			certificateAuthentication.Enabled = caEnabled

		} else {
			caEnabledNull := types.BoolNull()
			certificateAuthentication.Enabled = caEnabledNull
		}

		//    UseLdap
		if _, ok := ca["useLdap"].(types.Bool); ok {

			useLdap, err := utils.ExtractBoolAttr(ca, "useLdap")
			if diags.HasError() {
				diags.Append(err...)
			}

			certificateAuthentication.UseLdap = useLdap

		} else {
			useLdapNull := types.BoolNull()
			certificateAuthentication.UseLdap = useLdapNull
		}

		//    UseOcsp
		if _, ok := ca["useOcsp"].(types.Bool); ok {

			useOcsp, err := utils.ExtractBoolAttr(ca, "useOcsp")
			if diags.HasError() {
				diags.Append(err...)
			}

			certificateAuthentication.UseOcsp = useOcsp

		} else {
			useOcspNull := types.BoolNull()
			certificateAuthentication.UseOcsp = useOcspNull
		}

		//    OcspResponderUrl
		if _, ok := ca["ocspResponderUrl"].(types.String); ok {

			ocspResponderUrl, err := utils.ExtractStringAttr(ca, "ocspResponderUrl")
			if diags.HasError() {
				diags.Append(err...)
			}

			certificateAuthentication.OcspResponderUrl = ocspResponderUrl

		} else {
			ocspResponderUrlNull := types.StringNull()
			certificateAuthentication.OcspResponderUrl = ocspResponderUrlNull
		}

		//    ClientRootCaCertificate
		if crca, ok := rawResp["clientRootCaCertificate"].(map[string]interface{}); ok {
			var clientRootCaCertificate CaCertificate

			// Contents
			if _, ok := crca["contents"].(types.String); ok {

				contents, err := utils.ExtractStringAttr(ca, "contents")
				if diags.HasError() {
					diags.Append(err...)
				}

				clientRootCaCertificate.Contents = contents

			} else {
				contentsNull := types.StringNull()
				clientRootCaCertificate.Contents = contentsNull
			}
		}

		certificateAuthenticationObj, err := types.ObjectValueFrom(context.Background(), certificateAuthenticationAttr, certificateAuthentication)
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
		}

		localRadius.CertificateAuthentication = certificateAuthenticationObj

	} else {
		certificateAuthenticationObjNull := types.ObjectNull(certificateAuthenticationAttr)
		localRadius.CertificateAuthentication = certificateAuthenticationObjNull
	}

	outputObj, err := types.ObjectValueFrom(context.Background(), localRadiusAttrs, localRadius)
	if err.HasError() {
		diags.Append(err...)
	}

	return outputObj, diags
}

func NetworksWirelessSsidStateLdap(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var ldap LDAP

	serverAttr := map[string]attr.Type{
		"host": types.StringType,
		"port": types.Int64Type,
	}

	credentialsAttr := map[string]attr.Type{
		"distinguished_name": types.StringType,
		"password":           types.StringType,
	}

	contentsAttr := map[string]attr.Type{
		"contents": types.StringType,
	}

	ldapAttrs := map[string]attr.Type{
		"base_distinguished_name": types.StringType,
		"servers":                 types.ListType{ElemType: types.ObjectType{AttrTypes: serverAttr}},
		"credentials":             types.ObjectType{AttrTypes: credentialsAttr},
		"server_ca_certificate":   types.ObjectType{AttrTypes: contentsAttr},
	}

	if l, ok := httpResp["ldap"].(map[string]interface{}); ok {

		// baseDistinguishedName
		baseDistinguishedName, err := utils.ExtractStringAttr(l, "baseDistinguishedName")
		if err.HasError() {
			diags.AddError("baseDistinguishedName Attribute", fmt.Sprintf("%s", err.Errors()))
		}
		ldap.BaseDistinguishedName = baseDistinguishedName

		// credentials
		if credsMap, ok := l["credentials"].(map[string]interface{}); ok {
			var creds LdapCredentials

			// loginName
			DistinguishedNameObj, err := utils.ExtractStringAttr(credsMap, "DistinguishedName")
			if err.HasError() {
				diags.AddError("DistinguishedName Attribute", fmt.Sprintf("%s", err.Errors()))
			}
			creds.DistinguishedName = DistinguishedNameObj

			// Password
			passwordObj, err := utils.ExtractStringAttr(credsMap, "password")
			if err.HasError() {
				diags.AddError("password Attribute", fmt.Sprintf("%s", err.Errors()))
			}
			creds.Password = passwordObj

			credsObj, err := types.ObjectValueFrom(context.Background(), credentialsAttr, creds)
			if err.HasError() {
				diags.AddError("credentials object Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			ldap.Credentials = credsObj

		} else {
			credsNull := types.ObjectNull(credentialsAttr)
			ldap.Credentials = credsNull
		}

		// serverCaCertificate
		if serverCaCertificateMap, ok := l["serverCaCertificate"].(map[string]interface{}); ok {
			var serverCaCertificate LdapServerCaCertificate

			// contents
			contents, err := utils.ExtractStringAttr(serverCaCertificateMap, "contents")
			if err.HasError() {
				diags.AddError("contents Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			serverCaCertificate.Contents = contents

			ServerCaCertObj, err := types.ObjectValueFrom(context.Background(), contentsAttr, serverCaCertificate)
			if err.HasError() {
				diags.AddError("serverCaCertificate object Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			ldap.ServerCaCertificate = ServerCaCertObj

		} else {
			ServerCaCertObjNull := types.ObjectNull(contentsAttr)
			ldap.ServerCaCertificate = ServerCaCertObjNull
		}

		// servers
		if listMapArray, ok := l["servers"].([]map[string]interface{}); ok {

			var serversArray []types.Object

			for _, listMap := range listMapArray {
				var server ActiveDirectoryServer

				// host
				host, err := utils.ExtractStringAttr(listMap, "host")
				if err.HasError() {
					diags.AddError("host Attribute", fmt.Sprintf("%s", err.Errors()))
				}
				server.Host = host

				// port
				port, err := utils.ExtractInt64Attr(listMap, "port")
				if err.HasError() {
					diags.AddError("port Attribute", fmt.Sprintf("%s", err.Errors()))
				}
				server.Port = port

				serverObj, err := types.ObjectValueFrom(context.Background(), serverAttr, server)
				if err.HasError() {
					diags.AddError("server Attribute", fmt.Sprintf("%s", err.Errors()))
				}

				serversArray = append(serversArray, serverObj)
			}

			// returns a populated or empty list instead of a null value
			if serversArray != nil {
				// servers Array
				serversArrayObj, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: serverAttr}, serversArray)
				if err.HasError() {
					diags.AddError("servers array Attribute", fmt.Sprintf("%s", err.Errors()))
				}

				ldap.Servers = serversArrayObj
			} else {
				serversArrayObj, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: serverAttr}, []attr.Value{})
				if err.HasError() {
					diags.AddError("servers array Attribute", fmt.Sprintf("%s", err.Errors()))
				}

				ldap.Servers = serversArrayObj
			}

		} else {
			lisObjNull := types.ListNull(types.ObjectType{AttrTypes: serverAttr})
			ldap.Servers = lisObjNull

		}

	} else {
		ldapObjNull := types.ObjectNull(ldapAttrs)
		return ldapObjNull, diags
	}

	ldapObject, err := types.ObjectValueFrom(context.Background(), ldapAttrs, ldap)
	if err.HasError() {
		diags.AddError("ldap object Attribute", fmt.Sprintf("%s", err.Errors()))
	}

	return ldapObject, diags
}

func NetworksWirelessSsidStateActiveDirectory(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var activeDirectory ActiveDirectory

	serverAttr := map[string]attr.Type{
		"host": types.StringType,
		"port": types.Int64Type,
	}

	credentialsAttr := map[string]attr.Type{
		"login_name": types.StringType,
		"password":   types.StringType,
	}

	activeDirectoryAttrs := map[string]attr.Type{
		"servers": types.ListType{
			ElemType: types.ObjectType{AttrTypes: serverAttr},
		},
		"credentials": types.ObjectType{
			AttrTypes: credentialsAttr,
		},
	}

	if ad, ok := httpResp["activeDirectory"].(map[string]interface{}); ok {

		// servers
		if listMapArray, ok := ad["servers"].([]map[string]interface{}); ok {

			var serversArray []types.Object

			for _, listMap := range listMapArray {
				var server ActiveDirectoryServer

				// host
				host, err := utils.ExtractStringAttr(listMap, "host")
				if err.HasError() {
					diags.AddError("host Attribute", fmt.Sprintf("%s", err.Errors()))
				}
				server.Host = host

				// port
				port, err := utils.ExtractInt64Attr(listMap, "port")
				if err.HasError() {
					diags.AddError("port Attribute", fmt.Sprintf("%s", err.Errors()))
				}
				server.Port = port

				serverObj, err := types.ObjectValueFrom(context.Background(), serverAttr, server)
				if err.HasError() {
					diags.AddError("server Attribute", fmt.Sprintf("%s", err.Errors()))
				}

				serversArray = append(serversArray, serverObj)
			}

			// returns a populated or empty list instead of a null value
			if serversArray != nil {
				// servers Array
				serversArrayObj, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: serverAttr}, serversArray)
				if err.HasError() {
					diags.AddError("servers array Attribute", fmt.Sprintf("%s", err.Errors()))
				}

				activeDirectory.Servers = serversArrayObj
			} else {
				// servers Array
				serversArrayObj, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: serverAttr}, []attr.Value{})
				if err.HasError() {
					diags.AddError("servers array Attribute", fmt.Sprintf("%s", err.Errors()))
				}

				activeDirectory.Servers = serversArrayObj
			}

		} else {
			lisObjNull := types.ListNull(types.ObjectType{AttrTypes: serverAttr})
			activeDirectory.Servers = lisObjNull

		}

		// credentials
		if credsMap, ok := ad["credentials"].(map[string]interface{}); ok {
			var creds AdCredentials

			// loginName
			loginNameObj, err := utils.ExtractStringAttr(credsMap, "loginName")
			if err.HasError() {
				diags.AddError("loginName Attribute", fmt.Sprintf("%s", err.Errors()))
			}
			creds.LoginName = loginNameObj

			// Password
			passwordObj, err := utils.ExtractStringAttr(credsMap, "password")
			if err.HasError() {
				diags.AddError("password Attribute", fmt.Sprintf("%s", err.Errors()))
			}
			creds.Password = passwordObj

			credsObj, err := types.ObjectValueFrom(context.Background(), credentialsAttr, creds)
			if err.HasError() {
				diags.AddError("credentials object Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			activeDirectory.Credentials = credsObj

		} else {
			credsNull := types.ObjectNull(credentialsAttr)
			activeDirectory.Credentials = credsNull
		}

	} else {
		activeDirectoryObjNull := types.ObjectNull(activeDirectoryAttrs)
		return activeDirectoryObjNull, diags
	}

	activeDirectoryObj, err := types.ObjectValueFrom(context.Background(), activeDirectoryAttrs, activeDirectory)
	if err.HasError() {
		diags.AddError("Active Directory Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return activeDirectoryObj, diags
}

func NetworksWirelessSsidStateApTagsAndVlanIds(httpResp map[string]interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	apTagsAndVlanIdsAttr := map[string]attr.Type{
		"tags":    types.ListType{ElemType: types.StringType},
		"vlan_id": types.Int64Type,
	}

	apTagsAndVlanIdsAttrs := types.ObjectType{AttrTypes: apTagsAndVlanIdsAttr}

	apTagsAndVlanIdsList, err := utils.ExtractListAttr(httpResp, "apTagsAndVlanIds", apTagsAndVlanIdsAttrs)
	if err.HasError() {
		tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
	}

	return apTagsAndVlanIdsList, diags
}

func NetworksWirelessSsidStateGre(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var gre GRE

	concentratorAttrs := map[string]attr.Type{
		"host": types.StringType,
	}
	greAttrs := map[string]attr.Type{
		"concentrator": types.ObjectType{AttrTypes: concentratorAttrs},
		"key":          types.Int64Type,
	}

	if g, ok := httpResp["gre"].(map[string]interface{}); ok {

		// key
		gre.Key, diags = utils.ExtractInt64Attr(httpResp, "key")

		// concentrator
		if c, ok := g["concentrator"].(map[string]interface{}); ok {
			var concentrator GreConcentrator

			concentrator.Host, diags = utils.ExtractStringAttr(c, "host")

			concentratorObj, err := types.ObjectValueFrom(context.Background(), concentratorAttrs, concentrator)
			if err.HasError() {
				diags.Append(err...)
			}

			gre.Concentrator = concentratorObj
		} else {
			concentratorObjNull := types.ObjectNull(concentratorAttrs)
			gre.Concentrator = concentratorObjNull
		}

	} else {
		greObjNull := types.ObjectNull(greAttrs)
		return greObjNull, diags
	}

	greObj, err := types.ObjectValueFrom(context.Background(), greAttrs, gre)
	if err.HasError() {
		diags.Append(err...)
	}

	return greObj, diags
}

func NetworksWirelessSsidStateDnsRewrite(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var dnsRewrite DnsRewrite

	dnsCustomNameServersAttrs := types.ListType{ElemType: types.StringType}

	dnsRewriteAttrs := map[string]attr.Type{
		"enabled":                 types.BoolType,
		"dns_custom_name_servers": dnsCustomNameServersAttrs,
	}

	dns, ok := httpResp["dnsRewrite"].(map[string]interface{})
	if ok {

		// enabled
		enabled, err := utils.ExtractBoolAttr(dns, "enabled")
		if err.HasError() {
			diags.AddError("enabled Attr", fmt.Sprintf("%s", err.Errors()))
		}

		dnsRewrite.Enabled = enabled

		// dns custom Name Servers
		dnsCustomNameservers, err := utils.ExtractListStringAttr(dns, "dnsCustomNameServers")
		if err.HasError() {
			diags.AddError("dnsCustomNameservers Attr", fmt.Sprintf("%s", err.Errors()))
		}

		dnsRewrite.DnsCustomNameservers = dnsCustomNameservers

	} else {
		dnsRewriteObjNull := types.ObjectNull(dnsRewriteAttrs)
		return dnsRewriteObjNull, diags
	}

	// dnsRewrite Terraform types Object
	dnsRewriteObj, dnsRewriteDiags := types.ObjectValueFrom(context.Background(), dnsRewriteAttrs, dnsRewrite)
	if dnsRewriteDiags.HasError() {
		diags.AddError("dnsRewriteObject Attr", fmt.Sprintf("%s", dnsRewriteDiags.Errors()))
	}

	return dnsRewriteObj, diags
}

func NetworksWirelessSsidStateSpeedBurst(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var speedBurst SpeedBurst

	speedBurstAttrs := map[string]attr.Type{
		"enabled": types.BoolType,
	}

	sb, ok := httpResp["speedBurst"].(map[string]interface{})
	if ok {
		enabled, err := utils.ExtractBoolAttr(sb, "enabled")
		if err.HasError() {
			diags.AddError("enabled Attr", fmt.Sprintf("%s", err.Errors()))
		}

		speedBurst.Enabled = enabled
	} else {
		speedBurstObjNull := types.ObjectNull(speedBurstAttrs)
		return speedBurstObjNull, diags
	}

	speedBurstObj, speedBurstDiags := types.ObjectValueFrom(context.Background(), speedBurstAttrs, speedBurst)
	if speedBurstDiags.HasError() {
		diags.AddError("enabled Attr", fmt.Sprintf("%s", speedBurstDiags.Errors()))
	}

	return speedBurstObj, diags
}

func NetworksWirelessSsidStateNamedVlans(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var namedVlans NamedVlans

	byApTagsAttrs := map[string]attr.Type{
		"tags":      types.ListType{ElemType: types.StringType},
		"vlan_name": types.StringType,
	}

	taggingAttrs := map[string]attr.Type{
		"enabled":           types.BoolType,
		"default_vlan_name": types.StringType,
		"by_ap_tags":        types.ListType{ElemType: types.ObjectType{AttrTypes: byApTagsAttrs}},
	}

	guestVlanAttrs := map[string]attr.Type{
		"enabled": types.BoolType,
		"name":    types.StringType,
	}

	radiusAttrs := map[string]attr.Type{
		"guest_vlan": types.ObjectType{AttrTypes: guestVlanAttrs},
	}

	namedVlansAttrs := map[string]attr.Type{
		"tagging": types.ObjectType{AttrTypes: taggingAttrs},
		"radius":  types.ObjectType{AttrTypes: radiusAttrs},
	}

	nv, ok := httpResp["namedVlans"].(map[string]interface{})
	if ok {

		// tagging
		t, ok := nv["tagging"].(map[string]interface{})
		if ok {
			var tagging Tagging

			// Enabled
			enabled, err := utils.ExtractBoolAttr(t, "enabled")
			if err.HasError() {
				diags.AddError("enabled Attr", fmt.Sprintf("%s", err.Errors()))
			}
			tagging.Enabled = enabled

			// DefaultVlanName
			defaultVlanName, err := utils.ExtractStringAttr(t, "defaultVlanName")
			if err.HasError() {
				diags.AddError("defaultVlanName Attr", fmt.Sprintf("%s", err.Errors()))
			}
			tagging.DefaultVlanName = defaultVlanName

			var byApTagsArray []types.Object

			// ByApTags
			bat, ok := nv["byApTags"].([]interface{})
			if ok {

				for _, ba := range bat {
					if b, ok := ba.(map[string]interface{}); ok {
						var byApTags ByApTag

						// tags
						tags, err := utils.ExtractListStringAttr(b, "tags")
						if err.HasError() {
							diags.AddError("tags Attr", fmt.Sprintf("%s", err.Errors()))
						}
						byApTags.Tags = tags

						// vlanName
						vlanName, err := utils.ExtractStringAttr(b, "vlanName")
						if err.HasError() {
							diags.AddError("vlanName Attr", fmt.Sprintf("%s", err.Errors()))
						}
						byApTags.VlanName = vlanName

						byApTagsObj, err := types.ObjectValueFrom(context.Background(), byApTagsAttrs, byApTags)
						if err.HasError() {
							diags.AddError("byApTags Object Attr", fmt.Sprintf("%s", err.Errors()))
						}

						byApTagsArray = append(byApTagsArray, byApTagsObj)

					}
				}

				// returns a populated or empty list instead of a null value
				if byApTagsArray != nil {
					byApTagsObjArray, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: byApTagsAttrs}, byApTagsArray)
					if err.HasError() {
						diags.AddError("byApTags Array Attr", fmt.Sprintf("%s", err.Errors()))
					}

					tagging.ByApTags = byApTagsObjArray
				} else {
					byApTagsObjArray, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: byApTagsAttrs}, []attr.Value{})
					if err.HasError() {
						diags.AddError("byApTags Array Attr", fmt.Sprintf("%s", err.Errors()))
					}

					tagging.ByApTags = byApTagsObjArray
				}

			} else {
				byApTagsArrayNull := types.ListNull(types.ObjectType{AttrTypes: byApTagsAttrs})
				tagging.ByApTags = byApTagsArrayNull
			}

			taggingObj, err := types.ObjectValueFrom(context.Background(), taggingAttrs, tagging)
			if err.HasError() {
				diags.AddError("tagging Object Attr", fmt.Sprintf("%s", err.Errors()))
			}
			namedVlans.Tagging = taggingObj

		} else {
			taggingObjNull := types.ObjectNull(taggingAttrs)
			namedVlans.Tagging = taggingObjNull
		}

		// radius
		r, ok := nv["radius"].(map[string]interface{})
		if ok {
			var radius Radius

			g, ok := r["guestVlan"].(map[string]interface{})
			if ok {
				var guestVlans RadiusGuestVlan

				// enabled
				enabled, err := utils.ExtractBoolAttr(g, "enabled")
				if err.HasError() {
					diags.AddError("enabled Attr", fmt.Sprintf("%s", err.Errors()))
				}
				guestVlans.Enabled = enabled

				// name
				name, err := utils.ExtractStringAttr(g, "name")
				if err.HasError() {
					diags.AddError("name Attr", fmt.Sprintf("%s", err.Errors()))
				}
				guestVlans.Name = name

				guestVlansObj, err := types.ObjectValueFrom(context.Background(), guestVlanAttrs, guestVlans)
				if err.HasError() {
					diags.AddError("guestVlans Object Attr", fmt.Sprintf("%s", err.Errors()))
				}
				radius.GuestVlan = guestVlansObj

			} else {
				guestVlansObjNull := types.ObjectNull(guestVlanAttrs)
				radius.GuestVlan = guestVlansObjNull
			}

			radiusObj, err := types.ObjectValueFrom(context.Background(), radiusAttrs, radius)
			if err.HasError() {
				diags.AddError("radius object Attr", fmt.Sprintf("%s", err.Errors()))
			}
			namedVlans.Radius = radiusObj

		} else {
			radiusObjNull := types.ObjectNull(radiusAttrs)
			namedVlans.Radius = radiusObjNull
		}

	} else {
		namedVlansObjNull := types.ObjectNull(namedVlansAttrs)
		return namedVlansObjNull, diags
	}

	namedVlansObj, err := types.ObjectValueFrom(context.Background(), radiusAttrs, namedVlans)
	if err.HasError() {
		diags.AddError("namedVlans object Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return namedVlansObj, diags

}
