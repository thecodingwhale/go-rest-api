package main

import (
  "fmt"
  "database/sql"
  "log"
  "errors"
  // "reflect"

  _ "github.com/go-sql-driver/mysql"
  "github.com/go-ozzo/ozzo-validation"
  "github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
  Id        int       `json:"id"`
  Email     string    `json:"email"`
  Password  string    `json:"password"`
}

func (u User) validate() error {
  return validation.ValidateStruct(&u,
    validation.Field(&u.Email, validation.Required, is.Email),
    validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
  )
}

func (u User) isEmailExists(db *sql.DB) error {
  rows, err := db.Query("SELECT email FROM users WHERE email = ?", u.Email)
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer rows.Close()
  for rows.Next() {
    err := rows.Scan(&u.Email)
    if err != nil {
      log.Fatal(err)
      return err
    }
    return errors.New("email already exists")
  }
  err = rows.Err()
  if err != nil {
    log.Fatal(err)
    return err;
  }
  return nil
}

func getUsers(db *sql.DB) ([]User, error) {
  statement := fmt.Sprintf("SELECT id, email, password FROM users")
  rows, err := db.Query(statement)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  users := []User{}
  for rows.Next() {
    var u User
    if err := rows.Scan(&u.Id, &u.Email, &u.Password); err != nil {
      return nil, err
    }
    users = append(users, u)
  }
  return users, nil
}

func (u *User) createUser(db *sql.DB) error {
  query := `
    INSERT INTO users (email, password)
    VALUES (?, ?)
  `
  var err error
  _, err = db.Exec(query, u.Email, u.Password)
  if err != nil {
    return err
  }

  err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.Id)

  if err != nil {
    return err
  }

  return nil
}

func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}

