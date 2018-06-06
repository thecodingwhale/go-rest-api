package main

import (
  "fmt"
  "database/sql"

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
