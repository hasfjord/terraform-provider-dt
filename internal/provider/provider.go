// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &DTProvider{}
var _ provider.ProviderWithFunctions = &DTProvider{}

// DTProvider defines the provider implementation.
type DTProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *DTProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "disruptive-technologies"
	resp.Version = p.version
}

func (p *DTProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "The URL of the API server.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "The username to authenticate with.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password to authenticate with.",
				Sensitive:   true,
				Required:    true,
			},
		},
	}
}

// hashicupsProviderModel maps provider schema data to a Go type.
type dtProviderModel struct {
	URL      types.String `tfsdk:"url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *DTProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// retrieve provider data from configuration
	var config dtProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	url := os.Getenv("DT_API_URL")
	if url == "" {
		if config.URL.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("url"),
				"URL must be set",
				"The URL of the dt api server must be set",
			)
		} else {
			url = config.URL.String()
		}
	}

	username := os.Getenv("DT_API_USERNAME")
	if username == "" {
		if config.Username.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("username"),
				"Username must be set",
				"The username to authenticate with must be set",
			)
		} else {
			username = config.Username.String()
		}
	}

	password := os.Getenv("DT_API_PASSWORD")
	if password == "" && config.Password.IsUnknown() {
		if config.Password.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("password"),
				"Password must be set",
				"The password to authenticate with must be set",
			)
		} else {
			password = config.Password.String()
		}
	}

	// if there are any errors, return early
	if resp.Diagnostics.HasError() {
		return
	}

	client := dt.NewClient(url)
	client.WithBasicAuth(username, password)

	// make the client available to the rest of the provider
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *DTProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *DTProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *DTProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DTProvider{
			version: version,
		}
	}
}
