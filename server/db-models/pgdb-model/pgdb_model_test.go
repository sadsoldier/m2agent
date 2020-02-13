/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package pgdbModel

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
    xdb := PgDb{
        Name: "test123",
        Owner: "postgres",
    }
    model.Create(xdb)
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

    var dbs []PgDb
    page := Page{
        Limit: 500,
        PgDbs: &dbs,
    }
    err = model.List(&page)
    if err != nil {
        t.Error(err)
    }

    for num, item := range *page.PgDbs {
        if strings.HasPrefix(item.Name, "test123") {
            fmt.Println(num, item)
        }
    }
}

func TestUpdate(t *testing.T) {
    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)
    xdb := PgDb{
        Name: "test123",
        Owner: "postgres",
    }
    err = model.Update(xdb)
    if err != nil {
        t.Error(err)
    }
}

func TestIsExist(t *testing.T) {
    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)
    xdb := PgDb{
        Name: "test123",
    }
    res, err := model.IsExist(xdb)
    if err != nil {
        t.Error(err)
    }
    if res != true {
        t.Error(errors.New("db must be exist"))
    }
}


func TestIsExistNot(t *testing.T) {
    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)
    xdb := PgDb{
        Name: "1234567",
    }
    res, err := model.IsExist(xdb)
    if err != nil {
        t.Error(err)
    }
    if res == true {
        t.Error(errors.New("xxx"))
    }
}


func TestDelete(t *testing.T) {
    db, err := createDB()
    if err != nil {
        t.Error(err)
    }
    model := New(db)
    xdb := PgDb{
        Name: "test123",
    }
    model.Delete(xdb)
    if err != nil {
        t.Error(err)
    }
}
