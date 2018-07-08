package app

import (
    // "log"
    "encoding/json"
    "net/http"
    "go-rest-api/helpers"
)

func (app *App) Authenticate(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    if r.Method != "POST" {
        helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException("Method Not Allowed."))
        return
    }

    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&a); err != nil {
        helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException("Invalid request payload."))
        return
    }

    if err := app.Auth.Validate(a); err != nil {
        helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException(err))
        return
    }

    if err := app.Auth.IsUserExists(a); err != nil {
        helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException(err.Error()))
        return
    }

    // Get User Profile
    // app.User.GetByEmail(a.Email)
    // token, err := app.Auth.CreateToken(user)
    // if err != nil {
    //     helpers.HttpResponse(w, http.StatusNotFound, helpers.ResponseException(err.Error()))
    //     return
    // }

    // log.Println(token)

    // responseMessage := map[string]string{"token": token}
    // helpers.HttpResponse(w, http.StatusCreated, responseMessage)
}

