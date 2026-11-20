package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/GuyomeC/DevWebApiGo/db"
	"github.com/GuyomeC/DevWebApiGo/models"

	"github.com/google/uuid"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	usersData, err := db.GetUsersFromData()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Server could not get users", 500)
		return
	}

	msg, err := json.Marshal(usersData)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Server could not translate users", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", msg)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get user by id !")
	id := r.PathValue("UserId")
	err := uuid.Validate(id)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "User id is not a uuid", 400)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDto models.User
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, "User(id,username,password,grade) needed", http.StatusBadRequest)
		return
	}

	//validation du DataTransferObject DTO

	var errorsMsg []string

	if len(userDto.Username) > 50 || len(userDto.Username) < 8 {
		errorsMsg = append(errorsMsg, "username should be between 8 and 50 chars")
	}
	if len(userDto.Password) > 50 || len(userDto.Password) < 8 {
		errorsMsg = append(errorsMsg, "password should be between 8 and 50 chars")
	}
	if userDto.Grade > 20.0 || userDto.Grade < 0.0 {
		errorsMsg = append(errorsMsg, "grade should be between 0.0 and 20.0 chars")
	}
	if !strings.ContainsAny(userDto.Password, ".-_!?#") {
		errorsMsg = append(errorsMsg, "password must contains at least one special char . - _ ! ? # ")
	}
	if strings.Contains(userDto.Username, "putain") {
		errorsMsg = append(errorsMsg, "username cannot contain the forbidden word")
	}

	if len(errorsMsg) > 0 {
		encodedMsg, _ := json.Marshal(errorsMsg)
		http.Error(w, string(encodedMsg), http.StatusBadRequest)
		return
	}

	existing, err := db.GetUsersByUsername(userDto.Username)
	if existing != nil {
		http.Error(w, "username must be unique", http.StatusConflict)
		return
	}

	userDto.Id = uuid.New().String()

	err = db.InsertUser(userDto)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Server could not create user", 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
