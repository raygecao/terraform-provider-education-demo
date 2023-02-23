package education

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const school = "BUPT"

type schoolDataSource struct {
}

func (s *schoolDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_school"
}

func (s *schoolDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed: true,
		},
	}}
}

func (s *schoolDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	diagnostics := response.State.Set(ctx, schoolDataSourceModel{
		Name: types.StringValue(school),
	})
	response.Diagnostics.Append(diagnostics...)
}

func NewSchoolDataSource() datasource.DataSource {
	return &schoolDataSource{}
}

type schoolDataSourceModel struct {
	Name types.String `tfsdk:"name"`
}
