package viz

import "database/sql"

type VisualizationHandler struct {
	DB *sql.DB
}