// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

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
				),
			},
		},
	})
}
