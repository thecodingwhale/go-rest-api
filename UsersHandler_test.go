package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateUserHandlerToReturnInvalidRequestPayload(t *testing.T) {
	clearTable()
	body := ``
	payload := []byte(body)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	actualErrorMessage := m["error"].(map[string]interface{})["message"]
	expectedErrorMessage := "Invalid request payload."

	if actualErrorMessage != expectedErrorMessage {
		t.Errorf(`Expected reponse {"error": {"message": "`+expectedErrorMessage+`"}}. Got '%s'`, actualErrorMessage)
	}
}

func TestCreateUserHandlerToReturnInputRequestValidationForMissingPayload(t *testing.T) {
	clearTable()
	body := `{}`
	payload := []byte(body)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

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
	clearTable()
	inputEmail := "foo@email.com"
	createUserEmail(inputEmail)
	payload := []byte(`{"email":"` + inputEmail + `","name":"foobarbaz","password":"12345678"}`)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	errorMessage := m["error"].(map[string]interface{})["message"]
	expectedErrorMessage := "Email already exists."

	if errorMessage != expectedErrorMessage {
		t.Errorf(`Expected the 'error' key of the response to be set to "`+expectedErrorMessage+`". Got '%s'`, errorMessage)
	}
}

func TestCreateUserHandlerToReturnSuccessful(t *testing.T) {
	clearTable()
	payload := []byte(`{"email":"foo@email.com","name":"foobarbaz","password":"12345678"}`)
	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	if len(m) != 0 {
		t.Errorf("Expected reponse should return 0 or empty. Got '%d'", len(m))
	}
}
