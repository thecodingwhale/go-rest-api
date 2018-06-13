package main_test

import (
  "testing"
  "encoding/json"
  "net/http"
)

func TestEmptyJobs(t *testing.T) {
  clearTable()
  req, _ := http.NewRequest("GET", "/jobs", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)
  if body := response.Body.String(); body != "[]" {
    t.Errorf("Expescted an empty array. Got %s", body)
  }
}

func TestGetNonExistentJob(t *testing.T) {
  clearTable()

  req, _ := http.NewRequest("GET", "/jobs/1", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusNotFound, response.Code)

  var m map[string]string
  json.Unmarshal(response.Body.Bytes(), &m)
  if m["error"] != "No job post found." {
    t.Errorf("Expected the 'error' key of the response to be set to 'No job post found.'. Got '%s'", m["error"])
  }
}

func TestGetJob(t *testing.T) {
  clearTable()
  insertUserQuery := `INSERT INTO users (email, name, password)VALUES (?, ?, ?)`
  a.DB.Exec(insertUserQuery, "foo@email.com", "username", "password")

  insertJobQuery := `INSERT INTO jobs (user_id, post, location, company) VALUES (?, ?, ?, ?)`
  a.DB.Exec(insertJobQuery, "1", "post", "location", "company")

  req, _ := http.NewRequest("GET", "/jobs/1", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)
}
