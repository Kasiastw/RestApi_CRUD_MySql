package main

import (
	"github.com/getground/tech-tasks/backend/cmd/app/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	api.Run()
}

