// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceDataSource(t *testing.T) {
	t.Parallel()
	t.Log("TestAccDeviceDataSource")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "dt_device" "test" {name = "projects/cvinutal2ugc73b866v0/devices/emucvio050h6oic7398hljg"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.dt_device.test", "name", "projects/cvinutal2ugc73b866v0/devices/emucvio050h6oic7398hljg"),
					resource.TestCheckResourceAttr("data.dt_device.test", "device_id", "emucvio050h6oic7398hljg"),
					resource.TestCheckResourceAttr("data.dt_device.test", "project_id", "cvinutal2ugc73b866v0"),
					resource.TestCheckResourceAttr("data.dt_device.test", "type", "temperature"),
					resource.TestCheckResourceAttr("data.dt_device.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("data.dt_device.test", "labels.virtual-sensor", ""),
				),
			},
		},
	})
}
