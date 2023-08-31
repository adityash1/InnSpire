package main

import (
	"context"
	"flag"
	"github.com/adityash1/go-reservation-api/api/middleware"
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

	// handlers initialisation
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			User:    userStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookHandler(store)
	)

	app := fiber.New(config)
	auth := app.Group("api")
	apiv1 := app.Group("api/v1", middleware.JWTAuthentication(userStore))
	admin := apiv1.Group("/admin", middleware.AdminAuth)

	//auth
	auth.Post("/auth", authHandler.HandleAuth)

	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// rooms handlers
	apiv1.Post("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	//bookings handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBookings)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	err = app.Listen(*port)
	if err != nil {
		return
	}
}
