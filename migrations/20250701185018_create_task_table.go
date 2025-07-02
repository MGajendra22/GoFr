package migrations

import "gofr.dev/pkg/gofr/migration"

const createTaskTableSQL = `
CREATE TABLE IF NOT EXISTS tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    description TEXT,
    status BOOLEAN DEFAULT FALSE,
    userid INT NOT NULL
);`

const createUserTableSQL = `
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);`

func createTaskTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createTaskTableSQL)
			if err != nil {
				return err
			}

			_, err = d.SQL.Exec(createUserTableSQL)
			if err != nil {
				return err
			}

			return nil

		},
	}
}
