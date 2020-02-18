/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package dbdumpController

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "time"
    "encoding/json"
    "bytes"
    "fmt"

    "agent/server/db-models/pgdb-model"
    "github.com/jmoiron/sqlx"

    "agent/config"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
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

func TestDump(t *testing.T) {

    configuration := config.New()
    controller := New(configuration, nil)

    gin.SetMode(gin.TestMode)
    router := gin.Default()

    router.POST("/dump", controller.Dump)

    dumpReq := DumpRequest{
        ResourseName:       "postgres",
        ResourseType:   "pgsql",
        TransportType:  "s2",
        StorageURI:     "https://user1:12345@localhost:7001/dumps",
        ReportURI:      "http://localhost:7003/report",
        JobId:          "localhost.123",
        MagicCode:      "1234512345",
        Timestamp:      time.Now().Format(time.RFC3339),
    }

    data, err := json.Marshal(dumpReq)
    if err != nil {
        t.Error(err)
    }
    reader := bytes.NewReader(data)
    request, err := http.NewRequest(http.MethodPost, "dump", reader)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }
    request.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, request)

    assert.Equal(t, http.StatusOK, recorder.Code, "Not equal response code")

    fmt.Println(recorder.Body.String())

    response := Response{}
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.Equal(t, response.Error, false, "Model or controller error")

    time.Sleep(time.Second * 7)

}

func TestRestore(t *testing.T) {

    dbx, err := createDB()
    if err != nil {
        t.Error(err)
    }

    pgdb := pgdbModel.New(dbx)
    db := pgdbModel.PgDb{
        DbName: "tmp128",
    }
    defer pgdb.Delete(db)

    configuration := config.New()
    controller := New(configuration, dbx)

    gin.SetMode(gin.TestMode)
    router := gin.Default()

    router.POST("/restore", controller.Restore)

    restoreRequest := RestoreRequest{
        TransportType:  "s2",
        StorageURI:     "https://user1:12345@localhost:7001/dumps",
        DumpFilename:   "postgres--2020-01-28T11:05:55+02:00--thx.unix7.org.sqlz",

        ResourseType:   "pgsql",
        ResourseOwner:   "postgres",
        Destination:    "tmp128",

        ReportURI:      "http://localhost:7003/report",
        JobId:          "localhost.123",
        MagicCode:      "1234512345",
        Timestamp:      time.Now().Format(time.RFC3339),
    }

    data, err := json.Marshal(restoreRequest)
    if err != nil {
        t.Error(err)
    }
    reader := bytes.NewReader(data)
    request, err := http.NewRequest(http.MethodPost, "/restore", reader)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }
    request.Header.Set("Content-Type", "application/json")

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, request)

    assert.Equal(t, http.StatusOK, recorder.Code, "Not equal response code")

    fmt.Println(recorder.Body.String())

    response := Response{}
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.Equal(t, response.Error, false, "Model or controller error")

    time.Sleep(time.Second * 5)
}
