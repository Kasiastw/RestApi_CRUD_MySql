package api

import (
	"database/sql"
	"fmt"
	"github.com/getground/tech-tasks/backend/cmd/app/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Router   *mux.Router
	DB       *sql.DB
	Handlers *handlers.Post
}

func NewSerwer() *Server {
	return &Server{}
}


func (s *Server) Init(user, password, host, port, name string) {
	var err error
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name)
	s.DB, err = sql.Open("mysql", source)
	if err != nil {
		log.Fatal("cannot conntect to databasse", err)
	}
	pingErr := s.DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Println("DB connected!")

	s.Router = mux.NewRouter()
	s.Handlers = handlers.NewHandlerFunc(s.DB)
	s.Router.HandleFunc("/tables", s.Handlers.CreateTable).Methods("POST")
	s.Router.HandleFunc("/guest_list/{name}", s.Handlers.CreateGuestsListEntry).Methods("POST")
	s.Router.HandleFunc("/guests/{name}", s.Handlers.UpdateGuestsList).Methods("PUT")
	s.Router.HandleFunc("/guest_list", s.Handlers.GetGuestsList).Methods("GET")
	s.Router.HandleFunc("/guests", s.Handlers.GetArrivedGuests).Methods("GET")
	s.Router.HandleFunc("/seats_empty", s.Handlers.GetEmptySeats).Methods("GET")
	s.Router.HandleFunc("/guests/{name}", s.Handlers.GuestLeaves).Methods("DELETE")
}

func Run() {
	app:= NewSerwer()

	err := godotenv.Load("../../.env")
	if err!=nil {
		log.Printf("Couldn't load .env file %v", err)
		return
	}

	app.Init(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	http.ListenAndServe(":3000", app.Router)
}

