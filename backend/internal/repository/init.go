package repository

import (
	"database/sql"
)

func InitDatabase(db *sql.DB) error {
	// Create each table
	createUserTbl := `CREATE TABLE IF NOT EXISTS users 
					(id INTEGER NOT NULL PRIMARY KEY, 
					username TEXT NOT NULL UNIQUE,
					password_hash TEXT NOT NULL,
					created_at TEXT NOT NULL
					);`
	_, err := db.Exec(createUserTbl)
	if err != nil {
		return err
	}

	createDatasetTbl := `CREATE TABLE IF NOT EXISTS datasets
						(id INTEGER NOT NULL PRIMARY KEY,
						user_id INTEGER NOT NULL,
						filename TEXT NOT NULL,
						filepath TEXT NOT NULL,
						created_at TEXT NOT NULL,
						FOREIGN KEY (user_id) REFERENCES users (id)
						);`
	_, err = db.Exec(createDatasetTbl)
	if err != nil {
		return err
	}

	createVizTbl := `CREATE TABLE IF NOT EXISTS visualizations
						(id INTEGER NOT NULL PRIMARY KEY, 
						user_id INTEGER NOT NULL,
						dataset_id INTEGER NOT NULL,
						plot_type TEXT NOT NULL,
						config TEXT NOT NULL,
						created_at TEXT NOT NULL,
						FOREIGN KEY (user_id) REFERENCES users (id),
						FOREIGN KEY (dataset_id) REFERENCES datasets (id)
						);`
	_, err = db.Exec(createVizTbl)
	if err != nil {
		return err
	}

	return nil
}
