terraform {
  required_providers {
    disruptive-technologies = {
      source = "disruptive-technologies.com/api/dt"
    }
  }
}

provider "disruptive-technologies" {
  url      = "https://api.dev.disruptive-technologies.com/v2"
  username = "dummy"
  password = "dummy"
}
