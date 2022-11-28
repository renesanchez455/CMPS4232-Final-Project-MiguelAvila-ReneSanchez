package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

// Struct to hold a review
type Review struct {
	Review_id      int
	Insertion_date time.Time
	Author_name    string
	Book_name      string
	Review         string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to Book Review."))
}

func (app *application) createReviewForm(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/review_form_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (app *application) createReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/quote", http.StatusSeeOther)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w,
			http.StatusText(http.StatusBadRequest),
			http.StatusBadRequest)
		return
	}
	author := r.PostForm.Get("author_name")
	book_name := r.PostForm.Get("book_name")
	review := r.PostForm.Get("review")

	s := `
	INSERT INTO reviews(author_name, book_name, review)
	VALUES ($1, $2, $3)
	`
	_, err = app.db.Exec(s, author, book_name, review)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

func (app *application) displayReview(w http.ResponseWriter, r *http.Request) {
	// SQL statement
	readReview := `
		SELECT *
		FROM reviews
		LIMIT 5
	`
	rows, err := app.db.Query(readReview)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var r Review
		err = rows.Scan(&r.Review_id, &r.Insertion_date, &r.Author_name,
			&r.Book_name, &r.Review)

		if err != nil {
			log.Println(err.Error())
			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		reviews = append(reviews, r)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	// Display the reviews using a template
	ts, err := template.ParseFiles("./ui/html/show_page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, reviews)
	if err != nil {
		log.Println(err.Error())
		http.Error(w,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

}
