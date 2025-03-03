# Copyright (c) HashiCorp, Inc.

data "dt_device" "test_device" {
  provider = disruptive-technologies
  name     = "projects/your-project-id/devices/your-device-id"
}

output "device" {
  value = data.dt_device.test_device
}
