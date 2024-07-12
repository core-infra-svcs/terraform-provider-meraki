package wireless

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworksWirelessSsidsFirewallL3FirewallRulesResource,
		NewNetworksWirelessSsidsFirewallL7FirewallRulesResource,
		NewNetworksWirelessSsidsResource,
		NewNetworksWirelessSsidsSplashSettingsResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNetworksWirelessSsidsDataSource,
	}
}
