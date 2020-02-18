/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package pgdumpModel

import (
    "agent/config"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "errors"
    "time"

    "github.com/jmoiron/sqlx"

    "agent/server/db-models/pgdb-model"
    "agent/server/dbuser-models/pguser-model"
    "agent/clients/s2"

)

type Model struct{
    config *config.Config
    dbx *sqlx.DB
}

func Timestamp() string {
    return time.Now().Format(time.RFC3339)
}

func (this *Model) Get(storeUri, filename string) (string, error) {
    return s2client.Get(storeUri, filename, this.config.StoreDir)
}

func (this *Model) Put(storeUri, filename string) error {
    return s2client.Put(storeUri, filename)
}

func (this *Model) Dump(dbname string) (string, error) {

    timestamp := time.Now().Format(time.RFC3339)
    hostname, _ := os.Hostname()
    dumpname := fmt.Sprintf("%s--%s--%s.sqlz", dbname, timestamp, hostname)

    dumpPath := filepath.Join(this.config.StoreDir, dumpname)

    cmd := exec.Command("pg_dump", "--format=c", "--file=" + dumpPath, dbname)

    env := os.Environ()
    env = append(env, "PGHOST=" + this.config.DbHost)
    env = append(env, "PGUSER=" + this.config.DbUser)
    env = append(env, "PGPASSWORD=" + this.config.DbPass)
    cmd.Env = env
    output, err := cmd.CombinedOutput()

    if err != nil {
        return "", errors.New(string(output))
    }
    return dumpPath, nil
}

func (this *Model) Restore(filepath, dbname, dbowner string) error {
    var err error

    /* Check DB */
    pgdb := pgdbModel.New(this.dbx)
    db := pgdbModel.PgDb{
        DbName: dbname,
        Owner: dbowner,
    }
    exist, err := pgdb.IsExist(db)
    if err != nil {
        return err
    }
    if exist {
        return errors.New(fmt.Sprintf("database %s already exist", dbname))
    }

    /* Check db user */
    pguser := pguserModel.New(this.dbx)
    user := pguserModel.PgUser{
        Username: dbowner,
    }
    exist, err = pguser.IsExist(user)
    if err != nil {
        return err
    }
    if !exist {
        return errors.New(fmt.Sprintf("user %s not exist", dbowner))
    }

    /* Create empty database */
    err = pgdb.Create(db)
    if err != nil {
        return err
    }

    /* Restore */
    env := os.Environ()
    env = append(env, "PGHOST=" + this.config.DbHost)
    env = append(env, "PGUSER=" + this.config.DbUser)
    env = append(env, "PGPASSWORD=" + this.config.DbPass)

    restoreCmd := exec.Command("pg_restore", "--dbname=" + dbname, filepath)
    restoreCmd.Env = env
    output, err := restoreCmd.CombinedOutput()

    if err != nil {
        return errors.New(string(output))
    }
    return nil
}

func New(config *config.Config, dbx *sqlx.DB) *Model {
    return &Model{
        config: config,
        dbx: dbx,
    }
}
