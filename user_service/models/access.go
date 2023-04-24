package models

import (
	"github.com/princeparmar/go-helpers/mysqlmanager"
)

type Access struct {
	ID   int
	Name string
}

func (a *Access) Create() error {
	query := "INSERT INTO access (access_name, created_date, updated_date) VALUES (?, NOW(), NOW())"
	result, err := mysqlmanager.Exec(query, a.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = int(id)

	return nil
}

func GetAccess(id int) (*Access, error) {
	query := "SELECT access_id, access_name FROM access WHERE access_id = ?"
	row := mysqlmanager.QueryRow(query, id)
	access := &Access{}
	err := row.Scan(&access.ID, &access.Name)
	if err != nil {
		return nil, err
	}
	return access, nil
}

func (a *Access) Update() error {
	query := "UPDATE access SET access_name = ?, updated_date = NOW() WHERE access_id = ?"
	_, err := mysqlmanager.Exec(query, a.Name, a.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAccess(id int) error {
	query := "DELETE FROM access WHERE access_id = ?"
	_, err := mysqlmanager.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func GetAccesses() ([]*Access, error) {
	query := "SELECT access_id, access_name FROM access"
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accesses, nil
}

func CreateAccessTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS access (
            access_id INT AUTO_INCREMENT PRIMARY KEY,
            access_name VARCHAR(255) NOT NULL UNIQUE,
			created_date DATETIME NOT NULL DEFAULT NOW(),
			updated_date DATETIME NOT NULL DEFAULT NOW(),
			)`
	_, err := mysqlmanager.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
