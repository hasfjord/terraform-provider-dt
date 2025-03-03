# Copyright (c) HashiCorp, Inc.

data "dt_project" "test_project" {
  provider = disruptive-technologies
  name     = "projects/ccol8iuk9smqiha4e8l0"
}

output "project" {
  value = data.dt_project.test_project
}
