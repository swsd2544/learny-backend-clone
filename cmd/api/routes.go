package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app application) routes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/v1/students", app.registerStudentHandler).Methods(http.MethodPost)
	r.HandleFunc("/v1/students/login", app.loginStudentHandler).Methods(http.MethodPut)
	return r
}
