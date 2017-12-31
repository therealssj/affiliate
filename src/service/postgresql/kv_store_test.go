package postgresql

import (
	"github.com/spaco/affiliate/src/config"
	"github.com/spaco/affiliate/src/service/db"
	"testing"
)

func TestKvStore(t *testing.T) {
	config := config.GetServerConfig()
	dbo := db.OpenDb(&config.Db)
	defer dbo.Close()
	tx, _ := dbo.Begin()
	defer tx.Rollback()
	var (
		key          = "testkey"
		intVal int64 = 12312
		strVal       = ""
	)
	SaveKvStore(tx, key, intVal, strVal)
	i, s, found := GetKvStore(tx, key)
	if !found {
		t.Errorf("Failed. not found")
	}
	if i != intVal {
		t.Errorf("Failed. expect:%d, actual:%d", intVal, i)
	}
	if s != strVal {
		t.Errorf("Failed. expect:%s, actual:%s", strVal, s)
	}
	intVal = 0
	strVal = "testval"
	SaveKvStore(tx, key, intVal, strVal)
	i, s, found = GetKvStore(tx, key)
	if !found {
		t.Errorf("Failed. not found")
	}
	if i != intVal {
		t.Errorf("Failed. expect:%d, actual:%d", intVal, i)
	}
	if s != strVal {
		t.Errorf("Failed. expect:%s, actual:%s", strVal, s)
	}

}
