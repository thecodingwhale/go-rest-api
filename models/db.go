package models

import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"

)

type DB struct {
    *sql.DB
}

func NewDB(user, password, name string) (*DB, error) {
    connectionString := fmt.Sprintf("%s:%s@/%s?parseTime=true", user, password, name)
    db, err := sql.Open("mysql", connectionString)
    if err != nil {
        return nil, err
    }
    if err = db.Ping(); err != nil {
        return nil, err
    }
    return &DB{db}, nil
}
