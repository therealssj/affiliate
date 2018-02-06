package config

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestGetServerConfig(t *testing.T) {
	writeConfigFile("")
	defer removeConfigFile()
	if GetServerConfig() == nil {
		t.Errorf("Failed. got nil")
	}
}

func TestGetApiForTellerConfig(t *testing.T) {
	initApiForTellerConfig = false
	writeConfigFile("")
	defer removeConfigFile()
	if GetApiForTellerConfig() == nil {
		t.Errorf("Failed. got nil")
	}
	rc := GetApiForTellerConfig().RewardConfig
	if len(rc.LadderLine) != 1 || rc.LadderLine[0] != 0 {
		t.Errorf("Failed. LadderLine wrong")
	}
	if len(rc.PromoterRatio) != 1 || rc.PromoterRatio[0] != 0.05 {
		t.Errorf("Failed. PromoterRatio wrong")
	}
	if len(rc.SuperiorPromoterRatio) != 1 || rc.SuperiorPromoterRatio[0] != 0.03 {
		t.Errorf("Failed. SuperiorPromoterRatio wrong")
	}
}

func TestGetApiForTellerConfig2(t *testing.T) {
	initApiForTellerConfig = false
	writeConfigFile(`[RewardConfig]
LadderLine      =[0,1000,10000]
PromoterRatio   =[0.05,0.06,0.07]
SuperiorPromoterRatio =[0.03,0.04,0.05]`)
	defer removeConfigFile()
	if GetApiForTellerConfig() == nil {
		t.Errorf("Failed. got nil")
	}
	rc := GetApiForTellerConfig().RewardConfig
	if len(rc.LadderLine) != 3 || rc.LadderLine[0] != 0 || rc.LadderLine[1] != 1000 || rc.LadderLine[2] != 10000 {
		t.Errorf("Failed. LadderLine wrong")
	}
	if len(rc.PromoterRatio) != 3 || rc.PromoterRatio[0] != 0.05 || rc.PromoterRatio[1] != 0.06 || rc.PromoterRatio[2] != 0.07 {
		t.Errorf("Failed. PromoterRatio wrong")
	}
	if len(rc.SuperiorPromoterRatio) != 3 || rc.SuperiorPromoterRatio[0] != 0.03 || rc.SuperiorPromoterRatio[1] != 0.04 || rc.SuperiorPromoterRatio[2] != 0.05 {
		t.Errorf("Failed. SuperiorPromoterRatio wrong")
	}
}

const config_filename = "config.toml"

func removeConfigFile() {
	err := os.Remove(config_filename)
	if err != nil {
		fmt.Printf("%s", err)
	}
}
func writeConfigFile(content string) {
	outputFile, err := os.OpenFile(config_filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	outputWriter.WriteString(content)
	outputWriter.Flush()
}
