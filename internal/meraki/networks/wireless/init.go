package wireless

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/meraki/networks/wireless/ssids"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworksWirelessSsidsFirewallL3FirewallRulesResource,
		NewNetworksWirelessSsidsFirewallL7FirewallRulesResource,
		ssids.NewNetworksWirelessSsidsResource,
		NewNetworksWirelessSsidsSplashSettingsResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNetworksWirelessSsidsDataSource,
	}
}
