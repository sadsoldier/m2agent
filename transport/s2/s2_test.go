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

func TestPutGet(t *testing.T) {

    buffer := make([]byte, 1024)
    tmpname := "file.test"
    ioutil.WriteFile(tmpname, buffer, 0666)

    resourse := "https://user1:12345@localhost:7001/tmp"
    err := Put(resourse, tmpname)
    if err != nil {
        t.Error(err)
    }
    cwd, _ := os.Getwd()
    filename, err := Get(resourse, tmpname, cwd)
    if err != nil {
        t.Error(err)
    }
    fmt.Println(filename)
    os.Remove(filename)
}
