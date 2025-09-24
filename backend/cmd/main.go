package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/glebarez/go-sqlite"

	"connectrpc.com/connect"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"chart-organizer/backend/gen/contracts/auth/v1/authv1connect"
	"chart-organizer/backend/gen/contracts/dataset/v1/datasetv1connect"
	"chart-organizer/backend/gen/contracts/viz/v1/vizv1connect"

	"chart-organizer/backend/internal/handlers/auth"
	"chart-organizer/backend/internal/handlers/dataset"
	"chart-organizer/backend/internal/handlers/viz"
	"chart-organizer/backend/internal/interceptors"
	"chart-organizer/backend/internal/repository"
)

func getAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return "0.0.0.0:" + port
	}
	if addr := os.Getenv("ADDR"); addr != "" {
		return addr
	}
	return "0.0.0.0:8080"
}

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	slog.Info(".env successfully loaded")

	// Get database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./storage/chart-organizer.db"
	}

	// Start the database
	db, err := sql.Open("sqlite", dbPath)
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

	// Configure interceptors
	interceptors.JwtKey = []byte(os.Getenv("JWT_KEY"))

	// Create interceptors
	debugInterceptor := interceptors.NewDebugInterceptor()
	authInterceptor := interceptors.NewAuthInterceptor()

	// Configure Connect options with interceptors
	connectOptions := connect.WithInterceptors(debugInterceptor, authInterceptor)

	// Adding routes with interceptors
	mux := http.NewServeMux()

	authPath, authHandler := authv1connect.NewAuthServiceHandler(&auth.AuthHandler{DB: db}, connectOptions)
	datasetPath, datasetHandler := datasetv1connect.NewDatasetServiceHandler(&dataset.DatasetHandler{DB: db}, connectOptions)
	vizPath, vizHandler := vizv1connect.NewDashboardServiceHandler(&viz.VisualizationHandler{DB: db}, connectOptions)

	mux.Handle(authPath, authHandler)
	mux.Handle(datasetPath, datasetHandler)
	mux.Handle(vizPath, vizHandler)

	// Get server address
	addr := getAddr()

	// Starting the server
	slog.Info(fmt.Sprintf("Server hosting at %s", addr))

	http.ListenAndServe(
		addr,
		// Use h2c so we can serve HTTP/2 without TLS, with CORS support
		interceptors.CORSHandler(h2c.NewHandler(mux, &http2.Server{})),
	)
}
