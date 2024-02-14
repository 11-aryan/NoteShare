package response

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

func Wrap(c *fiber.Ctx) (a *model) {
	a = &model{c}
	return
}

func ErrorHandler(c *fiber.Ctx, input error) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c.SetUserContext(ctx)
	log.Error(input)
	code := fiber.StatusInternalServerError
	if e, ok := input.(*fiber.Error); ok {
		code = e.Code
	}
	response := &Response{
		Ok:    false,
		Error: input.Error(),
	}
	err = c.Status(code).JSON(response)
	if err != nil {
		log.Error(err)
		return
	}
	return
}
