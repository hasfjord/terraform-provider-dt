// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSafeProjectDataSource(t *testing.T) {
	t.Parallel()
	t.Log("TestAccDeviceDataSource")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "dt_project" "test" {name = "projects/cvinutal2ugc73b866v0"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.dt_project.test", "name", "projects/cvinutal2ugc73b866v0"),
					resource.TestCheckResourceAttr("data.dt_project.test", "id", "cvinutal2ugc73b866v0"),
					resource.TestCheckResourceAttr("data.dt_project.test", "display_name", "manual"),
					resource.TestCheckResourceAttr("data.dt_project.test", "inventory", "false"),
					resource.TestCheckResourceAttr("data.dt_project.test", "organization", "organizations/cvinmt9aq9sc738g6eog"),
					resource.TestCheckResourceAttr("data.dt_project.test", "organization_display_name", "Terraform Provider Acceptance Test Org"),
					resource.TestCheckResourceAttr("data.dt_project.test", "sensor_count", "1"),
					resource.TestCheckResourceAttr("data.dt_project.test", "cloud_connector_count", "0"),
				),
			},
		},
	})
}
