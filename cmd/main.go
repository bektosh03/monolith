package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
	"github.com/bektosh03/monolith/api"
	"github.com/bektosh03/monolith/api/handlers"
)

const (
	dbname  = "crud"
	dbuser  = "bektosh"
	dbpass  = "12345"
	dbhost  = "localhost"
	dbport  = 5432
	sslmode = "disable"
)

var connString = fmt.Sprintf(
	"host=%s user=%s password=%s port=%d dbname=%s sslmode=%s",
	dbhost, dbuser, dbpass, dbport, dbname, sslmode,
)

func initDb(dbString string) *sql.DB {
	db, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func main() {
	db := initDb(connString)
	defer db.Close()
	handler := handlers.NewHandler(db)
	app := api.New(handler)
	err := app.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
