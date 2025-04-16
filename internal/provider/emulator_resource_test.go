// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccEmulatorResourceExample(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: providerConfig + readTestFile(t, "../../examples/resources/dt_emulator/resource.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_emulator.my_emulator", "display_name", "Terraform created emulator"),
					resource.TestCheckResourceAttr("dt_emulator.my_emulator", "type", "touch"),
				),
			},
			// Destroy only the emulator resource
			{
				Config:  providerConfig + readTestFile(t, "../../examples/resources/dt_emulator/resource.tf"),
				Destroy: true,
			},
		},
	})
}

func TestAccEmulatorResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and read testing
			{
				Config: providerConfig +
					readTestFile(t, "../../testdata/emulator/with_labels.tf"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dt_emulator.test", "display_name", "Emulator with custom labels"),
					resource.TestCheckResourceAttr("dt_emulator.test", "type", "temperature"),
					resource.TestCheckResourceAttr("dt_emulator.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("dt_emulator.test", "labels.foo", "bar"),
					resource.TestCheckResourceAttr("dt_emulator.test", "labels.bar", "baz"),
				),
			},
			// Import testing
			{
				ResourceName:                         "dt_emulator.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["dt_emulator.test"].Primary.Attributes["name"], nil
				},
			},
			// Destroy only the emulator resource
			{
				Config: providerConfig +
					readTestFile(t, "../../testdata/emulator/with_labels.tf"),
				Destroy: true,
			},
		},
	})
}
