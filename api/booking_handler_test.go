package api

import (
	"fmt"
	"github.com/adityash1/go-reservation-api/db/fixtures"
	"testing"
	"time"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.Store, "aditya", "sharma", true)
	hotel := fixtures.AddHotel(tdb.Store, "shere punjab", "india", 3, nil)
	room := fixtures.AddRoom(tdb.Store, "small", true, 4.3, hotel.ID)
	booking := fixtures.AddBooking(tdb.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println(booking)
}
