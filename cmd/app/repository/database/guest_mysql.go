package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/getground/tech-tasks/backend/cmd/app/models"
	"github.com/getground/tech-tasks/backend/cmd/app/repository"
	"log"
	"time"
)

type mysqlGuestRepo struct {
	Conn *sql.DB
}

func NewSQLGuestRepo(Conn *sql.DB) repository.GuestRepo {
	return &mysqlGuestRepo{
		Conn: Conn,
	}
}

func (m *mysqlGuestRepo) CreateTableId(ctx context.Context, table models.Table) (int64, error){

	stmt, err := m.Conn.PrepareContext(
		ctx,
		"INSERT INTO tables(capacity, booked_seats, available_seats) VALUES(?, ?, ?);")
	if err != nil {
		return -1, err
	}
	res, err := stmt.ExecContext(ctx, table.Capacity, 0, table.Capacity)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	tableId, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return tableId, nil
}

func (m *mysqlGuestRepo) CreateGuestReservationID (ctx context.Context, guest *models.GuestsReservation) error {
 	var err error
	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ok, err:=m.checkIfTableAvailable(ctx, tx, guest.AccompanyingGuests, guest.TableId)
	if err !=nil {
		return err
	}
	if !ok {
		log.Printf("not enough seats, tableId=%v", guest.TableId)
		return errors.New("not enough seats")
	}

	_, err = tx.ExecContext(
		ctx,
		"UPDATE tables SET booked_seats= booked_seats + ?, available_seats = available_seats-? where id = ?",
		guest.AccompanyingGuests, guest.AccompanyingGuests, guest.TableId)
	if err != nil {
		return err
	}
	res, err := tx.ExecContext(
		ctx,
		"INSERT INTO guestsList(table_id, name, accompanying_guests, status) VALUES (?, ?, ?, ?);",
		guest.TableId, guest.Name, guest.AccompanyingGuests, models.Status(models.Upcoming))
	if err != nil {
		return err
	}
	tableId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	log.Printf("New reservation id=%v was added", tableId)
	return nil
}

func (m *mysqlGuestRepo) CheckAvailableSeats (ctx context.Context, guest *models.GuestsReservation) error {
	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	reservationId, tableId, err := m.getReservationId(tx, guest.Name)
	if err!=nil {
		return err
	}

	var diffGuestsNumber int64
	diffGuestsNumber, err= m.getAvailableSeatsAmount(tx, ctx, guest, reservationId)
	if err!=nil {
		return err
	}

	switch {
		case diffGuestsNumber==0:
			_, err = tx.ExecContext(
				ctx,
				"UPDATE guestsList SET status = ? where id=?", models.Attended, reservationId)
			if err != nil {
				return err
			}
		case diffGuestsNumber>0:
			ok, err:=m.checkIfTableAvailable(ctx, tx, diffGuestsNumber, tableId)
				if err !=nil {
					return err
				}
			if !ok {
				log.Println("not enough seats, tableId=%v", tableId)
				return err
			}
			err = m.updateSeatsAmount(ctx, tx, diffGuestsNumber, guest, reservationId, tableId)
			if err != nil {
				return err
			}
		case diffGuestsNumber<0:
			err := m.updateSeatsAmount(ctx, tx, diffGuestsNumber, guest, reservationId, tableId)
			if err != nil {
				return err
			}
		}
	if err = tx.Commit(); err != nil {
		return err
	}
	log.Printf("the guests: %s (reservationId=%v) arrived", guest.Name , reservationId)
	return nil
}
func (m *mysqlGuestRepo) getReservationId(tx *sql.Tx, name string) (int32, int32, error) {
	var reservationId int32
	var tableId int32
	var nTable sql.NullInt32
	var nReservation sql.NullInt32
	err := tx.QueryRow(
		"SELECT g.id, g.table_id FROM guestsList g where g.name=? ", name).Scan(&nReservation, &nTable)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, -1, err
		}
		return -1, -1, err
	}
	if nTable.Valid && nReservation.Valid {
		tableId = nTable.Int32
		reservationId = nReservation.Int32
	}
	return reservationId, tableId, nil
}

