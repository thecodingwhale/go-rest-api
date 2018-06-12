package main_test

import (
  main "go-rest-api"
  "os"
  "log"
  "testing"
  "encoding/json"
  "net/http"
  "net/http/httptest"

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
  a.Routes()

  ensureTableExists()
  code := m.Run()
  clearTable()
  os.Exit(code)
}

func clearTable() {
  a.DB.Exec("DELETE FROM jobs")
  a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func ensureTableExists() {
  if _, err := a.DB.Exec(tableCreationQuery); err != nil {
    log.Fatal(err)
  }
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS jobs (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  user_id int(11) NOT NULL,
  post varchar(255) DEFAULT '',
  location varchar(255) DEFAULT NULL,
  company varchar(255) DEFAULT NULL,
  created_date datetime DEFAULT CURRENT_TIMESTAMP,
  updated_date datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;`

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

func TestEmptyJobs(t *testing.T) {
  clearTable()
  req, _ := http.NewRequest("GET", "/jobs", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)
  if body := response.Body.String(); body != "[]" {
    t.Errorf("Expected an empty array. Got %s", body)
  }
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
