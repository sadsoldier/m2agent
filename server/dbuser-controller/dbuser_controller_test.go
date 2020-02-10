/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package dbuserController

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

    "agent/server/pguser-model"
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

func TestListAll(t *testing.T) {

    dbx, err := createDbx()
    if err != nil {
        t.Error(err)
    }

    configuration := config.New()
    controller := New(configuration, dbx)

    gin.SetMode(gin.TestMode)

    router := gin.Default()
    router.GET("/pgusers/listall", controller.ListAll)

    request, err := http.NewRequest(http.MethodGet, "/pgusers/listall", nil)
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
    router.POST("/pguser/create", controller.Create)

    pguser := pguserModel.PgUser{
        Username: "qwerty",
        Password: "123456",
    }

    data, err := json.Marshal(pguser)
    if err != nil {
        t.Error(err)
    }

    reader := bytes.NewReader(data)
    request, err := http.NewRequest(http.MethodPost, "/pguser/create", reader)
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
    router.POST("/pguser/update", controller.Update)

    pguser := pguserModel.PgUser{
        Username: "qwerty",
        Password: "876543",
    }

    data, err := json.Marshal(pguser)
    if err != nil {
        t.Error(err)
    }

    reader := bytes.NewReader(data)
    request, err := http.NewRequest(http.MethodPost, "/pguser/update", reader)
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


func TestList(t *testing.T) {

    dbx, err := createDbx()
    if err != nil {
        t.Error(err)
    }

    configuration := config.New()
    controller := New(configuration, dbx)

    gin.SetMode(gin.TestMode)

    router := gin.Default()
    router.POST("/pguser/list", controller.List)

    page := pguserModel.Page{
        Offset: 0,
        Limit: 7,
        Pattern: "qwerty",
    }
    data, err := json.Marshal(page)
    if err != nil {
        t.Error(err)
    }


    reader := bytes.NewReader(data)
    request, err := http.NewRequest(http.MethodPost, "/pguser/list", reader)
    if err != nil {
        t.Fatalf("Couldn't create request: %s\n", err)
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
    router.POST("/pguser/delete", controller.Delete)

    pguser := pguserModel.PgUser{
        Username: "qwerty",
    }

    data, err := json.Marshal(pguser)
    if err != nil {
        t.Error(err)
    }

    reader := bytes.NewReader(data)
    request, err := http.NewRequest(http.MethodPost, "/pguser/delete", reader)
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
