package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func setUpDB() (*sql.DB, error) {
	// Provide the credentials for the database
	const (
		host     = "localhost"
		port     = 5432
		user     = "bookreview"
		password = "$swordfish$"
		dbname   = "bookreview"
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Establish a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Testing connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Dependencie Injection
type application struct {
	db *sql.DB
}

func main() {
	var db, err = setUpDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := &application{
		db: db,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/review", app.createReviewForm)
	mux.HandleFunc("/review-add", app.createReview)
	mux.HandleFunc("/show", app.displayReview)
	log.Println("Starting server on port :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
