package main

import (
  // "fmt"
  "strconv"
  "net/http"
  "encoding/json"

  "github.com/gorilla/context"
  "github.com/gorilla/mux"
  "github.com/dgrijalva/jwt-go"
  "github.com/mitchellh/mapstructure"
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
  token, err := u.isUserExists(a.DB)
  if err != nil {
    response, _ := json.Marshal(map[string]map[string]string{"error": { "email": err.Error() }})
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write(response)
    return
  }

  defer r.Body.Close()

  // 2. return token
  responseJson(w, http.StatusCreated, map[string]string{"token": token})
}

func (a *App) createJob(w http.ResponseWriter, r *http.Request) {
  var j Job
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&j); err != nil {
    responseJsonErr(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  // add validation
  if err := j.validate(); err != nil {
    response, _ := json.Marshal(map[string]error{"error": err})
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write(response)
    return
  }

  // decode token from headers
  var u User
  decoded := context.Get(r, "decoded")
  mapstructure.Decode(decoded.(jwt.MapClaims), &u)

  // create new job post
  if err := j.createJob(a.DB, u.Id); err != nil {
    responseJsonErr(w, http.StatusInternalServerError, err.Error())
    return
  }

  defer r.Body.Close()

  // 2. throw empty string json object.
  responseJson(w, http.StatusCreated, map[string]string{})
}

func (a *App) getJob(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    responseJsonErr(w, http.StatusBadRequest, "Invalid user ID")
    return
  }

  j := Job{
    Id: id,
  }

  job, err := j.getJob(a.DB)
  if err != nil {
    responseJsonErr(w, http.StatusInternalServerError, err.Error())
    return
  }
  responseJson(w, http.StatusOK, job)
}

func (a *App) getJobs(w http.ResponseWriter, r *http.Request) {
  base := 10
  limit, _ := strconv.Atoi(r.FormValue("limit"))
  offset, _ := strconv.Atoi(r.FormValue("offset"))

  if limit > base || limit < 1 {
    limit = base
  }
  if offset < 0 {
    offset = 0
  }

  var j Job
  jobs, err := j.getJobs(a.DB, offset, limit)
  if err != nil {
    responseJsonErr(w, http.StatusInternalServerError, err.Error())
    return
  }
  responseJson(w, http.StatusOK, jobs)
}

func (a *App) updateJob(w http.ResponseWriter, r *http.Request) {
  var j Job
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    responseJsonErr(w, http.StatusBadRequest, "Missing job id")
    return
  }

  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&j); err != nil {
    responseJsonErr(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  // add validation
  if err := j.validate(); err != nil {
    response, _ := json.Marshal(map[string]error{"error": err})
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    w.Write(response)
    return
  }

  // decode token from headers
  var u User
  decoded := context.Get(r, "decoded")
  mapstructure.Decode(decoded.(jwt.MapClaims), &u)

  // update job base on user id and jobId
  updatedJob, err := j.updateJob(a.DB, u.Id, id);
  if err != nil {
    responseJsonErr(w, http.StatusInternalServerError, err.Error())
    return
  }

  defer r.Body.Close()

  // 2. throw empty string json object.
  responseJson(w, http.StatusCreated, updatedJob)
}
