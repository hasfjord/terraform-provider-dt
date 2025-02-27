// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt"
	"github.com/disruptive-technologies/terraform-provider-dt/internal/dt/oidc"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	resp.TypeName = "dt"
	resp.Version = p.version
}

func (p *DTProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "The URL of the API server.",
				// Can use either environment variables or configuration, therefore optional: true
				Optional: true,
			},
			"key_id": schema.StringAttribute{
				Description: "The key ID from the service account.",
				// Can use either environment variables or configuration, therefore optional: true
				Optional: true,
			},
			"key_secret": schema.StringAttribute{
				Description: "The key secret from the service account.",
				Sensitive:   true,
				// Can use either environment variables or configuration, therefore optional: true
				Optional: true,
			},
			"token_endpoint": schema.StringAttribute{
				Description: "The token endpoint for the OIDC provider.",
				// Can use either environment variables or configuration, therefore optional: true
				Optional: true,
			},
			"email": schema.StringAttribute{
				Description: "The email address used to authenticate with the OIDC provider.",
				// Can use either environment variables or configuration, therefore optional: true
				Optional: true,
			},
		},
	}
}

// hashicupsProviderModel maps provider schema data to a Go type.
type dtProviderModel struct {
	URL           types.String `tfsdk:"url"`
	ClientID      types.String `tfsdk:"key_id"`
	ClientSecret  types.String `tfsdk:"key_secret"`
	TokenEndpoint types.String `tfsdk:"token_endpoint"`
	Email         types.String `tfsdk:"email"`
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
			url = config.URL.ValueString()
		}
	}

	keyID := os.Getenv("DT_API_KEY_ID")
	if keyID == "" {
		if config.ClientID.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("key_id"),
				"Key ID must be set",
				"The key ID to authenticate with must be set",
			)
		} else {
			keyID = config.ClientID.ValueString()
		}
	}

	keySecret := os.Getenv("DT_API_KEY_SECRET")
	if keySecret == "" {
		if config.ClientSecret.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("key_secret"),
				"key secret must be set",
				"The secret to authenticate with must be set",
			)
		} else {
			keySecret = config.ClientSecret.ValueString()
		}
	}
	tokenEndpoint := os.Getenv("DT_OIDC_TOKEN_ENDPOINT")
	if tokenEndpoint == "" {
		if config.TokenEndpoint.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("token_endpoint"),
				"Token endpoint must be set",
				"The token endpoint for the OIDC provider must be set",
			)
		} else {
			tokenEndpoint = config.TokenEndpoint.ValueString()
		}
	}
	email := os.Getenv("DT_OIDC_EMAIL")
	if email == "" {
		if config.Email.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("email"),
				"Email must be set",
				"The email address used to authenticate with the OIDC provider must be set",
			)
		} else {
			email = config.Email.ValueString()
		}
	}

	// if there are any errors, return early
	if resp.Diagnostics.HasError() {
		for _, diag := range resp.Diagnostics {
			tflog.Error(ctx, diag.Summary())
		}
		return
	}

	ctx = tflog.SetField(ctx, "url", url)
	ctx = tflog.SetField(ctx, "key_id", keyID)
	ctx = tflog.SetField(ctx, "token_endpoint", tokenEndpoint)
	ctx = tflog.SetField(ctx, "email", email)
	tflog.Debug(ctx, "provider parameters")

	client := dt.NewClient(dt.Config{
		URL: url,
		Oidc: oidc.Config{
			TokenEndpoint: tokenEndpoint,
			ClientID:      keyID,
			ClientSecret:  keySecret,
			Email:         email,
		},
	})

	// make the client available to the rest of the provider
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources defines the resources implemented in the provider.
func (p *DTProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewDataConnectorResource,
	}
}

func (p *DTProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
	}
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
