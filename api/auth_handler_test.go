package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adityash1/go-reservation-api/db/fixtures"
	"github.com/adityash1/go-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAuthSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := fixtures.AddUser(tdb.Store, "aditya", "sharma", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuth)

	params := types.AuthParams{
		Email:    "aditya@sharma.com",
		Password: "aditya_sharma",
	}
	b, _ := json.Marshal(params)
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
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == " " {
		t.Fatalf("expected the JWT token to be in auth response")
	}
	// user also contains encrypted password so before comparing reset this
	// because we do not return it in any JSON response
	insertedUser.EncryptedPassword = ""
	fmt.Println("inserted user ->", insertedUser)
	fmt.Println("auth user ->", authResp.User)
	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("expected the user to be inserted user")
	}
}

func TestAuthWithWrongPasswordFailure(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	fixtures.AddUser(tdb.Store, "aditya", "sharma", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuth)

	params := types.AuthParams{
		Email:    "aditya@sharma.com",
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
