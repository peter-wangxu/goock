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

	"github.com/peter-wangxu/goock/pkg/connector"
	"github.com/peter-wangxu/goock/pkg/model"
)

var iscsiConnector = connector.NewISCSIConnector()

// SetISCSIConnector sets the iSCSI connector
// This will help when doing mock testing
func SetISCSIConnector(iscsi connector.ISCSIInterface) {
	iscsiConnector = iscsi
}

// Session2ConnectionProperty converts a session to an ConnectionProperty
func Session2ConnectionProperty(sessions []model.ISCSISession, lun int) connector.ConnectionProperty {
	conn := connector.ConnectionProperty{}
	var portals []string
	var iqns []string
	var lunIDs []int
	for _, session := range sessions {
		portals = append(portals, session.TargetPortal)
		iqns = append(iqns, session.TargetIqn)
		lunIDs = append(lunIDs, lun)
	}
	conn.TargetIqns = iqns
	conn.TargetPortals = portals
	conn.TargetLuns = lunIDs
	return conn
}

// HandleISCSIConnect connects the iSCSI target via iscsiadm
func HandleISCSIConnect(args ...string) error {
	var err error
	if len(args) == 1 {
		err = fmt.Errorf("currently Target IP is not supported")
		log.Error("Target IP and LUN ID(s) is required.")
		//log.Info("LUN ID is not specified, will query all LUNs on target IP: %s", args[0])
		//targetIP := args[0]
		//// TODO Need to login and find all possible LUN IDs
		//volumeInfo, conErr := FetchVolumeInfo(targetIP, 4)
		//if conErr != nil {
		//	err = conErr
		//} else {
		//	BeautifyVolumeInfo(volumeInfo)
		//}
	} else {
		log.Debugf("Trying to validate the target IP : %s, LUN ID: %s", args[0], args[1:])
		var lunIDs []int
		lunIDs, err = ValidateLunID(args[1:])
		if err == nil {
			targetIP := args[0]
			sessions := iscsiConnector.DiscoverPortal(targetIP)
			for _, lun := range lunIDs {
				volumeInfo, _ := FetchVolumeInfo(sessions, lun)
				BeautifyVolumeInfo(volumeInfo)
			}
		}
	}

	return err
}

// HandleISCSIDisconnect disconnects the volume devices from local host
// Accessible format likes follows:
//  /dev/sdb
//  sdb
//  <Target IP> <LUN ID>
func HandleISCSIDisconnect(args ...string) error {
	var err error
	if len(args) == 1 {
		// TODO Support the device name removal
		err = fmt.Errorf("currently device name is not supported")
	} else if len(args) >= 2 {

		targetIP := args[0]
		lunIDs, err := ValidateLunID(args[1:])
		if err == nil {
			sessions := iscsiConnector.DiscoverPortal(targetIP)
			for _, lun := range lunIDs {
				connectionProperty := Session2ConnectionProperty(sessions, lun)
				err = iscsiConnector.DisconnectVolume(connectionProperty)
			}
		}

	}
	if err != nil {
		log.WithError(err).Error("Unable to proceed.")
	}
	return err
}

// HandleISCSIExtend extends the iscsi block device
func HandleISCSIExtend(args ...string) error {
	targetIP := args[0]
	lunIDs, err := ValidateLunID(args[1:])

	sessions := iscsiConnector.DiscoverPortal(targetIP)
	if err == nil {
		for _, lun := range lunIDs {
			property := Session2ConnectionProperty(sessions, lun)
			iscsiConnector.ExtendVolume(property)
		}
	}
	return err
}

// FetchVolumeInfo fetches the volume information via iscsiConnector.
func FetchVolumeInfo(sessions []model.ISCSISession, lun int) (connector.VolumeInfo, error) {
	connectionProperty := Session2ConnectionProperty(sessions, lun)
	return iscsiConnector.ConnectVolume(connectionProperty)

}

// BeautifyVolumeInfo output the volume information to stdout.
func BeautifyVolumeInfo(info connector.VolumeInfo) {
	beautifiedPaths := ""
	for _, path := range info.Paths {
		beautifiedPaths += fmt.Sprintf("  %s\n", path)
	}
	fmt.Printf(fmt.Sprintf(VolumeFormat, info.Multipath, beautifiedPaths,
		info.MultipathId, info.Wwn))
}
