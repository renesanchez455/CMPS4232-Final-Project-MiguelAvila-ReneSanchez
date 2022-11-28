package main

import (
	"net/url"

	"rsma.net/review/pkg/models"
)

type templateData struct {
	Reviews        []*models.Review
	ErrorsFromForm map[string]string
	FormData       url.Values
}
