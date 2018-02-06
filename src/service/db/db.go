package db

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spolabs/affiliate/src/config"
	"strconv"
	"time"
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

func InClause(count int, first int) string {
	if count == 1 {
		return fmt.Sprintf("($%d)", first)
	} else if count == 0 {
		panic("count can't be zero")
	}
	var buffer bytes.Buffer
	buffer.WriteString("($")
	buffer.WriteString(strconv.Itoa(first))
	for i := 1; i < count; i++ {
		buffer.WriteString(", $")
		first++
		buffer.WriteString(strconv.Itoa(first))
	}
	buffer.WriteString(")")
	return buffer.String()
}

type Time time.Time

const (
	timeFormart = "2006-01-02 15:04:05"
)

func (self *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*self = Time(now)
	return
}

func (self Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(self).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (self Time) String() string {
	return time.Time(self).Format(timeFormart)
}

type NullTime struct {
	Time  Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (self *NullTime) Scan(value interface{}) error {
	var t time.Time
	t, self.Valid = value.(time.Time)
	self.Time = Time(t)
	return nil
}

// Value implements the driver Valuer interface.
func (self NullTime) Value() (driver.Value, error) {
	if !self.Valid {
		return nil, nil
	}
	return self.Time, nil
}
