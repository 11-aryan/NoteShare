package search

import (
	"NotesApp/Utils/cache"

	"github.com/gofiber/fiber/v2/log"
)

// add adds documents to the Index.
func (idx Index) Add(docs []Document, userID string) {
	indexMap := make(map[string][]string)
	for _, doc := range docs {
		for _, token := range analyze(doc.GetSearchString()) {
			ids := indexMap[token]
			if ids != nil && ids[len(ids)-1] == doc.GetID() {
				// Don't add same ID twice.
				continue
			}
			indexMap[token] = append(ids, doc.GetID())
		}
	}
	idx.hash = cache.RClient
	err := idx.hash.SetHash("searchIndex", userID, indexMap)
	if err != nil {
		log.Error(err)
		return
	}
}

// search queries the Index for the given text.
func (idx Index) Search(text string, userID string) (result []string, err error) {
	idx.hash = cache.RClient
	indexMap, err := idx.hash.GetHash("searchIndex", userID)
	if err != nil {
		log.Errorf("Cannot get hash")
		log.Error(err)
	}
	for _, token := range analyze(text) {
		if ids, ok := indexMap[token]; ok {
			if result == nil {
				result = ids
			} else {
				result = Intersection(result, ids)
			}
		} else {
			// Token doesn't exist.
			return
		}
	}
	return
}
