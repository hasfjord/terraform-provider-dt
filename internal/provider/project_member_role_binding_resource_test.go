// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMemberRoleBindingResource(t *testing.T) {
	t.Parallel()
	t.Log("TestAccDeviceDataSource")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create testing
			{
				Config: providerConfig + readTestFile(t, "../../testdata/member/three_projects.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project_member_role_bindings.test", "email", "d0hjenj24tsg00b24tb0@cvinmt9aq9sc738g6ep0.serviceaccount.d21s.com"),
					resource.TestCheckResourceAttr("dt_project_member_role_bindings.test", "organization", "organizations/cvinmt9aq9sc738g6eog"),
					resource.TestCheckResourceAttr("dt_project_member_role_bindings.test", "projects.#", "3"),
					resource.TestCheckResourceAttr("dt_project_member_role_bindings.test", "role", "roles/project.user"),
				),
			},
			// Import testing
			{
				ResourceName:                         "dt_project_member_role_bindings.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources["dt_project_member_role_bindings.test"].Primary.Attributes["name"], nil
				},
			},
			// Update testing
			{
				Config: providerConfig + readTestFile(t, "../../testdata/member/three_projects_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dt_project_member_role_bindings.test", "email", "d0hjenj24tsg00b24tb0@cvinmt9aq9sc738g6ep0.serviceaccount.d21s.com"),
					resource.TestCheckResourceAttr("dt_project_member_role_bindings.test", "organization", "organizations/cvinmt9aq9sc738g6eog"),
					resource.TestCheckResourceAttr("dt_project_member_role_bindings.test", "projects.#", "3"),
					resource.TestCheckResourceAttr("dt_project_member_role_bindings.test", "role", "roles/project.admin"),
				),
			},
		},
	})
}
