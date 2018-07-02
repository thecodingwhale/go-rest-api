package main_test

import (
	"bytes"
	"encoding/json"
	helpers "go-rest-api"
	main "go-rest-api"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/icrowley/fake"
	"github.com/joho/godotenv"
)

var a main.App

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	a = main.App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	ensureTableExists()

	a.Routes()
	code := m.Run()

	clearTable()
	os.Exit(code)
}

func clearTable() {
	a.DB.Exec("TRUNCATE TABLE jobs")
	a.DB.Exec("TRUNCATE TABLE users")
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableJobsCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(tableUsersCreationQuery); err != nil {
		log.Fatal(err)
	}
}

const tableJobsCreationQuery = `CREATE TABLE IF NOT EXISTS jobs (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  user_id int(11) NOT NULL,
  post varchar(255) DEFAULT '',
  location varchar(255) DEFAULT NULL,
  company varchar(255) DEFAULT NULL,
  created_date datetime DEFAULT CURRENT_TIMESTAMP,
  updated_date datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;`

const tableUsersCreationQuery = `CREATE TABLE IF NOT EXISTS users (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  email varchar(255) DEFAULT NULL,
  name varchar(255) DEFAULT NULL,
  password varchar(255) DEFAULT NULL,
  created_date datetime DEFAULT CURRENT_TIMESTAMP,
  updated_date datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func createUser() map[string]interface{} {
	query := `INSERT INTO users (email, name, password) VALUES (?, ?, ?)`
	email := strings.ToLower(fake.EmailAddress())
	name := fake.FullName()
	password, _ := helpers.HashPassword("password")

	res, err := a.DB.Exec(query, email, name, password)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return map[string]interface{}{
		"id":       id,
		"email":    email,
		"name":     name,
		"password": password,
	}
}

func createUserEmail(e string) map[string]interface{} {
	query := `INSERT INTO users (email, name, password) VALUES (?, ?, ?)`
	email := e
	name := fake.FullName()
	password, _ := helpers.HashPassword("password")

	res, err := a.DB.Exec(query, email, name, password)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return map[string]interface{}{
		"id":       id,
		"email":    email,
		"name":     name,
		"password": password,
	}
}

func createJob(userId interface{}) map[string]interface{} {
	query := `INSERT INTO jobs (user_id, post, location, company) VALUES (?, ?, ?, ?)`
	post := fake.JobTitle()
	location := fake.State()
	company := fake.Company()

	res, err := a.DB.Exec(query, userId, post, location, company)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return map[string]interface{}{
		"id":       id,
		"post":     post,
		"location": location,
		"company":  company,
	}

}

func createToken(t *testing.T, inputEmail string) (string, map[string]interface{}) {
	password := "password"
	user := createUserEmail(inputEmail)

	payload := []byte(`{"email":"` + inputEmail + `","password":"` + password + `"}`)
	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	return m["token"], user
}

func TestInvalidTypeParameter(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/XXXXXX", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Not Found." {
		t.Errorf("Expected the 'error' key of the response to be set to 'Not Found.'. Got '%s'", m["error"])
	}
}
