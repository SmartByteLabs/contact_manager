package models

import (
	"github.com/princeparmar/go-helpers/mysqlmanager"
)

type User struct {
	ID       int
	UserName string
	Mobile   string
	EmailID  string
}

func (u *User) Create() error {
	query := "INSERT INTO users (user_name, mobile, created_date, updated_date, email_id) VALUES (?, ?, NOW(), NOW(), ?)"
	result, err := mysqlmanager.Exec(query, u.UserName, u.Mobile, u.EmailID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = int(id)

	return nil
}

func GetUser(id int) (*User, error) {
	query := "SELECT user_id, user_name, mobile, email_id FROM users WHERE user_id = ?"
	row := mysqlmanager.QueryRow(query, id)
	user := &User{}
	err := row.Scan(&user.ID, &user.UserName, &user.Mobile, &user.EmailID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Update() error {
	query := "UPDATE users SET user_name = ?, mobile = ?, updated_date = NOW(), email_id = ? WHERE user_id = ?"
	_, err := mysqlmanager.Exec(query, u.UserName, u.Mobile, u.EmailID, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id int) error {
	query := "DELETE FROM users WHERE user_id = ?"
	_, err := mysqlmanager.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func GetUsers() ([]*User, error) {
	query := "SELECT user_id, user_name, mobile, email_id FROM users"
	rows, err := mysqlmanager.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.UserName, &user.Mobile, &user.EmailID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func CreateUserTable() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            user_id INT AUTO_INCREMENT PRIMARY KEY,
            user_name VARCHAR(255) NOT NULL UNIQUE,
            mobile VARCHAR(10) NOT NULL,
            email_id VARCHAR(255) NOT NULL,
			created_date DATETIME NOT NULL DEFAULT NOW(),
			updated_date DATETIME NOT NULL DEFAULT NOW(),
			)`
	_, err := mysqlmanager.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
