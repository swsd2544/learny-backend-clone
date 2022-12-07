package main

import (
	"errors"
	"github.com/swsd2544/learny-backend-clone/internal/entity"
	"github.com/swsd2544/learny-backend-clone/internal/repository"
	"github.com/swsd2544/learny-backend-clone/internal/validator"
	"net/http"
)

func (app application) registerStudentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &entity.User{
		Username:    input.Username,
		Firstname:   input.Firstname,
		Lastname:    input.Lastname,
		Email:       input.Email,
		Coin:        0,
		Role:        entity.RoleStudent,
		CharacterID: 1,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if entity.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.repositories.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = writeJSON(w, http.StatusCreated, envelope{"user": *user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app application) loginStudentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	entity.ValidateEmail(v, input.Email)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.repositories.Users.GetUserWithEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	ok, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !ok {
		app.invalidCredentialsResponse(w, r)
		return
	}

	err = writeJSON(w, http.StatusCreated, envelope{"user": *user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
