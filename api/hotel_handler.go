package api

import (
	"fmt"
	"github.com/adityash1/go-reservation-api/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	HotelStore db.HotelStore
	roomStore  db.RoomStore
}

type HotelQueryParams struct {
	Rooms bool
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		HotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	fmt.Println("qparams", qparams)
	hotels, err := h.HotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
