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
package connector

import (
	"errors"
	"fmt"
	"github.com/peter-wangxu/goock/pkg/exec"
	"github.com/peter-wangxu/goock/pkg/linux"
	"github.com/peter-wangxu/goock/pkg/model"
	goockutil "github.com/peter-wangxu/goock/pkg/util"
	"path/filepath"
)

const (
	ISCSIPathPattern = "/dev/disk/by-path/ip-%s-iscsi-%s-lun-%s"
)

type OPERATION_ENUM StringEnum

const (
	OperationNew           OPERATION_ENUM = "new"
	OperationDelete        OPERATION_ENUM = "new"
	OperationUpdate        OPERATION_ENUM = "update"
	OperationShow          OPERATION_ENUM = "show"
	OperationNonPersistent OPERATION_ENUM = "nonpersistent"
)

type ISCSIConnector struct {
	exec exec.Interface
}

func NewISCSIConnector() ISCSIInterface {
	return &ISCSIConnector{exec: executor}
}

// Returns host information regarding iSCSI and FC
func (iscsi *ISCSIConnector) GetHostInfo() (HostInfo, error) {
	return GetHostInfo()
}

// Get all logged-in sessions
func (iscsi *ISCSIConnector) getIscsiSessions() []model.ISCSISession {
	// parse the output from iscsiadm
	// lines are in the format of
	// tcp: [1] 192.168.121.250:3260,1 iqn.2010-10.org.openstack:volume-
	iscsiSession := model.NewISCSISession()
	return iscsiSession

}

// Get all possible volume paths under /dev/disk/by-path/
func (iscsi *ISCSIConnector) getVolumePaths(connectionProperty ConnectionProperty) []string {
	target_iqns := connectionProperty.TargetIqns
	target_portals := connectionProperty.TargetPortals
	target_luns := FormatLuns(connectionProperty.TargetLuns...)
	var potential_paths []string
	for i, iqn := range target_iqns {
		path := fmt.Sprintf(ISCSIPathPattern, target_portals[i], iqn, target_luns[i])
		potential_paths = append(potential_paths, path)
	}
	return potential_paths

}

func (iscsi *ISCSIConnector) validateIfaceTransport(transportIface string) string {
	// TODO need to support multiple transports?
	return "default"
}

// Discover all target portals
func (iscsi *ISCSIConnector) DiscoverPortal(targetPortal ...string) []model.ISCSISession {
	// Parse output like 10.64.76.253:3260,1 iqn.1992-04.com.emc:cx.fcnch097ae5ef3.h1
	iscsiSessions := model.DiscoverISCSISession(targetPortal)
	return iscsiSessions

}

// Login all target portals if needed
// TODO(peter) consider using goroutine to login concurrently?
func (iscsi *ISCSIConnector) LoginPortal(targetPortal string, targetIqn string) error {
	sessions := iscsi.getIscsiSessions()
	// If already logged in, skipped
	var loggedIn = false
	var err error
	for _, session := range sessions {
		if session.TargetIqn == targetIqn && session.TargetPortal == targetPortal {
			log.Debugf("Target %s, %s is already logged in. skip login.", targetPortal,
				targetIqn)
			loggedIn = true
			err = nil
			break
		}
	}
	if loggedIn != true {
		_, errLogin := iscsi.exec.Command("iscsiadm", "-m", "node", "-p", targetPortal, "-T", targetIqn, "--login").CombinedOutput()
		err = errLogin
	}
	return err
}

// Set the node to 'node.startup = automatic', it will login the portal
// automatically after reboot
func (iscsi *ISCSIConnector) SetNode2Auto(targetPortal string, targetIqn string) error {
	operations := iscsi.composeISCSIOperation(targetPortal, targetIqn, OperationUpdate, "node.startup", "automatic")
	_, err := iscsi.exec.Command("iscsiadm", operations...).CombinedOutput()
	return err
}

func (iscsi *ISCSIConnector) composeISCSIOperation(targetPortal string, targetIqn string,
	operation OPERATION_ENUM, key string, value string) []string {
	return []string{
		"-m", "node", "-p", targetPortal, "-T", targetIqn,
		"--op", string(operation), "-n", key, "-v", value,
	}
}

func (iscsi *ISCSIConnector) rescanISCSI() {
	iscsi.exec.Command("iscsiadm", "-m", "session", "--rescan").CombinedOutput()

}

// Return not logged portals for discovery
func (iscsi *ISCSIConnector) filterTargets(sessions []model.ISCSISession, connectionProperty ConnectionProperty) []string {
	var currPortals []string
	for _, session := range sessions {
		currPortals = append(currPortals, session.TargetPortal)
	}

	targetPortals := connectionProperty.TargetPortals
	var notLogged []string
	for _, portal := range targetPortals {
		if !goockutil.Contains(portal, currPortals) {
			notLogged = append(notLogged, portal)
		}
	}
	return notLogged
}

