package models

import (
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("models: no matching record found")

// Struct to hold a review
type Review struct {
	Review_id   int
	Created_at  time.Time
	Author_name string
	Book_name   string
	Review      string
}
