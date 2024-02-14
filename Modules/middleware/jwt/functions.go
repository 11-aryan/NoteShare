package jwt

import (
	"NotesApp/Utils/response"
	"context"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)


const SecretKey = "secret"

// ValidateJWT is a middleware function to check if the request has a valid JWT token stored in the cookies
func ValidateJWTCookie(c *fiber.Ctx) error {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	// Get the JWT token stored in the cookie
	cookie := c.Cookies("jwt")
	if cookie == "" {
		return resp.Error(errors.New("Unauthorized"))
	}
	// Parse the JWT token
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrForbidden
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return resp.Error(errors.New("Unauthorized"))
	}
	// Verify if the token is valid
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return c.Next()
	}
	return resp.Error(errors.New("Unauthorized"))
}
