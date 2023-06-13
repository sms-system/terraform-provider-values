package provider

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/maps"
)

var _ resource.Resource = &DiffResource{}

func NewDiffResource() resource.Resource {
	return &DiffResource{}
}

type DiffResource struct{}

func (n *DiffResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_diff"
}

func (n *DiffResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
			The resource detects changes in values using their identifiers.
		`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier of resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"is_initiated": schema.BoolAttribute{
				MarkdownDescription: "`true` on resource creation",
				Computed:            true,
			},

			"values": schema.MapAttribute{
				Description: "Items for tracking differences. The keys are here to identify a unique element",
				Required:    true,
				ElementType: types.StringType,
			},

			"last_values": schema.MapAttribute{
				Description: "Items from previous state",
				ElementType: types.StringType,
				Computed:    true,
			},

			"created": schema.ListAttribute{
				Description: "New added items",
				ElementType: types.StringType,
				Computed:    true,
			},

			"updated": schema.ListAttribute{
				Description: "Items whose value has been changed",
				ElementType: types.StringType,
				Computed:    true,
			},

			"deleted": schema.ListAttribute{
				Description: "Deleted items",
				ElementType: types.StringType,
				Computed:    true,
			},

			"commit_exp": schema.StringAttribute{
				MarkdownDescription: "JS expression. If it returns `true`, `last_values` will be updated. Aviable global variables: `values`, `last_values`, `created`, `updated`, `deleted`, `is_initiated`. Default: `\"true\"`",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("true"),
			},

			"is_value_commited": schema.BoolAttribute{
				MarkdownDescription: "`true` on successful `last_values` update",
				Computed:            true,
			},
		},
	}
}

type diffModel struct {
	Id types.String `tfsdk:"id"`

	IsInitiated types.Bool `tfsdk:"is_initiated"`

	Values     types.Map `tfsdk:"values"`
	LastValues types.Map `tfsdk:"last_values"`

	Created types.List `tfsdk:"created"`
	Updated types.List `tfsdk:"updated"`
	Deleted types.List `tfsdk:"deleted"`

	CommitExp       types.String `tfsdk:"commit_exp"`
	IsValueCommited types.Bool   `tfsdk:"is_value_commited"`
}

func (r *DiffResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *diffModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue("diff")

	data.IsInitiated = types.BoolValue(true)

	data.Created, _ = types.ListValueFrom(ctx, types.StringType, maps.Keys(data.Values.Elements()))
	data.Updated, _ = types.ListValueFrom(ctx, types.StringType, []string{})
	data.Deleted, _ = types.ListValueFrom(ctx, types.StringType, []string{})
	data.LastValues, _ = types.MapValueFrom(ctx, types.StringType, map[string]string{})

	isValueCommited, err := canCommitValue(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			errors.Unwrap(err).Error(),
		)
		return
	}

	data.IsValueCommited = types.BoolValue(isValueCommited)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DiffResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *DiffResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *diffModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.IsInitiated = types.BoolValue(false)

	current := data.Values
	previous := state.Values
	if !state.IsValueCommited.ValueBool() {
		previous = state.LastValues
	}

	currentItems := current.Elements()
	previousItems := previous.Elements()

	created, updated, deleted := calculateDiff(currentItems, previousItems)

	data.LastValues, _ = types.MapValueFrom(ctx, types.StringType, previousItems)
	data.Created, _ = types.ListValueFrom(ctx, types.StringType, created)
	data.Updated, _ = types.ListValueFrom(ctx, types.StringType, updated)
	data.Deleted, _ = types.ListValueFrom(ctx, types.StringType, deleted)

	isValueCommited, err := canCommitValue(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			errors.Unwrap(err).Error(),
		)
		return
	}

	data.IsValueCommited = types.BoolValue(isValueCommited)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DiffResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func calculateDiff(currentItems, previousItems map[string]attr.Value) ([]string, []string, []string) {
	created := []string{}
	updated := []string{}
	deleted := []string{}

	// TODO: Extract to func
	for k, v := range currentItems {
		val, ok := previousItems[k]
		if !ok {
			created = append(created, k)
		} else if v != val {
			updated = append(updated, k)
		}
	}

	for k := range previousItems {
		_, ok := currentItems[k]
		if !ok {
			deleted = append(deleted, k)
		}
	}

	return created, updated, deleted
}
