package main

import (
  "os"
  "fmt"
  "time"
  "log"
  "errors"
  "database/sql"

  _ "github.com/go-sql-driver/mysql"
  "github.com/go-ozzo/ozzo-validation"
  "github.com/go-ozzo/ozzo-validation/is"
  "github.com/joho/godotenv"
  "github.com/dgrijalva/jwt-go"
)

type User struct {
  Id int `json:"id"`
  Email string `json:"email"`
  Name string `json:"name"`
  Password string `json:"password"`
  CreatedDate time.Time `json:"created_date,string"`
  UpdatedDate time.Time `json:"updated_date,string"`
}

func (u User) validate() error {
  return validation.ValidateStruct(&u,
    validation.Field(&u.Email, validation.Required, is.Email),
    validation.Field(&u.Name, validation.Required, validation.Length(8, 50)),
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
    if err := rows.Scan(&u.Email); err != nil {
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
    "name": u.Name,
    "email": u.Email,
  })
  tokenString, err := token.SignedString([]byte(os.Getenv("APP_SECRET")))
  if err != nil {
    errors.New("token error")
  }
  return tokenString
}

func (u User) isUserExists(db *sql.DB) (tokenString string, err error) {
  requestPassword := u.Password
  rows, err := db.Query("SELECT id, email, name, password FROM users WHERE email=?", u.Email)
  if err != nil {
    log.Fatal(err)
    return "", err
  }
  defer rows.Close()

  for rows.Next() {
    rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password)
    if (CheckPasswordHash(requestPassword, u.Password)) {
      tokenString := u.getToken(db)
      return tokenString, nil
    }
    return "", errors.New("invalid password")
  }
  return "", errors.New("credentials not found")

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
    return "", err;
  }
  return "", nil
}

func getUsers(db *sql.DB) ([]User, error) {
  statement := fmt.Sprintf("SELECT id, email, name, password, created_date, updated_date FROM users")
  rows, err := db.Query(statement)
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  users := []User{}
  for rows.Next() {
    var u User
    if err := rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.CreatedDate, &u.UpdatedDate); err != nil {
      return nil, err
    }
    users = append(users, u)
  }
  return users, nil
}

func (u *User) createUser(db *sql.DB) error {
  query := `
    INSERT INTO users (email, name, password)
    VALUES (?, ?, ?)
  `
  var err error
  hashPassword, _ := HashPassword(u.Password)
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

func (u *User) updateUser(db *sql.DB, userId int) (map[string]interface{}, error) {
  stmt, err := db.Prepare(`
    UPDATE
      users
    SET
      email = ?,
      name = ?,
      password = ?,
      updated_date = NOW()
    WHERE
      id = ?
  `)
  hashPassword, _ := HashPassword(u.Password)
  res, err := stmt.Exec(u.Email, u.Name, hashPassword, userId)
  if err != nil {
    log.Fatal(err)
    return nil, err
  }

  count, err := res.RowsAffected()
  if err != nil {
    log.Fatal(err)
  }

  if count != 0 {
    return map[string]interface{}{}, nil
  }

  return nil, errors.New("No user found.")
}

func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}

