package user

import (
	notesPkg "NotesApp/Modules/notes"
	"NotesApp/Utils/mongodb"
	"NotesApp/Utils/response"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Create a new user
func SignUp(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user User
	resp := response.Wrap(c)
	//validating the request body
	err := c.BodyParser(&user)
	if err != nil {
		return resp.Error(err)
	}
	//using the validator library to validate required fields
    err = validate.Struct(&user); 
	if err != nil {
        return resp.Error(err)
    }
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return resp.Error(err)
	}
	newUser := User {
		Name: user.Name,
		Email: user.Email,
		Password: hashedPassword,
	}
	result, err := userCollection.InsertOne(ctx, newUser)
    if err != nil {
        return resp.Error(err)
    }
    return resp.Data(result)
}

// Get a User with given user_id
func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	userId := c.Params("id")
	var user User 
	// Converting the string user_id to mongoDB objectID 
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return resp.Error(err)
	}
	err = userCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return resp.Error(err)
	}
	return resp.Data(user)
}

func GetNotes(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c) 
	userID := c.Params("id")
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return resp.Error(err)
	}
	notesCollection := mongodb.GetCollection(mongodb.DB, "notes")
	filter := bson.M{
		"created.userID": userObjID,
	}
	var notes []notesPkg.Note
	cursor, err := notesCollection.Find(ctx, filter)
	if err != nil {
		log.Error(err)
		return resp.Error(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var note notesPkg.Note
		if err = cursor.Decode(&note); err != nil {
			return resp.Error(err)
		}
		notes = append(notes, note)
	}
	if err := cursor.Err(); err != nil {
        return resp.Error(err)
    }
	return resp.Data(notes)
}

//Function to hash the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}