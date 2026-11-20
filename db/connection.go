package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver   = "mysql"
	host     = "localhost"
	port     = 3306
	user     = "root"
	password = "SQLbase1304!esgi"
	dbname   = "test"
)

var Connection *sql.DB

func NewDB() *sql.DB {

	sqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		user, password, host, port, dbname)

	connection, err := sql.Open(driver, sqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection initiated")

	return connection
}
