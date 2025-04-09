// Copyright (c) HashiCorp, Inc.

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	// providerConfig is the configuration for the provider that will be used
	providerConfig = `provider "dt" {
		url            = "https://api.disruptive-technologies.com"
  		token_endpoint = "https://identity.disruptive-technologies.com/oauth2/token"
	}
	
	`
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"dt": providerserver.NewProtocol6WithError(New("test")()),
	}
)

// notificationActionExample is a helper for reading the test .tf file
func readTestFile(t *testing.T, filePath string) string {
	t.Helper()
	// Read the example file
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read example file: %v", err)
	}
	return string(content)
}
