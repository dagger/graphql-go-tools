package tools

import (
	"github.com/dagger/graphql"
	"github.com/dagger/graphql/language/kinds"
)

// Resolver interface to a resolver configuration
type Resolver interface {
	getKind() string
}

// ResolverMap a map of resolver configurations.
// Accept generic interfaces and identify types at build
type ResolverMap map[string]any

// internal resolver map
type resolverMap map[string]Resolver

// FieldResolveMap map of field resolve functions
type FieldResolveMap map[string]*FieldResolve

// FieldResolve field resolver
type FieldResolve struct {
	Resolve   graphql.FieldResolveFn
	Subscribe graphql.FieldResolveFn
}

// ObjectResolver config for object resolver map
type ObjectResolver struct {
	IsTypeOf graphql.IsTypeOfFn
	Fields   FieldResolveMap
}

// GetKind gets the kind
func (c *ObjectResolver) getKind() string {
	return kinds.ObjectDefinition
}

// ScalarResolver config for a scalar resolve map
type ScalarResolver struct {
	Serialize    graphql.SerializeFn
	ParseValue   graphql.ParseValueFn
	ParseLiteral graphql.ParseLiteralFn
}

// GetKind gets the kind
func (c *ScalarResolver) getKind() string {
	return kinds.ScalarDefinition
}

// InterfaceResolver config for interface resolve
type InterfaceResolver struct {
	ResolveType graphql.ResolveTypeFn
	Fields      FieldResolveMap
}

// GetKind gets the kind
func (c *InterfaceResolver) getKind() string {
	return kinds.InterfaceDefinition
}

// UnionResolver config for interface resolve
type UnionResolver struct {
	ResolveType graphql.ResolveTypeFn
}

// GetKind gets the kind
func (c *UnionResolver) getKind() string {
	return kinds.UnionDefinition
}

// EnumResolver config for enum values
type EnumResolver struct {
	Values map[string]any
}

// GetKind gets the kind
func (c *EnumResolver) getKind() string {
	return kinds.EnumDefinition
}
