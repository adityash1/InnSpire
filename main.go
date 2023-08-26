package main

import (
	"context"
	"flag"
	"log"

	"github.com/adityash1/go-reservation-api/api"
	"github.com/adityash1/go-reservation-api/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	port := flag.String("port", ":8080", "The listen address of the api server")
	flag.Parse()

	// connection to db
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	// handler initialisation
	var (
		userHandler  = api.NewUserHandler(db.NewMongoUserStore(client, db.DB_NAME))
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
	)

	app := fiber.New(config)
	apiv1 := app.Group("api/v1")

	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)

	err = app.Listen(*port)
	if err != nil {
		return
	}
}
