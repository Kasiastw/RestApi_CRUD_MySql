package tests

import (
	"github.com/getground/tech-tasks/backend/cmd/app/handlers.go"
	"github.com/getground/tech-tasks/backend/cmd/app/server"
	"os"
	"testing"
)

var app = api.NewSerwer()

func TestRun(t *testing.T) {
	app.Init(
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	app.Handlers = handlers.NewHandlerFunc(app.DB)
	clearTable()
	ensureTableExists()
}
