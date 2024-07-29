package graphql

import (
	"github.com/graphql-go/graphql"
)

// RootQuery represents the root GraphQL query.
var RootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"banking360": &graphql.Field{
				Type: Banking360QueriesType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &Banking360Queries{}, nil
				},
			},
			// Add other queries as needed
		},
	},
)

// RootMutation represents the root GraphQL query.
var RootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"banking360Mutations": &graphql.Field{
				Type: Banking360MutationsType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return &Banking360Mutations{}, nil
				},
			},
			// Add other queries as needed
		},
	},
)