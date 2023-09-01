package main

import (
	"context"
	"fmt"
	"github.com/adityash1/go-reservation-api/api"
	"github.com/adityash1/go-reservation-api/db"
	"github.com/adityash1/go-reservation-api/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DB_NAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
		Hotel:   hotelStore,
	}
	user := fixtures.AddUser(store, "aditya", "sharma", false)
	fmt.Println("aditya ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(store, "delight", "india", 4, nil)
	room := fixtures.AddRoom(store, "large", true, 88.44, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ->", booking.ID)
}
