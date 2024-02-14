package main

import (
	"NotesApp/Modules/middleware/rateLimiting"
	"NotesApp/Modules/notes"
	"NotesApp/Utils/response"
	"NotesApp/Modules/user"
	"NotesApp/Modules/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"NotesApp/Utils/mongodb"
)

func main() {
	app := fiber.New(
		fiber.Config{
			ErrorHandler: response.ErrorHandler,
			StrictRouting: true,
			CaseSensitive: true,
		},
	)
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	app.Use(ratelimiting.PerClientRateLimiter())

	mongodb.ConnectDB()

	auth.AuthRoutes(app)
	user.UserRoutes(app)
	notes.NoteRoutes(app)

	app.Listen(":8000")
}