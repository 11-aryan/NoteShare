package user

import (
	"NotesApp/Modules/middleware/jwt"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	group := app.Group("/api/users", jwt.ValidateJWTCookie)
	group.Get("/:id", GetUser)
	app.Post("/signup", SignUp)
}

