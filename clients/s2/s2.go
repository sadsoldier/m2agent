/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package s2client

import (
    "time"
    "strconv"
    "path/filepath"
    "os"
    "net/url"
    "net/http"
    "mime/multipart"
    "io/ioutil"
    "io"
    "fmt"
    "errors"
    "encoding/json"
    "crypto/tls"
    "bytes"
    "log"
)

type URI struct {
    Scheme      string
    //Username    string
    //Password    string
    Hostname    string
    Userinfo    string
    Port        int
    Path        string
}

const (
    bufferSize int64 = 1024 * 128
    partSuffix string = ".part"
    defaultPort int = 7001
)

func UriParse(compactURI string) (URI, error) {

    var err error
    var uri URI

    parsedURI, err := url.Parse(compactURI)
    if err != nil {
        return uri, err
    }

    uri.Scheme = parsedURI.Scheme
    uri.Hostname = parsedURI.Host

    uri.Port = defaultPort
    if len(parsedURI.Port()) != 0 {
        uri.Port, err = strconv.Atoi(parsedURI.Port())
        if err != nil {
            return uri, err
        }
    }

    //uri.Username = parsedURI.User.Username()
    //password, exists := parsedURI.User.Password()
    //if exists {
        //uri.Password = password
    //}

    uri.Userinfo = parsedURI.User.String()
    uri.Path = parsedURI.Path

    return uri, nil
}

type File struct {
    Name        string
    Size        int64
    ModTime     time.Time
}

type S2File struct {
    Name        string      `json:"name"`
    Size        int64       `json:"size"`
    ModTime     string      `json:"modtime"`
}

type S2ListResponse struct {
    Error       bool        `json:"error"`
    Message     string      `json:"message"`
    Result      []S2File    `json:"result"`
}

type ListForm struct {
    Bucket      string  `json:"bucket"`
    Pattern     string  `json:"pattern"`
}

type S2DeleteResponse struct {
    Error       bool        `json:"error"`
    Message     string      `json:"message"`
}


func List(storeURI string) (*[]File, error) {
    var err error
    var response S2ListResponse
    var files []File

    parsedURI, err := UriParse(storeURI)
    if err != nil {
        return &files, err
    }

    s2url := fmt.Sprintf("%s://%s@%s/api/v1/file/list",
        parsedURI.Scheme,
        parsedURI.Userinfo,
        parsedURI.Hostname)

    form := ListForm{
        Bucket: parsedURI.Path,
        Pattern: "*",
    }

    transCfg := &http.Transport{
         TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: transCfg}

    data, _ := json.Marshal(form)
    reader := bytes.NewReader([]byte(data))

    resp, err := client.Post(s2url, "application/json", reader)
    if err != nil {
        return &files, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return &files, err
    }

    err = json.Unmarshal(body, &response)
    if err != nil {
        return &files, err
    }

    if response.Error == true {
        log.Println(response.Message)
        return &files, errors.New(response.Message)
    }

    /* Map result to file array */
    for _, item := range response.Result {
        var file File
        file.Name = item.Name
        file.Size = item.Size
        file.ModTime, _ = time.Parse(time.RFC3339, item.ModTime)
        files = append(files, file)
    }

    return &files, nil
}

func Get(storeURI, filename, destDir string) (string, error) {
    var err error
    var destFilepath string

    parsedURI, err := UriParse(storeURI)
    if err != nil {
        return destFilepath, err
    }

    s2uri := fmt.Sprintf("%s://%s@%s/api/v1/file/down/%s/%s",
        parsedURI.Scheme,
        parsedURI.Userinfo,
        parsedURI.Hostname,
        parsedURI.Path,
        filename)

    transCfg := &http.Transport{
         TLSClientConfig: &tls.Config{ InsecureSkipVerify: true },
    }
    client := &http.Client{ Transport: transCfg }

    response, err := client.Get(s2uri)
    if err != nil {
        return destFilepath, err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return destFilepath, errors.New(fmt.Sprintf("Wrong status code: %d", response.StatusCode))
    }

    file, err := ioutil.TempFile(destDir, filename + partSuffix)
    if err != nil {
        return destFilepath, err
    }
    defer file.Close()

    buffer := make([]byte,  bufferSize)
    _, err = io.CopyBuffer(file, response.Body, buffer)
    if err != nil {
        os.Remove(file.Name())
        return destFilepath, err
    }
    destFilepath = filepath.Join(destDir, filename)

    err = os.Rename(file.Name(), destFilepath)
    if err != nil {
        os.Remove(file.Name())
        return destFilepath, err
    }

    return destFilepath, nil
}

func Put(storeURI, filename string) error {
    var err error

    parsedURI, err := UriParse(storeURI)
    if err != nil {
        return err
    }

    s2url := fmt.Sprintf("%s://%s@%s/api/v1/file/put",
        parsedURI.Scheme,
        parsedURI.Userinfo,
        parsedURI.Hostname)

    pipeOut, pipeIn := io.Pipe()
    writer := multipart.NewWriter(pipeIn)

    go func() {
        defer pipeIn.Close()
        defer writer.Close()

        _ = writer.WriteField("filename", filepath.Base(filename))
        _ = writer.WriteField("bucket", parsedURI.Path)

        part, err := writer.CreateFormFile("file", filepath.Base(filename))
        if err != nil {
            return
        }
        file, err := os.Open(filename)
        if err != nil {
            return
        }
        defer file.Close()
        if _, err = io.Copy(part, file); err != nil {
            return
        }
    }()

    transCfg := &http.Transport{
         TLSClientConfig: &tls.Config{ InsecureSkipVerify: true },
    }
    client := &http.Client{ Transport: transCfg }
    resp, err := client.Post(s2url, writer.FormDataContentType(), pipeOut)
    if err != nil {
        return err
    }
    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    return nil
}


type deleteForm struct {
    Filename    string  `form:"filename" json:"filename" binding:"required" `
    Bucketname  string  `form:"bucket"   json:"bucket"`
}


func Delete(storeURI, filename string) error {
    var err error
    var response S2DeleteResponse

    parsedURI, err := UriParse(storeURI)
    if err != nil {
        return err
    }

    s2url := fmt.Sprintf("%s://%s@%s/api/v1/file/delete",
        parsedURI.Scheme,
        parsedURI.Userinfo,
        parsedURI.Hostname)

    form := deleteForm{
        Bucketname: parsedURI.Path,
        Filename: filename,
    }

    transCfg := &http.Transport{
         TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: transCfg}

    data, _ := json.Marshal(form)
    reader := bytes.NewReader([]byte(data))

    resp, err := client.Post(s2url, "application/json", reader)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    err = json.Unmarshal(body, &response)
    if err != nil {
        return err
    }

    if response.Error == true {
        log.Println(response.Message)
        return errors.New(response.Message)
    }

    return nil
}
