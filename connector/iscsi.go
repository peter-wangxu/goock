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
	"github.com/peter-wangxu/goock/exec"
	"github.com/peter-wangxu/goock/linux"
	"github.com/peter-wangxu/goock/model"
	goockutil "github.com/peter-wangxu/goock/util"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

var executor = exec.New()

const (
	ISCSIPathPattern = "/dev/disk/by-path/ip-%s-iscsi-%s-lun-%d"
)

func SetExecutor(e exec.Interface) {
	executor = e
}

type ISCSIConnector struct {
	exec exec.Interface
}

func New() Interface {
	return &ISCSIConnector{exec: executor}
}

func (iscsi *ISCSIConnector) GetHostInfo(args []string) (HostInfo, error) {
	filePath := "/etc/iscsi/initiatorname.iscsi"
	cmd := iscsi.exec.Command("cat", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Log waring
		log.WithError(err).Debugf("Unable to fetch iscsi iqn under %s, iscsi is not installed?", filePath)
		return HostInfo{}, err
	}
	var info HostInfo
	pattern, err := regexp.Compile("InitiatorName=(?P<name>.*)\n$")
	matches := pattern.FindStringSubmatch(string(out))
	if len(matches) >= 2 {
		info.Initiator = matches[1]
	}
	osName := runtime.GOOS
	hostName, _ := os.Hostname()
	info.OSType = osName
	info.Hostname = hostName
	return info, err

}
func (iscsi *ISCSIConnector) getSearchPath() string {
	return "/dev/disk/by-path/"
}

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
	target_luns := connectionProperty.TargetLuns
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
// TODO(peter) set the session/node to "automatic"
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
		output, errLogin := iscsi.exec.Command("iscsiadm", "-m", "node", "-T", targetIqn, "-p", targetPortal, "--login").CombinedOutput()
		log.Debug("Login target with output: ", output)
		err = errLogin
	}
	return err
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

	//targetIqns := connectionProperty.TargetIqns
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

	paths := iscsi.getVolumePaths(connectionProperty)
	for _, path := range paths {
		linux.ExtendDevice(path)
	}
	mpathId := linux.GetWWN(paths[0])
	linux.ResizeMpath(mpathId)
	return nil
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
		for _, newSession := range discovered {
			iscsi.LoginPortal(newSession.TargetPortal, newSession.TargetIqn)
		}

	}
	iscsi.rescanISCSI()
	info := VolumeInfo{}
	possiblePaths := iscsi.getVolumePaths(connectionProperty)
	accessiblePath, err := goockutil.WaitForAnyPath(possiblePaths)
	if err != nil {
		log.WithError(err).Errorf("Unable to find any existing path within %s", possiblePaths)
		return info, err
	}
	wwn := linux.GetWWN(accessiblePath)
	if linux.IsMultipathEnabled() == true {
		// for multipath, returns the multipath descriptor
		log.Info("Multipath discovery for iSCSI enabled.")
		mPath := linux.FindMpathByWwn(wwn)
		info.Wwn = wwn
		info.Multipath = mPath
		info.Paths = possiblePaths
		if connectionProperty.AccessMode == READWRITE {
			log.Debugf("Checing to see if multipath %s is writable.", mPath)
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
	log.Debug("ConnectVolume returning %s", info)
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
