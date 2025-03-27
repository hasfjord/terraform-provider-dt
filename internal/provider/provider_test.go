// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/sirupsen/logrus"
)

var (
	// providerConfig is the configuration for the provider that will be used
	providerConfig string
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"dt": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func TestMain(m *testing.M) {
	providerConfig = `provider "dt" {
		url            = "https://api.disruptive-technologies.com"
  		token_endpoint = "https://identity.disruptive-technologies.com/oauth2/token"
	}
	
	`
	logrus.WithField("providerConfig", providerConfig).Info("providerConfig")

	// Run the tests
	resource.TestMain(m)
}
