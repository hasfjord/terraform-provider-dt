// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Setup separate project for the test.
// There can only be 10 data connectors per project.
var dataConnectorProviderConfig = providerConfig + `
resource "dt_project" "test" {
	display_name = "data connector Acceptance Test Project"
	organization = "organizations/cvinmt9aq9sc738g6eog"
	location = {
		time_location = "Europe/Oslo"
	}
}

`

func TestAccDataConnectorProvider(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create pubsub data connector with no labels or events
			{
				Config: dataConnectorProviderConfig + `
				resource "dt_data_connector" "test" {
					display_name = "data connector Acceptance Test"
					type = "GOOGLE_CLOUD_PUBSUB"
					project = dt_project.test.id
					pubsub_config = {
						topic    = "projects/your-project-id/topics/your-topic"
						audience = "//iam.googleapis.com/projects/12345689/locations/europe-west1/workloadIdentityPools/my-pool-id/providers/my-provider-id"
					}
				}
				
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_data_connector.test", "display_name", "data connector Acceptance Test"),
					resource.TestCheckResourceAttr("dt_data_connector.test", "type", "GOOGLE_CLOUD_PUBSUB"),
					resource.TestCheckResourceAttr("dt_data_connector.test", "pubsub_config.topic", "projects/your-project-id/topics/your-topic"),
					resource.TestCheckResourceAttr("dt_data_connector.test", "pubsub_config.audience", "//iam.googleapis.com/projects/12345689/locations/europe-west1/workloadIdentityPools/my-pool-id/providers/my-provider-id"),
					resource.TestCheckResourceAttr("dt_data_connector.test", "labels.%", "0"),
					resource.TestCheckResourceAttr("dt_data_connector.test", "events.%", "0"),
				),
			},
			// Import testing
			{
				ResourceName:                         "dt_data_connector.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["dt_data_connector.test"].Primary.Attributes["name"], nil
				},
			},
		},
	})
}
