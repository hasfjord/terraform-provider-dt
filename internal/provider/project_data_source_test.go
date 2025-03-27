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
				Config: providerConfig + `data "dt_project" "test" {name = "projects/cvimk1al2ugc73b866qg"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.dt_project.test", "name", "projects/cvimk1al2ugc73b866qg"),
					resource.TestCheckResourceAttr("data.dt_project.test", "id", "cvimk1al2ugc73b866qg"),
					resource.TestCheckResourceAttr("data.dt_project.test", "display_name", "Terraform Provider DT Acceptance Tests"),
					resource.TestCheckResourceAttr("data.dt_project.test", "inventory", "false"),
					resource.TestCheckResourceAttr("data.dt_project.test", "organization", "organizations/c4nif0cqjh0g02hh4t10"),
					resource.TestCheckResourceAttr("data.dt_project.test", "organization_display_name", "Thomas Hasfjord"),
					resource.TestCheckResourceAttr("data.dt_project.test", "sensor_count", "1"),
					resource.TestCheckResourceAttr("data.dt_project.test", "cloud_connector_count", "0"),
				),
			},
		},
	})
}
