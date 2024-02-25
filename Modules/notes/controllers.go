package notes

import (
	"NotesApp/Utils/response"
	"NotesApp/Utils/search"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Creates a new Note for a user
func CreateNote(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	var note Note
	//Decoding the request body into the note sturct
	err := c.BodyParser(&note)
	if err != nil {
		resp.Error(err)
	}
	//using the validator library to validate required fields
	err = validate.Struct(&note)
	if err != nil {
		resp.Error(err)
	}
	currentUserId := note.Users[0]
	currentUserObjID, err := primitive.ObjectIDFromHex(currentUserId)
	if err != nil {
		resp.Error(err)
	}
	currentTime := time.Now()
	note.Created.Time = currentTime
	note.Created.UserId = currentUserObjID
	note.Updated.Time = currentTime
	note.Updated.UserId = currentUserObjID
	// Initializing a new Note object with values retrieved from request body
	newNote := Note{
		Title:   note.Title,
		Users:   note.Users,
		Content: note.Content,
		Created: note.Created,
		Updated: note.Updated,
	}
	// Inserting the new note into the notes collection
	result, err := notesCollection.InsertOne(ctx, newNote)
	if err != nil {
		resp.Error(err)
	}
	return resp.Data(result)
}

// Search a note based on title
func Search(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	searchKey := c.Query("search_key")
	filter := bson.M{
		"title": bson.M{
			"$regex":   searchKey,
			"$options": "i",
		},
	}
	var notes []Note
	cursor, err := notesCollection.Find(ctx, filter)
	if err != nil {
		log.Error(err)
		return resp.Error(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var note Note
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

func CreateSearchIndex(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	userID := c.Params("userID")
	filter := bson.M{
		"users": bson.M{
			"$in": []string{userID},
		},
	}
	var notes []Note
	cursor, err := notesCollection.Find(ctx, filter)
	if err != nil {
		log.Error(err)
		return resp.Error(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var note Note
		if err = cursor.Decode(&note); err != nil {
			return resp.Error(err)
		}
		notes = append(notes, note)
	}
	if err := cursor.Err(); err != nil {
		return resp.Error(err)
	}
	documents := []search.Document{}
	for _, note := range notes {
		documents = append(documents, note)
	}
	idx := search.Index{}
	idx.Add(documents, userID)
	return resp.Message("Successfully created index")
}

// Search notes from the cache using note title
func SearchCache(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	searchKey := c.Query("search_key")
	userID := c.Query("user_id")
	idx := search.Index{}
	matchedIDs, err := idx.Search(searchKey, userID)
	if err != nil {
		log.Error(err)
		return err
	}
	return resp.Data(matchedIDs)
}

func GetNoteByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	noteID := c.Params("id")
	noteObjID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return resp.Error(err)
	}
	filter := bson.M{
		"_id": noteObjID,
	}
	var note Note
	err = notesCollection.FindOne(ctx, filter).Decode(&note)
	if err != nil {
		return resp.Error(err)
	}
	return resp.Data(note)
}

// Update a note with given ID
func UpdateNote(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	noteID := c.Params("id")
	noteObjID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return resp.Error(err)
	}
	var note Note
	//Decoding the request body into the note sturct
	err = c.BodyParser(&note)
	if err != nil {
		return resp.Error(err)
	}
	//using the validator library to validate required fields
	err = validate.Struct(&note)
	if err != nil {
		return resp.Error(err)
	}
	currentTime := time.Now()
	note.Updated.Time = currentTime
	// Initializing a new Note object with values retrieved from request body
	newNote := Note{
		Id:      noteObjID,
		Title:   note.Title,
		Users:   note.Users,
		Content: note.Content,
		Updated: note.Updated,
	}
	filter := bson.M{
		"_id": noteObjID,
	}
	update := bson.M{
		"$set": bson.M{
			"title":   newNote.Title,
			"content": newNote.Content,
			"users":   newNote.Users,
			"updated": newNote.Updated,
		},
	}
	_, err = notesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resp.Error(err)
	}
	return resp.Message("Successfully Updated Note")
}

func DeleteNote(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	noteID := c.Params("id")
	noteObjID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return resp.Error(err)
	}
	var data map[string]interface{}
	//Decoding the request body into the note sturct
	err = c.BodyParser(&data)
	if err != nil {
		return resp.Error(err)
	}
	// Get userID of the user who sent the delete request
	userID := data["user_id"].(string)
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return resp.Error(err)
	}
	currentTime := time.Now()
	updated := Updated{
		UserId: userObjID,
		Time:   currentTime,
	}
	deleted := Deleted{
		Ok:     true,
		UserId: userObjID,
		Time:   currentTime,
	}
	filter := bson.M{
		"_id": noteObjID,
	}
	update := bson.M{
		"$set": bson.M{
			"deleted": deleted,
			"updated": updated,
		},
	}
	_, err = notesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resp.Error(err)
	}
	return resp.Message("Successfully Deleted Note")
}

func ShareNote(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	noteID := c.Params("id")
	noteObjID, err := primitive.ObjectIDFromHex(noteID)
	if err != nil {
		return resp.Error(err)
	}
	var data map[string]interface{}
	err = c.BodyParser(&data)
	if err != nil {
		return resp.Error(err)
	}
	usersInterface := data["users"].([]interface{})
	var users []string
	// Iterate over each element and convert it to a string
	for _, user := range usersInterface {
		if userString, ok := user.(string); ok {
			users = append(users, userString)
		}
	}
	userID := data["user_id"].(string)
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return resp.Error(err)
	}
	currentTime := time.Now()
	updated := Updated{
		UserId: userObjID,
		Time:   currentTime,
	}
	// Update the users field of the note with new users
	filter := bson.M{
		"_id": noteObjID,
	}
	update := bson.M{
		"$set": bson.M{
			"users":   users,
			"updated": updated,
		},
	}
	_, err = notesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return resp.Error(err)
	}
	return resp.Message("Successfully Shared Note")
}
