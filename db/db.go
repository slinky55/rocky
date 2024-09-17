package db

import "github.com/jmoiron/sqlx"

var Conn *sqlx.DB

func Init() (err error) {
	Conn, err = sqlx.Connect("sqlite3", "./rocky.db")
	if err != nil {
		return err
	}

	return nil
}
