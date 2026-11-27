package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/GuyomeC/DevWebApiGo/db"
	"github.com/GuyomeC/DevWebApiGo/models"
	"github.com/GuyomeC/DevWebApiGo/utils"

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
	id, err := uuid.Parse(r.PathValue("UserId"))
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "User id is not a uuid", http.StatusBadRequest)
		return
	}
	user, err := db.GetUserById(id)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Server could not get user", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	msg, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Server could not translate user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", msg)
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

func Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	} // décode le body pour les donnees du client

	existing, err := db.GetUsersByUsername(creds.Username)
	if err != nil {
		http.Error(w, "Error requestiong user by username", http.StatusInternalServerError)
		return
	} // recup l'username dans la DB

	if existing == nil || creds.Username != existing.Username ||
		creds.Password != existing.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	} //verif que le username et password soient egaux si oui token

	token, err := utils.GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	encodedToken, _ := json.Marshal(token)
	fmt.Fprintf(w, "%s", encodedToken)
}
