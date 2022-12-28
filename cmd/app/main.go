package main

import (
	"fmt"
	"github.com/getground/tech-tasks/backend/cmd/app/models"
	"github.com/getground/tech-tasks/backend/cmd/app/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"log"
)

func main() {
	log.Println(fmt.Sprintf("xxx %+v", models.GuestList{Guests: nil}))
	api.Run()
}

