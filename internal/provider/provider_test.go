// Copyright (c) HashiCorp, Inc.

package provider

import (
	"net/http/httptest"
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
	// Start the test server
	testDTServer := httptest.NewServer(newTestHandler())
	defer testDTServer.Close()

	providerConfig = `provider "dt" {
email = "myServiceAccount@myProject.serviceaccount.d21s.com"
key_id = "myKey"
key_secret = "mySecret"
token_endpoint = "` + testDTServer.URL + `/token"
url = "` + testDTServer.URL + `"
}

`
	logrus.WithField("providerConfig", providerConfig).Info("providerConfig")

	// Run the tests
	resource.TestMain(m)
}