// Update the local kernel's size information
func (iscsi *ISCSIConnector) ExtendVolume(connectionProperty ConnectionProperty) error {
	var err error
	paths := iscsi.getVolumePaths(connectionProperty)
	paths, _ = goockutil.FilterPath(paths)

	if len(paths) > 0 {
		// Flush size of each single path
		for _, path := range paths {
			linux.ExtendDevice(path)
		}
		// Flush size for multipath descriptor
		mpathId := linux.GetWWN(paths[0])
		err = linux.ResizeMpath(mpathId)
	} else {
		err = fmt.Errorf("Unable to find any path to extend.")
	}
	return err
}

// Attach the volume from the remote to the local
// 1. Need to login/discover iscsi sessions from the ConnectionProperty
// 2. Get all possible paths for the targets
// 3. If multipath is enabled, return multipath_id as well, returned properties
//   ScsiWwn: <scsi wwn>
//   MultipathId: <multipath id>
//   Path: single path device description
func (iscsi *ISCSIConnector) ConnectVolume(connectionProperty ConnectionProperty) (VolumeInfo, error) {
	currSessions := iscsi.getIscsiSessions()
	notLogged := iscsi.filterTargets(currSessions, connectionProperty)
	if len(notLogged) > 0 {
		log.Debugf("Discovering the target(s) by iscsiadm...")
		discovered := iscsi.DiscoverPortal(notLogged...)
		// login to the session as needed
		// TODO(peter) can be accelerated by goroutine?
		// but the os-brick says parallel login can crash open-iscsi
		for _, newSession := range discovered {
			iscsi.LoginPortal(newSession.TargetPortal, newSession.TargetIqn)
			iscsi.SetNode2Auto(newSession.TargetPortal, newSession.TargetIqn)

		}

	}
	iscsi.rescanISCSI()
	info := VolumeInfo{}
	possiblePaths := iscsi.getVolumePaths(connectionProperty)
	accessiblePath, err := goockutil.WaitForAnyPath(possiblePaths, nil)
	if err != nil {
		log.WithError(err).Errorf("Unable to find any existing path within %s", possiblePaths)
		return info, err
	}
	wwn := linux.GetWWN(accessiblePath)
	log.Debugf("Found wwn [%s] for path %s.", wwn, accessiblePath)
	if linux.IsMultipathEnabled() == true {
		// for multipath, returns the multipath descriptor
		log.Info("Multipath discovery for iSCSI enabled.")
		mPath := linux.FindMpathByWwn(wwn)
		info.Wwn = wwn
		info.MultipathId = wwn
		info.Multipath = mPath
		info.Paths, _ = goockutil.FilterPath(possiblePaths)
		if connectionProperty.AccessMode == ReadWrite {
			log.Debugf("Checking to see if multipath %s is writable.", mPath)
			linux.CheckReadWrite(accessiblePath, wwn)
		}
	} else {
		// for single path, returns any of the found path
		log.Debug("Multipath discovery for iSCSI disabled.")
		newPath, _ := goockutil.FilterPath(possiblePaths)
		info.Wwn = wwn
		info.Paths = newPath
		info.Multipath = ""
		info.MultipathId = ""

	}
	log.Debugf("ConnectVolume returning %s", info)
	return info, nil

}

func (iscsi *ISCSIConnector) DisconnectVolume(connectProperty ConnectionProperty) error {

	possiblePaths := iscsi.getVolumePaths(connectProperty)
	possiblePaths, _ = goockutil.FilterPath(possiblePaths)
	if linux.IsMultipathEnabled() {
		log.Info("Multipath discovery for iSCSI enabled.")
		if len(possiblePaths) > 0 {
			accessiblePath := possiblePaths[0]
			wwn := linux.GetWWN(accessiblePath)
			multipath := linux.FindMultipathByWwn(wwn)
			if multipath.Wwn == "" {
				// Sometimes the multipath is not found under specific path,
				// We need to find it out by "multipath -l"
				log.Info("No any multipath path found for targets.")
				return errors.New("Multipath is not found, skip the deletion.")
			}
			// First, remove the multipath descriptor
			linux.FlushPath(multipath.Wwn)
			// Secondary, remove every single path from scsi bus
			for _, single := range multipath.Paths {
				linux.RemoveSCSIDevice(single.DevNode)
			}

		} else {
			log.Info("No path found for targets")
		}
	} else {
		log.Info("Multipath discovery for iSCSI disabled.")
		for _, path := range possiblePaths {
			path, _ = filepath.EvalSymlinks(path)
			linux.RemoveSCSIDevice(path)
		}

	}
	left := goockutil.WaitForPathRemoval(possiblePaths, 10)
	if len(left) > 0 {
		log.Warnf("Paths still exist on system: %s.", left)
		return errors.New(fmt.Sprintf("Paths %s are not removed from system.", left))
	}
	//TODO(peter) Need to check paths and logout targets that no path exists
	return nil
}
