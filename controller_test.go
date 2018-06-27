package main_test

import (
  "log"
  "testing"
  "encoding/json"
  "net/http"
  "bytes"
)

func TestEmptyJobs(t *testing.T) {
  clearTable()
  req, _ := http.NewRequest("GET", "/jobs", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)
  if body := response.Body.String(); body != "[]" {
    t.Errorf("Expected an empty array. Got %s", body)
  }
}

func TestGetJobsWithUrlParameters(t *testing.T) {
  clearTable()
  numberOfUsers := 13
  user := createUser()
  for i := 0; i < numberOfUsers; i++ {
    createJob(user["id"])
  }

  req, _ := http.NewRequest("GET", "/jobs", nil)

  q := req.URL.Query()
  q.Add("limit", "13")
  q.Add("offset", "10")
  req.URL.RawQuery = q.Encode()

  req.URL.RawQuery = q.Encode()

  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)

  var settings []map[string]interface{}
  if err := json.NewDecoder(response.Body).Decode(&settings); err != nil {
    log.Fatal(err)
  }

  expectedLength := 3
  if len(settings) != expectedLength {
    t.Errorf("Expected number of response should be '%d'. Got '%d'", expectedLength, len(settings))
  }
}

func TestGetJobs(t *testing.T) {
  clearTable()
  numberOfUsers := 13
  user := createUser()
  for i := 0; i < numberOfUsers; i++ {
    createJob(user["id"])
  }

  req, _ := http.NewRequest("GET", "/jobs", nil)

  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)

  var settings []map[string]interface{}
  if err := json.NewDecoder(response.Body).Decode(&settings); err != nil {
    log.Fatal(err)
  }

  expectedLength := 10
  if len(settings) != expectedLength {
    t.Errorf("Expected number of response should be '%d'. Got '%d'", expectedLength, len(settings))
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
  user := createUser()
  createJob(user["id"])

  req, _ := http.NewRequest("GET", "/jobs/1", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCreateUserSuccesfully(t *testing.T) {
  clearTable();
  payload := []byte(`{"email":"foo@email.com","name":"foobarbaz","password":"12345678"}`)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
  response := executeRequest(req)
  checkResponseCode(t, http.StatusCreated, response.Code)

  var m map[string]string
  json.Unmarshal(response.Body.Bytes(), &m)

  if len(m) != 0 {
    t.Errorf("Expected reponse should return 0 or empty. Got '%d'", len(m))
  }
}

func TestCreateUserInvalidRequestPayloadFailed(t *testing.T) {
  clearTable();
  body := ``
  payload := []byte(body)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
  response := executeRequest(req)
  checkResponseCode(t, http.StatusNotFound, response.Code)

  var m map[string]string
  json.Unmarshal(response.Body.Bytes(), &m)
  if m["error"] != "Invalid request payload" {
    t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
  }
}


func TestCreateUserEmptyBodyFailed(t *testing.T) {
  clearTable();
  payload := []byte(`{}`)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
  response := executeRequest(req)

  checkResponseCode(t, http.StatusNotFound, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  email := m["error"].(map[string]interface{})["email"]
  password := m["error"].(map[string]interface{})["password"]
  name := m["error"].(map[string]interface{})["name"]

  if email != "cannot be blank" {
    t.Errorf("expected error '%s'", "cannot be blank")
  }
  if password != "cannot be blank" {
    t.Errorf("expected error '%s'", "cannot be blank")
  }
  if name != "cannot be blank" {
    t.Errorf("expected error '%s'", "cannot be blank")
  }
}

func TestCreateUserEmailAlreadyExistsFailed(t *testing.T) {
  clearTable();
  inputEmail := "foo@email.com"
  createUserEmail(inputEmail)
  payload := []byte(`{"email":"` + inputEmail + `","name":"foobarbaz","password":"12345678"}`)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))
  response := executeRequest(req)
  checkResponseCode(t, http.StatusNotFound, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)
  email := m["error"].(map[string]interface{})["email"]

  if email != "email already exists" {
    t.Errorf("Expected the 'error' key of the response to be set to 'email already exists'. Got '%s'", email)
  }
}
