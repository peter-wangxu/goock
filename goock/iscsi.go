package iscsi

import (
	"exec"
	"fmt"
	"regexp"
	"strings"
)

type ConnectionProperty struct {
	target_iqns     []string
	target_portals  []string
	target_luns     []int
	storge_protocal string
}

const (
	iscsiadm = "iscsiadm"
)

type Interface interface {
}

type ISCSIConnector struct {
	exec exec.Interface
}

func New(exec exec.Interface) {
	executor := exec.New()
	return &ISCSIConnector{
		exec: executor}
}

func (c *ISCSIConnector) getConnectorProperties(args []string) (string, error) {
	file_path := "/etc/iscsi/initiatorname.iscsi"
	cmd := c.exec.Cmd("cat", file_path)
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
func (c *ISCSIConnector) getSearchPath() string {
	return "/dev/disk/by-path/"
}

func (c *ISCSIConnector) getIscsiSessions() []map[string]string {
	cmd := c.exec.Command(iscsiadm, "-m", "session")
	out, err := cmd.Output()
	// parse the output from iscsiadm
	// lines are in the format of
	// tcp: [1] 192.168.121.250:3260,1 iqn.2010-10.org.openstack:volume-
	if err != nil {
		// log warning
	}
	var session_maps []map[string]string
	sessions := strings.Fileds(out)
	for i, s := range sessions {
		se := strings.Fields(s)
		m = map[string]string{
			"target_portal": strings.Split(se[2], ",")[0],
			"target_iqn":    se[3],
		}
		append(session_maps, m)
	}
	return session_maps
}

func (c *ISCSIConnector) getVolumePaths(connection_properties ConnectionProperty) []string {
	target_iqns := connection_properties["target_iqns"]
	target_portals := connection_properties["target_luns"]
	target_luns := connection_properties["target_luns"]
	var potentail_paths []string
	for i, iqn := range target_iqns {
		path := fmt.Sprintf("/dev/disk/by-path/ip-%s-iscsi-%s-lun-%s",
			target_portals[i], iqn, target_luns[i])
		append(potentail_paths, path)
	}
	return potentail_paths

}

// func (c *ISCSIConnector) set_execute(self, execute){
// }
// func (c *ISCSIConnector) _validate_iface_transport(self, transport_iface){
// }
// func (c *ISCSIConnector) _get_transport(self){
// }
// func (c *ISCSIConnector) _discover_iscsi_portals(self, connection_properties){
// }
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
