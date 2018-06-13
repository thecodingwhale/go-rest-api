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
