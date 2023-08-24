package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/adityash1/go-reservation-api/db"
	"github.com/adityash1/go-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http/httptest"
	"testing"
)

type testdb struct {
	UserStore db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) error {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
	return nil
}

func setup(_ *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.TEST_DB_URI))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.TEST_DB_NAME),
	}
}

func TestPostUser(t *testing.T) {
	tbd := setup(t)
	defer tbd.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tbd.UserStore)
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
		t.Error(err)
	}
	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return
	}
	if len(user.ID) == 0 {
		t.Errorf("expecting a user to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting encryptedpassword not to be included in json")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected username %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}