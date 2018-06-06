package main

import (
  "fmt"
  "strings"
  "os"
  "log"
  "net/http"
  "encoding/json"

  "github.com/joho/godotenv"
  "github.com/gorilla/mux"
  "github.com/gorilla/context"
  "github.com/dgrijalva/jwt-go"
  "github.com/mitchellh/mapstructure"
)

type Exception struct {
  Message string `json:"message"`
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    authorizationHeader := req.Header.Get("authorization")
    if authorizationHeader != "" {
      bearerToken := strings.Split(authorizationHeader, " ")
      if len(bearerToken) == 2 {
        token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
          if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("There was an error")
          }
          return []byte(os.Getenv("APP_SECRET")), nil
        })
        if error != nil {
          json.NewEncoder(w).Encode(Exception{Message: error.Error()})
          return
        }
        if token.Valid {
          context.Set(req, "decoded", token.Claims)
          next(w, req)
        } else {
          json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
        }
      }
    } else {
      json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
    }
  })
}

func TestEndpoint(w http.ResponseWriter, req *http.Request) {
  decoded := context.Get(req, "decoded")
  var user User
  mapstructure.Decode(decoded.(jwt.MapClaims), &user)
  json.NewEncoder(w).Encode(user)
}

func (app *App) routes() {
  app.Router = mux.NewRouter()
  app.Router.HandleFunc("/users/all", app.getUsers).Methods("GET")

  app.Router.HandleFunc("/users", app.createUser).Methods("POST")
  app.Router.HandleFunc("/authenticate", app.createTokenEndpoint).Methods("POST")

  app.Router.HandleFunc("/test", ValidateMiddleware(TestEndpoint)).Methods("GET")

  app.Router.HandleFunc("/jobs/{id:[0-9]+}", app.getJob).Methods("GET")
  app.Router.HandleFunc("/jobs", ValidateMiddleware(app.createJob)).Methods("POST")
}
