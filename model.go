package main

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type User struct {
  ID        int       `json:"id"`
  Email     string    `json:"email"`
  Password  string    `json:"password"`
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
    if err := rows.Scan(&u.ID, &u.Email, &u.Password); err != nil {
      return nil, err
    }
    users = append(users, u)
  }
  return users, nil
}

func (u *User) createUser(db *sql.DB) error {
  statement := fmt.Sprintf("INSERT INTO users(email, password) VALUES('%s', '%s')", u.Email, u.Password)
  _, err := db.Exec(statement)

  if err != nil {
    return err
  }

  err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

  if err != nil {
    return err
  }

  return nil
}
