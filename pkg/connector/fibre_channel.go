package connector

import (
	"fmt"
	"github.com/peter-wangxu/goock/pkg/exec"
	"github.com/peter-wangxu/goock/pkg/linux"
	"github.com/peter-wangxu/goock/pkg/model"
	goockutil "github.com/peter-wangxu/goock/pkg/util"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	// Pattern to match the fibre channel device
	FibreChannelPathPattern = "/dev/disk/by-path/pci-%s-fc-%s-lun-%s"
)

// Connector for Fibre Channel
type FibreChannelConnector struct {
	exec exec.Interface
}

// Constructor for FibreChannelConnector
func NewFibreChannelConnector() FibreChannelInterface {
	return &FibreChannelConnector{exec: executor}
}

// Get Fibre Channel host information
func (fc *FibreChannelConnector) GetHostInfo() (HostInfo, error) {
	return GetHostInfo()
}

// Connect/Discover a FC device
func (fc *FibreChannelConnector) ConnectVolume(connectionProperty ConnectionProperty) (VolumeInfo, error) {

	var volumeInfo VolumeInfo
	hostPaths := fc.getVolumePaths(connectionProperty)

	if len(hostPaths) <= 0 {
		return volumeInfo, fmt.Errorf("unable to locate any Fibre Channel devices")
	}

	existedPath, err := goockutil.WaitForAnyPath(
		hostPaths, fc.wrapperRescanHosts(
			connectionProperty.TargetWwns, connectionProperty.TargetLun))

	if err != nil {
		log.WithError(err).Error("Unable to find any Fibre Channel devices.")
		return volumeInfo, err
	}
	lunWwn := linux.GetWWN(existedPath)
	log.Debugf("Found wwn [%s] for path %s.", lunWwn, existedPath)
	mPath := linux.FindMpathByWwn(lunWwn)

	volumeInfo.Wwn = lunWwn
	volumeInfo.MultipathId = lunWwn
	volumeInfo.Multipath = mPath
	volumeInfo.Paths, _ = goockutil.FilterPath(hostPaths)
	log.Debugf("ConnectVolume returning %s", volumeInfo)

	return volumeInfo, nil
}

// DisconnectVolume disconnect/remove an already-connected FC device
func (fc *FibreChannelConnector) DisconnectVolume(connectionProperty ConnectionProperty) error {
	return nil
}

// Extend the volume attributes when changes are made on storage side
func (fc *FibreChannelConnector) ExtendVolume(connectionProperty ConnectionProperty) error {
	return nil
}

// Get all possible fc devices from connection property
func (fc *FibreChannelConnector) getVolumePaths(connectionProperty ConnectionProperty) []string {
	pcis := fc.getPciNums()
	wwns := fc.formatWwns(connectionProperty.TargetWwns)
	combined := fc.combinePciWithWwn(pcis, wwns)

	formattedLun := FormatLuns(connectionProperty.TargetLun)[0]

	var possiblePaths []string
	for _, pciWwn := range combined {
		possiblePaths = append(possiblePaths, fmt.Sprintf(FibreChannelPathPattern, pciWwn[0], pciWwn[1], formattedLun))
	}

	return possiblePaths
}

// Try to get HBA channel and SCSI target to use as filters
// This could largely avoid unintended presence of the same target
func (fc *FibreChannelConnector) wrapperRescanHosts(wwpns []string, lunID int) func() {

	fcTargets := model.NewFibreChannelTarget()
	var connectedTargets [][]int
	for _, target := range fcTargets {
		for _, wwpn := range wwpns {
			if matched := strings.Contains(target.PortName, wwpn); matched == true {
				hct, _ := target.GetHostChannelTarget()
				connectedTargets = append(connectedTargets, append(hct, lunID))
			}
		}
	}
	log.WithFields(logrus.Fields{"Targets": connectedTargets, "lun": lunID}).Debug("Found connected targets.")
	return func() {
		linux.RescanHosts(connectedTargets, lunID)
	}
}

// Extract pci number from device path: /sys/devices/pci0000:00/0000:00:03.0/0000:05:00.1/host9/fc_host/host9
// 0000:05:00.1 is the correct pci number
func (fc *FibreChannelConnector) getPciNums() []string {
	var pciNums []string
	hbas := model.NewHBA()
	for _, hba := range hbas {
		pciNums = append(pciNums, strings.Split(hba.DevicePath, "/")[5])
	}
	return pciNums
}

// Insert "0x" before any wwns
func (fc *FibreChannelConnector) formatWwns(wwns []string) []string {
	var targets []string
	for _, wwn := range wwns {
		targets = append(targets, fmt.Sprintf("0x%s", wwn))
	}
	return targets
}

//Given one or more wwn  ports
//do the matrix math to figure out a set of pci device, wwn
//tuples that are potentially valid (they won't all be). This
//provides a search space for the device connection.
func (fc *FibreChannelConnector) combinePciWithWwn(pcis []string, wwns []string) [][2]string {
	var tuples [][2]string
	for _, pci := range pcis {
		if len(pci) >= 0 {
			for _, wwn := range wwns {
				tuples = append(tuples, [2]string{pci, wwn})
			}
		}
	}
	return tuples
}
