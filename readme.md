
# AuthorTodo
```
query {
  findTodosByAuthorName(name: "John Doe") {
    Id
    Text
    Done
    AuthorId
  }
}
```
```
mutation {
  deleteTodo(id: "TODO_ID_HERE") {
    Id
  }
}

```
```
mutation {
  createTodo(text: "Buy groceries", authorId: "YOUR_AUTHOR_ID_HERE") {
    Id
    Text
    Done
    AuthorId
  }
}


```
```
mutation {
  createTodo(text: "Buy groceries", authorId: "YOUR_AUTHOR_ID_HERE") {
    Id
    Text
    Done
    AuthorId
  }
}
```
```
query {
  getAuthorAndTodos(authorId: "YOUR_AUTHOR_ID_HERE") {
    author {
      Id
      Name
    }
    todos {
      Id
      Text
      Done
    }
  }
}

```
```
mutation {
  createAuthor(name: "Rahim") {
    Id
    Name
  }
}

```


```go
//single page implementation
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var todoCollection *mongo.Collection

type Todo struct {
	Id   string `json:"id" bson:"_id,omitempty"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"Id": &graphql.Field{
			Type: graphql.String,
		},
		"Text": &graphql.Field{
			Type: graphql.String,
		},
		"Done": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createTodo": &graphql.Field{
			Type:        todoType,
			Description: "Create new todo",
			Args: graphql.FieldConfigArgument{
				"text": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				text, _ := params.Args["text"].(string)

				newTodo := Todo{
					Text: text,
					Done: false,
				}

				_, err := todoCollection.InsertOne(context.Background(), newTodo)
				if err != nil {
					return nil, err
				}

				return newTodo, nil
			},
		},
		"updateTodo": &graphql.Field{
			Type:        todoType,
			Description: "Update existing todo, mark it Done or not Done",
			Args: graphql.FieldConfigArgument{
				"done": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				doneParam, _ := params.Args["done"].(bool)
				idParam, _ := params.Args["id"].(string)

				filter := bson.M{"_id": idParam}
				update := bson.M{"$set": bson.M{"done": doneParam}}

				_, err := todoCollection.UpdateOne(context.Background(), filter, update)
				if err != nil {
					return nil, err
				}

				var updatedTodo Todo
				err = todoCollection.FindOne(context.Background(), filter).Decode(&updatedTodo)
				if err != nil {
					return nil, err
				}

				return updatedTodo, nil
			},
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"todo": &graphql.Field{
			Type:        todoType,
			Description: "Get single todo",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(string)
				if !isOK {
					return nil, nil
				}

				var todo Todo
				err := todoCollection.FindOne(context.Background(), bson.M{"_id": idQuery}).Decode(&todo)
				if err != nil {
					return nil, err
				}

				return todo, nil
			},
		},
		"todoList": &graphql.Field{
			Type:        graphql.NewList(todoType),
			Description: "List of todos",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var all []Todo
				cursor, err := todoCollection.Find(context.Background(), bson.M{})
				if err != nil {
					return nil, err
				}
				defer cursor.Close(context.Background())

				err = cursor.All(context.Background(), &all)
				if err != nil {
					return nil, err
				}

				return all, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Disconnect(ctx)

	todoCollection = client.Database("todoDB").Collection("todos")

	router := gin.Default()

	router.POST("/graphql", func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		query, exists := requestBody["query"].(string)
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query not provided"})
			return
		}

		params := graphql.Params{
			Schema:        schema,
			RequestString: query,
		}

		result := graphql.Do(params)
		c.JSON(http.StatusOK, result)
	})

	fmt.Println("Now server is running on port 8081")
	fmt.Println("Create new todo: Send a POST request to 'http://localhost:8081/graphql' with the following JSON body:")
	fmt.Println(`{
  "query": "mutation { createTodo(text: \"My new todo\") { Id Text Done } }"
}`)
	fmt.Println("Update todo: Send a POST request to 'http://localhost:8081/graphql' with the following JSON body:")
	fmt.Println(`{
  "query": "mutation { updateTodo(id: \"your-todo-id\", done: true) { Id Text Done } }"
}`)
	fmt.Println("Get single todo: Send a POST request to 'http://localhost:8081/graphql' with the following JSON body:")
	fmt.Println(`{
  "query": "{ todo(id: \"your-todo-id\") { Id Text Done } }"
}`)
	fmt.Println("Load todo list: Send a POST request to 'http://localhost:8081/graphql' with the following JSON body:")
	fmt.Println(`{
  "query": "{ todoList { Id Text Done } }"
}`)
	fmt.Println("Access the web app via a browser at 'http://localhost:8081'")

	router.Run(":8081")
}

```
