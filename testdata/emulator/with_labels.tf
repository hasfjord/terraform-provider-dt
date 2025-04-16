# Copyright (c) HashiCorp, Inc.

resource "dt_emulator" "test" {
  display_name = "Emulator with custom labels"
  project_id   = "cvvrosbeetdc738h9r0g"
  type         = "temperature"
  labels = {
    foo = "bar"
    bar = "baz"
  }
}
