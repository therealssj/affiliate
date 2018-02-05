package service

import (
	"encoding/json"
	//	"fmt"
	"testing"

	"github.com/spolabs/affiliate/src/config"
	"github.com/spolabs/affiliate/src/service/db"
)

func TestBuildRewardRecord(t *testing.T) {
	conf := config.GetDaemonConfig()
	dbo := db.OpenDb(&conf.Db)
	defer dbo.Close()
	tx, _ := dbo.Begin()
	defer tx.Rollback()
	rewardConfig := config.GetDaemonConfig().RewardConfig
	rewardRecords := make([]db.RewardRecord, 0, 2)
	remainMap := make(map[string]uint64, 8)
	changedRemainMap := make(map[string]uint64, 8)
	jsonStr := `{
                "seq": 1,
                "update_at": 1514726238,
                "address": "2g4AGfrjU91cQHAQhWZLnc9DmDDfmPVR8o9",
                "deposit_address": "2do3K1YLMy3Aq6EcPMdncEurP5BfAUdFPJj",
                "txid": "2a6ccb1dbfd9ed65bcd74ee7a2c7877d110812ba7cf496cbcdb034bd671e1490",
                "sent": 10000000,
                "rate": 100,
                "coin_type": "skycoin",
                "deposit_value": 100000,
                "height": 7455
            }`
	dr := new(db.DepositRecord)
	err := json.Unmarshal([]byte(jsonStr), &dr)
	if err != nil {
		panic(err)
	}
	dr.MappingId, dr.RefAddr, dr.SuperiorRefAddr = 110, "2Fo3oqg4c3ugewzi6ZUrFnAmp4AFLQ4bucQ", "mXY1Tu3Gb4tQ1wHUmbFUbbGQmBSZeuzFJc"
	if len(dr.RefAddr) > 0 {
		rewardRecords = appendBuyerRewardRecord(tx, rewardRecords, dr, &rewardConfig, remainMap, changedRemainMap)
		rewardRecords = appendPromoterRewardRecord(tx, rewardRecords, dr, &rewardConfig, remainMap, changedRemainMap)
		if len(dr.SuperiorRefAddr) > 0 {
			rewardRecords = appendSuperiorPromoterRewardRecord(tx, rewardRecords, dr, &rewardConfig, remainMap, changedRemainMap)
		}
	}
	if len(rewardRecords) != 3 {
		t.Errorf("Failed. Got len(rewardRecords)=%d, expected %d", len(rewardRecords), 3)
	}
	if len(changedRemainMap) != 3 {
		t.Errorf("Failed. Got len(changedRemainMap)=%d, expected %d", len(changedRemainMap), 3)
	}
	//	fmt.Println(rewardRecords)
	//	fmt.Println(changedRemainMap)
	rewardConfig.BuyerRatio = 0.2
	rewardConfig.PromoterRatio = []float64{0.5, 0.7}
	rewardConfig.SuperiorPromoterRatio = []float64{0.3, 0.5}
	jsonStr = `{
		              "seq": 2,
		              "update_at": 1514726616,
		              "address": "2g4AGfrjU91cQHAQhWZLnc9DmDDfmPVR8o9",
		              "deposit_address": "2do3K1YLMy3Aq6EcPMdncEurP5BfAUdFPJj",
		              "txid": "fa92485d739e64e55f7a4beab9f5d7e6e23aa6f7e289bd5a5e7597f5e1fa4cf9",
		              "sent": 20000000,
		              "rate": 100,
		              "coin_type": "skycoin",
		              "deposit_value": 200000,
		              "height": 7456
		          }`
	dr = new(db.DepositRecord)
	err = json.Unmarshal([]byte(jsonStr), &dr)
	if err != nil {
		panic(err)
	}
	dr.MappingId, dr.RefAddr, dr.SuperiorRefAddr = 110, "2Fo3oqg4c3ugewzi6ZUrFnAmp4AFLQ4bucQ", "mXY1Tu3Gb4tQ1wHUmbFUbbGQmBSZeuzFJc"
	if len(dr.RefAddr) > 0 {
		rewardRecords = appendBuyerRewardRecord(tx, rewardRecords, dr, &rewardConfig, remainMap, changedRemainMap)
		rewardRecords = appendPromoterRewardRecord(tx, rewardRecords, dr, &rewardConfig, remainMap, changedRemainMap)
		if len(dr.SuperiorRefAddr) > 0 {
			rewardRecords = appendSuperiorPromoterRewardRecord(tx, rewardRecords, dr, &rewardConfig, remainMap, changedRemainMap)
		}
	}
	if len(rewardRecords) != 6 {
		t.Errorf("Failed. Got len(rewardRecords)=%d, expected %d", len(rewardRecords), 6)
	}
	if len(changedRemainMap) != 3 {
		t.Errorf("Failed. Got len(changedRemainMap)=%d, expected %d", len(changedRemainMap), 3)
	}
	//	fmt.Println(rewardRecords)
	//	fmt.Println(changedRemainMap)
	rewardConfig.BuyerRatio = 0
	rewardConfig.PromoterRatio = []float64{0, 0}
	rewardConfig.SuperiorPromoterRatio = []float64{0, 0}
	if len(dr.RefAddr) > 0 {
		rewardRecords = appendBuyerRewardRecord(tx, rewardRecords, dr, &rewardConfig, remainMap, changedRemainMap)
		rewardRecords = appendPromoterRewardRecord(tx, rewardRecords, dr, &rewardConfig, remainMap, changedRemainMap)
		if len(dr.SuperiorRefAddr) > 0 {
			rewardRecords = appendSuperiorPromoterRewardRecord(tx, rewardRecords, dr, &rewardConfig, remainMap, changedRemainMap)
		}
	}
	if len(rewardRecords) != 6 {
		t.Errorf("Failed. Got len(rewardRecords)=%d, expected %d", len(rewardRecords), 6)
	}
	if len(changedRemainMap) != 3 {
		t.Errorf("Failed. Got len(changedRemainMap)=%d, expected %d", len(changedRemainMap), 3)
	}
}

