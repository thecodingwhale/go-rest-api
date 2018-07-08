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

    "github.com/go-ozzo/ozzo-validation"
    "github.com/go-ozzo/ozzo-validation/is"
)

type authenticateMockDB struct{}


func (mdb *authenticateMockDB) Validate(u models.Auth) error {
  return validation.ValidateStruct(&a,
    validation.Field(&a.Email, validation.Required, is.Email),
    validation.Field(&a.Password, validation.Required, validation.Length(8, 50)),
  )
}

func (mdb *authenticateMockDB) IsUserExists(a models.Auth) error {
  return errors.New("error")
}

func (mdb *authenticateMockDB) CreateToken(u models.User) (token string, err error) {
  return "asdadasd", nil
}

func TestAuthenticateHandlerToReturnMethodNotAllowed(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(``)
  req, _ := http.NewRequest("GET", "/authenticate", bytes.NewBuffer(payload))

  var a = App{Auth: &authenticateMockDB{}}
  http.HandlerFunc(a.Authenticate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusNotFound, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  actualErrorMessage := m["error"].(map[string]interface{})["message"]
  expectedErrorMessage := "Method Not Allowed."

  if actualErrorMessage != expectedErrorMessage {
    t.Errorf(`Expected reponse {"error": {"message": "`+expectedErrorMessage+`"}}. Got '%s'`, actualErrorMessage)
  }
}

func TestAuthenticateHandlerToReturnInvalidRequestPayload(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(``)
  req, _ := http.NewRequest("POST", "/authenticate", bytes.NewBuffer(payload))

  var a = App{Auth: &authenticateMockDB{}}
  http.HandlerFunc(a.Authenticate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusNotFound, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  actualErrorMessage := m["error"].(map[string]interface{})["message"]
  expectedErrorMessage := "Invalid request payload."

  if actualErrorMessage != expectedErrorMessage {
    t.Errorf(`Expected reponse {"error": {"message": "`+expectedErrorMessage+`"}}. Got '%s'`, actualErrorMessage)
  }
}

func TestAuthenticateHandlerToReturnInputRequestValidationForMissingPayload(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(`{}`)
  req, _ := http.NewRequest("POST", "/authenticate", bytes.NewBuffer(payload))

  var a = App{Auth: &authenticateMockDB{}}

  http.HandlerFunc(a.Authenticate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusNotFound, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  errorMessage := m["error"].(map[string]interface{})["message"]
  emailErrorMessage := errorMessage.(map[string]interface{})["email"]
  passwordErrorMessage := errorMessage.(map[string]interface{})["password"]
  expectedErrorMessage := "cannot be blank"

  if emailErrorMessage != expectedErrorMessage {
    t.Errorf("expected error '%s'", expectedErrorMessage)
  }
  if passwordErrorMessage != expectedErrorMessage {
    t.Errorf("expected error '%s'", expectedErrorMessage)
  }
}

func TestAuthenticateHandlerToReturnErrorIfUserDoesNotExists(t *testing.T) {
  rec := httptest.NewRecorder()
  payload := []byte(`{"email":"foo@email.com","password":"12345678"}`)
  req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(payload))

  var a = App{
    Auth: &authenticateMockDB{},
  }
  http.HandlerFunc(a.Authenticate).ServeHTTP(rec, req)

  CheckResponseCode(t, http.StatusNotFound, rec.Code)

  var m map[string]interface{}
  json.Unmarshal(rec.Body.Bytes(), &m)

  errorMessage := m["error"].(map[string]interface{})["message"]
  expectedErrorMessage := "error"

  if errorMessage != expectedErrorMessage {
    t.Errorf("expected error '%s'", expectedErrorMessage)
  }
}
