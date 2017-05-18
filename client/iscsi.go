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

package client

import (
	"fmt"
	"github.com/peter-wangxu/goock/connector"
	"github.com/peter-wangxu/goock/model"
	"regexp"
	"strconv"
	"strings"
)

var iscsiConnector = connector.New()

func SetISCSIConnector(iscsi connector.Interface) {
	iscsiConnector = iscsi
}

func Session2ConnectionProperty(sessions []model.ISCSISession, lun int) connector.ConnectionProperty {
	conn := connector.ConnectionProperty{}
	var portals []string
	var iqns []string
	var lunIds []int
	for _, session := range sessions {
		portals = append(portals, session.TargetPortal)
		iqns = append(iqns, session.TargetIqn)
		lunIds = append(lunIds, lun)
	}
	conn.TargetIqns = iqns
	conn.TargetPortals = portals
	conn.TargetLuns = lunIds
	return conn
}

func HandleISCSIDisconnect(args ...string) error {
	return nil
}

func HandleISCSIConnect(args ...string) error {
	var err error
	if len(args) <= 0 {
		log.Error("Target IP is required.")
		err = fmt.Errorf("Target IP is required.")
	} else if len(args) == 1 {
		log.Info("LUN ID is not specified, will query all LUNs on target IP: %s", args[0])
		targetIP := args[0]
		// TODO Need to login and find all possible LUN IDs
		volumeInfo, conErr := FetchVolumeInfo(targetIP, 4)
		if conErr != nil {
			err = conErr
		} else {
			BeautifyVolumeInfo(volumeInfo)
		}
	} else {
		log.Debugf("Trying to validate the target IP : %s, LUN ID: %s", args[0], args[1:])
		re, _ := regexp.Compile("\\d+")
		for _, lun := range args[1:] {
			if re.MatchString(lun) == false {
				err = fmt.Errorf("%s does not look like a LUN ID.", lun)
				break
			}
		}
		if err == nil {
			targetIP := args[0]
			for _, lun := range args[1:] {
				i, _ := strconv.Atoi(lun)
				volumeInfo, _ := FetchVolumeInfo(targetIP, i)
				BeautifyVolumeInfo(volumeInfo)
			}
		}
	}

	return err
}

func FetchVolumeInfo(targetIP string, lun int) (connector.VolumeInfo, error) {
	sessions := iscsiConnector.DiscoverPortal(targetIP)
	connectionProperty := Session2ConnectionProperty(sessions, lun)
	return iscsiConnector.ConnectVolume(connectionProperty)

}

var VOLUME_FORMAT = `Volume Information:
Multipath    : %s
Single paths :
%s
Multipath ID : %s
WWN          : %s
`

func BeautifyVolumeInfo(info connector.VolumeInfo) {
	var beautiPath []string
	for _, path := range info.Paths {
		path = "  [*] " + path
		beautiPath = append(beautiPath, path)
	}
	fmt.Printf(fmt.Sprintf(VOLUME_FORMAT, info.Multipath, strings.Join(beautiPath, "\n"),
		info.MultipathId, info.Wwn))
}
