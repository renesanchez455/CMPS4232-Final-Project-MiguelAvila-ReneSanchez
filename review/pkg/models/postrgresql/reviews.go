package postgresql

import (
	"database/sql"
	"errors"

	"rsma.net/review/pkg/models"
)

type ReviewModel struct {
	DB *sql.DB
}

func (m *ReviewModel) Insert(author, book_name, review string) (int, error) {
	var id int

	s := `
	INSERT INTO reviews(author_name, book_name, review)
	VALUES ($1, $2, $3)
	RETURNING review_id
	`
	err := m.DB.QueryRow(s, author, book_name, review).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *ReviewModel) Read() ([]*models.Review, error) {
	s := `
		SELECT author_name, book_name, review
		FROM reviews
		LIMIT 20
	`
	rows, err := m.DB.Query(s)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []*models.Review{}

	for rows.Next() {
		r := &models.Review{}
		err = rows.Scan(&r.Author_name, &r.Book_name, &r.Review)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}

func (m *ReviewModel) Get(id int) (*models.Review, error) {
	s := `
		SELECT author_name, book_name, review
		FROM reviews
		WHERE review_id = $1
	`
	r := &models.Review{}
	err := m.DB.QueryRow(s, id).Scan(&r.Author_name, &r.Book_name, &r.Review)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrRecordNotFound
		} else {
			return nil, err
		}
	}
	return r, nil
}
