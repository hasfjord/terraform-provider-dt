# Copyright (c) HashiCorp, Inc.

resource "dt_emulator" "my_emulator" {
  display_name = "Terraform created emulator"
  project_id   = "d0ito5m62hus73ae3lr0"
  type         = "touch"
}
