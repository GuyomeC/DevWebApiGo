package models

type User struct {
	Id       string  `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Grade    float64 `json:"grade"`
}
