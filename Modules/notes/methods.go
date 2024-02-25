package notes

import "go.mongodb.org/mongo-driver/bson/primitive"

func (n Note) GetSearchString() string {
	return n.Title
}

func (n Note) GetID() primitive.ObjectID {
	return n.Id
}

