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
		Serialize: func(value any) any {
			valStr := fmt.Sprintf("%v", value)
			return valStr == "true" || valStr == "1"
		},
		ParseValue: func(value any) any {
			b, ok := value.(bool)
			if !ok {
				return "false"
			} else if b {
				return "true"
			}
			return "false"
		},
		ParseLiteral: func(astValue ast.Value) any {
			value := astValue.GetValue()
			b, ok := value.(bool)
			if !ok {
				return "false"
			} else if b {
				return "true"
			}
			return "false"
		},
	},
)
