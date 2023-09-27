package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/saifuljnu/todo/config"
	db "github.com/saifuljnu/todo/db/mongo"
	"github.com/saifuljnu/todo/todoql"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	config.SetEnvionment()

}

func main() {
	// Initialize MongoDB
	client, _, err := db.SetupMongoDB()
	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	authorCon := todoql.NewAuthorController(client)
	//todoCon := todoql.NewTodoController(client)

	// router.POST("/graphql/todo", func(c *gin.Context) {
	// 	todoql.TodoGraphQLHandler(*todoCon).ServeHTTP(c.Writer, c.Request)
	// })

	router.POST("/graphql/author", func(c *gin.Context) {
		todoql.AuthorGraphQLHandler(*authorCon).ServeHTTP(c.Writer, c.Request)
	})

	// Start the server
	fmt.Println("Now the server is running on port 8081")
	router.Run(":8081")
}
