// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	t.Parallel()
	t.Log("TestAccDeviceDataSource")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "dt_project" "test" {name = "projects/your-project-id"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.dt_project.test", "name", "projects/your-project-id"),
					resource.TestCheckResourceAttr("data.dt_project.test", "id", "your-project-id"),
					resource.TestCheckResourceAttr("data.dt_project.test", "display_name", "Test Project"),
					resource.TestCheckResourceAttr("data.dt_project.test", "inventory", "true"),
					resource.TestCheckResourceAttr("data.dt_project.test", "organization", "organizations/your-organization-id"),
					resource.TestCheckResourceAttr("data.dt_project.test", "organization_display_name", "Test Organization"),
					resource.TestCheckResourceAttr("data.dt_project.test", "sensor_count", "10"),
					resource.TestCheckResourceAttr("data.dt_project.test", "cloud_connector_count", "1"),
					resource.TestCheckResourceAttr("data.dt_project.test", "location.latitude", "63.44539"),
					resource.TestCheckResourceAttr("data.dt_project.test", "location.longitude", "10.910202"),
					resource.TestCheckResourceAttr("data.dt_project.test", "location.time_location", "Europe/Oslo"),
				),
			},
		},
	})
}
