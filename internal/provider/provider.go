// Copyright 2024 Antimetal LLC
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/antimetal/terraform-provider-antimetal/internal/antimetal"
	amDataSource "github.com/antimetal/terraform-provider-antimetal/internal/provider/datasource"
	amResource "github.com/antimetal/terraform-provider-antimetal/internal/provider/resource"
)

const (
	providerType = "antimetal"
)

var (
	_ provider.Provider = &Antimetal{}
)

type AntimetalModel struct {
	URL types.String `tfsdk:"url"`
}

// Antimetal defines the provider implementation.
type Antimetal struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Antimetal{
			version: version,
		}
	}
}

func (p *Antimetal) Metadata(ctx context.Context,
	req provider.MetadataRequest, resp *provider.MetadataResponse) {

	resp.TypeName = providerType
	resp.Version = p.version
}

func (p *Antimetal) Schema(ctx context.Context,
	req provider.SchemaRequest, resp *provider.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional: true,
				Description: "The default is 'antimetal.com'. This is optional and shouldn't be " +
					"changed under normal circumstances.",
			},
		},
	}
}

func (p *Antimetal) Configure(ctx context.Context,
	req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var cfg AntimetalModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var clientOpts []antimetal.ClientOption

	if !cfg.URL.IsNull() {
		clientOpts = append(clientOpts, antimetal.WithURL(cfg.URL.ValueString()))
	}

	client, err := antimetal.NewClient(clientOpts...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Antimetal Client",
			fmt.Sprintf("... details ... %s", err),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *Antimetal) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		amResource.NewHandshake,
	}
}

func (p *Antimetal) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		amDataSource.NewAWSIAMPolicyDocument,
		amDataSource.NewAWSIAMAssumeRolePolicy,
	}
}
