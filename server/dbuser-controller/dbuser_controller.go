/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */


package dbuserController

import (
    "net/http"
    "fmt"
    "errors"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "agent/config"
    "agent/server/pguser-model"

)

type Controller struct {
    config *config.Config
    db *sqlx.DB
    pguser *pguserModel.Model
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
    var page pguserModel.Page
    _ = context.Bind(&page)
    this.pguser.List(&page)
    sendResult(context, &page)
}

func (this *Controller) ListAll(context *gin.Context) {
    var pgusers []pguserModel.PgUser
    this.pguser.ListAll(&pgusers, "")
    sendResult(context, &pgusers)
}

func (this *Controller) Create(context *gin.Context) {
    var pguser pguserModel.PgUser
    var err error
    err = context.Bind(&pguser)
    log.Println(pguser)
    if err != nil {
        sendError(context, err)
        return
    }

    err = this.pguser.Create(pguser)
    if err != nil {
        sendError(context, err)
        return
    }
    sendOk(context)
}

func (this *Controller) Update(context *gin.Context) {
    var pguser pguserModel.PgUser
    var err error
    err = context.Bind(&pguser)
    if err != nil {
        sendError(context, err)
        return
    }

    err = this.pguser.Update(pguser)
    if err != nil {
        sendError(context, err)
        return
    }
    sendOk(context)
}

func (this *Controller) Delete(context *gin.Context) {
    var pguser pguserModel.PgUser
    var err error
    err = context.Bind(&pguser)
    if err != nil {
        sendError(context, err)
        return
    }

    err = this.pguser.Delete(pguser)
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
        pguser: pguserModel.New(db),
    }
}
