package models

import (
	"errors"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := "INSERT INTO users(email, password) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword)
	if err != nil {
		return err
	}

	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func GetAllUsers() ([]User, error) {
	query := "	SELECT * FROM users"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, nil
}

func (u *User) ValidateCredentials() error {
	query := "SELECT ID, password FROM users WHERE email = ?"
	row := db.DB.QueryRow(query, u.Email)

	var retrivedPassword string
	err := row.Scan(&u.ID, &retrivedPassword)
	if err != nil {
		return errors.New("Can't find user!")
	}

	passwordIsValid := utils.CheckPassword(retrivedPassword, u.Password)

	if !passwordIsValid {
		return errors.New("Invalid credentials!")
	}

	return nil
}
