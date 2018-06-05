package main

import "github.com/gorilla/mux"

func (app *App) routes() {
  app.Router = mux.NewRouter()
  app.Router.HandleFunc("/users/all", app.getUsers).Methods("GET")
  app.Router.HandleFunc("/users", app.createUser).Methods("POST")
  app.Router.HandleFunc("/authenticate", app.createTokenEndpoint).Methods("POST")
}
