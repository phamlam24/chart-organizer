package repository

import (
	"database/sql"
)

func InitDatabase(db *sql.DB) error {
	// Create each table
	createUserTbl := `CREATE TABLE IF NOT EXISTS users 
					(id TEXT NOT NULL PRIMARY KEY, 
					username TEXT NOT NULL UNIQUE,
					password_hash TEXT NOT NULL,
					created_at TEXT NOT NULL
					);`
	_, err := db.Exec(createUserTbl)
	if err != nil {
		return err
	}

	createDatasetTbl := `CREATE TABLE IF NOT EXISTS datasets
						(id TEXT NOT NULL PRIMARY KEY,
						user_id TEXT NOT NULL,
						name TEXT NOT NULL,
						created_at TEXT NOT NULL,
						FOREIGN KEY (user_id) REFERENCES users (id)
						);`
	_, err = db.Exec(createDatasetTbl)
	if err != nil {
		return err
	}

	createDashboardTbl := `CREATE TABLE IF NOT EXISTS dashboards
						(id TEXT NOT NULL PRIMARY KEY, 
						dataset_id TEXT NOT NULL,
						visualizations TEXT NOT NULL,
						created_at TEXT NOT NULL,
						FOREIGN KEY (dataset_id) REFERENCES datasets (id)
						);`
	_, err = db.Exec(createDashboardTbl)
	if err != nil {
		return err
	}

	return nil
}
