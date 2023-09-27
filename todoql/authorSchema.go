package todoql

import (
	"github.com/graphql-go/graphql"
)

var authorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Author",
	Fields: graphql.Fields{
		"Id": &graphql.Field{
			Type: graphql.String,
		},
		"Name": &graphql.Field{
			Type: graphql.String,
		},
	},
})
var authorTodoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthorTodo",
	Fields: graphql.Fields{
		"Id": &graphql.Field{
			Type: graphql.String,
		},
		"Text": &graphql.Field{
			Type: graphql.String,
		},
		"Done": &graphql.Field{
			Type: graphql.String,
		},
		"AuthorId": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var authorRootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createAuthor": &graphql.Field{
			Type:        authorType,
			Description: "Create new author",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveCreateAuthor(params)
			},
		},
		"createTodo": &graphql.Field{
			Type: authorTodoType,
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"authorId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveCreateAuthorTodo(params)
			},
		},
		"updateTodo": &graphql.Field{
			Type:        authorTodoType,
			Description: "Update existing todo, mark it Done or not Done",
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"done": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveUpdateAuthorTodo(params)
			},
		},
		"deleteTodo": &graphql.Field{
			Type:        authorTodoType,
			Description: "Delete a todo by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String), // Change the argument type to ID
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveDeleteAuthorTodo(params)
			},
		},
	},
})

var authorRootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{

		"authorList": &graphql.Field{
			Type:        graphql.NewList(authorType),
			Description: "List of author",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolveAuthorList(p)
			},
		},

		"getAuthorAndTodos": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "AuthorAndTodos",
				Fields: graphql.Fields{
					"author": &graphql.Field{
						Type: authorType,
					},
					"todos": &graphql.Field{
						Type: graphql.NewList(authorTodoType),
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				"authorId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolveGetAuthorAndTodos(p)
			},
		},
		"findTodosByAuthorName": &graphql.Field{
			Type:        graphql.NewList(authorTodoType),
			Description: "Find todos by author's name",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return resolveFindTodosByAuthorName(params)
			},
		},

		"getAuthorAndTodoss": &graphql.Field{
			Type: graphql.NewObject(graphql.ObjectConfig{
				Name: "AuthorAndTodos",
				Fields: graphql.Fields{
					"author": &graphql.Field{
						Type: authorType,
					},
					"todos": &graphql.Field{
						Type: graphql.NewList(authorTodoType),
					},
				},
			}),
			Args: graphql.FieldConfigArgument{
				"authorId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return resolveGetAuthorAndTodoss(p)
			},
		},
	},
})

var authorSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    authorRootQuery,
	Mutation: authorRootMutation,
})

func NewAuthorSchema() graphql.Schema {
	return authorSchema
}
