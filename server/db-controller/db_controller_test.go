/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package dbController

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "encoding/json"
    "bytes"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
    "github.com/stretchr/testify/assert"

    "agent/server/pgdb-model"
    "agent/config"
)

func createDbx() (*sqlx.DB, error) {
    dburi := fmt.Sprintf("postgres://postgres@localhost/postgres?sslmode=disable")
    dbx, err := sqlx.Open("pgx", dburi)
    if err != nil {
        return nil, err
    }
    err = dbx.Ping()
    if err != nil {
        return nil, err
    }
    return dbx, nil
}


func TestList(t *testing.T) {
  dbx, err := createDbx()
    if err != nil {
        t.Error(err)
    }

    configuration := config.New()
    controller := New(configuration, dbx)

    gin.SetMode(gin.TestMode)

    router := gin.Default()
    router.GET("/list", controller.ListAll)

    request, err := http.NewRequest(http.MethodGet, "/list", nil)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }

    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, request)

    assert.Equal(t, http.StatusOK, recorder.Code, "Not equal response code")
    response := Response{}
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.Equal(t, response.Error, false, "Model or controller error")
}

func TestCreate(t *testing.T) {
    dbx, err := createDbx()
    if err != nil {
        t.Error(err)
    }

    configuration := config.New()
    controller := New(configuration, dbx)

    gin.SetMode(gin.TestMode)

    router := gin.Default()
    router.POST("/create", controller.Create)

    pgdb := pgdbModel.PgDb{
        Name: "qwerty",
        Owner: "postgres",
    }

    data, err := json.Marshal(pgdb)
    if err != nil {
        t.Error(err)
    }

    reader := bytes.NewReader(data)
    request, err := http.NewRequest(http.MethodPost, "/create", reader)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }
    request.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, request)

    assert.Equal(t, http.StatusOK, recorder.Code, "Not equal response code")
    response := Response{}
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.Equal(t, response.Error, false, "Model or controller error")
}

func TestUpdate(t *testing.T) {
    dbx, err := createDbx()
    if err != nil {
        t.Error(err)
    }

    configuration := config.New()
    controller := New(configuration, dbx)

    gin.SetMode(gin.TestMode)

    router := gin.Default()
    router.POST("/update", controller.Update)

    pgdb := pgdbModel.PgDb{
        Name: "qwerty",
        Owner: "postgres",
    }

    data, err := json.Marshal(pgdb)
    if err != nil {
        t.Error(err)
    }

    reader := bytes.NewReader(data)
    request, err := http.NewRequest(http.MethodPost, "/update", reader)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }
    request.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, request)

    assert.Equal(t, http.StatusOK, recorder.Code, "Not equal response code")
    response := Response{}
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.Equal(t, response.Error, false, "Model or controller error")
}

func TestDelete(t *testing.T) {

    dbx, err := createDbx()
    if err != nil {
        t.Error(err)
    }

    configuration := config.New()
    controller := New(configuration, dbx)

    gin.SetMode(gin.TestMode)

    router := gin.Default()
    router.POST("/delete", controller.Delete)

    pgdb := pgdbModel.PgDb{
        Name: "qwerty",
    }

    data, err := json.Marshal(pgdb)
    if err != nil {
        t.Error(err)
    }

    reader := bytes.NewReader(data)
    request, err := http.NewRequest(http.MethodPost, "/delete", reader)
    if err != nil {
        t.Fatalf("Couldn't create request: %v\n", err)
    }
    request.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, request)

    assert.Equal(t, http.StatusOK, recorder.Code, "response code")
    response := Response{}
    err = json.Unmarshal(recorder.Body.Bytes(), &response)
    assert.Equal(t, response.Error, false, "Model or controller error")
}
