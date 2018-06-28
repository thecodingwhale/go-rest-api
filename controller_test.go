package main_test

import (
  "log"
  "testing"
  // "os"
  "bytes"
  "encoding/json"
  "net/http"

  // "github.com/joho/godotenv"
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

func TestUpdateUserEmptyBodyFailed(t *testing.T) {
  clearTable();
  token, _ := createToken(t, "foo@email.com");
  payload := []byte(`{}`)
  req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusBadRequest, response.Code)

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

func TestUpdateUserInvalidRequestPayloadFailed(t *testing.T) {
  clearTable();
  inputEmail := "foo@email.com"
  token, _ := createToken(t, inputEmail)
  payload := []byte(``)
  req, _ := http.NewRequest("PUT", "/users/2", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusBadRequest, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)
  if m["error"] != "Invalid request payload" {
    t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
  }
}

func TestUpdateUserNotFoundFailed(t *testing.T) {
  clearTable();
  token, _ := createToken(t, "foo@email.com");
  payload := []byte(``)
  req, _ := http.NewRequest("PUT", "/users/XX", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusNotFound, response.Code)

  var m map[string]string
  json.Unmarshal(response.Body.Bytes(), &m)
  if m["error"] != "Not Found." {
    t.Errorf("Expected the 'error' key of the response to be set to 'Not Found.'. Got '%s'", m["error"])
  }
}

func TestUpdateUserSuccesfully(t *testing.T) {
  clearTable();
  token, _ := createToken(t, "foo@email.com");
  payload := []byte(`{"email":"foo@email.com","name":"foobarbaz","password":"12345678"}`)
  req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusCreated, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  if len(m) != 0 {
    t.Errorf("Expected reponse should return 0 or empty. Got '%d'", len(m))
  }
}

func TestCreateJobSuccesfully(t *testing.T) {
  clearTable();
  token, _ := createToken(t, "foo@email.com");
  payload := []byte(`{"company": "sample company", "location": "sample location", "post": "sample post"}`)
  req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusCreated, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  if len(m) != 0 {
    t.Errorf("Expected reponse should return 0 or empty. Got '%d'", len(m))
  }
}

func TestCreateJobEmptyBodyFailed(t *testing.T) {
  clearTable();
  token, _ := createToken(t, "foo@email.com");
  payload := []byte(`{}`)
  req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusBadRequest, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  company := m["error"].(map[string]interface{})["company"]
  location := m["error"].(map[string]interface{})["location"]
  post := m["error"].(map[string]interface{})["post"]

  if company != "cannot be blank" {
    t.Errorf("expected error '%s'", "cannot be blank")
  }
  if location != "cannot be blank" {
    t.Errorf("expected error '%s'", "cannot be blank")
  }
  if post != "cannot be blank" {
    t.Errorf("expected error '%s'", "cannot be blank")
  }
}

func TestCreateJobInvalidRequestPayloadFailed(t *testing.T) {
  clearTable();
  token, _ := createToken(t, "foo@email.com");
  payload := []byte(``)
  req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusBadRequest, response.Code)

  var m map[string]string
  json.Unmarshal(response.Body.Bytes(), &m)
  if m["error"] != "Invalid request payload" {
    t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
  }
}

