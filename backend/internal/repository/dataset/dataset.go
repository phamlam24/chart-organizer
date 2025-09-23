package dataset

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const datasetStorageFilePath = "../storage/datasets"

// Add a new dataset into our database.
// First, an ID and timestamp is generated.
// Then, add the file bytes into our dataset storage.
// Dataset storage is in `datasetStorageFilePath`
// The file name when stored should be the generated id + ".csv"
// Finally, insert the dataset into our SQL database. Refer to init.go for the schema
func AddNewDataset(db *sql.DB, userId string, name string, file []byte) (string, error) {
	// Generate a UUID4 for the dataset ID
	id := uuid.New().String()

	// Get the current time
	currentTime := time.Now().Format(time.RFC3339)

	// Create the storage directory if it doesn't exist
	err := os.MkdirAll(datasetStorageFilePath, 0755)
	if err != nil {
		return "", err
	}

	// Create the file path with the generated ID + ".csv"
	fileName := id + ".csv"
	filePath := filepath.Join(datasetStorageFilePath, fileName)

	// Write the file bytes to storage
	err = os.WriteFile(filePath, file, 0644)
	if err != nil {
		return "", err
	}

	// Insert the dataset into our SQL database
	_, err = db.Exec("INSERT INTO datasets (id, user_id, name, created_at) VALUES (?, ?, ?, ?)", id, userId, name, currentTime)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetDataset(db *sql.DB, userId, id string) ([]byte, error) {
	var datasetId string
	err := db.QueryRow("SELECT id FROM datasets WHERE id = ? AND user_id = ?", id, userId).Scan(&datasetId)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(datasetStorageFilePath, id+".csv")
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

type DatasetInfo struct {
	ID   string
	Name string
}

func GetAllDatasetsFromUser(db *sql.DB, userId string) ([]DatasetInfo, error) {
	rows, err := db.Query("SELECT id, name FROM datasets WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var datasets []DatasetInfo
	for rows.Next() {
		var info DatasetInfo
		if err := rows.Scan(&info.ID, &info.Name); err != nil {
			return nil, err
		}
		datasets = append(datasets, info)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return datasets, nil
}
