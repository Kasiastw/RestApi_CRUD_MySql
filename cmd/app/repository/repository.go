package repository

import (
	"context"
	"github.com/getground/tech-tasks/backend/cmd/app/models"
)

type GuestRepo interface {
	CreateTableId(ctx context.Context, table models.Table) (int64, error)
	CreateGuestReservationID (ctx context.Context, guest *models.GuestsReservation) error
	CheckAvailableSeats (ctx context.Context, guest *models.GuestsReservation) error
	GetGuestsList() (*models.GuestList, error)
	GetArrivedGuests() (*models.GuestList, error)
	GetEmptySeats() (*models.Seats, error)
	GuestLeaves(ctx context.Context, name string) error
}
