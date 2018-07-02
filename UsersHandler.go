package main

import (
	"encoding/json"
	"net/http"
)

func ResponseExemption(e interface{}) map[string]Exception {
	return map[string]Exception{"error": Exception{Message: e}};
}

func (app *App) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		httpResponse(w, http.StatusNotFound, ResponseExemption("Invalid request payload."))
		return
	}

	if err := u.validate(); err != nil {
		httpResponse(w, http.StatusNotFound, ResponseExemption(err))
		return
	}

	if err := u.isEmailExists(app.DB); err != nil {
		httpResponse(w, http.StatusNotFound, ResponseExemption("Email already exists."))
		return
	}

	if err := u.CreateUser(app.DB); err != nil {
		httpResponse(w, http.StatusNotFound, ResponseExemption("Something went wrong."))
		return
	}

	defer r.Body.Close()

	httpResponse(w, http.StatusCreated, map[string]string{})
}

