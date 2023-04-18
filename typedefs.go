package tools

import (
	"fmt"

	"github.com/dagger/graphql/language/ast"
	"github.com/dagger/graphql/language/parser"
	"github.com/dagger/graphql/language/source"
)

// ConcatenateTypeDefs combines one ore more typeDefs into an ast Document
func (c *ExecutableSchema) ConcatenateTypeDefs() (*ast.Document, error) {
	switch defs := c.TypeDefs.(type) {
	case string:
		return c.concatenateTypeDefs([]string{defs})
	case []string:
		return c.concatenateTypeDefs(defs)
	case func() []string:
		return c.concatenateTypeDefs(defs())
	}
	return nil, fmt.Errorf("unsupported TypeDefs value. Must be one of string, []string, or func() []string")
}

// appends all type definitions together into one document
func (c *ExecutableSchema) concatenateTypeDefs(typeDefs []string) (*ast.Document, error) {
	doc := ast.NewDocument(nil)

	for _, defs := range typeDefs {
		sub, err := parser.Parse(parser.ParseParams{
			Source: &source.Source{
				Body: []byte(defs),
				Name: "GraphQL",
			},
		})
		if err != nil {
			return nil, err
		}

		doc.Definitions = append(doc.Definitions, sub.Definitions...)
	}

	return doc, nil
}
