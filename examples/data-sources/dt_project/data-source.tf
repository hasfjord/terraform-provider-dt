# Copyright (c) HashiCorp, Inc.

data "dt_project" "test_project" {
  provider = disruptive-technologies
  name     = "projects/your-project-id"
}

output "project" {
  value = data.dt_project.test_project
}
