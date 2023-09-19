package scalars

import (
	"fmt"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
)

// ScalarBoolString converts boolean to a string
var ScalarBoolString = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "BoolString",
		Description: "BoolString converts a boolean to/from a string",
		Serialize: func(value any) (any, error) {
			valStr := fmt.Sprintf("%v", value)
			return valStr == "true" || valStr == "1", nil
		},
		ParseValue: func(value any) (any, error) {
			b, ok := value.(bool)
			if !ok {
				return "false", nil
			} else if b {
				return "true", nil
			}
			return "false", nil
		},
		ParseLiteral: func(astValue ast.Value) (any, error) {
			value := astValue.GetValue()
			b, ok := value.(bool)
			if !ok {
				return "false", nil
			} else if b {
				return "true", nil
			}
			return "false", nil
		},
	},
)
