package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Ok      bool        `json:"ok" bson:"ok"`
	Data    interface{} `json:"data,omitempty" bson:"data,omitempty"`
	Error   string      `json:"error,omitempty" bson:"error,omitempty"`
	Message string      `json:"message,omitempty" bson:"message,omitempty"`
}

type model struct {
	*fiber.Ctx
}
