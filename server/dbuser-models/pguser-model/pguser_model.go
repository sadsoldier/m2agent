/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package pguserModel

import (
    "log"
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v4/stdlib"
)

type Model struct {
    db *sqlx.DB
}

type PgUser struct {
    Username     string  `db:"usename"   json:"username"`
    Password     string  `db:"passwd"    json:"password"`
    Superuser    bool    `db:"usesuper"  json:"superuser"`
}

type Page struct {
    Total       int         `json:"total"`
    Offset      int         `json:"offset"`
    Limit       int         `json:"limit"`
    Pattern     string      `json:"pattern"`
    PgUsers     *[]PgUser   `json:"users,omitempty"`
}

//func QuoteString(str string) string {
//        return "'" + strings.Replace(str, "'", "''", -1) + "'"
//}

func (this *Model) List(page *Page) (error) {
    var request string
    var err error
    var total int

    pattern := "%" + page.Pattern + "%"
    request = `SELECT COUNT(usename) as total FROM pg_catalog.pg_user WHERE usename LIKE $1`
    err = this.db.QueryRow(request, pattern).Scan(&total)
    if err != nil {
        log.Println(err)
        return err
    }
    page.Total = total

    var users []PgUser
    request = `SELECT usename, usesuper, '' as passwd FROM pg_catalog.pg_user
                WHERE usename LIKE $1
                ORDER BY usename LIMIT $2 OFFSET $3`
    err = this.db.Select(&users, request, pattern, page.Limit, page.Offset)
    if err != nil {
        log.Println(err)
        return err
    }
    page.PgUsers = &users
    return nil
}

func (this *Model) ListAll(pattern string) (*[]PgUser, error) {
    var err error
    pattern = "%" + pattern + "%"
    request := `SELECT usename, usesuper, '' as passwd FROM pg_catalog.pg_user
                WHERE usename LIKE $1
                ORDER BY usename`
    var users []PgUser
    err = this.db.Select(&users, request, pattern)
    if err != nil {
        log.Println(err)
        //return &users, err
    }
    return &users, err
}

func (this *Model) IsExist(user PgUser) (bool, error) {
    var err error
    var users []PgUser
    request := `SELECT usename FROM pg_catalog.pg_user
                    WHERE usename = $1
                    LIMIT 1`
    err = this.db.Select(&users, request, user.Username)
    if err != nil {
        log.Println(err)
        return false, err
    }
    if len(users) == 0 {
        return false, nil
    }
    return true, nil
}


func (this *Model) Create(user PgUser) error {
    request := fmt.Sprintf(`CREATE USER "%s" WITH PASSWORD '%s'`,
                                user.Username,
                                user.Password)
    if user.Superuser {
        request = request + " SUPERUSER"
    } else {
        request = request + " NOSUPERUSER"
    }
    _, err := this.db.Exec(request)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (this *Model) Update(user PgUser) error {
    request := fmt.Sprintf(`ALTER USER "%s" WITH PASSWORD '%s'`,
                            user.Username,
                            user.Password)
    if user.Superuser {
        request = request + " SUPERUSER"
    } else {
        request = request + " NOSUPERUSER"
    }
    _, err := this.db.Exec(request)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (this *Model) Delete(user PgUser) error {
    request := fmt.Sprintf(`DROP USER "%s"`, user.Username)
    _, err := this.db.Exec(request)
    if err != nil {
        log.Println(err)
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
