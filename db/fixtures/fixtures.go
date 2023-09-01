package fixtures

import (
	"context"
	"fmt"
	"github.com/adityash1/go-reservation-api/db"
	"github.com/adityash1/go-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"
)

func AddHotel(store *db.Store, name string, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	roomIDs := rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rating:   rating,
		Rooms:    roomIDs,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store, fn, ln string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		SeaSide: ss,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBooking(store *db.Store, uid, rid primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}
