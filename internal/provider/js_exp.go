package provider

import (
	"context"

	"github.com/dop251/goja"
)

func canCommitValue(ctx context.Context, data *diffStateItemsModel) (bool, string, error) {
	prog, err := goja.Compile("", data.CommitExp.ValueString(), true)
	if err != nil {
		return false, "JavaScript Syntax Error in commit_exp", err
	}

	vm := goja.New()

	vm.Set("is_initiated", data.IsInitiated.ValueBool())

	var values map[string]string
	data.Values.ElementsAs(ctx, &values, false)
	vm.Set("values", values)

	var last_values map[string]string
	data.LastValues.ElementsAs(ctx, &last_values, false)
	vm.Set("last_values", last_values)

	var created []string
	data.Created.ElementsAs(ctx, &created, false)
	vm.Set("created", created)

	var updated []string
	data.Updated.ElementsAs(ctx, &updated, false)
	vm.Set("updated", updated)

	var deleted []string
	data.Deleted.ElementsAs(ctx, &deleted, false)
	vm.Set("deleted", deleted)

	result, err := vm.RunProgram(prog)
	if err != nil {
		return false, "JavaScript Runtime Error in commit_exp", err
	}

	return result.Export() == true, "", nil
}
