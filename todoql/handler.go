// handler.go
package todoql

import (
	"net/http"

	"github.com/graphql-go/handler"
)

// This returns an HTTP handler for the Author GraphQL schema.
func AuthorGraphQLHandler(ac AuthorController) http.Handler {

	authorCollection := ac.db.Database("todoDB").Collection("author")
	authorTodoCollection := ac.db.Database("todoDB").Collection("todos")

	authorSchema := NewAuthorSchema()

	// Create a new GraphQL handler for the Author schema
	graphQLHandler := handler.New(&handler.Config{
		Schema: &authorSchema,
		Pretty: true,
		//GraphiQL: true, //GraphiQL for web interface
	})

	// Pass the authorCollection to the resolver functions
	SetAuthorCollection(authorCollection)
	SetAuthorTodoCollection(authorTodoCollection)
	return graphQLHandler
}

// func TodoGraphQLHandler(td TodoController) http.Handler {
// 	todoCollection := td.db.Database("todoDB").Collection("todos")
// 	todoSchema := NewTodoSchema()
// 	graphQLHandler := handler.New(&handler.Config{
// 		Schema: &todoSchema,
// 		Pretty: true,
// 		//GraphiQL: true,
// 	})
// 	SetTodoCollection(todoCollection)
// 	//SetAuthorTodoCollection(todoCollection)
// 	return graphQLHandler
// }
