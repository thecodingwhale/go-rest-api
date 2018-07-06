package app

import (
    "go-rest-api/models"
    "github.com/gorilla/mux"
)

type App struct {
    Router *mux.Router
    User models.UserDatastore
}
