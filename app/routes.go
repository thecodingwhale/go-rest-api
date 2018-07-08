package app

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

func (app *App) Routes() {
    app.Router = mux.NewRouter()
    app.Router.HandleFunc("/users", app.UsersCreate).Methods("POST")
    app.Router.HandleFunc("/authenticate", app.Authenticate).Methods("POST")
    log.Fatal(http.ListenAndServe("3000", app.Router))
}
