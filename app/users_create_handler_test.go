package app

import (
  // "log"
  "go-rest-api/models"
  "encoding/json"
  "net/http"
  "net/http/httptest"
  "bytes"
  "testing"
  "errors"
)

type userCreateMockDB struct{}

func (mdb *userCreateMockDB) Create(u models.User) error {
  if (u.Email == somethingWentWrongEmail) {
    return errors.New("Something went wrong.")
  }
  return nil
}

func (mdb *userCreateMockDB) IsEmailExists(u models.User) error {
  if (u.Email == baseEmail) {
    return errors.New("Email Already Exists.")
  }
  return nil
}

func CheckResponseCode(t *testing.T, expected, actual int) {
  if expected != actual {
    t.Errorf("Expected response code %d. Got %d\n", expected, actual)
  }
}

func TestUsersHandlerToReturnMethodNotAllowed(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(``)
  req, _ := http.NewRequest("GET", "/users", bytes.NewBuffer(payload))

  var a = App{User: &userCreateMockDB{}}
  http.HandlerFunc(a.UsersCreate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusNotFound, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  actualErrorMessage := m["error"].(map[string]interface{})["message"]
  expectedErrorMessage := "Method Not Allowed."

  if actualErrorMessage != expectedErrorMessage {
    t.Errorf(`Expected reponse {"error": {"message": "`+expectedErrorMessage+`"}}. Got '%s'`, actualErrorMessage)
  }
}

func TestUsersHandlerToReturnInvalidRequestPayload(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(``)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))

  var a = App{User: &userCreateMockDB{}}
  http.HandlerFunc(a.UsersCreate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusNotFound, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  actualErrorMessage := m["error"].(map[string]interface{})["message"]
  expectedErrorMessage := "Invalid request payload."

  if actualErrorMessage != expectedErrorMessage {
    t.Errorf(`Expected reponse {"error": {"message": "`+expectedErrorMessage+`"}}. Got '%s'`, actualErrorMessage)
  }
}

func TestUsersHandlerToReturnInputRequestValidationForMissingPayload(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(`{}`)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))

  var a = App{User: &userCreateMockDB{}}
  http.HandlerFunc(a.UsersCreate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusNotFound, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  errorMessage := m["error"].(map[string]interface{})["message"]
  emailErrorMessage := errorMessage.(map[string]interface{})["email"]
  nameErrorMessage := errorMessage.(map[string]interface{})["name"]
  passwordErrorMessage := errorMessage.(map[string]interface{})["password"]
  expectedErrorMessage := "cannot be blank"

  if emailErrorMessage != expectedErrorMessage {
    t.Errorf("expected error '%s'", expectedErrorMessage)
  }
  if nameErrorMessage != expectedErrorMessage {
    t.Errorf("expected error '%s'", expectedErrorMessage)
  }
  if passwordErrorMessage != expectedErrorMessage {
    t.Errorf("expected error '%s'", expectedErrorMessage)
  }
}

func TestCreateUserHandlerToReturnUserEmailAlreadyExists(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(`{"email":"` + baseEmail + `","name":"foobarbaz","password":"12345678"}`)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))

  var a = App{User: &userCreateMockDB{}}
  http.HandlerFunc(a.UsersCreate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusNotFound, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  errorMessage := m["error"].(map[string]interface{})["message"]
  expectedErrorMessage := "Email already exists."

  if errorMessage != expectedErrorMessage {
    t.Errorf("expected error '%s'", expectedErrorMessage)
  }
}

func TestCreateUserHandlerToReturnSomethingWentWrong(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(`{"email":"` + somethingWentWrongEmail + `","name":"foobarbaz","password":"12345678"}`)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))

  var a = App{User: &userCreateMockDB{}}
  http.HandlerFunc(a.UsersCreate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusNotFound, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  errorMessage := m["error"].(map[string]interface{})["message"]
  expectedErrorMessage := "Something went wrong."

  if errorMessage != expectedErrorMessage {
    t.Errorf("expected error '%s'", expectedErrorMessage)
  }
}

func TestCreateUserHandlerToReturnSuccess(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(`{"email":"foo@email.com","name":"foobarbaz","password":"12345678"}`)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))

  var a = App{User: &userCreateMockDB{}}
  http.HandlerFunc(a.UsersCreate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusCreated, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  if len(m) != 0 {
    t.Errorf("Expected reponse should return 0 or empty. Got '%d'", len(m))
  }
}
