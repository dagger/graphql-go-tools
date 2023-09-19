package scalars

import (
	"reflect"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
)

func ensureArray(value any) (any, error) {
	switch kind := reflect.TypeOf(value).Kind(); kind {
	case reflect.Slice, reflect.Array:
		return value, nil
	default:
		if reflect.ValueOf(value).IsNil() {
			return nil, nil
		}
		return []any{value}, nil
	}
}

func serializeStringSetFn(value any) (any, error) {
	switch kind := reflect.TypeOf(value).Kind(); kind {
	case reflect.Slice, reflect.Array:
		v := reflect.ValueOf(value)
		if v.Len() == 1 {
			return v.Index(0).Interface(), nil
		}
		return value, nil
	default:
		return []any{}, nil
	}
}

// ScalarStringSet allows string or array of strings
// stores as an array of strings
var ScalarStringSet = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "StringSet",
		Description: "StringSet allows either a string or list of strings",
		Serialize:   serializeStringSetFn,
		ParseValue:  ensureArray,
		ParseLiteral: func(astValue ast.Value) (any, error) {
			val, err := parseLiteralJSONFn(astValue)
			if err != nil {
				return nil, err
			}
			return ensureArray(val)
		},
	},
)
