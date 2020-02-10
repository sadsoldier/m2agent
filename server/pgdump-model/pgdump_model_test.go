/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */


package pgdumpModel

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "testing"

    "agent/config"
    "agent/server/pgdb-model"

    "github.com/jmoiron/sqlx"
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


func _TestDump(t *testing.T) {

    dbx, err := createDB()
    if err != nil {
        t.Error(err)
    }

    configuration := config.Config{
        StoreDir: "/var/tmp/",
        DbUser: "postgres",
        DbPass: "password",
        DbHost: "localhost",
    }

    agent := New(&configuration, dbx)
    path, err := agent.Dump("postgres")
    if err != nil {
        t.Error(err)
    }
    fmt.Println(path)
}

func TestPutGet(t *testing.T) {

    dbx, err := createDB()
    if err != nil {
        t.Error(err)
    }

    configuration := config.Config{
        StoreDir: "/var/tmp/",
        DbUser: "postgres",
        DbPass: "password",
        DbHost: "localhost",
    }
    agent := New(&configuration, dbx)

    buffer := make([]byte, 1024)
    tmpname := "file.test"
    ioutil.WriteFile(tmpname, buffer, 0666)
    defer os.Remove(tmpname)

    resourse := "https://user1:12345@localhost:7001/tmp"
    err = agent.Put(resourse, tmpname)
    if err != nil {
        t.Error(err)
    }

    filename, err := agent.Get(resourse, tmpname)
    defer os.Remove(filename)
    if err != nil {
        t.Error(err)
    }
    fmt.Println(filename)
}

func TestDumpPutGetRestore(t *testing.T) {

    dbx, err := createDB()
    if err != nil {
        t.Error(err)
    }


    configuration := config.Config{
        StoreDir: "/var/tmp/",
        DbUser: "postgres",
        DbPass: "password",
        DbHost: "localhost",
    }

    agent := New(&configuration, dbx)
    outpath, err := agent.Dump("postgres")
    if err != nil {
        t.Error(err)
    }
    fmt.Println(outpath)

    storeUri := "https://user1:12345@localhost:7001/dumps"
    err = agent.Put(storeUri, outpath)
    if err != nil {
        t.Error(err)
    }

    tmppath, err := agent.Get(storeUri, filepath.Base(outpath))
    if err != nil {
        t.Error(err)
    }
    defer os.Remove(tmppath)

    pgdb := pgdbModel.New(dbx)
    db := pgdbModel.PgDb{
        Name: "tmp127",
    }
    pgdb.Delete(db)

    err = agent.Restore(tmppath, "tmp127", "postgres")
    defer pgdb.Delete(db)
    if err != nil {
        t.Error(err)
    }

    fmt.Println(tmppath)
}

func _TestDumpRestore(t *testing.T) {

    dbx, err := createDB()
    if err != nil {
        t.Error(err)
    }

    configuration := config.Config{
        StoreDir: "/var/tmp/",
        DbUser: "postgres",
        DbPass: "password",
        DbHost: "localhost",
    }
    agent := New(&configuration, dbx)

    outpath, err := agent.Dump("postgres")
    if err != nil {
        t.Error(err)
    }
    fmt.Println(outpath)

    pgdb := pgdbModel.New(dbx)
    db := pgdbModel.PgDb{
        Name: "tmp128",
    }
    pgdb.Delete(db)

    err = agent.Restore(outpath, "tmp128", "postgres")
    if err != nil {
        t.Error(err)
    }

}
