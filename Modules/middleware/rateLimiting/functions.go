package ratelimiting

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

func PerClientRateLimiter() fiber.Handler {
	//Remove the clients from map if the client is inactive
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.LastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &Client{
				Limiter:  rate.NewLimiter(2, 4),
				LastSeen: time.Now(),
			}
		}
		//Return an error if the user has exceed the number of requests
		if !clients[ip].Limiter.Allow() {
			mu.Unlock()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, try again later",
			})
		}
		mu.Unlock()
		return c.Next()
	}
}
