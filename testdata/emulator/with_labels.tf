# Copyright (c) HashiCorp, Inc.

resource "dt_emulator" "test" {
  display_name = "Emulator with custom labels"
  project_id   = "d0ito5m62hus73ae3lr0"
  type         = "temperature"
  labels = {
    foo = "bar"
    bar = "baz"
  }
}
