package app

import (
    "encoding/json"
    "net/http"
    "go-rest-api/helpers"
)

func (app *App) UsersCreate(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    if r.Method != "POST" {
        helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException("Method Not Allowed."))
        return
    }

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&u); err != nil {
        helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException("Invalid request payload."))
        return
    }

    if err := u.Validate(); err != nil {
        helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException(err))
        return
    }

    if err := app.User.IsEmailExists(u); err != nil {
        helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException("Email already exists."))
        return
    }

    if err := app.User.Create(u); err != nil {
        helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException("Something went wrong."))
        return
    }

    helpers.HttpResponse(w, http.StatusCreated, map[string]string{})
}
