# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    disruptive-technologies = {
      source = "registry.terraform.io/hasfjord/dt"
    }
  }
}

provider "disruptive-technologies" {
  url            = "https://api.dev.disruptive-technologies.com/v2"
  token_endpoint = "https://identity.dev.disruptive-technologies.com/oauth2/token"
}

data "dt_device" "test_device" {
  provider = disruptive-technologies
  name     = "projects/ccol8iuk9smqiha4e8l0/devices/emucv0799gjrncc73fnv1dg"
}

output "device" {
  value = data.dt_device.test_device
}
