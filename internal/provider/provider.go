// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"terraform-provider-kickflow/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure KickflowProvider satisfies various provider interfaces.
var _ provider.Provider = &KickflowProvider{}
var _ provider.ProviderWithFunctions = &KickflowProvider{}

// KickflowProvider defines the provider implementation.
type KickflowProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KickflowProviderModel describes the provider data model.
type KickflowProviderModel struct {
	Endpoint    types.String `tfsdk:"endpoint"`
	AccessToken types.String `tfsdk:"access_token"`
	CallerID    types.String `tfsdk:"caller_id"`
}

func (p *KickflowProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kickflow"
	resp.Version = p.version
}

func (p *KickflowProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Kickflow API endpoint (default: https://api.kickflow.com)",
				Optional:            true,
			},
			"access_token": schema.StringAttribute{
				MarkdownDescription: "Kickflow API access token",
				Required:            true,
				Sensitive:           true,
			},
			"caller_id": schema.StringAttribute{
				MarkdownDescription: "ID to identify the API request caller",
				Required:            true,
			},
		},
	}
}

func (p *KickflowProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config KickflowProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the configuration values.
	if config.Endpoint.IsNull() {
		config.Endpoint = types.StringValue("https://api.kickflow.com")
	}

	if config.AccessToken.IsNull() {
		resp.Diagnostics.AddError(
			"Access token is not set",
			"Please set an access token to access the Kickflow API.",
		)
		return
	}

	if config.CallerID.IsNull() {
		resp.Diagnostics.AddError(
			"Caller ID is not set",
			"Please set a Caller ID to identify the API request caller.",
		)
		return
	}

	// Configure the Kickflow client
	kickflowClient := client.NewKickflowClient(
		config.Endpoint.ValueString(),
		config.AccessToken.ValueString(),
		config.CallerID.ValueString(),
	)
	resp.DataSourceData = kickflowClient
	resp.ResourceData = kickflowClient
}

func (p *KickflowProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *KickflowProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
		NewKickflowUserDataSource,
	}
}

func (p *KickflowProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KickflowProvider{
			version: version,
		}
	}
}
