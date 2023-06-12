package provider

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
)

type CommitExpError struct {
	Message string
	Err     error
}

func (e *CommitExpError) Error() string {
	return fmt.Sprintf("%s in `commit_exp`", e.Message)
}

func (e *CommitExpError) Unwrap() error {
	return e.Err
}

func canCommitValue(ctx context.Context, data *diffModel) (bool, error) {
	prog, err := goja.Compile("", data.CommitExp.ValueString(), true)
	if err != nil {
		return false, &CommitExpError{"JavaScript Syntax Error", err}
	}

	vm := goja.New()

	vm.Set("is_initiated", data.IsInitiated.ValueBool())

	vm.Set("values", asValue[stringMap](ctx, data.Values))
	vm.Set("last_values", asValue[stringMap](ctx, data.LastValues))

	vm.Set("created", asValue[stringList](ctx, data.Created))
	vm.Set("updated", asValue[stringList](ctx, data.Updated))
	vm.Set("deleted", asValue[stringList](ctx, data.Deleted))

	result, err := vm.RunProgram(prog)
	if err != nil {
		return false, &CommitExpError{"JavaScript Runtime Error", err}
	}

	return result.Export() == true, nil
}
