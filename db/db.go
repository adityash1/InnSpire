package db

import (
	"context"
)

const (
	MongoDBNameEnvName = "MONGO_DB_NAME"
	MongoDBUrlEnvName  = "MONGO_DB_URL"
	userCol            = "users"
	hotelCol           = "hotels"
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
