/*
Copyright 2017 The Goock Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package linux

import (
	"fmt"
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/util"
	"strings"
)

func IsFCSupport() bool {
	err := util.IsPathExists("/sys/class/fc_host")
	if nil != err {
		return false
	}
	return true
}

func GetFCHBA() []model.HBA {
	return model.NewHBA()
}

func GetFcWwpn() []string {
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

func GetFcWwnn() []string {
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

// Do a more specific scan instead of a wildcard
func RescanHosts(allHct [][]int, lunID int) {
	for _, hct := range allHct {
		path := fmt.Sprintf("/sys/class/scsi_host/host%d/scan", hct[0])
		ScanSCSIBus(path, fmt.Sprintf("%d %d %d", hct[1], hct[2], lunID))
	}
}

func IsFCDevice(device string) bool {
	return false
}
