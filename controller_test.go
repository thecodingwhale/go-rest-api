package main_test

import (
  "log"
  "testing"
  "encoding/json"
  "net/http"

  // "github.com/bitly/go-simplejson"
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
  for i := 0; i < numberOfUsers; i++ {
    user := createUser()
    createJob(user["id"])
  }

  req, _ := http.NewRequest("GET", "/jobs", nil)

  q := req.URL.Query()
  q.Add("limit", "13")
  q.Add("offset", "10")
  req.URL.RawQuery = q.Encode()

  req.URL.RawQuery = q.Encode()
  log.Println(req.URL.String())

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
  for i := 0; i < numberOfUsers; i++ {
    user := createUser()
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

/*
func TestPublishOK(t *testing.T) {
  msg := "Test message"
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    if r.Method != "POST" {
      t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
    }
    if r.URL.EscapedPath() != "/pub" {
      t.Errorf("Expected request to ‘/pub’, got ‘%s’", r.URL.EscapedPath())
    }
    r.ParseForm()
    topic := r.Form.Get("topic")
    if topic != "meaningful-topic" {
      t.Errorf("Expected request to have ‘topic=meaningful-topic’, got: ‘%s’", topic)
    }
    reqJson, err := simplejson.NewFromReader(r.Body)
    if err != nil {
      t.Errorf("Error while reading request JSON: %s", err)
    }
    lifeMeaning := reqJson.GetPath("meta", "lifeMeaning").MustInt()
    if lifeMeaning != 42 {
      t.Errorf("Expected request JSON to have meta/lifeMeaning = 42, got %d", lifeMeaning)
    }
    msgActual := reqJson.GetPath("data", "message").MustString()
    if msgActual != msg {
      t.Errorf("Expected request JSON to have data/message = ‘%s’, got ‘%s’", msg, msgActual)
    }
  }))
  defer ts.Close()
  nsqdUrl := ts.URL
  err := Publish(nsqdUrl, msg)
  if err != nil {
    t.Errorf("Publish() returned an error: %s", err)
  }
}

*/
