package goock

import (
	"fmt"
	"github.com/peter-wangxu/goock/exec"
	"github.com/peter-wangxu/goock/linux"
	"github.com/peter-wangxu/goock/model"
	"regexp"
	"github.com/Sirupsen/logrus"
	goockutil "github.com/peter-wangxu/goock/util"
)

type StringEnum string

const (
	READWRITE StringEnum = "rw"
	READONLY StringEnum = "ro"
)

const (
	ISCSI_PROTOCOL StringEnum = "iscsi"
	FC_PROTOCOL StringEnum = "fc"
)

type ConnectionProperty struct {
	TargetIqns      []string
	TargetPortals   []string
	TargetLuns      []int
	StorageProtocol string
	AccessMode      StringEnum
}

type DeviceInfo struct {
	MultipathId string
	paths       []string
	Wwn         string
	Multipath   string
}

type Interface interface {
}

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

func (iscsi *ISCSIConnector) getConnectorProperties(args []string) (map[string]string, error) {
	file_path := "/etc/iscsi/initiatorname.iscsi"
	cmd := iscsi.exec.Command("cat", file_path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Log waring
		return map[string]string{}, err
	}
	var props map[string]string
	pattern, err := regexp.Compile("InitiatorName=(?P<name>.*)\n$")
	matches := pattern.FindStringSubmatch(string(out))
	if len(matches) >= 2 {
		props["initiator"] = matches[1]
	}

	return props, nil

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
func (iscsi *ISCSIConnector) loginISCSIPortal(targetPortal string, targetIqn string) error {
	sessions := iscsi.getIscsiSessions()
	// If already logged in, skipped
	var loggedIn = false
	var err error
	for _, session := range (sessions) {
		if (session.TargetIqn == targetIqn && session.TargetPortal == targetPortal) {
			logrus.Debugf("Target %s, %s is already logged in. skip login.", targetPortal,
				targetIqn)
			loggedIn = true
			err = nil
			break
		}
	}
	if (loggedIn != true) {
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
func (iscsi *ISCSIConnector) filterTargets(sessions []model.ISCSISession, connectionProperty ConnectionProperty) []string {
	var currPortals []string
	for _, session := range (sessions) {
		currPortals = append(currPortals, session.TargetPortal)
	}

	//targetIqns := connectionProperty.TargetIqns
	targetPortals := connectionProperty.TargetPortals
	var notLogged []string
	for _, portal := range (targetPortals) {
		if (!goockutil.Contains(portal, currPortals)) {
			notLogged = append(notLogged, portal)
		}
	}
	return notLogged
}


// Update the local kernel's size information
func (iscsi *ISCSIConnector) ExtendVolume(connectionProperty ConnectionProperty) {

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
func (iscsi *ISCSIConnector) ConnectVolume(connectionProperty ConnectionProperty) (DeviceInfo, error) {
	currSessions := iscsi.getIscsiSessions()
	notLogged := iscsi.filterTargets(currSessions, connectionProperty)
	if (len(notLogged) > 0) {
		logrus.Debugf("Discovering the target(s) by iscsiadm...")
		discovered := iscsi.discoverISCSIPortals(notLogged)
		// login to the session as needed
		// TODO(peter) can be accelerated by goroutine?
		for _, newSession := range (discovered) {
			iscsi.loginISCSIPortal(newSession.TargetPortal, newSession.TargetIqn)
		}

	}
	iscsi.rescanISCSI()
	info := DeviceInfo{}
	possiblePaths := iscsi.getVolumePaths(connectionProperty)
	accessiblePath, err := goockutil.WaitForAnyPath(possiblePaths)
	if (err != nil) {
		logrus.WithError(err).Errorf("Unable to find any existing path in %s", possiblePaths)
		return info, err
	}
	wwn := linux.GetWWN(accessiblePath)
	if (linux.IsMultipathEnabled() == true) {
		// for multipath, returns the multipath descriptor
		logrus.Info("Multipath discovery for iSCSI enabled.")
		mPath := linux.FindMpathByWwn(wwn)
		info.Wwn = wwn
		info.Multipath = mPath
		info.paths = possiblePaths
		if (connectionProperty.AccessMode == READWRITE) {
			logrus.Debugf("Checing to see if multipath %s is writable.", mPath)
			linux.CheckReadWrite(accessiblePath, wwn)
		}
	} else {
		// for single path, returns any of the found path
		logrus.Debug("Multipath discovery for iSCSI disabled.")
		newPath, _ := goockutil.FilterPath(possiblePaths)
		info.Wwn = wwn
		info.paths = newPath
		info.Multipath = ""
		info.MultipathId = ""

	}
	logrus.Debug("ConnectVolume returning %s", info)
	return info, nil

}

func (iscsi *ISCSIConnector) DisconnectVolume(connectProperty ConnectionProperty) {

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
