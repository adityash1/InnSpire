package api

import (
	"github.com/adityash1/go-reservation-api/db"
	"github.com/adityash1/go-reservation-api/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelID": objId}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams types.HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), objId)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}
