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

// commmented code is to check things if they're working before
// proceeding with abstractions

// const dburi = "mongodb://localhost:27017"
// const dbname = "hotel-reservation"
// const userCol = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	port := flag.String("port", ":8080", "The listen address of the api server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}

	// ctx := context.Background()
	// col := client.Database(dbname).Collection(userCol)

	// user := types.User{
	// 	FirstName: "Aditya",
	// 	LastName:  "Sharma",
	// }

	// _, err = col.InsertOne(ctx, user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var aditya types.User
	// if err := col.FindOne(ctx, bson.M{}).Decode(&aditya); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(aditya)

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("api/v1")

	// apiv1.Get("/user", api.HandleGetUsers)
	// apiv1.Get("/user/:id", api.HandleGetUser)

	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	app.Listen(*port)
}