func TestGetPromoterRatioBySalesVolume(t *testing.T) {
	rewardConfig := config.RewardConfig{}
	rewardConfig.BuyerRatio = 0.2
	rewardConfig.LadderLine = []int{0, 1000}
	rewardConfig.PromoterRatio = []float64{0.5, 0.7}
	rewardConfig.SuperiorPromoterRatio = []float64{0.3, 0.5}
	pr, spr := getPromoterRatioBySalesVolume(&rewardConfig, 0)
	if pr != 0.5 || spr != 0.3 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 999)
	if pr != 0.5 || spr != 0.3 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 1000)
	if pr != 0.7 || spr != 0.5 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 1001)
	if pr != 0.7 || spr != 0.5 {
		t.Errorf("Failed.")
	}
	rewardConfig.LadderLine = []int{0}
	rewardConfig.PromoterRatio = []float64{0.2}
	rewardConfig.SuperiorPromoterRatio = []float64{0.1}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 0)
	if pr != 0.2 || spr != 0.1 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 999)
	if pr != 0.2 || spr != 0.1 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	rewardConfig.LadderLine = []int{0, 100, 10000}
	rewardConfig.PromoterRatio = []float64{0.2, 0.4, 0.6}
	rewardConfig.SuperiorPromoterRatio = []float64{0.1, 0.3, 0.5}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 0)
	if pr != 0.2 || spr != 0.1 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 1)
	if pr != 0.2 || spr != 0.1 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 99)
	if pr != 0.2 || spr != 0.1 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 100)
	if pr != 0.4 || spr != 0.3 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 101)
	if pr != 0.4 || spr != 0.3 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 9999)
	if pr != 0.4 || spr != 0.3 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 10000)
	if pr != 0.6 || spr != 0.5 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 10001)
	if pr != 0.6 || spr != 0.5 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
	rewardConfig.LadderLine = []int{10, 100, 10000}
	rewardConfig.PromoterRatio = []float64{0.2, 0.4, 0.6}
	rewardConfig.SuperiorPromoterRatio = []float64{0.1, 0.3, 0.5}
	pr, spr = getPromoterRatioBySalesVolume(&rewardConfig, 5)
	if pr != 0.2 || spr != 0.1 {
		t.Errorf("Failed.pr:%g,spr:%g", pr, spr)
	}
}
