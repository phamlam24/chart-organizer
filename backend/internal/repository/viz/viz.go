package viz

import (
	vizv1 "chart-organizer/backend/gen/contracts/viz/v1"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

func AddNewDashboard(db *sql.DB, userId string, datasetId string, visualizations []*vizv1.Visualization) (string, error) {
	// Generate a UUID4 for the dataset ID
	id := uuid.New().String()

	// Get the current time
	currentTime := time.Now().Format(time.RFC3339)

	// Marshal each visualization using protojson
	var jsonVizs [][]byte
	for _, viz := range visualizations {
		vizJson, err := protojson.Marshal(viz)
		if err != nil {
			return "", err
		}
		jsonVizs = append(jsonVizs, vizJson)
	}

	// Combine into a JSON array
	visualizationsJson := "["
	for i, vizJson := range jsonVizs {
		if i > 0 {
			visualizationsJson += ","
		}
		visualizationsJson += string(vizJson)
	}
	visualizationsJson += "]"

	// Insert the dataset into our SQL database
	_, err := db.Exec("INSERT INTO dashboards (id, dataset_id, visualizations, created_at) VALUES (?, ?, ?, ?)", id, datasetId, visualizationsJson, currentTime)
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

	// First, unmarshal as raw JSON to get the array structure
	var rawVizs []json.RawMessage
	err = json.Unmarshal([]byte(visualizationsJson), &rawVizs)
	if err != nil {
		return nil, "", err
	}

	// Then unmarshal each visualization using protojson
	var visualizations []*vizv1.Visualization
	for _, rawViz := range rawVizs {
		viz := &vizv1.Visualization{}
		err = protojson.Unmarshal(rawViz, viz)
		if err != nil {
			return nil, "", err
		}
		visualizations = append(visualizations, viz)
	}

	return visualizations, datasetId, nil
}
