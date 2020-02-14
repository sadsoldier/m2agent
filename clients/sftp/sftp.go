/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 *
 */

/* https://godoc.org/github.com/pkg/sftp */

package sftpClient



import (
    "fmt"
    "log"
    "time"
    "net/url"
    "strconv"

    "os"
    "io"
    "path/filepath"

    "golang.org/x/crypto/ssh"
    "github.com/pkg/sftp"
)

type File struct {
    Name        string      `json:"name"`
    Size        int64       `json:"size"`
    ModTime     string      `json:"modtime"`
}

type URI struct {
    Scheme      string
    Username    string
    Password    string
    Hostname    string
    Port        int
    Path        string
}

const (
    timeout time.Duration = 10
    bufferSize int64 = 1024 * 128
    partSuffix string = ".part"
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

    uri.Port = 22
    if len(parsedURI.Port()) != 0 {
        uri.Port, err = strconv.Atoi(parsedURI.Port())
        if err != nil {
            return uri, err
        }
    }

    uri.Username = parsedURI.User.Username()
    password, exists := parsedURI.User.Password()
    if exists {
        uri.Password = password
    }

    uri.Path = parsedURI.Path

    return uri, nil
}

func Connect(compactURI string) (*ssh.Client, error) {

    var err error
    var client *ssh.Client

    uri, err := UriParse(compactURI)
    if err != nil {
        log.Println(err)
        return client, err
    }

    config := ssh.ClientConfig{
        User: uri.Username,
        Auth: []ssh.AuthMethod{
            ssh.Password(uri.Password),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout: time.Duration(timeout) * time.Second,
    }

    connect, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", uri.Hostname, uri.Port), &config)
    if err != nil {
        log.Println(err)
        return client, err
    }
    return connect, nil
}

func List(compactURI string) (*[]File, error) {

    var files []File
    var err error

    connect, err := Connect(compactURI)
    if err != nil {
        log.Println(err)
        return &files, err
    }
    defer connect.Close()

    client, err := sftp.NewClient(connect)
    if err != nil {
       log.Println(err)
       return &files, err
    }
    defer client.Close()

    uri, err := UriParse(compactURI)
    if err != nil {
        log.Println(err)
        return &files, err
    }

    list, err := client.ReadDir(uri.Path)
    if err != nil {
       log.Println(err)
       return &files, err
    }

    for _, item := range list {
        var file File
        file.Name = item.Name()
        file.Size = item.Size()
        file.ModTime = item.ModTime().Format(time.RFC3339)
        files = append(files, file)
    }

    return &files, err
}

func Get(compactURI, srcFilename, destDir string) (string, error) {

    var err error
    var destFilepath string

    uri, err := UriParse(compactURI)
    if err != nil {
        log.Println(err)
        return destFilepath, err
    }

    connect, err := Connect(compactURI)
    if err != nil {
        log.Println(err)
        return destFilepath, err
    }
    defer connect.Close()

    client, err := sftp.NewClient(connect)
    if err != nil {
       log.Println(err)
       return destFilepath, err
    }
    defer client.Close()

    src, err := client.Open(filepath.Join(uri.Path, srcFilename))
    if err != nil {
       log.Println(err)
       return destFilepath, err
    }
    defer src.Close()

    destPartial := filepath.Join(destDir, srcFilename + partSuffix)
    destFilepath = filepath.Join(destDir, srcFilename)

    dest, err := os.Create(destPartial)
    if err != nil {
       log.Println(err)
       return destPartial, err
    }
    defer dest.Close()

    buffer := make([]byte, bufferSize)
    _, err = io.CopyBuffer(dest, src, buffer)
    if err != nil {
       log.Println(err)
       return destPartial, err
    }

    err = dest.Sync()
    if err != nil {
        log.Println(err)
        return destFilepath, err
    }

    err = os.Rename(destPartial, destFilepath)
    if err != nil {
        os.Remove(destPartial)
        return destFilepath, err
    }

    return destFilepath, nil
}


func Put(compactURI, srcFilepath string) (error) {

    var err error

    uri, err := UriParse(compactURI)
    if err != nil {
        log.Println(err)
        return err
    }

    connect, err := Connect(compactURI)
    if err != nil {
        log.Println(err)
        return err
    }
    defer connect.Close()

    client, err := sftp.NewClient(connect)
    if err != nil {
       log.Println(err)
       return err
    }
    defer client.Close()

    src, err := os.Open(srcFilepath)
    if err != nil {
       log.Println(err)
       return err
    }
    defer src.Close()

    destFilepath := filepath.Join(uri.Path, filepath.Base(srcFilepath))
    partialFilepath := destFilepath + partSuffix

    dest, err := client.Create(partialFilepath)
    if err != nil {
       log.Println(err)
       return err
    }
    defer dest.Close()

    buffer := make([]byte, bufferSize)
    _, err = io.CopyBuffer(dest, src, buffer)
    if err != nil {
       log.Println(err)
       return err
    }

    err = client.PosixRename(partialFilepath, destFilepath)
    if err != nil {
        client.Remove(partialFilepath)
        return err
    }

    return nil
}

func Remove(compactURI, xfilepath string) (error) {

    var err error

    uri, err := UriParse(compactURI)
    if err != nil {
        log.Println(err)
        return err
    }

    connect, err := Connect(compactURI)
    if err != nil {
        log.Println(err)
        return err
    }
    defer connect.Close()

    client, err := sftp.NewClient(connect)
    if err != nil {
       log.Println(err)
       return err
    }
    defer client.Close()

    client.Remove(filepath.Join(uri.Path, xfilepath))
    if err != nil {
       log.Println(err)
       return err
    }

    return nil
}
