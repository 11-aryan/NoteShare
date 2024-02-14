package response

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (resp *model) Data(data interface{}) error {
	response := &Response{
		Ok: true,
	}
	if data != nil {
		response.Data = data
	}
	return resp.Status(fiber.StatusOK).JSON(response)
}

func (resp *model) Error(err interface{}) error {
	response := &Response{
		Ok: false,
	}
	if err != nil {
		response.Error = fmt.Sprintf("%+[1]v", err)
	}
	return resp.Status(fiber.StatusOK).JSON(response)
}

func (resp *model) Message(message interface{}) error {
	response := &Response{
		Ok: true,
	}
	if message != nil {
		response.Message = fmt.Sprintf("%+[1]v", message)
	}
	return resp.Status(fiber.StatusOK).JSON(response)
}

