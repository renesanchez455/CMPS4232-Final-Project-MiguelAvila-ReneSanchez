package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	postgresql "rsma.net/review/pkg/models/postrgresql"
)

func setUpDB(dsn string) (*sql.DB, error) {
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

// Dependency Injection
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	reviews  *postgresql.ReviewModel
}

func main() {
	// Create a command line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", os.Getenv("BOOKREVIEW_DB_DSN"), "PostgreSQL DSN (Data Source Name)")
	flag.Parse()
	// Create a logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var db, err = setUpDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		reviews: &postgresql.ReviewModel{
			DB: db,
		},
	}
	// Create a custom web server
	srv := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errorLog,
	}

	// Start our server
	infoLog.Printf("Starting server on port %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
