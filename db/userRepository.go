package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/GuyomeC/DevWebApiGo/models"

	"github.com/google/uuid"
)

func GetUsersByUsername(username string) (*models.User, error) {
	var user models.User
	row := Connection.QueryRow("SELECT id, username, password, grade FROM user WHERE username = ?", username)
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Grade)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserByUsername(db) (username : %s) : %s", username, err.Error())
	}
	return &user, nil
}

func GetUserById(id uuid.UUID) (*models.User, error) {
	var user models.User
	row := Connection.QueryRow("SELECT id, username, password, grade FROM user WHERE id = ?", id)
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Grade)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserById(db) (id : %s) : %s", id, err.Error())
	}
	return &user, nil
}

func GetUsersFromData() ([]models.User, error) {
	var users []models.User

	rows, err := Connection.Query("SELECT id, username, password, grade FROM test.user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Id, &user.Username, &user.Password, &user.Grade)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func InsertUser(user models.User) error {
	_, err := Connection.Exec("INSERT INTO test.user(id, username, password, grade) VALUES (?,?,?,?)",
		user.Id, user.Username, user.Password, user.Grade)

	if err != nil {
		return err
	}

	return nil
}
