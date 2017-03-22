package linux

import (
	"os"
	"github.com/peter-wangxu/goock/model"
	"strings"
	"fmt"
)


func IsFCSupport() bool {
	_, err := os.Stat("/sys/class/fc_host")
	if(nil != err){
		return false
	}
	return true
}

func GetFCHBA() []model.HBA {
	return model.NewHBA().Parse()
}


func GetFCWWPN() []string{
	hbas := GetFCHBA()
	var wwpns []string
	var index = 0
	for _, hba := range hbas{
		if hba.PortState == "Online" {
			wwpns[index] = strings.Replace(hba.PortName, "0x", "", -1)
			index++
		}
	}
	return wwpns
}


func GetFCWWNN() []string{
	hbas := GetFCHBA()
	var wwpns []string
	var index = 0
	for _, hba := range hbas{
		if hba.PortState == "Online" {
			wwpns[index] = strings.Replace(hba.NodeName, "0x", "", -1)
			index++
		}
	}
	return wwpns
}

func RescanHosts() {
	hbas := GetFCHBA()
	for _, hba := range hbas{
		path := fmt.Sprintf("/sys/class/scsi_host/%s/scan", hba.Name)
		ScanSCSIBus(path, "")
	}
}
