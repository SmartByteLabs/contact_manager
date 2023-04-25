package models

import (
	"github.com/princeparmar/go-helpers/mysqlmanager"
)

// Access represents an access object with ID and Name fields
type Access struct {
	ID   int
	Name string
}

// Create creates a new access object in the database and sets its ID
func (a *Access) Create() error {
	// Prepare the query to insert a new access object
	query := "INSERT INTO access (access_name, created_date, updated_date) VALUES (?, NOW(), NOW())"
	// Execute the query with the access name parameter
	result, err := mysqlmanager.Exec(query, a.Name)
	if err != nil {
		return err
	}

	// Get the ID of the newly inserted access object
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID of the access object
	a.ID = int(id)

	return nil
}

// GetAccess retrieves an access object with the given ID from the database
func GetAccess(id int) (*Access, error) {
	// Prepare the query to select an access object by ID
	query := "SELECT access_id, access_name FROM access WHERE access_id = ?"
	// Execute the query with the ID parameter
	row := mysqlmanager.QueryRow(query, id)
	access := &Access{}
	err := row.Scan(&access.ID, &access.Name)
	return access, err
}

// Update updates an access object in the database with the new data
func (a *Access) Update() error {
	// Prepare the query to update an access object by ID
	query := "UPDATE access SET access_name = ?, updated_date = NOW() WHERE access_id = ?"
	// Execute the query with the access name and ID parameters
	_, err := mysqlmanager.Exec(query, a.Name, a.ID)

	return err
}

// DeleteAccess deletes an access object with the given ID from the database
func DeleteAccess(id int) error {
	// Prepare the query to delete an access object by ID
	query := "DELETE FROM access WHERE access_id = ?"
	// Execute the query with the ID parameter
	_, err := mysqlmanager.Exec(query, id)

	return err
}

// GetAccesses retrieves all access objects from the database
func GetAccesses() ([]*Access, error) {
	// Prepare the query to select all access objects
	query := "SELECT access_id, access_name FROM access"
	// Execute the query
	rows, err := mysqlmanager.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accesses := []*Access{}

	for rows.Next() {
		access := &Access{}
		err := rows.Scan(&access.ID, &access.Name)
		if err != nil {
			return nil, err
		}
		accesses = append(accesses, access)
	}

	return accesses, rows.Err()
}

// CreateAccessTable creates the access table in the database
func CreateAccessTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS access (
            access_id INT AUTO_INCREMENT PRIMARY KEY,
            access_name VARCHAR(255) NOT NULL UNIQUE,
			created_date DATETIME NOT NULL DEFAULT NOW(),
			updated_date DATETIME NOT NULL DEFAULT NOW(),
			)`
	// Execute the query to create the access table
	_, err := mysqlmanager.Exec(query)
	return err
}
