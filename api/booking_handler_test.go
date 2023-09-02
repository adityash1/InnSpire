package api

import (
	"encoding/json"
	"fmt"
	"github.com/adityash1/go-reservation-api/api/middleware"
	"github.com/adityash1/go-reservation-api/db/fixtures"
	"github.com/adityash1/go-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser      = fixtures.AddUser(tdb.Store, "admin", "admin", true)
		user           = fixtures.AddUser(tdb.Store, "aditya", "sharma", false)
		hotel          = fixtures.AddHotel(tdb.Store, "shere punjab", "india", 3, nil)
		room           = fixtures.AddRoom(tdb.Store, "small", true, 4.3, hotel.ID)
		booking        = fixtures.AddBooking(tdb.Store, adminUser.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(tdb.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("not authorized, response is %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
	if bookings[0].ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookings[0].ID)
	}
	if bookings[0].UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookings[0].UserID)
	}

	// test for non admin acccess to all bookings info
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non 200 status code but got %d", resp.StatusCode)
	}
}

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		nonAuthUser    = fixtures.AddUser(tdb.Store, "not", "auth", false)
		user           = fixtures.AddUser(tdb.Store, "aditya", "sharma", false)
		hotel          = fixtures.AddHotel(tdb.Store, "shere punjab", "india", 3, nil)
		room           = fixtures.AddRoom(tdb.Store, "small", true, 4.3, hotel.ID)
		booking        = fixtures.AddBooking(tdb.Store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 2))
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(tdb.User))
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	route.Get("/", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	fmt.Println(bookingResp)
	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResp.UserID)
	}

	// test for non-authorised user to get a booking info
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected non 200 status code but got %d", resp.StatusCode)
	}
}
