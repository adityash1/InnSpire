package main

import (
	"context"
	"fmt"
	"github.com/adityash1/go-reservation-api/api"
	"github.com/adityash1/go-reservation-api/db"
	"github.com/adityash1/go-reservation-api/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	var (
		ctx         = context.Background()
		MongoDBUrl  = os.Getenv(db.MongoDBUrlEnvName)
		MongoDBName = os.Getenv(db.MongoDBNameEnvName)
	)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MongoDBUrl))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(MongoDBName).Drop(ctx); err != nil {
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
	fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
