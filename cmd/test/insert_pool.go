package main

// schema we can use along with some select statements
// create table test ( gopher_id int, created timestamp );
// select * from test order by created asc limit 1;
// select * from test order by created desc limit 1;
// select count(created) from test;

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	gophers = 10
	entries = 1000
	//host     = "192.168.99.100"
	host     = "localhost"
	port     = 54320
	user     = "postgres"
	password = "postgres"
	dbname   = "ddoor_db"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {

	defer timeTrack(time.Now(), "insert_pool")

	var wg sync.WaitGroup

	// lazily open db (doesn't truly open until first request)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB Successfully connected!")

	// create string to pass
	var sStmt string = "insert into test_pool (gopher_id, created) values ($1, $2)"

	// run the insert function using 10 go routines
	for i := 1; i <= gophers; i++ {
		wg.Add(1)
		// spin up a gopher
		go gopher(&wg, i, sStmt, db)
	}

	fmt.Println("Main: Waiting for workers to finish")
	wg.Wait()
	fmt.Println("Main: Completed")

}

func gopher(wg *sync.WaitGroup, gopher_id int, sStmt string, db *sql.DB) {

	defer wg.Done()
	fmt.Printf("Gopher Id: %v || StartTime: %v\n", gopher_id, time.Now())

	stmt, err := db.Prepare(sStmt)
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= entries; i++ {
		res, err := stmt.Exec(gopher_id, time.Now())
		if err != nil || res == nil {
			log.Fatal(err)
		}
	}
	// close prepared stmt
	stmt.Close()

	fmt.Printf("Gopher Id: %v || StopTime: %v\n", gopher_id, time.Now())

}
