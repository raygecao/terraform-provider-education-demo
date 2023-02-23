package education

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-education-demo/mockplatform"
	"terraform-provider-education-demo/mockplatform/model"
)

type teacherResource struct {
	cli *mockplatform.Platform
}

func NewTeacherResource() resource.Resource {
	return &teacherResource{}
}

func (t *teacherResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData != nil {
		t.cli = request.ProviderData.(*mockplatform.Platform)
	}
}

func (t *teacherResource) Metadata(ctx context.Context, request resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_teacher"
}

func (t *teacherResource) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"subject": schema.StringAttribute{
				Required: true,
			},
			"salary": schema.Int64Attribute{
				Computed: true,
				Optional: true,
			},
			"organ": schema.StringAttribute{
				Optional: true,
				//Computed: true,
			},
		},
	}
}

func (t *teacherResource) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	var plan teacherResourceModel
	diags := request.Config.Get(ctx, &plan)
	response.Diagnostics = append(response.Diagnostics, diags...)
	if response.Diagnostics.HasError() {
		return
	}

	teacher := &model.Teacher{
		ID:      int(plan.ID.ValueInt64()),
		Name:    plan.Name.ValueString(),
		Subject: plan.Subject.ValueString(),
		Salary:  int(plan.Salary.ValueInt64()),
		Organ:   plan.Organ.ValueString(),
	}
	teacher, err := t.cli.CreateTeacher(teacher)
	if err != nil {
		response.Diagnostics.AddError("fail to create teacher", err.Error())
		return
	}
	state := teacherResourceModel{
		ID:      types.Int64Value(int64(teacher.ID)),
		Name:    types.StringValue(teacher.Name),
		Subject: types.StringValue(teacher.Subject),
		Salary:  types.Int64Value(int64(teacher.Salary)),
		Organ:   types.StringValue(teacher.Organ),
	}
	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
	tflog.Warn(ctx, "CALL createTeacherResource")
}

func (t *teacherResource) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	var plan teacherResourceModel
	diags := request.State.Get(ctx, &plan)
	response.Diagnostics = append(response.Diagnostics, diags...)
	if response.Diagnostics.HasError() {
		return
	}
	teacher, err := t.cli.GetTeacher(int(plan.ID.ValueInt64()))
	if err != nil {
		response.Diagnostics.AddError("fail to get teacher", err.Error())
		return
	}
	plan = teacherResourceModel{
		ID:      types.Int64Value(int64(teacher.ID)),
		Name:    types.StringValue(teacher.Name),
		Subject: types.StringValue(teacher.Subject),
		Salary:  types.Int64Value(int64(teacher.Salary)),
		Organ:   types.StringValue(teacher.Organ),
	}
	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	tflog.Warn(ctx, "CALL readTeacherResource")
}

func (t *teacherResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	var plan teacherResourceModel
	diags := request.Config.Get(ctx, &plan)
	response.Diagnostics = append(response.Diagnostics, diags...)
	if response.Diagnostics.HasError() {
		return
	}

	teacher := &model.Teacher{
		ID:      int(plan.ID.ValueInt64()),
		Name:    plan.Name.ValueString(),
		Subject: plan.Subject.ValueString(),
		Salary:  int(plan.Salary.ValueInt64()),
		Organ:   plan.Organ.ValueString(),
	}

	teacher, err := t.cli.UpdateTeacher(teacher)

	if err != nil {
		response.Diagnostics.AddError("fail to update teacher", err.Error())
		return
	}
	plan = teacherResourceModel{
		ID:      types.Int64Value(int64(teacher.ID)),
		Name:    types.StringValue(teacher.Name),
		Subject: types.StringValue(teacher.Subject),
		Salary:  types.Int64Value(int64(teacher.Salary)),
		Organ:   types.StringValue(teacher.Organ),
	}
	diags = response.State.Set(ctx, plan)
	response.Diagnostics.Append(diags...)
	tflog.Warn(ctx, "CALL updateTeacherResource")
}

func (t *teacherResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	var plan teacherResourceModel
	diags := request.State.Get(ctx, &plan)
	response.Diagnostics = append(response.Diagnostics, diags...)
	if response.Diagnostics.HasError() {
		return
	}
	if err := t.cli.DeleteTeacher(int(plan.ID.ValueInt64())); err != nil {
		response.Diagnostics.AddError("fail to delete teacher", err.Error())
	}
	tflog.Warn(ctx, "CALL deleteTeacherResource")
}

type teacherResourceModel struct {
	ID      types.Int64  `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Subject types.String `tfsdk:"subject"`
	Salary  types.Int64  `tfsdk:"salary"`
	Organ   types.String `tfsdk:"organ"`
}
