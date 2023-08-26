package api

import (
	"github.com/adityash1/go-reservation-api/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	HotelStore db.HotelStore
	roomStore  db.RoomStore
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		HotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelID": objId}
	rooms, err := h.roomStore.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	hotels, err := h.HotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
