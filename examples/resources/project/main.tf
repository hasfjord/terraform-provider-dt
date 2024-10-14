terraform {
  required_providers {
    disruptive-technologies = {
      source = "disruptive-technologies.com/api/dt"
    }
  }
}

provider "disruptive-technologies" {
  url      = "https://api.dev.disruptive-technologies.com/v2"
  username = "cs5v9pr24td000b24tp0"
}

resource "dt_project" "provider_test_project" {
  provider     = disruptive-technologies
  display_name = "Provider Test Project"
  organization = "organizations/dt"
  location = {
    latitude      = 63.44539
    longitude     = 10.910202
    time_location = "Europe/Oslo"
  }
}
