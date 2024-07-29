package main

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"banking360/config"
	gql "banking360/internal/graphql"
	
)

func main() {
	cfg := config.NewConfig()
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: gql.RootQuery,
		Mutation:     gql.RootMutation,
		
	})
	if err != nil {
		panic(err)
	}
	// Create a GraphQL handler for HTTP requests
	graphqlHandler := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false, // Disable GraphiQL for subscriptions endpoint
		Playground: true,
	})

	// Serve GraphQL API at /graphql endpoint
	http.Handle("/graphql", graphqlHandler)

	// Create the server
	server := &http.Server{
		Addr: fmt.Sprintf(":%v", cfg.GraphQLPort),
	}

	// Start the HTTP server
	fmt.Printf("Server is running at http://localhost:%v/graphql\n", cfg.GraphQLPort)
	server.ListenAndServe()
}


