package api

import (
	"github.com/adityash1/go-reservation-api/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return ErrUnauthorized()
	}
	if !user.IsAdmin {
		return ErrUnauthorized()
	}
	return c.Next()
}
