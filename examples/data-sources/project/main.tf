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

data "dt_project" "thomas_test_project" {
  provider = disruptive-technologies
  name     = "projects/ccol8iuk9smqiha4e8l0"
}

output "project" {
  value = data.dt_project.thomas_test_project
}
