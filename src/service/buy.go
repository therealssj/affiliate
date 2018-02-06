package service

import (
	"github.com/spolabs/affiliate/src/service/db"
	pg "github.com/spolabs/affiliate/src/service/postgresql"
	client "github.com/spolabs/affiliate/src/teller_client"
)

func AllCryptocurrencyMap() map[string]db.CryptocurrencyInfo {
	slice := AllCryptocurrency()
	m := make(map[string]db.CryptocurrencyInfo, len(slice))
	for _, info := range slice {
		m[info.ShortName] = info
	}
	return m
}
func GetCryptocurrency(currencyType string) *db.CryptocurrencyInfo {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	info := pg.GetCryptocurrency(tx, currencyType)
	checkErr(tx.Commit())
	commit = true
	return info
}
func AllCryptocurrency() []db.CryptocurrencyInfo {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	all := pg.AllCryptocurrency(tx)
	checkErr(tx.Commit())
	commit = true
	return all
}

func MappingDepositAddr(address string, currencyType string, ref string) (string, bool, error) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	buyAddrMapping := pg.QueryMappingDepositAddr(tx, address, currencyType)
	if buyAddrMapping != nil {
		checkErr(tx.Commit())
		commit = true
		return buyAddrMapping.DepositAddr, false, nil
	}
	depositAddr, err := client.Bind(currencyType, address)
	if err != nil {
		return "", true, err
	}
	pg.SaveDepositAddrMapping(tx, address, currencyType, ref, depositAddr)
	checkErr(tx.Commit())
	commit = true
	return depositAddr, true, nil
}

func CheckMappingAddr(address string, currencyType string) bool {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	mapping := pg.QueryMappingDepositAddr(tx, address, currencyType)
	checkErr(tx.Commit())
	commit = true
	return mapping != nil
}

func QueryDepositRecord(address string, currencyType string) []db.DepositRecord {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	res := pg.QueryDepositRecord(tx, address, currencyType)
	checkErr(tx.Commit())
	commit = true
	return res
}

func SyncCryptocurrency(slice []db.CryptocurrencyInfo) {
	currencyMap := AllCryptocurrencyMap()
	newCurrency := make([]db.CryptocurrencyInfo, 0, 4)
	updateCur := make([]db.CryptocurrencyInfo, 0, 4)
	for _, info := range slice {
		if old, ok := currencyMap[info.ShortName]; ok {
			if old.Rate != info.Rate || old.Enabled != info.Enabled {
				updateCur = append(updateCur, info)
			}
		} else {
			newCurrency = append(newCurrency, info)
		}
	}
	syncCryptocurrency(newCurrency, updateCur)
}

func syncCryptocurrency(newCurrency []db.CryptocurrencyInfo, updateCur []db.CryptocurrencyInfo) {
	tx, commit := db.BeginTx()
	defer db.Rollback(tx, &commit)
	pg.AddBatchCryptocurrency(tx, newCurrency)
	pg.UpdateBatchRateAndEnabled(tx, updateCur)
	checkErr(tx.Commit())
	commit = true
}
