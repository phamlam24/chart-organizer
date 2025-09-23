package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/glebarez/go-sqlite"

	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"chart-organizer/backend/gen/contracts/auth/v1/authv1connect"
	"chart-organizer/backend/gen/contracts/dataset/v1/datasetv1connect"

	"chart-organizer/backend/internal/handlers/auth"
	"chart-organizer/backend/internal/handlers/dataset"
	"chart-organizer/backend/internal/middleware"
	"chart-organizer/backend/internal/repository"
)

const addr = "localhost:8080"

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	slog.Info(".env successfully loaded")
	
	// Setting up middleware
	middleware.JwtKey = []byte(os.Getenv("JWT_KEY"))

	// Start the database
	db, err := sql.Open("sqlite", "../storage/chart-organizer.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repository.InitDatabase(db)

	var sqliteVersion string
	err = db.QueryRow("SELECT sqlite_version()").Scan(&sqliteVersion)
	if err != nil {
		fmt.Println(err)
		return
	}
	slog.Info(fmt.Sprintf("SQLite version %s loaded", sqliteVersion))

	// Adding routes
	mux := http.NewServeMux()

	authPath, authHandler := authv1connect.NewAuthServiceHandler(&auth.AuthHandler{DB: db})
	datasetPath, datasetHandler := datasetv1connect.NewDatasetServiceHandler(&dataset.DatasetHandler{DB: db})
	
	mux.Handle(authPath, authHandler)
	mux.Handle(datasetPath, middleware.AuthMiddleware(datasetHandler))

	// Starting the server
	slog.Info(fmt.Sprintf("Server hosting at %s", addr))
	http.ListenAndServe(
		addr,
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
