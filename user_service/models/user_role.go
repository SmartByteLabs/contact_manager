package models

import (
	"time"

	"github.com/princeparmar/go-helpers/mysqlmanager"
)

type UserRole struct {
	UserID     int
	RoleID     int
	ExpiryDate time.Time
}

func (ur *UserRole) Create() error {
	query := "INSERT INTO user_roles (user_id, role_id, expiry_date) VALUES (?, ?, ?)"
	result, err := mysqlmanager.Exec(query, ur.UserID, ur.RoleID, ur.ExpiryDate)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func GetUserRole(userID, roleID int) (*UserRole, error) {
	query := "SELECT user_id, role_id, expiry_date FROM user_roles WHERE user_id = ? AND role_id = ?"
	row := mysqlmanager.QueryRow(query, userID, roleID)
	userRole := &UserRole{}
	err := row.Scan(&userRole.UserID, &userRole.RoleID, &userRole.ExpiryDate)
	if err != nil {
		return nil, err
	}
	return userRole, nil
}

func (ur *UserRole) Update() error {
	query := "UPDATE user_roles SET expiry_date = ? WHERE user_id = ? AND role_id = ?"
	_, err := mysqlmanager.Exec(query, ur.ExpiryDate, ur.UserID, ur.RoleID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserRole(userID, roleID int) error {
	query := "DELETE FROM user_roles WHERE user_id = ? AND role_id = ?"
	_, err := mysqlmanager.Exec(query, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserRoles() ([]*UserRole, error) {
	query := "SELECT user_id, role_id, expiry_date FROM user_roles"
	rows, err := mysqlmanager.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	userRoles := []*UserRole{}

	for rows.Next() {
		userRole := &UserRole{}
		err := rows.Scan(&userRole.UserID, &userRole.RoleID, &userRole.ExpiryDate)
		if err != nil {
			return nil, err
		}
		userRoles = append(userRoles, userRole)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userRoles, nil
}

func GetUserRolesForUser(userID int) ([]*Role, error) {
	query := "SELECT r.role_id, r.role_name FROM roles r INNER JOIN user_roles ur ON r.role_id = ur.role_id WHERE ur.user_id = ?"
	rows, err := mysqlmanager.Query(query, userID)
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

func CreateUserRolesTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS user_roles (
            user_id INT NOT NULL,
            role_id INT NOT NULL,
            expiry_date DATETIME,
			created_date DATETIME NOT NULL DEFAULT NOW(),
			updated_date DATETIME NOT NULL DEFAULT NOW(),
			PRIMARY KEY (user_id, role_id),
            FOREIGN KEY (user_id) REFERENCES users(user_id),
            FOREIGN KEY (role_id) REFERENCES roles(role_id)
        )`
	_, err := mysqlmanager.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
