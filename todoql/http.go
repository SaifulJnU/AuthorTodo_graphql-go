package todoql

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// //////////for author controller///////////////
type TodoController struct {
	db *mongo.Client
}

func NewTodoController(db *mongo.Client) *TodoController {
	return &TodoController{
		db: db,
	}
}

type AuthorController struct {
	db *mongo.Client
}

func NewAuthorController(db *mongo.Client) *AuthorController {
	return &AuthorController{
		db: db,
	}
}
