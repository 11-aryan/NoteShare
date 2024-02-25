package main

import (
	"NotesApp/Modules/auth"
	"NotesApp/Modules/middleware/rateLimiting"
	"NotesApp/Modules/notes"
	"NotesApp/Modules/user"
	"NotesApp/Utils/cache"
	"NotesApp/Utils/mongodb"
	"NotesApp/Utils/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(
		fiber.Config{
			ErrorHandler:  response.ErrorHandler,
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

	cache.InitRedis()

	app.Listen(":8000")
}
