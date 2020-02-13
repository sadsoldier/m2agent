/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package pguserModel

import (
    "github.com/jmoiron/sqlx"
    "testing"
    "fmt"
    "strings"
    "errors"
)

func createDB() (*sqlx.DB, error) {
    db, err := sqlx.Open("pgx", "postgres://pgsql@localhost/postgres?sslmode=disable")
    if err != nil {
        return db, err
    }
    err = db.Ping()
    if err != nil {
        return db, err
    }
    return db, nil
}


func TestCreate(t *testing.T) {

    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)

    var user1 = PgUser{
        Username: "user1",
        Password: "12345",
        Superuser: false,
    }

    err = model.Create(user1)
    if err != nil {
        t.Error(err)
    }

    var user2 = PgUser{
        Username: "user2",
        Password: "12345",
        Superuser: true,
    }

    err = model.Create(user2)
    if err != nil {
        t.Error(err)
    }
}

func TestList(t *testing.T) {
    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)

    var users []PgUser
    page := Page{
        Limit: 50,
        PgUsers: &users,
    }
    err = model.List(&page)
    if err != nil {
        t.Error(err)
    }

    for num, item := range *page.PgUsers {
        if strings.HasPrefix(item.Username, "user") {
            fmt.Println(num, item)
        }
    }
}

func TestIsExist(t *testing.T) {
    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)

    var user = PgUser{
        Username: "user1",
    }
    res, err := model.IsExist(user)
    if err != nil {
        t.Error(err)
    }
    if res != true {
        t.Error(errors.New("user must be exist"))
    }
}


func TestUpdate(t *testing.T) {
    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)

    var user = PgUser{
        Username: "user1",
        Password: "12345",
        Superuser: false,
    }

    err = model.Update(user)
    if err != nil {
        t.Error(err)
    }

    user = PgUser{
        Username: "user1",
        Password: "12345",
        Superuser: true,
    }

    err = model.Update(user)
    if err != nil {
        t.Error(err)
    }
}

func TestListAll(t *testing.T) {
    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)

    var users []PgUser
    err = model.ListAll(&users, "u")
    if err != nil {
        t.Error(err)
    }
    for num, item := range users {
        if strings.HasPrefix(item.Username, "user") {
            fmt.Println(num, item)
        }
    }
}

func TestDelete(t *testing.T) {
    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)

    var user1 = PgUser{
        Username: "user1",
    }

    err = model.Delete(user1)
    if err != nil {
        t.Error(err)
    }

    var user2 = PgUser{
        Username: "user2",
    }

    err = model.Delete(user2)
    if err != nil {
        t.Error(err)
    }
}
