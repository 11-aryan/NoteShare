package notes

import (
	"NotesApp/Utils/mongodb"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var notesCollection *mongo.Collection = mongodb.GetCollection(mongodb.DB, "notes")
var validate = validator.New()

type Note struct {
	Id      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Title   string             `json:"title" bson:"title" validate:"required"`
	Users   []string           `json:"users" bson:"users"`
	Content string             `json:"content" bson:"content" validate:"required"`
	Created Created            `json:"created" bson:"created"`
	Updated Updated            `json:"updated" bson:"updated"`
	Deleted Deleted            `json:"deleted" bson:"deleted,omitempty"`
}

type Created struct {
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	Time   time.Time          `json:"time" bson:"time"`
}

type Updated struct {
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	Time   time.Time          `json:"time" bson:"time"`
}

type Deleted struct {
	Ok     bool               `json:"ok" bson:"ok,omitempty"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	Time   time.Time          `json:"time" bson:"time,omitempty"`
}
