package api

import (
	"errors"
	"fmt"
	"github.com/adityash1/go-reservation-api/db"
	"github.com/adityash1/go-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var params types.AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}
	fmt.Println("authenticated ->", user)
	return nil
}
