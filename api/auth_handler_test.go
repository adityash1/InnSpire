package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/adityash1/go-reservation-api/db"
	"github.com/adityash1/go-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     "some@test.com",
		FirstName: "Aditya",
		LastName:  "Sharma",
		Password:  "password12345",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := insertTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuth)

	params := types.AuthParams{
		Email:    "some@test.com",
		Password: "password12345",
	}
	b, err := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status 200 but got %d", resp.StatusCode)
	}
	var authResp types.AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		t.Fatal(err)
	}
	if authResp.Token == " " {
		t.Fatalf("expected the JWT token to be in auth response")
	}
	// user also contains encrypted password so before comparing reset this
	// because we do not return it in any JSON response
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be inserted user")
	}
}

func TestAuthWithWrongPasswordFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuth)

	params := types.AuthParams{
		Email:    "some@test.com",
		Password: "password",
	}
	b, err := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status 400 but got %d", resp.StatusCode)
	}
	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error")
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response msg to be <invalid credential> but got %s", genResp.Msg)
	}
}
