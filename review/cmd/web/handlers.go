package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"rsma.net/review/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	q, err := app.reviews.Read()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// an instance of templeData
	data := &templateData{
		Reviews: q,
	}

	// Display the reviews using a template
	ts, err := template.ParseFiles("./ui/html/show_page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, data)

	if err != nil {
		app.serverError(w, err)
		return
	}

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
		app.serverError(w, err)
	}
}

func (app *application) createReview(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	author := r.PostForm.Get("author_name")
	book_name := r.PostForm.Get("book_name")
	review := r.PostForm.Get("review")

	errors := make(map[string]string)

	if strings.TrimSpace(author) == "" {
		errors["author_name"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(author) > 100 {
		errors["author_name"] = "This field is too long (maximum is 100 characters)"
	}
	if strings.TrimSpace(book_name) == "" {
		errors["book_name"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(book_name) > 75 {
		errors["book_name"] = "This field is too long (maximum is 75 characters)"
	}
	if strings.TrimSpace(review) == "" {
		errors["review"] = "This field cannot be left blank"
	} else if utf8.RuneCountInString(review) > 300 {
		errors["review"] = "This field is too long (maximum is 300 characters)"
	}

	if len(errors) > 0 {
		ts, err := template.ParseFiles("./ui/html/review_form_page.tmpl")
		if err != nil {
			app.serverError(w, err)
			return
		}
		err = ts.Execute(w, &templateData{
			ErrorsFromForm: errors,
			FormData:       r.PostForm,
		})
		if err != nil {
			log.Println(err.Error())
			app.serverError(w, err)
			return
		}
		return
	}
	// Insert a review
	id, err := app.reviews.Insert(author, book_name, review)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/review/%d", id), http.StatusSeeOther)
}

func (app *application) showReview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	q, err := app.reviews.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%v", q)
}
