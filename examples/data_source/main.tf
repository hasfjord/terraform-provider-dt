terraform {
  required_providers {
    disruptive-technologies = {
      source = "disruptive-technologies.com/api/dt"
    }
  }
}

provider "disruptive-technologies" {
  url      = "https://api.dev.disruptive-technologies.com/v2"
  username = "cs5v9pr24td000b24tp0"
}

data "dt_project" "thomas_test_project" {
  provider = disruptive-technologies
  name     = "projects/ccol8iuk9smqiha4e8l0"
}

output "project" {
  value = data.dt_project.thomas_test_project
}
