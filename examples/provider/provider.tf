# Copyright (c) HashiCorp, Inc.

provider "dt" {
  url            = "https://api.disruptive-technologies.com"
  emulator_url   = "https://emulator.disruptive-technologies.com"
  token_endpoint = "https://identity.disruptive-technologies.com/oauth2/token"
}
