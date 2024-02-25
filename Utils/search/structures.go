package search

import "go.mongodb.org/mongo-driver/bson/primitive"

// Index is an inverted index. It maps tokens to document IDs.
type Index map[string][]primitive.ObjectID

var IndexMap = make(Index)

type Document interface {
    GetSearchString() string
    GetID() primitive.ObjectID
}