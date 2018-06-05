package main

import (
  // "fmt"
  "net/http"
  "encoding/json"
  // "strings"
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
  if err := u.validate(); err != nil {
    response, _ := json.Marshal(map[string]error{"error": err})
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write(response)
    return
  }

  // 2. check if email is already registered.
  if err := u.isEmailExists(a.DB); err != nil {
    response, _ := json.Marshal(map[string]map[string]string{"error": { "email": err.Error() }})
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write(response)
    return
  }

  if err := u.createUser(a.DB); err != nil {
    responseJsonErr(w, http.StatusInternalServerError, err.Error())
    return
  }

  defer r.Body.Close()

  // 3. throw empty string json object.
  responseJson(w, http.StatusCreated, map[string]string{})
}

func (a *App) createTokenEndpoint(w http.ResponseWriter, r *http.Request) {
  var u User
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&u); err != nil {
    responseJsonErr(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  // 1. check if the user exists
  if err := u.isUserExists(a.DB); err != nil {
    response, _ := json.Marshal(map[string]map[string]string{"error": { "email": err.Error() }})
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write(response)
    return
  }

  defer r.Body.Close()

  // 2. return token
  responseJson(w, http.StatusCreated, map[string]string{"token": u.getToken(a.DB)})
}
