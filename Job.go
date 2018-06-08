package main

import (
  "fmt"
  "log"
  "database/sql"
  "errors"

  _ "github.com/go-sql-driver/mysql"
  "github.com/go-ozzo/ozzo-validation"
  // "github.com/go-ozzo/ozzo-validation/is"
)

type Job struct {
  Id int `json:"id"`
  UserId int `json:"user_id"`
  Post string `json:"post"`
  Location string `json:"location"`
  Company string `json:"company"`
}

func (j Job) validate() error {
  return validation.ValidateStruct(&j,
    validation.Field(&j.Post, validation.Required),
    validation.Field(&j.Location, validation.Required),
    validation.Field(&j.Company, validation.Required),
  )
}

func (j *Job) getJob(db *sql.DB) (map[string]interface{}, error) {
  rows, err := db.Query(`
    SELECT
      jobs.id,
      jobs.post,
      jobs.location,
      jobs.company,
      users.name
    FROM
      jobs
    JOIN
      users
      ON jobs.user_id = users.id
    WHERE
      jobs.id = ?
  `, j.Id)
  if err != nil {
    log.Fatal(err)
    return nil, err
  }
  defer rows.Close()
  var u User
  for rows.Next() {
    rows.Scan(&j.Id, &j.Post, &j.Location, &j.Company, &u.Name)
    return map[string]interface{}{
      "id" : j.Id,
      "post": j.Post,
      "location": j.Location,
      "company": j.Company,
      "name": u.Name,
    }, nil
  }
  return nil, errors.New("No job post found.")

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
    return map[string]interface{}{}, err;
  }
  return map[string]interface{}{}, nil
}

func (j *Job) createJob(db *sql.DB, userId int) error {
  fmt.Println(userId)

  query := `
    INSERT INTO jobs (user_id, post, location, company)
    VALUES (?, ?, ?, ?)
  `
  var err error
  _, err = db.Exec(query, userId, j.Post, j.Location, j.Company)
  if err != nil {
    return err
  }

  err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&j.Id)

  if err != nil {
    return err
  }

  return nil
}

func (j *Job) getJobs(db *sql.DB, offset, limit int) ([]map[string]interface{}, error) {
  rows, err := db.Query(`
    SELECT
      jobs.id,
      jobs.post,
      jobs.location,
      jobs.company,
      users.name
    FROM
      jobs
    JOIN
      users
      ON jobs.user_id = users.id
    LIMIT ?
    OFFSET ?
  `, limit, offset)
  if err != nil {
    return nil, err
  }
  jobs := []Job{}
  allJobs := make([]map[string]interface{}, len(jobs))
  for rows.Next() {
    var j Job
    var u User
    if err := rows.Scan(&j.Id, &j.Post, &j.Location, &j.Company, &u.Name); err != nil {
      return nil, err
    }
    allJobs = append(allJobs, map[string]interface{}{
      "id" : j.Id,
      "post": j.Post,
      "location": j.Location,
      "company": j.Company,
      "name": u.Name,
    })
  }
  return allJobs, nil
}

func (j *Job) updateJob(db *sql.DB, userId int, jobId int) (map[string]interface{}, error) {
  stmt, err := db.Prepare(`
    UPDATE
      jobs
    SET
      post = ?,
      location = ?,
      company = ?
    WHERE
      id = ?
    AND
      user_id = ?
  `)
  res, err := stmt.Exec(j.Post, j.Location, j.Company, jobId, userId)
  if err != nil {
    log.Fatal(err)
    return nil, err
  }

  count, err := res.RowsAffected()
  if err != nil {
    log.Fatal(err)
  }

  if count != 0 {
    return map[string]interface{}{
      "post": j.Post,
      "location": j.Location,
      "company": j.Company,
    }, nil
  }

  return nil, errors.New("No job post found.")
}
