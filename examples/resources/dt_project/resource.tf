# Copyright (c) HashiCorp, Inc.

resource "dt_project" "my_project" {
  display_name = "Terraform created project"
  organization = "organizations/cvinmt9aq9sc738g6eog"
  location = {
    // Hell, Norway
    latitude      = 63.44539
    longitude     = 10.910202
    time_location = "Europe/Oslo"
  }
}
