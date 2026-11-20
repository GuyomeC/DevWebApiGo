package main

import (
	"fmt"
	"net/http"

	"github.com/GuyomeC/DevWebApiGo/application"
	"github.com/GuyomeC/DevWebApiGo/db"
	"github.com/GuyomeC/DevWebApiGo/models"
)

var usersData []models.User

func healthCheck(w http.ResponseWriter, r *http.Request) {
	err := db.Connection.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ping successful")
}

func main() {

	db.Connection = db.NewDB()

	http.HandleFunc("GET /health", healthCheck)
	http.HandleFunc("GET /users", application.GetAllUsers)
	http.HandleFunc("GET /users/{UserId}", application.GetUserById)
	http.HandleFunc("POST /users", application.CreateUser)

	http.ListenAndServe(":8080", nil)
}
