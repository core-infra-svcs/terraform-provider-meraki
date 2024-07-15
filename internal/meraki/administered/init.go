package administered

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var administeredResources []func() resource.Resource
var administeredDataSources []func() datasource.DataSource

func init() {

	// Initialize with all administered resources //

	//administeredResources = append(administeredResources, subFolder.Resources()...)

	// package level resources
	//administeredResources = append(administeredResources, )

	// Initialize with all administered data sources //

	//administeredDataSources = append(administeredDataSources, subFolder.DataSources()...)

	// package level data sources
	administeredDataSources = append(administeredDataSources,
		NewAdministeredIdentitiesMeDataSource,
	)

}

func Resources() []func() resource.Resource {
	return administeredResources
}

func DataSources() []func() datasource.DataSource {
	return administeredDataSources
}
