# Copyright (c) HashiCorp, Inc.

terraform {
  backend "local" {
    path = "./terraform.tfstate"
  }
  required_providers {
    disruptive-technologies = {
      source = "registry.terraform.io/hasfjord/dt"
    }
  }
}

provider "disruptive-technologies" {
  url            = "https://api.dev.disruptive-technologies.com/v2"
  key_id         = "ct48i7r24te000b24trg"
  token_endpoint = "https://identity.dev.disruptive-technologies.com/oauth2/token"
  email          = "cs50l4324te000b24v0g@ccol8iuk9smqiha4e8l0.serviceaccount.d21s.com"
}

resource "dt_project" "provider_test_project" {
  provider     = disruptive-technologies
  display_name = "Provider Test"
  organization = "organizations/dt"
  location = {
    // Hell, Norway
    latitude      = 63.44539
    longitude     = 10.910202
    time_location = "Europe/Oslo"
  }
}
