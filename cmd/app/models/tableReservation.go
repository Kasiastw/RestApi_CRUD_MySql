package models

import (
	"strconv"
	"time"
)

const (
	Upcoming       	Status = 0
	Attended 		Status = 1
	Archived       	Status = 2
)
type (
	Status      			int
	timestamp 				uint64
	Table struct {
		Id 					int64 			`json:"id"`
		Capacity			int				`json:"capacity"`
		BookedSeats 		int				`json:"booked_seats"`
		AvailableSeats 		int 			`json:"available_seats"`
	}
	GuestsReservation struct {
		Id 					int64 			`json:"id"`
		TableId				int32			`json:"table_id"`
		AccompanyingGuests	int64			`json:"accompanying_guests"`
		Status				Status			`json:"status"`
		Name				string			`json:"name"`
		ArrivalTime       	timestamp		`json:"time_arrived"`
	}
	Seats struct {
		SeatsEmpty 			int32 			`json:"seats_empty"`
	}
	GuestList struct {
		Guests 				[]GuestsReservation `json:"guests"`
	}
	GuestDto struct {
		Name 				string 			`json:"name"`
	}
)

func GuestDtoFromEntity(guestEntity GuestsReservation) GuestDto {
	return GuestDto{Name: guestEntity.Name}
}

func (ts timestamp) MarshalJSON() (data []byte, _ error) {
		layout := "2006-01-02 15:04:05"
		x:= time.Unix(int64(ts), 0).Format(layout)
		return strconv.AppendQuote(data, x), nil
}




