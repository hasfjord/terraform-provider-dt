# Copyright (c) HashiCorp, Inc.

resource "dt_project" "provider_test_project" {
  display_name = "Provider Test"
  organization = "organizations/dt"
  location = {
    // Hell, Norway
    latitude      = 63.44539
    longitude     = 10.910202
    time_location = "Europe/Oslo"
  }
}
