package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/getground/tech-tasks/backend/cmd/app/models"
	"github.com/getground/tech-tasks/backend/cmd/app/repository"
	"github.com/getground/tech-tasks/backend/cmd/app/repository/database"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


type Post struct {
	repo repository.GuestRepo
}

func NewHandlerFunc(db *sql.DB) *Post {
	return &Post{
		repo: database.NewSQLGuestRepo(db),
	}
}

func (s *Post) CreateTable(w http.ResponseWriter, r *http.Request) {
	var table models.Table
	err := json.NewDecoder(r.Body).Decode(&table)
	if err != nil {
		models.RespondwithJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	if table.Capacity <=0 {
		models.RespondWithError(w, http.StatusBadRequest, "Invalid table capacity")
		log.Printf("Invalid table capacity, %v", table.Capacity)
		return
	}

	tableId, err := s.repo.CreateTableId(r.Context(), table)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}
	log.Printf("New table with id=%v was added", tableId)

	t:=&models.Table{
		Id:       tableId,
		Capacity: table.Capacity,
	}
	models.RespondwithJSON(w, http.StatusOK, t)
}

func (s *Post) CreateGuestsListEntry(w http.ResponseWriter, r *http.Request) {
	params:= mux.Vars(r)
	name := params["name"]
	if name == "" {
		models.RespondWithError(w, http.StatusBadRequest, "Invalid guest name")
		log.Fatalln("Invalid guest name")
		return
	}

	guestsReservation:= models.GuestsReservation{Name: name}
	err := json.NewDecoder(r.Body).Decode(&guestsReservation)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		log.Println("There was an error decoding the request body into the struct")
		return
	}
	defer r.Body.Close()

	if guestsReservation.AccompanyingGuests <=0 {
		models.RespondWithError(w, http.StatusBadRequest, "Invalid guest number")
		log.Printf("Invalid guest number %v", guestsReservation.AccompanyingGuests)
		return
	}

	err = s.repo.CreateGuestReservationID(r.Context(), &guestsReservation)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	mappedResult:= models.GuestDtoFromEntity(models.GuestsReservation(guestsReservation))
	models.RespondwithJSON(w, http.StatusOK, mappedResult)
}

func (s *Post) UpdateGuestsList(w http.ResponseWriter, r *http.Request) {
	params:= mux.Vars(r)
	name := params["name"]
	if name == "" {
		models.RespondWithError(w, http.StatusBadRequest, "Invalid guest name")
		log.Println("Invalid guest name")
		return
	}

	guestsReservation:= models.GuestsReservation{Name: name}
	err := json.NewDecoder(r.Body).Decode(&guestsReservation)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		log.Println("There was an error decoding the request body into the struct")
		return
	}
	defer r.Body.Close()

	if guestsReservation.AccompanyingGuests <=0 {
		models.RespondWithError(w, http.StatusBadRequest, "Invalid guest number")
		log.Printf("Accompanying guests number %v should be greater than zero", guestsReservation.AccompanyingGuests)
		return
	}

	err = s.repo.CheckAvailableSeats(r.Context(), &guestsReservation)
	if err != nil {
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	mappedResult:= models.GuestDtoFromEntity(models.GuestsReservation(guestsReservation))
	models.RespondwithJSON(w, http.StatusOK, mappedResult)
}

func (s *Post) GetGuestsList(w http.ResponseWriter, r *http.Request) {
	guests, err:= s.repo.GetGuestsList()
	if err!=nil {
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	models.RespondwithJSON(w, http.StatusOK, guests)
}

func (s *Post) GetArrivedGuests(w http.ResponseWriter, r *http.Request) {
	guests, err:= s.repo.GetArrivedGuests()
	if err!=nil {
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	models.RespondwithJSON(w, http.StatusOK, guests)
}

func (s *Post) GetEmptySeats(w http.ResponseWriter, r *http.Request) {
	emptySeats, err:= s.repo.GetEmptySeats()
	if err!=nil {
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	models.RespondwithJSON(w, http.StatusOK, emptySeats)
}

func (s *Post) GuestLeaves(w http.ResponseWriter, r *http.Request) {
	params:= mux.Vars(r)
	name := params["name"]
	if name == "" {
		models.RespondWithError(w, http.StatusBadRequest, "Invalid guest name")
		log.Println("Invalid guest name")
		return
	}

	err:= s.repo.GuestLeaves(r.Context(), name)
	if err!=nil {
		models.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	models.RespondwithJSON(w, http.StatusNoContent, name)
}

