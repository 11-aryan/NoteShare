package notes

import (
	"NotesApp/Modules/middleware/jwt"

	"github.com/gofiber/fiber/v2"
)

func NoteRoutes(app *fiber.App) {
	group := app.Group("/api/notes", jwt.ValidateJWTCookie)
	group.Post("/", CreateNote)
	group.Get("/", GetNotes)
	group.Get("/search", Search)
	group.Post("/:id/share", ShareNote)
	group.Get("/:id", GetNoteByID)
	group.Put("/:id", UpdateNote)
	group.Delete("/:id", DeleteNote)
}
