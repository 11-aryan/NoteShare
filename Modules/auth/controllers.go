package auth

import (
	"NotesApp/Utils/response"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func login(c *fiber.Ctx) error {
	fmt.Println("Login Called")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		return resp.Error(err)
	}
	var user User
	// Get the user from DB using email
	err = userCollection.FindOne(ctx, bson.M{"email": data["email"]}).Decode(&user)
	if err != nil {
		return resp.Error(err)
	}
	if user.Id == primitive.NilObjectID {
		c.Status(fiber.StatusNotFound)
		return resp.Error(errors.New("User Not Found"))
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		return resp.Error(err)
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Id.Hex(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		return resp.Error(err)
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return resp.Message("Login Successful")
}

func loggedInUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp := response.Wrap(c)
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return resp.Error(err)
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var user User
	userId, err := primitive.ObjectIDFromHex(claims.Issuer)
	if err != nil {
		return resp.Error(err)
	}
	err = userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return resp.Error(err)
	}
	return resp.Data(user)
}

func refresh(c *fiber.Ctx) error {
	resp := response.Wrap(c)
	refreshToken := c.Cookies("jwt")
	if refreshToken == "" {
		return resp.Error(errors.New("Unauthorized: Refresh token not provided"))
	}
	// Parse and verify the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrForbidden
		}
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		return resp.Error(errors.New("Unauthorized: Invalid refresh token"))
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return resp.Error(errors.New("Internal Server Error"))
	}
	// Check remaining time before expiration
	// Refresh only if 30 seconds or less are left before expiration
	expTime := time.Unix(int64(claims["exp"].(float64)), 0)
	timeUntilExp := time.Until(expTime)
	if timeUntilExp > 30*time.Second {
		return resp.Error(errors.New("Token does not need refreshing"))
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims["exp"] = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return resp.Error(err)
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  expirationTime,
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return resp.Message("Token Successfully Refreshed")
}

func logout(c *fiber.Ctx) error {
	resp := response.Wrap(c)
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return resp.Message("Logout Successful")
}
