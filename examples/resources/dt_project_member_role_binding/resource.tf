# Copyright (c) HashiCorp, Inc.

data "dt_project" "test1" { name = "projects/d0hj3ndaoups738bc8og" }
data "dt_project" "test2" { name = "projects/d0hj3qdaoups738bc8pg" }
data "dt_project" "test3" { name = "projects/d0hj3s5aoups738bc8qg" }

resource "dt_project_member_role_bindings" "test" {
  email        = "d0hjenj24tsg00b24tb0@cvinmt9aq9sc738g6ep0.serviceaccount.d21s.com"
  organization = "organizations/cvinmt9aq9sc738g6eog"
  projects = [
    data.dt_project.test1.name,
    data.dt_project.test2.name,
    data.dt_project.test3.name,
  ]
  role = "roles/project.user"
}
