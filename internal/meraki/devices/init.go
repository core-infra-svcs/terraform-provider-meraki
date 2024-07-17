package devices

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var devicesResources []func() resource.Resource
var devicesDataSources []func() datasource.DataSource

func init() {

	// Initialize with all devices resources //

	//devicesResources = append(devicesResources, subFolder.Resources()...)

	// package level resources
	devicesResources = append(devicesResources,
		NewDevicesCellularSimsResource,
		NewDevicesResource,
		NewDevicesTestAccDevicesManagementInterfaceResourceResource,
	)

	// Initialize with all devices data sources //

	//devicesDataSources = append(devicesDataSources, subFolder.DataSources()...)

	// package level data sources
	devicesDataSources = append(devicesDataSources,
		NewNetworkDevicesDataSource,
		NewDevicesManagementInterfaceDatasource,
	)

}

func Resources() []func() resource.Resource {
	return devicesResources
}

func DataSources() []func() datasource.DataSource {
	return devicesDataSources
}
