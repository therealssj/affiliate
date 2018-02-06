package service

import (
	"github.com/spolabs/affiliate/src/service/db"
	pg "github.com/spolabs/affiliate/src/service/postgresql"
	"github.com/spolabs/affiliate/src/tracking_code"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetTrackingCodeOrGenerate(address string, refCode string) uint64 {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	id, _ := pg.GetTrackingCode(tx, address)
	if id == 0 {
		var refAddress string
		if len(refCode) > 0 {
			id := tracking_code.GetId(refCode)
			if id != 0 {
				refAddress, _ = pg.GetAddrById(tx, id)
				if refAddress == address {
					refAddress = ""
				}
			}
		}
		id = pg.GenerateTrackingCode(tx, address, refAddress)
	}
	checkErr(tx.Commit())
	commit = true
	return id
}

func QueryRewardRecord(address string) []db.RewardRecord {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	res := pg.QueryRewardRecord(tx, address)
	checkErr(tx.Commit())
	commit = true
	return res
}

func QueryRewardRemain(address string) uint64 {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	m := pg.QueryRewardRemain(tx, address)
	checkErr(tx.Commit())
	commit = true
	if rm, ok := m[address]; ok {
		return rm
	}
	return 0
}
