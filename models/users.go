package models

import (
  "log"
  "errors"
  "time"
  "go-rest-api/helpers"

  "github.com/go-ozzo/ozzo-validation"
  "github.com/go-ozzo/ozzo-validation/is"
)

type UserDatastore interface {
    Create(User) error
    ReadByEmail(email string) (User, error)
    IsEmailExists(User) error
}

type User struct {
    Id          int       `json:"id"`
    Email       string    `json:"email"`
    Name        string    `json:"name"`
    Password    string    `json:"password"`
    CreatedDate time.Time `json:"created_date,string"`
    UpdatedDate time.Time `json:"updated_date,string"`
}

func (u User) Validate() error {
  return validation.ValidateStruct(&u,
    validation.Field(&u.Email, validation.Required, is.Email),
    validation.Field(&u.Name, validation.Required, validation.Length(8, 50)),
    validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
  )
}

func (db *DB) IsEmailExists(u User) error {
  rows, err := db.Query("SELECT email FROM users WHERE email = ?", u.Email)
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer rows.Close()

  for rows.Next() {
    if err := rows.Scan(&u.Email); err != nil {
      return err
    }
    return errors.New("email already exists")
  }
  err = rows.Err()
  if err != nil {
    log.Fatal(err)
    return err
  }
  return nil
}


func (db *DB) Create(u User) error {
  query := `
    INSERT INTO users (email, name, password)
    VALUES (?, ?, ?)
  `
  var err error
  hashPassword, _ := helpers.HashPassword(u.Password)
  _, err = db.Exec(query, u.Email, u.Name, hashPassword)
  if err != nil {
    return err
  }

  err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.Id)

  if err != nil {
    return err
  }

  return nil
}

func (db *DB) ReadByEmail(email string) (User, error) {
  var u User
  rows, err := db.Query("SELECT id, email, name FROM users WHERE email = ?", email)
  if err != nil {
    log.Fatal(err)
    return User{}, err
  }
  defer rows.Close()

  for rows.Next() {
    rows.Scan(&u.Id, &u.Email, &u.Name)
    return User{
      Id: u.Id,
      Email: u.Email,
      Name: u.Name,
    }, nil
  }
  return User{}, err

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
    return User{}, err
  }
  return User{}, nil
}
