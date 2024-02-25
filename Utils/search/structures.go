package search

import (
	"NotesApp/Utils/cache"
)

// Index is an inverted index. It maps tokens to document IDs.
// type Index map[string][]primitive.ObjectID

type Index struct {
	hash cache.RedisClient
}

// var IndexMap = make(Index)

type Document interface {
	GetSearchString() string
	GetID() string
}
