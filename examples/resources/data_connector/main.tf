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
  token_endpoint = "https://identity.dev.disruptive-technologies.com/oauth2/token"
}

data "dt_project" "provider_test_project" {
  provider = disruptive-technologies
  name     = "projects/ct5m0ndfb7rvtcogrl0g"
}

resource "dt_data_connector" "pub_sub_data_connector" {
  provider     = disruptive-technologies
  display_name = "Pub/Sub Data Connector"
  type         = "GOOGLE_CLOUD_PUBSUB"
  project      = data.dt_project.provider_test_project.id
  //labels       = ["name", "location"]
  //events = ["temperature", "humidity"]
  pubsub_config = {
    topic    = "projects/dt-dev-169909/topics/data-connector-test"
    audience = "//iam.googleapis.com/projects/392692775543/locations/global/workloadIdentityPools/dt-data-connector-oidc-pool/providers/dt-data-connector-oidc-provider"
  }
}

