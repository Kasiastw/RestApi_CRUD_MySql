package tests

import (
	"bytes"
	"encoding/json"
	"github.com/getground/tech-tasks/backend/cmd/app/models"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestEmptyTable(t *testing.T) {
	ensureTableExists()

	req, _ := http.NewRequest("GET", "/guest_list", nil)
	response := executeRequest(req)
	var guestReservations []models.GuestsReservation
	emptyJson, _:= json.MarshalIndent(models.GuestList{Guests: guestReservations}, "", "")
	checkResponseCode(t, http.StatusOK, response.Code)
	if body := response.Body.String(); body != string(emptyJson) {
		t.Errorf("Expected an empty object of array. Got %s, %s", body, string(emptyJson))
	}
}

func TestCreateTable(t *testing.T)  {
	clearTable()
	ensureTableExists()

	tests:= []struct{
		name 		string
		args 		string
		want 		int
	}{
		{
			name: "test if the table was added",
			args: `{"capacity":10}`,
			want: 200,
		},
		{
			name: "test if the table won't be added due to invalid capacity",
			args: `{"capacity":-10}`,
			want: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jsonStr = []byte(tt.args)
			req, _ := http.NewRequest("POST", "/tables", bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req)
			checkResponseCode(t, tt.want, response.Code)
			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)
			if got :=response.Code; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(" Expected response code %v, want %v", tt.want, got)
			}
		})
	}
}

func TestCreateGuestsListEntry(t *testing.T)  {

	tests:= []struct{
		name 		string
		guestName 	string
		args 		string
		want 		int
	}{
		{
			name: "test if the guests were added to table",
			args: `{"accompanying_guests":10, "table_id":1}`,
			guestName: "Tom",
			want: 200,
		},
		{
			name: "test if the guests won't be added to a table due to the unavailable seats",
			guestName: "oli",
			args: `{"accompanying_guests":100, "table_id":1}`,
			want: 500,
		},
		{
			name: "test if the guests won't be added to table due to invalid number of accompanying guests",
			guestName: "oli",
			args: `{"accompanying_guests":0, "table_id":1}`,
			want: 400,
		},
		{
			name: "test when incorrect table_is is given",
			guestName: "oli",
			args: `{"accompanying_guests":2, "table_id":111}`,
			want: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jsonStr = []byte(tt.args)
			req, _ := http.NewRequest("POST", "/guest_list/"+tt.guestName, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req)
			checkResponseCode(t, tt.want, response.Code)
			var m map[string]string
			json.Unmarshal(response.Body.Bytes(), &m)
			if got :=response.Code; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(" Expected response code %v, want %v", tt.want, got)
			}
			if response.Code==200 && m["name"] != tt.guestName {
					t.Errorf("Expected name to be %s. Got %s", tt.guestName, m["name"])
			}
		})
	}
}

func TestUpdateGuestsList(t *testing.T)  {
	tests:= []struct{
		name 		string
		guestName 	string
		args 		string
		want 		int
	}{
		{
			name: "test if the guests were added to table",
			args: `{"accompanying_guests":10, "table_id":1}`,
			guestName: "Tom",
			want: 200,
		},
		{
			name: "test if the guests were added to table",
			args: `{"accompanying_guests":10, "table_id":1}`,
			guestName: "Tom",
			want: 200,
		},
		{
			name: "test if the guests were added to table",
			args: `{"accompanying_guests":2, "table_id":1}`,
			guestName: "Tom",
			want: 200,
		},
		{
			name: "test when invalid guest number is given",
			args: `{"accompanying_guests":-20, "table_id":1}`,
			guestName: "Tom",
			want: 400,
		},
		{
			name: "test when invalid guest name is given",
			args: `{"accompanying_guests":2, "table_id":1}`,
			guestName: "John",
			want: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var jsonStr = []byte(tt.args)
			req, _ := http.NewRequest("PUT", "/guests/"+tt.guestName, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req)
			checkResponseCode(t, tt.want, response.Code)
			var m map[string]string
			json.Unmarshal(response.Body.Bytes(), &m)
			if got :=response.Code; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(" Expected response code %v, want %v", tt.want, got)
			}
			if response.Code==200 && m["name"] != tt.guestName {
				t.Errorf("Expected name to be %s. Got %s", tt.guestName, m["name"])
			}
		})
	}
}

func TestGetGuestsList(t *testing.T)  {
	tests:= []struct{
		name 		string
		want 		int
	}{
		{
			name: "test if we get the list",
			want: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/guests", nil)
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req)
			checkResponseCode(t, tt.want, response.Code)
			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)
			if got :=response.Code; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(" Expected response code %v, want %v", tt.want, got)
			}
		})
	}
}

func TestGetEmptySeats(t *testing.T)  {
	tests:= []struct{
		name 		string
		want 		int
		wantedSeats	int
	}{
		{
			name: "test if the number of available seats is correct",
			want: 200,
			wantedSeats: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/seats_empty", nil)
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req)
			checkResponseCode(t, tt.want, response.Code)
			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)
			if got :=response.Code; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(" Expected response code %v, want %v", tt.want, got)
			}
			if response.Code==200 && m["seats_empty"] != float64(tt.wantedSeats) {
				t.Errorf("Expected name to be %v. Got %v", tt.wantedSeats, m["seats_empty"])
			}
		})
	}
}

func TestGuestLeaves(t *testing.T)  {
	tests:= []struct{
		name 		string
		want 		int
		guestName	string
	}{
		{
			name: "test if the reservation was archived",
			want: 204,
			guestName: "Tom",
		},
		{
			name: "test if the reservation couldn't be archived because of invalid name",
			want: 500,
			guestName: "Anthony",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/guests/"+tt.guestName, nil)
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req)
			checkResponseCode(t, tt.want, response.Code)
			var m map[string]interface{}
			json.Unmarshal(response.Body.Bytes(), &m)
			if got :=response.Code; !reflect.DeepEqual(got, tt.want) {
				t.Errorf(" Expected response code %v, want %v", tt.want, got)
			}
		})
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func ensureTableExists() {
	if _, err := app.DB.Exec(createTableTables); err != nil {
		log.Fatal(err)
	}

	if _, err := app.DB.Exec(createTableGuestList); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	if _, err := app.DB.Exec("DROP TABLE IF EXISTS getground.guestsList"); err != nil {
		log.Fatal(err)
	}
	if _, err := app.DB.Exec("DROP TABLE IF EXISTS getground.tables"); err != nil {
		log.Fatal(err)
	}
}

const createTableTables = `CREATE TABLE IF NOT EXISTS tables
(
	id INT NOT NULL auto_increment, PRIMARY KEY (id),
                          capacity int,
                          booked_seats int,
                          available_seats int
)`

const createTableGuestList = `CREATE TABLE IF NOT EXISTS guestsList
(
	id INT NOT NULL auto_increment,
                         PRIMARY KEY (id),
                         table_id INT,
                         name VARCHAR(100) NOT NULL,
                         accompanying_guests INT,
                         status int,
                         arrival_time bigint,
                         FOREIGN KEY (table_id) REFERENCES tables(id)
)`