func (m *mysqlGuestRepo) getAvailableSeatsAmount(tx *sql.Tx, ctx context.Context,  guests *models.GuestsReservation,
	reservationId int32) (int64, error)  {
	var updateGuestNumber int64
	var n sql.NullInt64
	err := tx.QueryRowContext(ctx, "SELECT ?- accompanying_guests from guestsList where id = ?",
		guests.AccompanyingGuests, reservationId).Scan(&n)
	if err != nil {
		log.Printf("no such reservation id=%v", guests.Id)
		return -1, err
	}
	if n.Valid {
		updateGuestNumber = n.Int64
	}
	return updateGuestNumber, err
}

func (m *mysqlGuestRepo) checkIfTableAvailable(ctx context.Context, tx *sql.Tx, val int64, tableId int32) (bool, error) {
	var enough bool
	if err := tx.QueryRowContext(ctx, "SELECT (available_seats >= ?) from tables where id = ?",
		val, tableId).Scan(&enough); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no such table with id=%v", tableId)
			return false, fmt.Errorf("no such table_id=%v", tableId)
		}
		return false, fmt.Errorf("checkIfTableAvailable %d: %v", tableId)
	}
	return enough, nil
}


func (m *mysqlGuestRepo) updateSeatsAmount(ctx context.Context, tx *sql.Tx,
											diffGuestNumber int64, guestReservation *models.GuestsReservation,
											reservationId int32, tableId int32) error {
	tArrival := time.Now().UTC().Unix()
	_, err := tx.ExecContext(
		ctx,
		"UPDATE tables SET booked_seats= booked_seats + ?, available_seats = available_seats - ? where id = ?",
		diffGuestNumber, diffGuestNumber, tableId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE guestsList SET accompanying_guests = ?, status = ?, arrival_time=? where id=?",
		guestReservation.AccompanyingGuests, models.Attended, tArrival, reservationId)
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlGuestRepo) GetGuestsList() (*models.GuestList, error) {

	rows, err := m.Conn.Query("SELECT g.table_id, g.name, g.accompanying_guests FROM guestsList g")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guestReservations []models.GuestsReservation
	for rows.Next() {
		var r models.GuestsReservation
		if err := rows.Scan(&r.TableId, &r.Name, &r.AccompanyingGuests); err != nil {
			log.Printf("DB: Error during sql statement to get arrived guest , error=%v", err)
			return nil, err
		}
	guestReservations = append(guestReservations, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &models.GuestList{Guests: guestReservations}, nil
}


func (m *mysqlGuestRepo) GetArrivedGuests() (*models.GuestList, error) {
	rows, err := m.Conn.Query(
		"SELECT g.name, g.accompanying_guests, g.arrival_time FROM guestsList g where g.status=1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guestReservations []models.GuestsReservation
	for rows.Next() {
		var r models.GuestsReservation
		if err := rows.Scan(&r.Name, &r.AccompanyingGuests, &r.ArrivalTime); err != nil {
			log.Printf("DB: Error during sql statement to get arrived guest , error=%v", err)
			return nil, err
		}
		guestReservations = append(guestReservations, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &models.GuestList{Guests: guestReservations}, nil
}

func (m *mysqlGuestRepo) GetEmptySeats() (*models.Seats, error) {
	var emptySeats int32
	var n sql.NullInt32
	err := m.Conn.QueryRow("SELECT SUM(available_seats) FROM tables").Scan(&n)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	if n.Valid {
		emptySeats = n.Int32
	}
	return &models.Seats{SeatsEmpty: emptySeats}, nil
}

func (m *mysqlGuestRepo) GuestLeaves(ctx context.Context, name string) error {
	tx, err := m.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var id int32
	var tableId int32
	var guestAmount int32
	var nId sql.NullInt32
	var nTableId sql.NullInt32
	var nGuestAmount sql.NullInt32
	err = tx.QueryRowContext(ctx,
		"SELECT g.id, g.table_id, g.accompanying_guests from guestsList g where g.name = ?",
		name).Scan(&nId, &nTableId, &nGuestAmount)
	if err != nil {
		log.Printf("no such reservation for %s", name)
		return err
	}
	if nId.Valid && nTableId.Valid && nGuestAmount.Valid  {
		id = nId.Int32
		tableId = nTableId.Int32
		guestAmount = nGuestAmount.Int32
	}
	_, err = tx.ExecContext(
		ctx,
		"UPDATE tables SET booked_seats= booked_seats + ?, available_seats = available_seats - ? where id = ?",
		guestAmount, guestAmount, tableId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,  "UPDATE guestsList SET accompanying_guests = ?, status = ? where id=?",
		guestAmount, models.Archived, id)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	log.Printf("the guests with id=%v left", id)
	return nil
}
