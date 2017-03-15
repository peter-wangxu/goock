package linux

import (
	"fmt"
	"github.com/peter-wangxu/goock/model"
	"os"
	"strings"
)

func IsFCSupport() bool {
	_, err := os.Stat("/sys/class/fc_host")
	if nil != err {
		return false
	}
	return true
}

func GetFCHBA() []model.HBA {
	return model.NewHBA().Parse()
}

func GetFCWWPN() []string {
	hbas := GetFCHBA()
	wwpns := make([]string, len(hbas))
	var index = 0
	for _, hba := range hbas {
		if hba.PortState == "Online" {
			wwpns[index] = strings.Replace(hba.PortName, "0x", "", -1)
			index++
		}
	}
	return wwpns
}

func GetFCWWNN() []string {
	hbas := GetFCHBA()
	wwnns := make([]string, len(hbas))
	var index = 0
	for _, hba := range hbas {
		if hba.PortState == "Online" {
			wwnns[index] = strings.Replace(hba.NodeName, "0x", "", -1)
			index++
		}
	}
	return wwnns
}

func RescanHosts() {
	hbas := GetFCHBA()
	for _, hba := range hbas {
		path := fmt.Sprintf("/sys/class/scsi_host/%s/scan", hba.Name)
		ScanSCSIBus(path, "")
	}
}
