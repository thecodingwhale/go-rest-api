package app

import (
    "log"
    "net/http"
    "go-rest-api/models"

    "github.com/gorilla/mux"
)

type App struct {
    DB models.Datastore
    Router *mux.Router
}

func (app *App) Routes() {
  app.Router = mux.NewRouter()
  app.Router.HandleFunc("/books", app.BooksIndex).Methods("GET")
  log.Fatal(http.ListenAndServe("3000", app.Router))
}
