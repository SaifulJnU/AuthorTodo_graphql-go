// handler.go
package todoql

import (
	"net/http"

	"github.com/graphql-go/handler"
)

// AuthorGraphQLHandler returns an HTTP handler for the Author GraphQL schema.
func AuthorGraphQLHandler(ac AuthorController) http.Handler {
	// Create a new schema using the NewAuthorSchema function from schema.go
	authorCollection := ac.db.Database("todoDB").Collection("author")
	authorTodoCollection := ac.db.Database("todoDB").Collection("todos")

	authorSchema := NewAuthorSchema()

	// Create a new GraphQL handler for the Author schema
	graphQLHandler := handler.New(&handler.Config{
		Schema: &authorSchema,
		Pretty: true,
		//GraphiQL: true, // enable the GraphiQL web interface
	})

	// Pass the authorCollection to the resolver functions
	SetAuthorCollection(authorCollection)
	SetAuthorTodoCollection(authorTodoCollection)
	return graphQLHandler
}

// TodoGraphQLHandler returns an HTTP handler for the Todo GraphQL schema.
// func TodoGraphQLHandler(td TodoController) http.Handler {
// 	// Create a new schema using the NewSchema function from schema.go
// 	todoCollection := td.db.Database("todoDB").Collection("todos")
// 	todoSchema := NewTodoSchema()

// 	// Create a new GraphQL handler for the Todo schema
// 	graphQLHandler := handler.New(&handler.Config{
// 		Schema: &todoSchema,
// 		Pretty: true, // Set this to true for pretty JSON responses in the playground
// 		//GraphiQL: true, // Set this to true to enable the GraphiQL web interface
// 	})

// 	// Pass the todoCollection to the resolver functions
// 	SetTodoCollection(todoCollection)
// 	//SetAuthorTodoCollection(todoCollection)

// 	return graphQLHandler
// }
