/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package s2client

import (
    "net/url"
    "crypto/tls"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "errors"
)

func Get(storeUri, filename, destDir string) (string, error) {
    var err error

    resourse, err := url.Parse(storeUri)
    if err != nil {
        return "", err
    }

    scheme := resourse.Scheme
    host := resourse.Host

    userinfo := resourse.User.String()
    bucket := resourse.Path

    s2uri := fmt.Sprintf("%s://%s@%s/api/v1/file/down/%s/%s", scheme, userinfo, host, bucket, filename)

    transCfg := &http.Transport{
         TLSClientConfig: &tls.Config{ InsecureSkipVerify: true },
    }
    client := &http.Client{ Transport: transCfg }

    response, err := client.Get(s2uri)
    if err != nil {
        return "", err
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return "", errors.New(fmt.Sprintf("Wrong status code: %d", response.StatusCode))
    }

    file, err := ioutil.TempFile(destDir, filename + ".part.*")
    if err != nil {
        return "", err
    }
    defer file.Close()

    buffer := make([]byte,  128 * 1024)
    _, err = io.CopyBuffer(file, response.Body, buffer)
    if err != nil {
        os.Remove(file.Name())
        return "", err
    }
    resPath := filepath.Join(destDir, filename)

    err = os.Rename(file.Name(), resPath)
    if err != nil {
        os.Remove(file.Name())
        return "", err
    }

    return resPath, nil
}

func Put(storeUri, filename string) error {
    var err error

    resourse, err := url.Parse(storeUri)
    if err != nil {
        return err
    }

    scheme := resourse.Scheme
    host := resourse.Host

    userinfo := resourse.User.String()
    bucket := resourse.Path

    s2url := fmt.Sprintf("%s://%s@%s/api/v1/file/put", scheme, userinfo, host)

    pipeOut, pipeIn := io.Pipe()
    writer := multipart.NewWriter(pipeIn)

    go func() {
        defer pipeIn.Close()
        defer writer.Close()

        _ = writer.WriteField("filename", filepath.Base(filename))
        _ = writer.WriteField("bucket", bucket)

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
