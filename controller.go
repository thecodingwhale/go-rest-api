package main

import (
  "net/http"
  "encoding/json"
)

func (app *App) getUsers(w http.ResponseWriter, r *http.Request) {
  users, err := getUsers(app.DB)
  if err != nil {
    responseJsonErr(w, http.StatusInternalServerError, err.Error())
    return
  }
  responseJson(w, http.StatusOK, users)
}

func (a *App) createUser(w http.ResponseWriter, r *http.Request) {
  var u User
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&u); err != nil {
    responseJsonErr(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  // 1. add validation
  //   - email should be valid
  //   - password minimum of 8 characters
/*
  {
    error: {
      email: "error email message",
      password: "error password message",
    }
  }
*/
  if err := u.Validate(); err != nil {
    response, _ := json.Marshal(err)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write(response)
    return
  }

  defer r.Body.Close()

  // 2. check if email is already registered.
  // if err := u.createUser(a.DB); err != nil {
  //   responseJsonErr(w, http.StatusInternalServerError, err.Error())
  //   return
  // }

  responseJson(w, http.StatusCreated, u)
}