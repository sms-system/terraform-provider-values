package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type mapValueOrListValue interface {
	ElementsAs(context.Context, interface{}, bool) diag.Diagnostics
}

type stringMap map[string]string
type stringList []string

func asValue[T stringMap | stringList](ctx context.Context, val mapValueOrListValue) T {
	var res T
	val.ElementsAs(ctx, &res, false)
	return res
}
