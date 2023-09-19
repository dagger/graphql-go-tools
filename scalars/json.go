package scalars

import (
	"fmt"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
	"github.com/dagger/graphql/language/kinds"
)

// ScalarJSON a scalar JSON type
var ScalarJSON = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "JSON",
		Description: "The `JSON` scalar type represents JSON values as specified by [ECMA-404](http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf)",
		Serialize: func(value any) (any, error) {
			return value, nil
		},
		ParseValue: func(value any) (any, error) {
			return value, nil
		},
		ParseLiteral: parseLiteralJSONFn,
	},
)

// recursively parse ast
func parseLiteralJSONFn(astValue ast.Value) (any, error) {
	switch kind := astValue.GetKind(); kind {
	// get value for primitive types
	case kinds.StringValue, kinds.BooleanValue, kinds.IntValue, kinds.FloatValue:
		return astValue.GetValue(), nil

	// make a map for objects
	case kinds.ObjectValue:
		obj := make(map[string]any)
		for _, v := range astValue.GetValue().([]*ast.ObjectField) {
			recur, err := parseLiteralJSONFn(v.Value)
			if err != nil {
				return nil, err
			}
			obj[v.Name.Value] = recur
		}
		return obj, nil

	// make a slice for lists
	case kinds.ListValue:
		list := make([]any, 0)
		for _, v := range astValue.GetValue().([]ast.Value) {
			recur, err := parseLiteralJSONFn(v)
			if err != nil {
				return nil, err
			}
			list = append(list, recur)
		}
		return list, nil

	default:
		return nil, fmt.Errorf("unknown kind %v", kind)
	}
}
