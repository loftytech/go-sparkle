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

func Query(query string) {
	_db.Query(query)
}

func Exec(query string) {
	_, err := _db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}
