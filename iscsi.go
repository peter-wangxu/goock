package main

import (
	"github.com/peter-wangxu/goock/exec"
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/linux"
	"fmt"
	"regexp"
)

type ConnectionProperty struct {
	TargetIqn       []string
	TargetPortals   []string
	TargetLuns      []int
	StorageProtocol string
}

const (
	iscsiadm = "iscsiadm"
)

type Interface interface {
}

type ISCSIConnector struct {
	exec exec.Interface
}

func New() {
	executor := exec.New()
	return &ISCSIConnector{
		exec: executor}
}

func (iscsi *ISCSIConnector) getConnectorProperties(args []string) (string, error) {
	file_path := "/etc/iscsi/initiatorname.iscsi"
	cmd := iscsi.exec.Command("cat", file_path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Log waring
		return "", err
	}
	var props map[string]string
	pattern, err := regexp.Compile("InitiatorName=(?P<name>.*)\n$")
	matches := pattern.FindStringSubmatch(out)
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
	iscsiSession :=  model.NewISCSISession()
	return iscsiSession.Parse()


}

func (iscsi *ISCSIConnector) getVolumePaths(connectionProperty ConnectionProperty) []string {
	target_iqns := connectionProperty["target_iqns"]
	target_portals := connectionProperty["target_luns"]
	target_luns := connectionProperty["target_luns"]
	var potential_paths []string
	for i, iqn := range target_iqns {
		path := fmt.Sprintf("/dev/disk/by-path/ip-%s-iscsi-%s-lun-%s",
			target_portals[i], iqn, target_luns[i])
		append(potential_paths, path)
	}
	return potential_paths

}

func (iscsi *ISCSIConnector) validateIfaceTransport(transportIface string) string {
	// TODO need to support multiple transports?
	return ""
}

// Discover all target portals
func (iscsi *ISCSIConnector) discoverISCSIPortal(targetPortal string) []model.ISCSISession {
	// Parse output like 10.64.76.253:3260,1 iqn.1992-04.com.emc:cx.fcnch097ae5ef3.h1
	//iscsiSessions := model.NewISCSISession()
	return nil

}

// Update the local kernel's size information
func (iscsi *ISCSIConnector) ExtendVolume(connectionProperty ConnectionProperty) {

	paths := iscsi.getVolumePaths(connectionProperty)
	for _, path := range paths {
		// TODO extend every path via scsi command
		linux.ExtendDevice(path)
	}
	// TODO extend multipath device again
	mpathId := linux.GetWWN(paths[0])
	linux.ResizeMpath(mpathId)
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
