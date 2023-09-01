package api

import (
	"bytes"
	"encoding/json"
	"github.com/adityash1/go-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"net/http/httptest"
	"testing"
)

func TestPostUser(t *testing.T) {
	tbd := setup(t)
	defer tbd.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tbd.User)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@test.com",
		FirstName: "John",
		LastName:  "Jacob",
		Password:  "dskdfsjbnfsdiodfnd",
	}
	b, err := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return
	}
	if len(user.ID) == 0 {
		t.Fatalf("expecting a user to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Fatalf("expecting encryptedpassword not to be included in json")
	}
	if user.FirstName != params.FirstName {
		t.Fatalf("expected username %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Fatalf("expected lastname %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Fatalf("expected email %s but got %s", params.Email, user.Email)
	}
}
