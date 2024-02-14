package auth

import (
	"NotesApp/Modules/middleware/jwt"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	group := app.Group("/api")
	group.Post("/login", login)
	group.Get("/loggedInUser", loggedInUser, jwt.ValidateJWTCookie)
	group.Post("/logout", logout)
	group.Post("/refresh", refresh)
}
