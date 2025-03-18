# Copyright (c) HashiCorp, Inc.

data "dt_project" "provider_test_project" {
  name = "projects/your-project-id"
}

data "google_iam_workload_identity_pool" "pool" {
  project = data.dt_project.provider_test_project.id
  name    = "your-pool-name"
}

data "google_iam_workload_identity_pool_provider" "provider" {
  project                            = data.dt_project.provider_test_project.id
  workload_identity_pool_provider_id = "your-provider-id"
}

resource "google_pubsub_topic_iam_member" "default" {
  project = "your-project-id"
  topic   = "projects/your-project-id/topics/your-topic"
  role    = "roles/pubsub.publisher"
  member  = "principal://iam.googleapis.com/${data.google_iam_workload_identity_pool.pool.name}/subject/your-dt-org-id"
}

resource "dt_data_connector" "pub_sub_data_connector" {
  display_name = "Pub/Sub Data Connector"
  type         = "GOOGLE_CLOUD_PUBSUB"
  project      = data.dt_project.provider_test_project.id
  labels       = ["name", "location"]
  events       = ["temperature", "humidity"]
  pubsub_config = {
    topic    = "projects/your-project-id/topics/your-topic"
    audience = "//iam.googleapis.com/${data.google_iam_workload_identity_pool_provider.provider.name}"
  }
}

