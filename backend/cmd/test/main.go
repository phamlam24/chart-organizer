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

	// Try to read and print the dashboards table
	rows, err := db.Query("SELECT id, visualizations FROM dashboards")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var visualizations string
		err := rows.Scan(&id, &visualizations)
		if err != nil {
			log.Fatal(err)
		}
		slog.Info("Viz details", "ID", id, "visualizations", visualizations)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Dropping tables
	// dropDatasetTbl := `
	// 				DROP TABLE IF EXISTS users;
	// 				DROP TABLE IF EXISTS datasets;
	// 				DROP TABLE IF EXISTS visualizations;
	// 				DROP TABLE IF EXISTS dashboards;`
	// _, err = db.Exec(dropDatasetTbl)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
