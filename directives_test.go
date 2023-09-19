package tools

import (
	"testing"

	"github.com/dagger/graphql"
)

func TestDirectives(t *testing.T) {
	typeDefs := `
directive @test(message: String) on FIELD_DEFINITION

type Foo {
	name: String!
	description: String
}

type Query {
	foos(
		name: String
	): [Foo] @test(message: "foobar")
}
`

	// create some data
	foos := []map[string]any{
		{
			"name":        "foo",
			"description": "a foo",
		},
	}

	// make the schema
	schema, err := MakeExecutableSchema(ExecutableSchema{
		TypeDefs: typeDefs,
		Resolvers: ResolverMap{
			"Query": &ObjectResolver{
				Fields: FieldResolveMap{
					"foos": &FieldResolve{
						Resolve: func(p graphql.ResolveParams) (any, error) {
							return foos, nil
						},
					},
				},
			},
		},
		SchemaDirectives: SchemaDirectiveVisitorMap{
			"test": &SchemaDirectiveVisitor{
				VisitFieldDefinition: func(v VisitFieldDefinitionParams) error {
					resolveFunc := v.Config.Resolve
					v.Config.Resolve = func(p graphql.ResolveParams) (any, error) {
						result, err := resolveFunc(p)
						if err != nil {
							return result, err
						}
						res := result.([]map[string]any)
						res0 := res[0]
						res0["description"] = v.Args["message"]
						return res, nil
					}

					return nil
				},
			},
		},
	})

	if err != nil {
		t.Error(err)
		return
	}

	// perform a query
	r := graphql.Do(graphql.Params{
		Schema: schema,
		RequestString: `query Query {
			foos(name:"foo") {
				name
				description
			}
		}`,
	})

	if r.HasErrors() {
		t.Error(r.Errors)
		return
	}

	d := r.Data.(map[string]any)
	fooResult := d["foos"]
	foos0 := fooResult.([]any)[0]
	foos0Desc := foos0.(map[string]any)["description"]
	if foos0Desc.(string) != "foobar" {
		t.Error("failed to set field with directive")
		return
	}
}
