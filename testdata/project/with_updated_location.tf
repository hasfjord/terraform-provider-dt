# Copyright (c) HashiCorp, Inc.

resource "dt_project" "test" {
  display_name = "Acceptance Test Project"
  organization = "organizations/cvinmt9aq9sc738g6eog"
  location = {
    latitude      = 59.910953
    longitude     = 10.639040
    time_location = "Europe/Oslo"
  }
}
