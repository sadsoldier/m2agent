package sftpClient

import (
    "io/ioutil"
    "testing"
    "log"
    "os"
)


const uri string = "sftp://root:bowie2@localhost/tmp/xxxx"
const filename string = "test3.bin"

func TestPut(t *testing.T) {
    os.Mkdir("/tmp/123", os.ModeDir)

    buffer := make([]byte, 1024 * 1024)
    ioutil.WriteFile(filename, buffer, 0660)

    err := Put(uri, filename)
    if err != nil {
        t.Error(err)
    }
}

func TestList1(t *testing.T) {

    files, err := List(uri)
    if err != nil {
        t.Error(err)
    }

    for n, item := range *files {
        log.Println(n, item)
    }
}

func TestGet(t *testing.T) {

    outFilepath, err := Get(uri, filename, "/tmp")
    if err != nil {
        t.Error(err)
    }
    log.Println(outFilepath)
}

func TestRemove(t *testing.T) {

    err := Remove(uri, filename)
    if err != nil {
        t.Error(err)
    }
}

func TestList2(t *testing.T) {

    files, err := List(uri)
    if err != nil {
        t.Error(err)
    }

    for n, item := range *files {
        log.Println(n, item)
    }
}
