# Copyright (c) HashiCorp, Inc.

resource "dt_project" "test" {
  display_name = "Empty Location Project"
  organization = "organizations/cvinmt9aq9sc738g6eog"
  location     = {}
}
