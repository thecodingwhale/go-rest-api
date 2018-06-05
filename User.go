package main

import (
  "os"
  "fmt"
  "database/sql"
  "log"
  "errors"
  // "reflect"

  _ "github.com/go-sql-driver/mysql"
  "github.com/go-ozzo/ozzo-validation"
  "github.com/go-ozzo/ozzo-validation/is"
  "github.com/joho/godotenv"
  "github.com/dgrijalva/jwt-go"
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
  // tuple - return multiple values
  rows, err := db.Query("SELECT email FROM users WHERE email = ?", u.Email)
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer rows.Close()

  for rows.Next() {
    // err := rows.Scan(&u.Email)
    if err := rows.Scan(&u.Email); err != nil {
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

func (u User) getToken(db *sql.DB) string {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "id": u.Id,
    "email": u.Email,
  })
  tokenString, err := token.SignedString([]byte(os.Getenv("APP_SECRET")))
  if err != nil {
    errors.New("token error")
  }
  return tokenString
}

func (u User) isUserExists(db *sql.DB) error {
  requestPassword := u.Password
  rows, err := db.Query("SELECT email, password FROM users WHERE email=?", u.Email)
  if err != nil {
    log.Fatal(err)
    return err
  }
  defer rows.Close()

  for rows.Next() {
    rows.Scan(&u.Email, &u.Password);
    if (CheckPasswordHash(requestPassword, u.Password)) {
      return nil
    }
    return errors.New("invalid password")
  }
  return errors.New("credentials not found")

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
  hashPassword, _ := HashPassword(u.Password)
  _, err = db.Exec(query, u.Email, hashPassword)
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

