package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*denvrProvider)(nil)

func New() func() provider.Provider {
	return func() provider.Provider {
		return &denvrProvider{}
	}
}

type denvrProvider struct{}

func (p *denvrProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {

}

func (p *denvrProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *denvrProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "denvr"
}

func (p *denvrProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *denvrProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVmResource,
	}
}
