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
				Config: providerConfig + `data "dt_device" "test" {name = "projects/your-project-id/devices/your-device-id"}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.dt_device.test", "name", "projects/your-project-id/devices/your-device-id"),
					resource.TestCheckResourceAttr("data.dt_device.test", "device_id", "your-device-id"),
					resource.TestCheckResourceAttr("data.dt_device.test", "project_id", "your-project-id"),
					resource.TestCheckResourceAttr("data.dt_device.test", "type", "temperature"),
					resource.TestCheckResourceAttr("data.dt_device.test", "labels.%", "1"),
					resource.TestCheckResourceAttr("data.dt_device.test", "labels.key", "value"),
				),
			},
		},
	})
}
