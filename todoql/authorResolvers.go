package todoql

import (
	"context"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/saifuljnu/todo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	authorCollection     *mongo.Collection
	authorTodoCollection *mongo.Collection
)

func SetAuthorCollection(collection *mongo.Collection) {
	authorCollection = collection
}

func SetAuthorTodoCollection(collection *mongo.Collection) {
	authorTodoCollection = collection
}

// Resolver functions with additional arguments
func resolveCreateAuthor(params graphql.ResolveParams) (interface{}, error) {
	name, _ := params.Args["name"].(string)

	newAuthor := models.Author{
		Name: name,
	}

	_, err := authorCollection.InsertOne(context.Background(), newAuthor)
	if err != nil {
		return nil, err
	}

	return newAuthor, nil
}

func resolveAuthorList(params graphql.ResolveParams) (interface{}, error) {
	var all []models.Author
	cursor, err := authorCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &all)
	if err != nil {
		return nil, err
	}

	return all, nil
}
func resolveGetAuthorAndTodos(params graphql.ResolveParams) (interface{}, error) {
	authorID, isOK := params.Args["authorId"].(string)
	if !isOK {
		fmt.Println("Author ID is missing or invalid")
		return nil, nil
	}
	authorHexID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		fmt.Println("Error parsing author ID:", err)
		return nil, err
	}
	// Fetch Author data
	var author models.Author
	authorQuery := bson.M{"_id": authorHexID}

	err = authorCollection.FindOne(context.Background(), authorQuery).Decode(&author)
	if err != nil {
		fmt.Println("Error fetching author:", err)
		return nil, err
	}

	// Fetch Todos data associated with the author
	fmt.Println(authorID, " ", authorHexID)

	todosQuery := bson.M{"authorId": authorHexID}

	var todos []models.AuthorTodo
	cursor, err := authorTodoCollection.Find(context.Background(), todosQuery)
	if err != nil {
		fmt.Println("Error fetching todos:", err)
		return nil, err
	}
	fmt.Println(cursor)

	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &todos)
	if err != nil {
		fmt.Println("Error decoding todos:", err)
		return nil, err
	}

	result := map[string]interface{}{
		"author": author,
		"todos":  todos,
	}
	return result, nil
}

func resolveCreateAuthorTodo(params graphql.ResolveParams) (interface{}, error) {
	text, _ := params.Args["text"].(string)
	authorID, _ := params.Args["authorId"].(string)
	authorHexID, _ := primitive.ObjectIDFromHex(authorID)
	newTodo := models.AuthorTodo{
		Text:     text,
		Done:     false,
		AuthorID: authorHexID,
	}

	_, err := authorTodoCollection.InsertOne(context.Background(), newTodo)
	if err != nil {
		return nil, err
	}

	return newTodo, nil
}

func resolveUpdateAuthorTodo(params graphql.ResolveParams) (interface{}, error) {
	textParam, textParamOK := params.Args["text"].(string)
	doneParam, doneParamOK := params.Args["done"].(bool)
	idParam, _ := params.Args["id"].(string)

	idParamHex, _ := primitive.ObjectIDFromHex(idParam)

	// Check if at least one of the update parameters is provided
	if !textParamOK && !doneParamOK {
		return nil, nil
	}

	filter := bson.M{"_id": idParamHex}

	// Define update operations based on provided parameters
	update := bson.M{}

	if textParamOK {
		update["$set"] = bson.M{"text": textParam}
	}

	if doneParamOK {
		update["$set"] = bson.M{"done": doneParam}
	}

	// Perform the update operation
	_, err := authorTodoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Fetch and return the updated todo item
	var updatedTodo models.AuthorTodo
	err = authorTodoCollection.FindOne(context.Background(), filter).Decode(&updatedTodo)
	if err != nil {
		return nil, err
	}

	return updatedTodo, nil
}

func resolveDeleteAuthorTodo(params graphql.ResolveParams) (interface{}, error) {
	idParam, _ := params.Args["id"].(string)

	fmt.Printf("Deleting todo with ID: %s\n", idParam)

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		fmt.Printf("Error converting ID: %v\n", err)
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	result, err := authorTodoCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Printf("Error deleting todo: %v\n", err)
		return nil, err
	}

	if result.DeletedCount == 0 {
		errMsg := fmt.Sprintf("Todo with ID %s not found", idParam)
		fmt.Println(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	fmt.Println("Todo deleted successfully")
	return nil, nil
}

func resolveFindTodosByAuthorName(params graphql.ResolveParams) (interface{}, error) {
	// Retrieve the author's name from the query arguments
	name, isOK := params.Args["name"].(string)
	if !isOK {
		return nil, nil
	}

	// Find the author with the given name
	var author models.Author
	authorQuery := bson.M{"name": name}
	err := authorCollection.FindOne(context.Background(), authorQuery).Decode(&author)
	if err != nil {
		return nil, err
	}

	// Find todos by author's ID
	todosQuery := bson.M{"authorId": author.ID}
	var todos []models.AuthorTodo
	cursor, err := authorTodoCollection.Find(context.Background(), todosQuery)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &todos)
	if err != nil {
		return nil, err
	}

	return todos, nil
}
