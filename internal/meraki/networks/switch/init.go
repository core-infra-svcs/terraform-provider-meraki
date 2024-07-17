package _switch

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func Resources() []func() resource.Resource {
	return []func() resource.Resource{
		NewNetworksSwitchDscpToCosMappingsResource,
		NewNetworksSwitchMtuResource,
		NewNetworksSwitchQosRuleResource,
		NewNetworksSwitchSettingsResource,
		NewDevicesSwitchPortResource,
		NewDevicesSwitchPortsCycleResource,
	}
}

func DataSources() []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNetworksSwitchMtuDataSource,
		NewNetworksSwitchQosRulesDataSource,
		NewDevicesSwitchPortsStatusesDataSource,
	}
}
