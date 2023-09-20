package todoql

import (
	"go.mongodb.org/mongo-driver/mongo"
)


type AuthorController struct {
	db *mongo.Client
}

func NewAuthorController(db *mongo.Client) *AuthorController {
	return &AuthorController{
		db: db,
	}
}
