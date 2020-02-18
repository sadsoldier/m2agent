/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */


package dbdumpController

import (
    "net/http"
    "fmt"
    "errors"
    "log"
    "os"
    "path/filepath"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "agent/config"
    "agent/server/dump-models/pgdump-model"
    "agent/clients/s2"
)

type Controller struct {
    config  *config.Config
    dbx     *sqlx.DB
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

type DumpRequest struct {
    ResourseType    string      `json:"resourseType"` // pgsql, mysql, file
    ResourseName    string      `json:"resourseName"`

    TransportType   string      `json:"transport"`  // s2, sftp
    StorageURI      string      `json:"storageURI"`

    JobId           string      `json:"jobId"`
    JobOrigin       string      `json:"jobOrigin"`
    Timestamp       string      `json:"timestamp"`
    ReportURI       string      `json:"reportURI"`
    MagicCode       string      `json:"magicCode"`
}

type RestoreRequest struct {
    TransportType   string      `json:"transportType"`
    StorageURI      string      `json:"storageURI"`

    DumpFilename  string        `json:"dumpFilename"`

    ResourseType    string      `json:"resourseType"`
    ResourseOwner   string      `json:"resourseOwner"`
    Destination     string      `json:"destination"`

    JobId           string      `json:"jobId"`
    JobOrigin       string      `json:"jobOrigin"`
    Timestamp       string      `json:"timestamp"`

    ReportURI       string      `json:"reportURI"`
    MagicCode       string      `json:"magicCode"`
}


func (this *Controller) Dump(context *gin.Context) {
    var request DumpRequest
    var err error
    err = context.Bind(&request)
    if err != nil {
        sendError(context, err)
        return
    }
    go this.dumpProcess(request)
    sendOk(context)
}

func (this *Controller) dumpProcess(request DumpRequest) {
    pgdump := pgdumpModel.New(this.config, this.dbx)

    /* Dump database to tmp file */
    log.Println("dump process start, jobid:", request.JobId)
    outpath, err := pgdump.Dump("postgres")
    if err != nil {
        log.Println("dump process error: ", request.JobId, err)
        os.Remove(outpath)
        return
    }
    defer os.Remove(outpath)
    log.Println("dump process done: ", request.JobId, filepath.Base(outpath))

    /* Put dumpfile to storage */
    log.Println("send process start:", request.JobId, filepath.Base(outpath))
    err = s2client.Put(request.StorageURI, outpath)
    if err != nil {
        log.Println("send process error:", request.JobId, err)
        return
    }
    log.Println("send process done:", request.JobId)
    /* Report */
}


func (this *Controller) Restore(context *gin.Context) {
    var request RestoreRequest
    var err error
    err = context.Bind(&request)
    if err != nil {
        sendError(context, err)
        return
    }
    go this.restoreProcess(request)
    sendOk(context)
}

func (this *Controller) restoreProcess(request RestoreRequest) {

    pgdump := pgdumpModel.New(this.config, this.dbx)

    log.Println("get process start:", request.JobId)
    tmppath, err := s2client.Get(request.StorageURI, request.DumpFilename, this.config.StoreDir)
    if err != nil {
        log.Println("get process error:", request.JobId, err)
        return
    }

    defer os.Remove(tmppath)
    log.Println("get process done:", request.JobId, tmppath)

    log.Println("restore process start:", request.JobId)
    err = pgdump.Restore(tmppath, request.Destination, request.ResourseOwner)
    if err != nil {
        log.Println("restore process error:", request.JobId, err)
        return
    }
    log.Println("restore process done:", request.JobId, request.Destination)
}

func New(config *config.Config, dbx *sqlx.DB) *Controller {
    return &Controller{
        config: config,
        dbx: dbx,
    }
}
