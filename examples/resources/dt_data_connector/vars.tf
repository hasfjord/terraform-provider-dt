# Copyright (c) HashiCorp, Inc.

variable "dt_project_id" {
  description = "The ID of the project to create the data connector in."
}
variable "pub_sub_topic" {
  description = "The Pub/Sub topic to send events to."
}
variable "pub_sub_audience" {
  description = "The audience to use for authentication to Pub/Sub."
}
