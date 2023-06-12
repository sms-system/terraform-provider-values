package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DiffStateItemsResource{}

func NewDiffStateItemsResource() resource.Resource {
	return &DiffStateItemsResource{}
}

type DiffStateItemsResource struct{}

func (n *DiffStateItemsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_items"
}

func (n *DiffStateItemsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"values": schema.MapAttribute{
				Description: "Items for tracking differences. The keys are here to identify a unique element",
				Required: true,
				ElementType: types.StringType,
			},

			"previous": schema.MapAttribute{
				Description: "Items from previous state",
				ElementType: types.StringType,
				Computed: true,
			},

			"new": schema.ListAttribute{
				Description: "New added items",
				ElementType: types.StringType,
				Computed: true,
			},

			"updated": schema.ListAttribute{
				Description: "Items whose value has been changed",
				ElementType: types.StringType,
				Computed: true,
			},

			"deleted": schema.ListAttribute{
				Description: "Deleted items",
				ElementType: types.StringType,
				Computed: true,
			},
		},
	}
}

type diffStateItemsModel struct {
	Id       types.String `tfsdk:"id"`
	Values   types.Map    `tfsdk:"values"`
	Previous types.Map    `tfsdk:"previous"`
	New      types.List   `tfsdk:"new"`
	Updated  types.List   `tfsdk:"updated"`
	Deleted  types.List   `tfsdk:"deleted"`
}

func (r *DiffStateItemsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *diffStateItemsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue("diff")

	data.New, _      = types.ListValueFrom(ctx, types.StringType, []string{})
	data.Updated, _  = types.ListValueFrom(ctx, types.StringType, []string{})
	data.Deleted, _  = types.ListValueFrom(ctx, types.StringType, []string{})
	data.Previous, _ = types.MapValueFrom(ctx, types.StringType, map[string]string{})
	
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DiffStateItemsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *DiffStateItemsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state *diffStateItemsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Values.Equal(state.Values) {
		current := data.Values.Elements()
		previous := state.Values.Elements()

		var new []string
		var updated []string
		var deleted []string

		for k, v := range current {
			val, ok := previous[k]
			if !ok {
				new = append(new, k)
			} else if v != val {
				updated = append(updated, k)
			}
		}

		for k, _ := range previous {
			_, ok := current[k]
			if !ok {
				deleted = append(deleted, k)
			}
		}

		data.Previous, _ = types.MapValueFrom(ctx, types.StringType, previous)
		data.New, _      = types.ListValueFrom(ctx, types.StringType, new)
		data.Updated, _  = types.ListValueFrom(ctx, types.StringType, updated)
		data.Deleted, _  = types.ListValueFrom(ctx, types.StringType, deleted)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DiffStateItemsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}