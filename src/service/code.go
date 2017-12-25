package service

import (
	"github.com/spaco/affiliate/src/service/db"
	pg "github.com/spaco/affiliate/src/service/postgresql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetTrackingCodeOrGenerate(address string, refAddress string) uint64 {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	id, _ := pg.GetTrackingCode(tx, address)
	if id == 0 {
		id = pg.GenerateTrackingCode(tx, address, refAddress)
	}
	checkErr(tx.Commit())
	commit = true
	return id
}
