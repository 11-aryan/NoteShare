package auth

import (
	"NotesApp/Utils/mongodb"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = mongodb.GetCollection(mongodb.DB, "users")
var validate = validator.New()

type User struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name" validate:"required"`
	Email    string             `json:"email" bson:"email" validate:"required"`
	Password string             `json:"-" bson:"password" validate:"required"`
}
