package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var dbDriver *sql.DB

// TestPrint Called from main
func TestPrint(str string) {

	fmt.Println("HI" + str)
}

// PgConn return conn handle
func PgConn(dataSourceName string) *sql.DB {
	var err error
	fmt.Println("connecting...")
	dbDriver, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}
	//defer db.Close()

	if err = dbDriver.Ping(); err != nil {
		log.Panic(err)
	}

	fmt.Println("#****Successfully connected*****#")
	return dbDriver
}
