package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spaco/affiliate/src/config"
)

var db *sql.DB

func OpenDb(dbConfig *config.Db) *sql.DB {
	conn_str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SslMode)
	//	fmt.Println(conn_str)
	var err error
	db, err = sql.Open("postgres", conn_str)
	checkErr(err)
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.Ping()
	return db
}

func CloseDb() {
	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func BeginTx() (*sql.Tx, bool) {
	tx, err := db.Begin()
	checkErr(err)
	return tx, false
}

func BeginReadTx() (*sql.Tx, bool) {
	tx, err := db.Begin()
	checkErr(err)
	return tx, false
}

func Rollback(tx *sql.Tx, commit *bool) {
	if !*commit {
		tx.Rollback()
	}
}
