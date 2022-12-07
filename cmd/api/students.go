package main

import (
	"errors"
	"github.com/swsd2544/learny-backend-clone/internal/entity"
	"github.com/swsd2544/learny-backend-clone/internal/repository"
	"net/http"
)

func (app application) registerStudentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username  string `json:"username" validate:"required"`
		Firstname string `json:"firstname" validate:"required"`
		Lastname  string `json:"lastname" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required,regexp="`
	}
	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.validate.Struct(input)
	if err != nil {
		app.failedValidationResponse(w, r, err)
	}

	user := &entity.User{
		Username:    input.Username,
		Firstname:   input.Firstname,
		Lastname:    input.Lastname,
		Email:       input.Email,
		Coin:        0,
		Role:        entity.STUDENT,
		CharacterID: 1,
	}
	err = user.SetPassword(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.repositories.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateEmail):
			app.failedValidationResponse(w, r, err)
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
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,regexp=^(?=.*[A-Za-z])(?=.*\\d)[A-Za-z\\d]{8,}$"`
	}
	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.validate.Struct(input)
	if err != nil {
		app.failedValidationResponse(w, r, err)
	}

	user, err := app.repositories.Users.GetUsersWithEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			app.failedValidationResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	ok, err := user.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	if !ok {
		app.invalidCredentialsResponse(w, r)
	}

	err = writeJSON(w, http.StatusCreated, envelope{"user": *user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
