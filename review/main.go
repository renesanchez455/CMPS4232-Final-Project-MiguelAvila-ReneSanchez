package main

import (
	// This provided by Go
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Provide the credentials for the database
const (
	host     = "localhost"
	port     = 5432
	user     = "bookreview"
	password = "$swordfish$"
	dbname   = "bookreview"
)

// dsn : data source name
func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Establish a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Testing connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	// Insert a review
	insertReview := `
	INSERT INTO reviews (author_name, book_name, review)
	VALUES ($1, $2, $3)
	RETURNING review_id, 
	`
	review_id := 0
	err = db.QueryRow(insertReview, "John Doe",
		"Deep Sea Exploration", "Awesome deep sea content").Scan(&review_id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The recently inserted record has review_id:", review_id)

}
