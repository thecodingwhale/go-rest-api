package main

import (
  "encoding/json"
  "net/http"
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
