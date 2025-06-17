# Copyright (c) HashiCorp, Inc.

data "dt_project" "test" {
  name = "projects/d18gf79mee4c73bk8lsg"
}

resource "dt_emulator" "test" {
  display_name = "Added emulator"
  project_id   = data.dt_project.test.id
  type         = "co2"
}
