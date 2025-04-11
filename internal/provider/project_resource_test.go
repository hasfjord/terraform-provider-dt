// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResourceExamples(t *testing.T) {
	t.Parallel()
	t.Log("TestAccProjectResourceExamples")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: providerConfig + readTestFile(t, "../../examples/resources/dt_project/resource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project.my_project", "display_name", "Terraform created project"),
					resource.TestCheckResourceAttr("dt_project.my_project", "location.latitude", "63.44539"),
					resource.TestCheckResourceAttr("dt_project.my_project", "location.longitude", "10.910202"),
					resource.TestCheckResourceAttr("dt_project.my_project", "location.time_location", "Europe/Oslo"),
					resource.TestCheckResourceAttr("dt_project.my_project", "inventory", "false"),
				),
			},
		},
	})
}

func TestAccProjectResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create and read testing
				Config: providerConfig + `resource "dt_project" "test" {
					display_name = "Acceptance Test Project"
					organization = "organizations/cvinmt9aq9sc738g6eog"
					location = {
						latitude = 63.44539
						longitude = 10.910202
						time_location = "Europe/Oslo"
					}
				}
				
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project.test", "display_name", "Acceptance Test Project"),
					resource.TestCheckResourceAttr("dt_project.test", "location.latitude", "63.44539"),
					resource.TestCheckResourceAttr("dt_project.test", "location.longitude", "10.910202"),
					resource.TestCheckResourceAttr("dt_project.test", "location.time_location", "Europe/Oslo"),
					resource.TestCheckResourceAttr("dt_project.test", "inventory", "false"),
				),
			},
		},
	})
}
