package helpers

import (
  "encoding/json"
  "net/http"
  "golang.org/x/crypto/bcrypt"
)

type Exception struct {
  Message interface{} `json:"message"`
}

func ResponseException(e interface{}) map[string]Exception {
  return map[string]Exception{"error": Exception{Message: e}};
}

func HttpResponse(w http.ResponseWriter, code int, res interface{}) {
  response, _ := json.Marshal(res)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}

func HashPassword(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}
