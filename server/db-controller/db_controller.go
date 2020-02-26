/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */


package dbController

import (
    "net/http"
    "fmt"
    "errors"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "agent/config"
    "agent/server/db-models/pgdb-model"
)

type Controller struct {
    config *config.Config
    db *sqlx.DB
    pgdb *pgdbModel.Model
}

type Response struct {
    Error       bool        `json:"error"`
    Message     string      `json:"message,omitempty"`
    Result      interface{} `json:"result,omitempty"`
}

func sendError(context *gin.Context, err error) {
    if err == nil {
        err = errors.New("undefined")
    }
    log.Printf("%s\n", err)
    response := Response{
        Error: true,
        Message: fmt.Sprintf("%s", err),
        Result: nil,
    }
    context.JSON(http.StatusOK, response)
}

func sendOk(context *gin.Context) {
    response := Response{
        Error: false,
        Message: "",
        Result: nil,
    }
    context.JSON(http.StatusOK, response)
}

func sendMessage(context *gin.Context, message string) {
    log.Printf("%s\n", message)
    response := Response{
        Error: false,
        Message: fmt.Sprintf("%s", message),
        Result: nil,
    }
    context.JSON(http.StatusOK, response)
}

func sendResult(context *gin.Context, result interface{}) {
    response := Response{
        Error: false,
        Message: "",
        Result: result,
    }
    context.JSON(http.StatusOK, &response)
}


func (this *Controller) List(context *gin.Context) {
    var page pgdbModel.Page
    err := context.Bind(&page)
    if err != nil {
        sendError(context, err)
        return
    }
    this.pgdb.List(&page)
    sendResult(context, &page)
}

type ListAllRequest struct {
    Pattern string `json:"pattern"`
}

func (this *Controller) ListAll(context *gin.Context) {
    var pgdbs []pgdbModel.PgDb
    var request ListAllRequest
    err := context.Bind(&request)
    if err != nil {
        sendError(context, err)
        return
    }
    this.pgdb.ListAll(&pgdbs, request.Pattern)
    sendResult(context, &pgdbs)
}

func (this *Controller) Create(context *gin.Context) {
    var pgdb pgdbModel.PgDb
    var err error
    err = context.Bind(&pgdb)
    log.Println(pgdb)
    if err != nil {
        sendError(context, err)
        return
    }

    err = this.pgdb.Create(pgdb)
    if err != nil {
        sendError(context, err)
        return
    }
    sendOk(context)
}

func (this *Controller) Update(context *gin.Context) {
    var pgdb pgdbModel.PgDb
    var err error
    err = context.Bind(&pgdb)
    if err != nil {
        sendError(context, err)
        return
    }

    err = this.pgdb.Update(pgdb)
    if err != nil {
        sendError(context, err)
        return
    }
    sendOk(context)
}

func (this *Controller) Delete(context *gin.Context) {
    var pgdb pgdbModel.PgDb
    var err error
    err = context.Bind(&pgdb)
    if err != nil {
        sendError(context, err)
        return
    }

    err = this.pgdb.Delete(pgdb)
    if err != nil {
        sendError(context, err)
        return
    }
    sendOk(context)
}

func New(config *config.Config, db *sqlx.DB) *Controller {
    return &Controller{
        config: config,
        db: db,
        pgdb: pgdbModel.New(db),
    }
}