func TestDeleteJobSuccesfully(t *testing.T) {
  clearTable();
  token, user := createToken(t, "foo@email.com")
  createJob(user["id"])

  payload := []byte(`{}`)
  req, _ := http.NewRequest("DELETE", "/jobs/1", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusCreated, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  if len(m) != 0 {
    t.Errorf("Expected reponse should return 0 or empty. Got '%d'", len(m))
  }
}

func TestDeleteJobInvalidRequestPayloadFailed(t *testing.T) {
  clearTable();
  token, user := createToken(t, "foo@email.com")
  createJob(user["id"])

  payload := []byte(``)
  req, _ := http.NewRequest("DELETE", "/jobs", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusMethodNotAllowed, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  if len(m) != 0 {
    t.Errorf("Expected reponse should return 0 or empty. Got '%d'", len(m))
  }
}

func TestDeleteJobNotFoundFailed(t *testing.T) {
  clearTable();
  token, user := createToken(t, "foo@email.com")
  createJob(user["id"])

  payload := []byte(`{}`)
  req, _ := http.NewRequest("DELETE", "/jobs/XX", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusNotFound, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  if m["error"] != "Not Found." {
    t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
  }
}

func TestDeleteJobNoJobPostFoundFailed(t *testing.T) {
  clearTable();
  token, user := createToken(t, "foo@email.com")
  createJob(user["id"])

  payload := []byte(`{}`)
  req, _ := http.NewRequest("DELETE", "/jobs/2", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusBadRequest, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  if m["error"] != "No job post found." {
    t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
  }
}

func TestUpdateJobInvalidRequestPayloadFailed(t *testing.T) {
  clearTable();
  token, user := createToken(t, "foo@email.com")
  createJob(user["id"])

  payload := []byte(``)
  req, _ := http.NewRequest("PUT", "/jobs/1", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusBadRequest, response.Code)

  var m map[string]string
  json.Unmarshal(response.Body.Bytes(), &m)

  if m["error"] != "Invalid request payload" {
    t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
  }
}

func TestUpdateEmptyBodyFailed(t *testing.T) {
  clearTable();
  token, _ := createToken(t, "foo@email.com");
  payload := []byte(`{}`)
  req, _ := http.NewRequest("PUT", "/jobs/1", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusBadRequest, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  company := m["error"].(map[string]interface{})["company"]
  location := m["error"].(map[string]interface{})["location"]
  post := m["error"].(map[string]interface{})["post"]

  if company != "cannot be blank" {
    t.Errorf("expected error '%s'", "cannot be blank")
  }
  if location != "cannot be blank" {
    t.Errorf("expected error '%s'", "cannot be blank")
  }
  if post != "cannot be blank" {
    t.Errorf("expected error '%s'", "cannot be blank")
  }
}

func TestUpdateJobNoJobPostFoundFailed(t *testing.T) {
  clearTable();
  token, user := createToken(t, "foo@email.com")
  createJob(user["id"])

  payload := []byte(`{"company": "sample company", "location": "sample location", "post": "sample post"}`)
  req, _ := http.NewRequest("PUT", "/jobs/2", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusBadRequest, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  if m["error"] != "No job post found." {
    t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
  }
}


func TestUpdateJobNonRelatedTokenToJobPostIdFailed(t *testing.T) {
  clearTable();
  _, firstUser := createToken(t, "foo@email.com")
  secondToken, _ := createToken(t, "bar@email.com")

  createJob(firstUser["id"])

  payload := []byte(`{"company": "sample company", "location": "sample location", "post": "sample post"}`)
  req, _ := http.NewRequest("PUT", "/jobs/1", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + secondToken)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusBadRequest, response.Code)

  var m map[string]interface{}
  json.Unmarshal(response.Body.Bytes(), &m)

  if m["error"] != "No job post found." {
    t.Errorf("Expected the 'error' key of the response to be set to 'Invalid request payload'. Got '%s'", m["error"])
  }
}

func TestUpdateJobSucessfully(t *testing.T) {
  clearTable();
  token, user := createToken(t, "foo@email.com")
  createJob(user["id"])

  payload := []byte(`{"company": "sample company", "location": "sample location", "post": "sample post"}`)
  req, _ := http.NewRequest("PUT", "/jobs/1", bytes.NewBuffer(payload))
  req.Header.Set("Authorization", "Bearer: " + token)
  response := executeRequest(req)
  checkResponseCode(t, http.StatusCreated, response.Code)

  var m map[string]string
  json.Unmarshal(response.Body.Bytes(), &m)

  if len(m) != 0 {
    t.Errorf("Expected reponse should return 0 or empty. Got '%d'", len(m))
  }
}
