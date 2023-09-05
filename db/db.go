package db

import "context"

const (
	DB_URI       = "mongodb://localhost:27017"
	TEST_DB_URI  = "mongodb://localhost:27017"
	DB_NAME      = "hotel-reservation"
	TEST_DB_NAME = "hotel-reservation-test"
	userCol      = "users"
	hotelCol     = "hotels"
)

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Dropper interface {
	Drop(ctx context.Context) error
}

type Map map[string]any
