package education

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-education-demo/mockplatform"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Provider struct{}

var _ provider.Provider = &Provider{}

func NewProvider() provider.Provider {
	return &Provider{}
}

type providerModel struct {
	User   types.String `tfsdk:"user"`
	Passwd types.String `tfsdk:"passwd"`
}

func (p *Provider) Metadata(ctx context.Context, request provider.MetadataRequest, response *provider.MetadataResponse) {
	response.TypeName = "education"
}

func (p *Provider) Schema(ctx context.Context, request provider.SchemaRequest, response *provider.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user": schema.StringAttribute{
				Required: true,
			},
			"passwd": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *Provider) Configure(ctx context.Context, request provider.ConfigureRequest, response *provider.ConfigureResponse) {
	var plan providerModel
	request.Config.Get(ctx, &plan)
	platform, err := mockplatform.NewPlatform(plan.User.ValueString(), plan.Passwd.ValueString())
	if err != nil {
		response.Diagnostics.AddError("failed to creat provider", err.Error())
		tflog.Warn(ctx, "connect to mock platform failed")
		return
	}
	tflog.Info(ctx, "connect to mock platform successfully")
	response.DataSourceData = platform
	response.ResourceData = platform
}

func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSchoolDataSource,
	}
}

func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTeacherResource,
	}
}
