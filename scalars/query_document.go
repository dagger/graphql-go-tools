package scalars

import (
	"encoding/json"
	"regexp"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
	"github.com/dagger/graphql/language/kinds"
)

var queryDocOperatorRx = regexp.MustCompile(`^\$`)
var storedQueryDocOperatorRx = regexp.MustCompile(`^_`)

func replacePrefixedKeys(obj any, prefixRx *regexp.Regexp, replacement string) any {
	switch obj.(type) {
	case map[string]any:
		result := map[string]any{}
		for k, v := range obj.(map[string]any) {
			newKey := prefixRx.ReplaceAllString(k, replacement)
			result[newKey] = replacePrefixedKeys(v, prefixRx, replacement)
		}
		return result

	case []any:
		result := []any{}
		for _, v := range obj.([]any) {
			result = append(result, replacePrefixedKeys(v, prefixRx, replacement))
		}
		return result

	default:
		return obj
	}
}

func serializeQueryDocFn(value any) any {
	return replacePrefixedKeys(value, storedQueryDocOperatorRx, "$")
}

func parseValueQueryDocFn(value any) any {
	return replacePrefixedKeys(value, queryDocOperatorRx, "_")
}

func parseLiteralQueryDocFn(astValue ast.Value) any {
	var val any
	switch astValue.GetKind() {
	case kinds.StringValue:
		bvalue := []byte(astValue.GetValue().(string))
		if err := json.Unmarshal(bvalue, &val); err != nil {
			return nil
		}
		return replacePrefixedKeys(val, queryDocOperatorRx, "_")
	case kinds.ObjectValue:
		return parseLiteralJSONFn(astValue)
	}
	return nil
}

// ScalarQueryDocument a mongodb style query document
var ScalarQueryDocument = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:         "QueryDocument",
		Description:  "MongoDB style query document",
		Serialize:    serializeQueryDocFn,
		ParseValue:   parseValueQueryDocFn,
		ParseLiteral: parseLiteralQueryDocFn,
	},
)
