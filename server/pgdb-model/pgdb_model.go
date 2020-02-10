/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package pgdbModel

import (
    "log"
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v4/stdlib"
)

type Model struct {
    db *sqlx.DB
}

type PgDb struct {
    DbName          string  `db:"name"          json:"dbname"`
    PrevDbName      string  `db:"-"             json:"prevdbname,omitempty"`
    Owner           string  `db:"owner"         json:"owner"`
    Size            int     `db:"size"          json:"size"`
    Numbackends     int     `db:"numbackends"   json:"numbackends"`
}

type Page struct {
    Total       int         `json:"total"`
    Offset      int         `json:"offset"`
    Limit       int         `json:"limit"`
    Pattern     string      `json:"pattern"`
    PgDbs     *[]PgDb       `json:"dbs,omitempty"`
}

func (this *Model) List(page *Page) (error) {
    var request string
    var err error
    var total int

    pattern := "%" + page.Pattern + "%"
    request = `SELECT COUNT(datname) FROM pg_catalog.pg_database WHERE datname LIKE $1`
    err = this.db.QueryRow(request, pattern).Scan(&total)
    if err != nil {
        log.Println(err)
        return err
    }
    page.Total = total

    var dbs []PgDb
    //request = `SELECT c.datname AS name,
    //                    pg_get_userbyid(c.datdba) AS owner,
    //                    pg_database_size(c.datname) AS size
    //                FROM pg_catalog.pg_database AS c
    //                WHERE c.datname LIKE $1
    //                ORDER BY c.datname
    //                LIMIT $2 OFFSET $3`
    request = `SELECT d.datname AS name,
                        pg_database_size(d.datname) AS size,
                        u.usename AS owner,
                        s.numbackends AS numbackends
                    FROM pg_database d, pg_user u, pg_stat_database s
                    WHERE d.datdba = u.usesysid AND d.datname = s.datname
                    AND d.datname LIKE $1
                    ORDER by d.datname
                    LIMIT $2 OFFSET $3`
    err = this.db.Select(&dbs, request, pattern, page.Limit, page.Offset)
    if err != nil {
        log.Println(err)
        return err
    }
    page.PgDbs = &dbs
    return nil
}

func (this *Model) ListAll(databases *[]PgDb, pattern string) error {
    var err error
    pattern = "%" + pattern + "%"
    //request := `SELECT c.datname AS name,
    //                    pg_get_userbyid(c.datdba) AS owner,
    //                    pg_database_size(c.datname) AS size
    //               FROM pg_catalog.pg_database AS c
    //                WHERE name LIKE $1
    //                ORDER BY name`
    request := `SELECT d.datname AS name,
                        pg_database_size(d.datname) AS size,
                        u.usename AS owner,
                        s.numbackends AS numbackends
                    FROM pg_database d, pg_user u, pg_stat_database s
                    WHERE d.datdba = u.usesysid
                        AND d.datname = s.datname
                        AND d.datname LIKE $1
                    ORDER by d.datname`
    err = this.db.Select(databases, request, pattern)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}


func (this *Model) IsExist(database PgDb) (bool, error) {
    var err error
    var databases []PgDb
    request := `SELECT datname AS name
                    FROM pg_catalog.pg_database
                    WHERE datname = $1
                    LIMIT 1`
    err = this.db.Select(&databases, request, database.DbName)
    if err != nil {
        log.Println(err)
        return false, err
    }
    if len(databases) == 0 {
        return false, nil
    }
    return true, nil
}


func (this *Model) Create(database PgDb) error {
    log.Println(database)
    request := fmt.Sprintf(`CREATE DATABASE "%s" OWNER "%s"`,
                                database.DbName,
                                database.Owner)
    _, err := this.db.Exec(request)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (this *Model) Update(database PgDb) error {
    request := fmt.Sprintf(`ALTER DATABASE "%s" OWNER TO "%s"`,
                            database.PrevDbName,
                            database.Owner)
    _, err := this.db.Exec(request)
    if err != nil {
        log.Println(err)
        return err
    }

    if database.PrevDbName != database.DbName {
        request := fmt.Sprintf(`ALTER DATABASE "%s" RENAME TO "%s"`,
                                database.PrevDbName,
                                database.DbName)
        _, err := this.db.Exec(request)
        if err != nil {
            log.Println(err)
            return err
        }
    }
    return nil
}

func (this *Model) Delete(database PgDb) error {
    request := fmt.Sprintf(`DROP DATABASE "%s"`, database.DbName)
    _, err := this.db.Exec(request)
    if err != nil {
        //log.Println(err)
        return err
    }
    return nil
}

func New(db *sqlx.DB) *Model {
    model := Model{
        db: db,
    }
    return &model
}
