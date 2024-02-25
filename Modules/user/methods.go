package user 

import "go.mongodb.org/mongo-driver/bson/primitive"

func (u User) GetSearchString() string {
	return u.Name
}

func (u User) GetID() primitive.ObjectID {
	return u.Id
}
