package main

import (
	"github.com/gorilla/schema"
	"github.com/swsd2544/learny-backend-clone/internal/validator"
	"net/http"
)

func (app application) getClassesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var filter struct {
		Offset     int64  `schema:"offset"`
		Limit      int64  `schema:"limit"`
		SortBy     string `schema:"sort_by"`
		Asc        bool   `schema:"asc"`
		IsEnrolled bool   `schema:"is_enrolled"`
	}
	if err := schema.NewDecoder().Decode(filter, r.Form); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(filter.Offset >= 0, "offset", "must be positive")
	v.Check(filter.Limit >= 0, "limit", "must be positive")
	v.Check(validator.PermittedValue(filter.SortBy, "id", "name", "created_at"), "sort_by", "must be either id, name, or created_at")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

}
