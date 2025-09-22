package main

import (
	"database/sql"
	"log"
	"log/slog"

	_ "github.com/glebarez/go-sqlite"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	slog.Info(".env successfully loaded")

	// Start the database
	db, err := sql.Open("sqlite", "../storage/chart-organizer.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


}