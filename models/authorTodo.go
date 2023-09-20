package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthorTodo struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Text     string             `json:"text" bson:"text"`
	Done     bool               `json:"done" bson:"done"`
	AuthorID primitive.ObjectID `json:"authorId" bson:"authorId"`
}
