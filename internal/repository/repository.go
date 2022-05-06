package repository

import (
	"time"

	"github.com/maxkobzin/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityByDatesForAllRooms(start, end time.Time) ([]models.Room, error)
}
