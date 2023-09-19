package tools

// taken from https://github.com/dagger/graphql/values.go
// since none of these functions are exported

import (
	"fmt"
	"math"
	"reflect"

	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/ast"
)

// Prepares an object map of argument values given a list of argument
// definitions and list of argument AST nodes.
func GetArgumentValues(argDefs []*graphql.Argument, argASTs []*ast.Argument, variableVariables map[string]any) (map[string]any, error) {

	argASTMap := map[string]*ast.Argument{}
	for _, argAST := range argASTs {
		if argAST.Name != nil {
			argASTMap[argAST.Name.Value] = argAST
		}
	}
	results := map[string]any{}
	for _, argDef := range argDefs {

		name := argDef.PrivateName
		var valueAST ast.Value
		if argAST, ok := argASTMap[name]; ok {
			valueAST = argAST.Value
		}

		value, err := valueFromAST(valueAST, argDef.Type, variableVariables)
		if err != nil || isNullish(value) {
			value = argDef.DefaultValue
		}

		// fix for checking that non nulls are not null
		typeString := argDef.Type.String()
		isNonNull := typeString[len(typeString)-1:] == "!"
		if isNonNull && isNullish(value) {
			return nil, fmt.Errorf("graphql input %q cannot be null", name)
		}

		if !isNullish(value) {
			results[name] = value
		}
	}
	return results, nil
}

// Returns true if a value is null, undefined, or NaN.
func isNullish(src any) bool {
	if src == nil {
		return true
	}
	value := reflect.ValueOf(src)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.String:
		// if src is ptr type and len(string)=0, it returns false
		if !value.IsValid() {
			return true
		}
	case reflect.Int:
		return math.IsNaN(float64(value.Int()))
	case reflect.Float32, reflect.Float64:
		return math.IsNaN(float64(value.Float()))
	}
	return false
}

/**
 * Produces a value given a GraphQL Value AST.
 *
 * A GraphQL type must be provided, which will be used to interpret different
 * GraphQL Value literals.
 *
 * | GraphQL Value        | JSON Value    |
 * | -------------------- | ------------- |
 * | Input Object         | Object        |
 * | List                 | Array         |
 * | Boolean              | Boolean       |
 * | String / Enum Value  | String        |
 * | Int / Float          | Number        |
 *
 */
func valueFromAST(valueAST ast.Value, ttype graphql.Input, variables map[string]any) (any, error) {
	if valueAST == nil {
		return nil, nil
	}
	// precedence: value > type
	if valueAST, ok := valueAST.(*ast.Variable); ok {
		if valueAST.Name == nil {
			return nil, fmt.Errorf("invalid variable")
		}

		var val any
		var found bool
		if variables != nil {
			val, found = variables[valueAST.Name.Value]
		}
		if !found {
			return nil, fmt.Errorf("missing variable: $%s", valueAST.Name.Value)
		}
		// Note: we're not doing any checking that this variable is correct. We're
		// assuming that this query has been validated and the variable usage here
		// is of the correct type.
		return val, nil
	}
	switch ttype := ttype.(type) {
	case *graphql.NonNull:
		return valueFromAST(valueAST, ttype.OfType, variables)
	case *graphql.List:
		values := []any{}
		if valueAST, ok := valueAST.(*ast.ListValue); ok {
			for _, itemAST := range valueAST.Values {
				recur, err := valueFromAST(itemAST, ttype.OfType, variables)
				if err != nil {
					return nil, err
				}
				values = append(values, recur)
			}
			return values, nil
		}
		recur, err := valueFromAST(valueAST, ttype.OfType, variables)
		if err != nil {
			return nil, err
		}
		return append(values, recur), nil
	case *graphql.InputObject:
		var (
			ok bool
			ov *ast.ObjectValue
			of *ast.ObjectField
		)
		if ov, ok = valueAST.(*ast.ObjectValue); !ok {
			return nil, fmt.Errorf("expected %T, found %T", ov, valueAST)
		}
		fieldASTs := map[string]*ast.ObjectField{}
		for _, of = range ov.Fields {
			if of == nil || of.Name == nil {
				continue
			}
			fieldASTs[of.Name.Value] = of
		}
		obj := map[string]any{}
		for name, field := range ttype.Fields() {
			var value any
			if of, ok = fieldASTs[name]; ok {
				recur, err := valueFromAST(of.Value, field.Type, variables)
				if err != nil {
					return nil, err
				}
				value = recur
			} else {
				value = field.DefaultValue
			}
			if !isNullish(value) {
				obj[name] = value
			}
		}
		return obj, nil
	case *graphql.Scalar:
		return ttype.ParseLiteral(valueAST)
	case *graphql.Enum:
		return ttype.ParseLiteral(valueAST)
	}

	return nil, fmt.Errorf("valueFromAST: unknown type %T", ttype)
}
