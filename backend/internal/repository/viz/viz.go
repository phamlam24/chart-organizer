package viz

import (
	vizv1 "chart-organizer/backend/gen/contracts/viz/v1"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

func AddNewDashboard(db *sql.DB, userId string, datasetId string, visualizations []*vizv1.Visualization) (string, error) {
	// Generate a UUID4 for the dataset ID
	id := uuid.New().String()

	// Get the current time
	currentTime := time.Now().Format(time.RFC3339)

	visualizationsJson, err := json.Marshal(visualizations)
	if err != nil {
		return "", err
	}

	// Insert the dataset into our SQL database
	_, err = db.Exec("INSERT INTO dashboards (id, dataset_id, visualizations, created_at) VALUES (?, ?, ?, ?)", id, datasetId, visualizationsJson, currentTime)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetDashboard(db *sql.DB, id string) ([]*vizv1.Visualization, string, error) {
	var visualizationsJson string
	var datasetId string
	row := db.QueryRow("SELECT visualizations, dataset_id FROM dashboards WHERE id = ?", id)
	err := row.Scan(&visualizationsJson, &datasetId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", nil
		}
		return nil, "", err
	}

	var visualizations []*vizv1.Visualization
	err = json.Unmarshal([]byte(visualizationsJson), &visualizations)
	if err != nil {
		return nil, "", err
	}

	return visualizations, datasetId, nil
}
