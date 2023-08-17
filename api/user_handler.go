package api

import (
	"github.com/adityash1/go-reservation-api/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUser(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "Aditya",
		LastName:  "Sharma",
	}
	return c.JSON(user)
}

func HandleGetUsers(c *fiber.Ctx) error {
	return c.JSON("Aditya")
}
