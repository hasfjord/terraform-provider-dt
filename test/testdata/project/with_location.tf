# Copyright (c) HashiCorp, Inc.

resource "dt_project" "test" {
  display_name = "Acceptance Test Project"
  organization = "organizations/cvinmt9aq9sc738g6eog"
  location = {
    latitude      = 63.44539
    longitude     = 10.910202
    time_location = "Europe/Oslo"
  }
}
