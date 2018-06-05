package main

import (
  "encoding/json"
  "net/http"

  "golang.org/x/crypto/bcrypt"
)

func responseJsonErr(w http.ResponseWriter, code int, message string) {
  responseJson(w, code, map[string]string{"error": message})
}

func responseJson(w http.ResponseWriter, code int, payload interface{}) {
  response, _ := json.Marshal(payload)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}

func HashPassword(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  return err == nil
}
