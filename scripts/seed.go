package main

import (
	"context"
	"fmt"
	"github.com/adityash1/go-reservation-api/db"
	"github.com/adityash1/go-reservation-api/db/fixtures"
	"github.com/adityash1/go-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func seedRoom(size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		SeaSide: seaside,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	resp, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking:", resp.ID)
}

func main() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DB_NAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
		Hotel:   hotelStore,
	}
	user := fixtures.AddUser(store, "aditya", "sharma", false)
	fmt.Println(user)
	hotel := fixtures.AddHotel(store, "delight", "india", 4, nil)
	fmt.Println(hotel)
	room := fixtures.AddRoom(store, "large", true, 88.44, hotel.ID)
	fmt.Println(room)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println(booking)
	return
	//aditya := seedUser(true, "Aditya", "Sharma", "aditya@gmail.com", "adi@123")
	//seedUser(false, "Yash", "Sharma", "yash123@gmail.com", "yash@123")
	//seedUser(false, "Tuntun", "Tiwari", "tuntun@gmail.com", "tuntun@123")
	//
	//seedHotel("Raffles Istanbul", "Turkey", 4)
	//seedHotel("The Driskil", "US", 3)
	//hotel := seedHotel("Taj", "India", 5)
	//
	//seedRoom("medium", true, 149.99, hotel.ID)
	//seedRoom("large", true, 299.99, hotel.ID)
	//room := seedRoom("small", false, 99.99, hotel.ID)
	//
	//seedBooking(aditya.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DB_URI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DB_NAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}
