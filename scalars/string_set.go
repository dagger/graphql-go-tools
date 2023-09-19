package scalars

import (
	"reflect"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
)

func ensureArray(value any) any {
	switch kind := reflect.TypeOf(value).Kind(); kind {
	case reflect.Slice, reflect.Array:
		return value
	default:
		if reflect.ValueOf(value).IsNil() {
			return nil
		}
		return []any{value}
	}
}

func serializeStringSetFn(value any) any {
	switch kind := reflect.TypeOf(value).Kind(); kind {
	case reflect.Slice, reflect.Array:
		v := reflect.ValueOf(value)
		if v.Len() == 1 {
			return v.Index(0).Interface()
		}
		return value
	default:
		return []any{}
	}
}

// ScalarStringSet allows string or array of strings
// stores as an array of strings
var ScalarStringSet = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "StringSet",
		Description: "StringSet allows either a string or list of strings",
		Serialize:   serializeStringSetFn,
		ParseValue: func(value any) any {
			return ensureArray(value)
		},
		ParseLiteral: func(astValue ast.Value) any {
			return ensureArray(parseLiteralJSONFn(astValue))
		},
	},
)
