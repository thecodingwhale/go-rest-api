package main

import (
  "fmt"
  "log"
  "database/sql"
  "net/http"

  _ "github.com/go-sql-driver/mysql"
  "github.com/gorilla/mux"
)

type App struct {
  Router *mux.Router
  DB     *sql.DB
}

func (app *App) database(driver, user, password, name string) {
  connectionString := fmt.Sprintf("%s:%s@/%s", user, password, name)
  var err error
  app.DB, err = sql.Open(driver, connectionString)
  if err != nil {
    log.Fatal(err)
  }
}

func (app *App) routes() {
  app.Router = mux.NewRouter()
  app.Router.HandleFunc("/users/all", app.getUsers).Methods("GET")
  app.Router.HandleFunc("/users", app.createUser).Methods("POST")
}

func (app *App) run(port string) {
  log.Fatal(http.ListenAndServe(port, app.Router))
}
