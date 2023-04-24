package models

import (
	"github.com/princeparmar/go-helpers/mysqlmanager"
)

type Role struct {
	ID   int
	Name string
}

func (r *Role) Create() error {
	query := "INSERT INTO roles (role_name, created_date, updated_date) VALUES (?, NOW(), NOW())"
	result, err := mysqlmanager.Exec(query, r.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	r.ID = int(id)

	return nil
}

func GetRole(id int) (*Role, error) {
	query := "SELECT role_id, role_name FROM roles WHERE role_id = ?"
	row := mysqlmanager.QueryRow(query, id)
	role := &Role{}
	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *Role) Update() error {
	query := "UPDATE roles SET role_name = ?, updated_date = NOW() WHERE role_id = ?"
	_, err := mysqlmanager.Exec(query, r.Name, r.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteRole(id int) error {
	query := "DELETE FROM roles WHERE role_id = ?"
	_, err := mysqlmanager.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func GetRoles() ([]*Role, error) {
	query := "SELECT role_id, role_name FROM roles"
	rows, err := mysqlmanager.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := []*Role{}

	for rows.Next() {
		role := &Role{}
		err := rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func CreateRoleTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS roles (
            role_id INT AUTO_INCREMENT PRIMARY KEY,
            role_name VARCHAR(255) NOT NULL UNIQUE,
			created_date DATETIME NOT NULL DEFAULT NOW(),
			updated_date DATETIME NOT NULL DEFAULT NOW()
			)`
	_, err := mysqlmanager.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
