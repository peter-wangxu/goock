package goock

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/peter-wangxu/goock/connector"
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

func SetExecutor(e exec.Interface) {
	executor = e
}

type ISCSIConnector struct {
	exec exec.Interface
}

func New() *ISCSIConnector {
	return &ISCSIConnector{
		exec: executor}
}

func (iscsi *ISCSIConnector) GetHostInfo(args []string) (connector.HostInfo, error) {
	filePath := "/etc/iscsi/initiatorname.iscsi"
	cmd := iscsi.exec.Command("cat", filePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Log waring
		logrus.WithError(err).Debugf("Unable to fetch iscsi iqn under %s, iscsi is not installed?", filePath)
		return connector.HostInfo{}, err
	}
	var info connector.HostInfo
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
func (iscsi *ISCSIConnector) getVolumePaths(connectionProperty connector.ConnectionProperty) []string {
	target_iqns := connectionProperty.TargetIqns
	target_portals := connectionProperty.TargetPortals
	target_luns := connectionProperty.TargetLuns
	var potential_paths []string
	for i, iqn := range target_iqns {
		path := fmt.Sprintf("/dev/disk/by-path/ip-%s-iscsi-%s-lun-%d",
			target_portals[i], iqn, target_luns[i])
		potential_paths = append(potential_paths, path)
	}
	return potential_paths

}

func (iscsi *ISCSIConnector) validateIfaceTransport(transportIface string) string {
	// TODO need to support multiple transports?
	return "default"
}

// Discover all target portals
func (iscsi *ISCSIConnector) discoverISCSIPortal(targetPortal string) []model.ISCSISession {
	// Parse output like 10.64.76.253:3260,1 iqn.1992-04.com.emc:cx.fcnch097ae5ef3.h1
	iscsiSessions := model.DiscoverISCSISession([]string{targetPortal})
	return iscsiSessions

}

// Discover all target portals
func (iscsi *ISCSIConnector) discoverISCSIPortals(targetPortals []string) []model.ISCSISession {
	// Parse output like 10.64.76.253:3260,1 iqn.1992-04.com.emc:cx.fcnch097ae5ef3.h1
	iscsiSessions := model.DiscoverISCSISession(targetPortals)
	return iscsiSessions

}

// Login all target portals if needed
// TODO(peter) consider using goroutine to login concurrently?
func (iscsi *ISCSIConnector) LoginISCSIPortal(targetPortal string, targetIqn string) error {
	sessions := iscsi.getIscsiSessions()
	// If already logged in, skipped
	var loggedIn = false
	var err error
	for _, session := range sessions {
		if session.TargetIqn == targetIqn && session.TargetPortal == targetPortal {
			logrus.Debugf("Target %s, %s is already logged in. skip login.", targetPortal,
				targetIqn)
			loggedIn = true
			err = nil
			break
		}
	}
	if loggedIn != true {
		output, errLogin := iscsi.exec.Command("iscsiadm", "-m", "node", "-T", targetPortal, "-p", targetIqn, "--login").CombinedOutput()
		logrus.Debug("Login target with output: ", output)
		err = errLogin
	}
	return err
}

func (iscsi *ISCSIConnector) rescanISCSI() {
	iscsi.exec.Command("iscsiadm", "-m", "session", "--rescan").CombinedOutput()

}

// Return not logged portals for discovery
func (iscsi *ISCSIConnector) filterTargets(sessions []model.ISCSISession, connectionProperty connector.ConnectionProperty) []string {
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
func (iscsi *ISCSIConnector) ExtendVolume(connectionProperty connector.ConnectionProperty) {

	paths := iscsi.getVolumePaths(connectionProperty)
	for _, path := range paths {
		linux.ExtendDevice(path)
	}
	mpathId := linux.GetWWN(paths[0])
	linux.ResizeMpath(mpathId)
}

// Attach the volume from the remote to the local
// 1. Need to login/discover iscsi sessions from the ConnectionProperty
// 2. Get all possible paths for the targets
// 3. If multipath is enabled, return multipath_id as well, returned properties
//   ScsiWwn: <scsi wwn>
//   MultipathId: <multipath id>
//   Path: single path device description
func (iscsi *ISCSIConnector) ConnectVolume(connectionProperty connector.ConnectionProperty) (connector.DeviceInfo, error) {
	currSessions := iscsi.getIscsiSessions()
	notLogged := iscsi.filterTargets(currSessions, connectionProperty)
	if len(notLogged) > 0 {
		logrus.Debugf("Discovering the target(s) by iscsiadm...")
		discovered := iscsi.discoverISCSIPortals(notLogged)
		// login to the session as needed
		// TODO(peter) can be accelerated by goroutine?
		for _, newSession := range discovered {
			iscsi.LoginISCSIPortal(newSession.TargetPortal, newSession.TargetIqn)
		}

	}
	iscsi.rescanISCSI()
	info := connector.DeviceInfo{}
	possiblePaths := iscsi.getVolumePaths(connectionProperty)
	accessiblePath, err := goockutil.WaitForAnyPath(possiblePaths)
	if err != nil {
		logrus.WithError(err).Errorf("Unable to find any existing path in %s", possiblePaths)
		return info, err
	}
	wwn := linux.GetWWN(accessiblePath)
	if linux.IsMultipathEnabled() == true {
		// for multipath, returns the multipath descriptor
		logrus.Info("Multipath discovery for iSCSI enabled.")
		mPath := linux.FindMpathByWwn(wwn)
		info.Wwn = wwn
		info.Multipath = mPath
		info.Paths = possiblePaths
		if connectionProperty.AccessMode == connector.READWRITE {
			logrus.Debugf("Checing to see if multipath %s is writable.", mPath)
			linux.CheckReadWrite(accessiblePath, wwn)
		}
	} else {
		// for single path, returns any of the found path
		logrus.Debug("Multipath discovery for iSCSI disabled.")
		newPath, _ := goockutil.FilterPath(possiblePaths)
		info.Wwn = wwn
		info.Paths = newPath
		info.Multipath = ""
		info.MultipathId = ""

	}
	logrus.Debug("ConnectVolume returning %s", info)
	return info, nil

}

func (iscsi *ISCSIConnector) DisconnectVolume(connectProperty connector.ConnectionProperty) error {

	possiblePaths := iscsi.getVolumePaths(connectProperty)
	possiblePaths, _ = goockutil.FilterPath(possiblePaths)
	if linux.IsMultipathEnabled() {
		logrus.Info("Multipath discovery for iSCSI enabled.")
		if len(possiblePaths) > 0 {
			accessiblePath := possiblePaths[0]
			wwn := linux.GetWWN(accessiblePath)
			multipath := linux.FindMultipathByWwn(wwn)
			if multipath.Wwn == "" {
				// Sometimes the multipath is not found under specific path,
				// We need to find it out by "multipath -l"
				logrus.Info("No any multipath path found for targets.")
				return errors.New("Multipath is not found, skip the deletion.")
			}
			// First, remove the multipath descriptor
			linux.FlushPath(multipath.Wwn)
			// Secondary, remove every single path from scsi bus
			for _, single := range multipath.Paths {
				linux.RemoveSCSIDevice(single.DevNode)
			}

		} else {
			logrus.Info("No path found for targets")
		}
	} else {
		logrus.Info("Multipath discovery for iSCSI disabled.")
		for _, path := range possiblePaths {
			path, _ = filepath.EvalSymlinks(path)
			linux.RemoveSCSIDevice(path)
		}

	}
	left := goockutil.WaitForPathRemoval(possiblePaths, 10)
	if len(left) > 0 {
		logrus.Warnf("Paths still exist on system: %s.", left)
		return errors.New(fmt.Sprintf("Paths %s are not removed from system.", left))
	}
	//TODO(peter) Need to check paths and logout targets that no path exists
	return nil
}

// func (c *ISCSIConnector) _run_iscsiadm_update_discoverydb(self, connection_properties,
// func (c *ISCSIConnector) extend_volume(self, connection_properties){
// }
// func (c *ISCSIConnector) connect_volume(self, connection_properties){
// }
// func (c *ISCSIConnector) disconnect_volume(self, connection_properties, device_info){
// }
// func (c *ISCSIConnector) _disconnect_volume_iscsi(self, connection_properties){
// }
// func (c *ISCSIConnector) _munge_portal(self, target){
// }
// func (c *ISCSIConnector) _get_device_path(self, connection_properties){
// }
// func (c *ISCSIConnector) get_initiator(self){
// }
// func (c *ISCSIConnector) _run_iscsiadm(self, connection_properties, iscsi_command, **kwargs){
// }
// func (c *ISCSIConnector) _iscsiadm_update(self, connection_properties, property_key,
// func (c *ISCSIConnector) _get_target_portals_from_iscsiadm_output(self, output){
// }
// func (c *ISCSIConnector) _disconnect_volume_multipath_iscsi(self, connection_properties,
// func (c *ISCSIConnector) _connect_to_iscsi_portal(self, connection_properties){
// }
// func (c *ISCSIConnector) _disconnect_from_iscsi_portal(self, connection_properties){
// }
// func (c *ISCSIConnector) _get_iscsi_devices(self){
// }
// func (c *ISCSIConnector) _disconnect_mpath(self, connection_properties, ips_iqns){
// }
// func (c *ISCSIConnector) _get_multipath_iqns(self, multipath_devices, mpath_map){
// }
// func (c *ISCSIConnector) _get_multipath_device_map(self){
// }
// func (c *ISCSIConnector) _run_iscsi_session(self){
// }
// func (c *ISCSIConnector) _run_iscsiadm_bare(self, iscsi_command, **kwargs){
// }
// func (c *ISCSIConnector) _run_multipath(self, multipath_command, **kwargs){
// }
// func (c *ISCSIConnector) _rescan_iscsi(self){
// }
