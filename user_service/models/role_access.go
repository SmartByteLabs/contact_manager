package models

import (
	"github.com/princeparmar/go-helpers/mysqlmanager"
)

type RoleAccess struct {
	RoleID   int
	AccessID int
}

func (ra *RoleAccess) Create() error {
	query := "INSERT INTO access_role (role_id, access_id, created_date, updated_date) VALUES (?, ?, NOW(), NOW())"
	result, err := mysqlmanager.Exec(query, ra.RoleID, ra.AccessID)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func GetRoleAccess(roleID, accessID int) (*RoleAccess, error) {
	query := "SELECT role_id, access_id FROM access_role WHERE role_id = ? AND access_id = ?"
	row := mysqlmanager.QueryRow(query, roleID, accessID)
	roleAccess := &RoleAccess{}
	err := row.Scan(&roleAccess.RoleID, &roleAccess.AccessID)
	if err != nil {
		return nil, err
	}
	return roleAccess, nil
}

func DeleteRoleAccess(roleID, accessID int) error {
	query := "DELETE FROM access_role WHERE role_id = ? AND access_id = ?"
	_, err := mysqlmanager.Exec(query, roleID, accessID)
	if err != nil {
		return err
	}

	return nil
}

func GetRoleAccesses() ([]*RoleAccess, error) {
	query := "SELECT role_id, access_id FROM access_role"
	rows, err := mysqlmanager.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roleAccesses := []*RoleAccess{}

	for rows.Next() {
		roleAccess := &RoleAccess{}
		err := rows.Scan(&roleAccess.RoleID, &roleAccess.AccessID)
		if err != nil {
			return nil, err
		}
		roleAccesses = append(roleAccesses, roleAccess)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roleAccesses, nil
}

func GetRoleAccessesForRole(roleID int) ([]*Access, error) {
	query := "SELECT a.access_id, a.access_name FROM access a INNER JOIN access_role ar ON a.access_id = ar.access_id WHERE ar.role_id = ?"
	rows, err := mysqlmanager.Query(query, roleID)
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

func GetAccessNamesForUser(userID int) ([]string, error) {
	// Execute the JOIN query
	query := `
		SELECT DISTINCT access.access_name
		FROM user_role
		JOIN role ON user_role.role_id = role.role_id
		JOIN access_role ON role.role_id = access_role.role_id
		JOIN access ON access_role.access_id = access.access_id
		WHERE user_role.user_id = ?
	`
	rows, err := mysqlmanager.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Extract the access names from the result set
	accessNames := make([]string, 0)
	for rows.Next() {
		var accessName string
		err := rows.Scan(&accessName)
		if err != nil {
			return nil, err
		}
		accessNames = append(accessNames, accessName)
	}

	return accessNames, nil
}

func CreateRoleAccessTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS access_role (
		role_id INT NOT NULL,
		access_id INT NOT NULL,
		created_date DATETIME NOT NULL DEFAULT NOW(),
		updated_date DATETIME NOT NULL DEFAULT NOW(),
		PRIMARY KEY (role_id, access_id),
		FOREIGN KEY (role_id) REFERENCES role(role_id) ON DELETE CASCADE,
		FOREIGN KEY (access_id) REFERENCES access(access_id) ON DELETE CASCADE
	)	
`
	_, err := mysqlmanager.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
