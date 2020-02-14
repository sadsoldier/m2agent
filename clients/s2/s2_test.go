/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */


package s2client

import (
    "fmt"
    "testing"
    "io/ioutil"
    "os"

)

const (
    uri string = "https://user1:12345@localhost:7001/12345"
    filename string = "test.bin"
)

func TestPut(t *testing.T) {
    buffer := make([]byte, 1024)
    ioutil.WriteFile(filename, buffer, 0666)

    err := Put(uri, filename)
    if err != nil {
        t.Error(err)
    }
    os.Remove(filename)
}

func TestList1(t *testing.T) {
    response, err := List(uri)
    if err != nil {
        t.Error(err)
    }
    fmt.Println(*response)
}

func TestGet(t *testing.T) {
    outname, err := Get(uri, filename, "/tmp")
    if err != nil {
        t.Error(err)
    }
    fmt.Println(outname)
    os.Remove(outname)
}


func TestDelete(t *testing.T) {
    err := Delete(uri, filename)
    if err != nil {
        t.Error(err)
    }
}

func TestList2(t *testing.T) {
    response, err := List(uri)
    if err != nil {
        t.Error(err)
    }
    fmt.Println(*response)
}
