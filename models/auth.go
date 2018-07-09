package models

import (
  "log"
  "os"
  "errors"
  "go-rest-api/helpers"

  "github.com/joho/godotenv"
  "github.com/dgrijalva/jwt-go"
  "github.com/go-ozzo/ozzo-validation"
  "github.com/go-ozzo/ozzo-validation/is"
)

type AuthDatastore interface {
    Validate(Auth) error
    IsUserExists(Auth) error
    CreateToken(User) (token string, err error)
}

type Auth struct {
    Email       string    `json:"email"`
    Password    string    `json:"password"`
}

func (db *DB) Validate(a Auth) error {
  return validation.ValidateStruct(&a,
    validation.Field(&a.Email, validation.Required, is.Email),
    validation.Field(&a.Password, validation.Required, validation.Length(8, 50)),
  )
}

func (db *DB) IsUserExists(a Auth) error {
  requestPassword := a.Password
  rows, err := db.Query("SELECT email, password FROM users WHERE email=?", a.Email)
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer rows.Close()

  for rows.Next() {
    rows.Scan(&a.Email, &a.Password)
    if helpers.CheckPasswordHash(requestPassword, a.Password) {
      return nil
    }
    return errors.New("Invalid Password.")
  }
  return errors.New("Credentials Not Found.")

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
    return err
  }
  return nil
}

func (db *DB) CreateToken(u User) (token string, tokenError error) {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
    return "", errors.New("Error loading .env file")
  }
  jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "id": u.Id,
    "email": u.Email,
    "name": u.Name,
  })
  token, jwtErr := jwtToken.SignedString([]byte(os.Getenv("APP_SECRET")))
  if jwtErr != nil {
    return "", errors.New("token error")
  }
  return token, nil
}

