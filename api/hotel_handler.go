package api

import (
	"github.com/adityash1/go-reservation-api/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	HotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		HotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.HotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
