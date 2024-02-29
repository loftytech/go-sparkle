package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var _db *sql.DB

func Init() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/golang")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	_db = db
}

func Query(query string) (*sql.Rows, error) {
	v, err :=  _db.Query(query)
	return v, err
}

func Exec(query string) (sql.Result, error) {
	v, err := _db.Exec(query)
	return v, err
}

func GetDBInstance() *sql.DB {
	return _db
}
