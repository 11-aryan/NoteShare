package search

import "go.mongodb.org/mongo-driver/bson/primitive"

// add adds documents to the Index.
func (idx Index) Add(docs []Document) {
	for _, doc := range docs {
		for _, token := range analyze(doc.GetSearchString()) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.GetID() {
				// Don't add same ID twice.
				continue
			}
			idx[token] = append(ids, doc.GetID())
		}
	}
}

// search queries the Index for the given text.
func (idx Index) Search(text string) []primitive.ObjectID {
	var result []primitive.ObjectID
	for _, token := range analyze(text) {
		if ids, ok := idx[token]; ok {
			if result == nil {
				result = ids
			} else {
				result = Intersection(result, ids)
			}
		} else {
			// Token doesn't exist.
			return nil
		}
	}
	return result
}