package main

import (
	"context"
	"github.com/swsd2544/learny-backend-clone/internal/entity"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func (app application) contextSetUser(r *http.Request, user *entity.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app application) contextGetUser(r *http.Request) *entity.User {
	user, ok := r.Context().Value(userContextKey).(*entity.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
