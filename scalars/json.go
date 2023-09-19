package scalars

import (
	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
	"github.com/dagger/graphql/language/kinds"
)

// ScalarJSON a scalar JSON type
var ScalarJSON = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "JSON",
		Description: "The `JSON` scalar type represents JSON values as specified by [ECMA-404](http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf)",
		Serialize: func(value any) any {
			return value
		},
		ParseValue: func(value any) any {
			return value
		},
		ParseLiteral: parseLiteralJSONFn,
	},
)

// recursively parse ast
func parseLiteralJSONFn(astValue ast.Value) any {
	switch kind := astValue.GetKind(); kind {
	// get value for primitive types
	case kinds.StringValue, kinds.BooleanValue, kinds.IntValue, kinds.FloatValue:
		return astValue.GetValue()

	// make a map for objects
	case kinds.ObjectValue:
		obj := make(map[string]any)
		for _, v := range astValue.GetValue().([]*ast.ObjectField) {
			obj[v.Name.Value] = parseLiteralJSONFn(v.Value)
		}
		return obj

	// make a slice for lists
	case kinds.ListValue:
		list := make([]any, 0)
		for _, v := range astValue.GetValue().([]ast.Value) {
			list = append(list, parseLiteralJSONFn(v))
		}
		return list

	// default to nil
	default:
		return nil
	}
}
