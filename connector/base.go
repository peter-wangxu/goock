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
	"fmt"
	"github.com/peter-wangxu/goock/exec"
	"github.com/peter-wangxu/goock/model"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"runtime"
)

type StringEnum string

const (
	ReadWrite StringEnum = "rw"
	ReadOnly  StringEnum = "ro"
)

const (
	IscsiProtocol StringEnum = "iscsi"
	FcProtocol    StringEnum = "fibre_channel"
)

type ConnectionProperty struct {
	// Only for iscsi
	TargetIqns    []string
	TargetPortals []string
	TargetLuns    []int
	// Only for fibre channel
	TargetWwns []string
	TargetLun  int
	// Shared by fibre change and iscsi
	StorageProtocol StringEnum
	AccessMode      StringEnum
}

var executor = exec.New()

func SetExecutor(e exec.Interface) {
	executor = e
}

// IsEmpty validates whether the ConnectionProperty is empty
func (prop ConnectionProperty) IsEmpty() error {
	if prop.StorageProtocol == IscsiProtocol {
		if len(prop.TargetPortals) == 0 || len(prop.TargetLuns) == 0 {
			return fmt.Errorf("An empty ConnectionProperty is specified, forget target IPs or LUN id?")
		}
	} else if prop.StorageProtocol == FcProtocol {
		if len(prop.TargetWwns) == 0 || len(prop.TargetLuns) == 0 {
			return fmt.Errorf("An empty ConnectionProperty is specified, forget target wwns or LUN id?")
		}
	} else {
		return fmt.Errorf("Unknown storage protocol specified.")
	}

	return nil
}

type HostInfo struct {
	Initiator string
	// Wwnns: the node name of the host HBA
	Wwnns []string
	// Wwpns: the port name of the HOST HBA
	Wwpns []string
	// TargetWwnns: the node name of connected targets
	TargetWwnns []string
	// TargetWwpns: the port name of connected targets
	TargetWwpns []string
	// Connected target portals
	TargetPortals []string
	TargetIqns    []string
	Ip            string
	Hostname      string
	OSType        string
}

type VolumeInfo struct {
	MultipathId string
	Paths       []string
	Wwn         string
	Multipath   string
}

// Defining these interfaces is mainly for unit testing
// Any caller of ISCSIConnector can implement this interface for testing purpose

type Interface interface {
	GetHostInfo() (HostInfo, error)
	ConnectVolume(connectionProperty ConnectionProperty) (VolumeInfo, error)
	DisconnectVolume(connectionProperty ConnectionProperty) error
	ExtendVolume(connectionProperty ConnectionProperty) error
}

type ISCSIInterface interface {
	GetHostInfo() (HostInfo, error)
	ConnectVolume(connectionProperty ConnectionProperty) (VolumeInfo, error)
	DisconnectVolume(connectionProperty ConnectionProperty) error
	ExtendVolume(connectionProperty ConnectionProperty) error
	LoginPortal(targetPortal string, targetIqn string) error
	SetNode2Auto(targetPortal string, targetIqn string) error
	DiscoverPortal(targetPortal ...string) []model.ISCSISession
}

type FibreChannelInterface interface {
	GetHostInfo() (HostInfo, error)
	ConnectVolume(connectionProperty ConnectionProperty) (VolumeInfo, error)
	DisconnectVolume(connectionProperty ConnectionProperty) error
	ExtendVolume(connectionProperty ConnectionProperty) error
}

var log *logrus.Logger = logrus.New()

func SetLogger(l *logrus.Logger) {
	log = l
}

// Common functions

// Specific handling for LUN ID.
// For lun id < 256, the return should be as original
// For lun id >= 256, return "0x" prefixed string
func FormatLuns(luns ...int) []string {
	var formated []string
	for _, lun := range luns {
		var s string
		if lun < 256 {
			s = fmt.Sprintf("%d", lun)
		} else {
			s = fmt.Sprintf("0x%04x%04x00000000", lun&0xffff, lun>>16&0xffff)
		}
		formated = append(formated, s)
	}
	return formated
}

// GetHostInfo returns host iscsi and fc related information
func GetHostInfo() (HostInfo, error) {
	var info HostInfo

	filePath := "/etc/iscsi/initiatorname.iscsi"
	cmd := executor.Command("cat", filePath)
	out, err := cmd.CombinedOutput()
	if err == nil {
		// Log warning
		pattern, _ := regexp.Compile("InitiatorName=(?P<name>.*)\n$")
		matches := pattern.FindStringSubmatch(string(out))
		if len(matches) >= 2 {
			info.Initiator = matches[1]
		}
	} else {
		log.WithError(err).Debugf("Unable to fetch iscsi iqn under %s, permission denied or open-iscsi is not installed?", filePath)
	}
	info.OSType = runtime.GOOS
	info.Hostname, _ = os.Hostname()
	hbas := model.NewHBA()
	for _, hba := range hbas {
		info.Wwnns = append(info.Wwnns, hba.NodeName)
		info.Wwpns = append(info.Wwpns, hba.PortName)
	}
	targets := model.NewFibreChannelTarget()
	for _, target := range targets {
		info.TargetWwnns = append(info.TargetWwnns, target.NodeName)
		info.TargetWwpns = append(info.TargetWwpns, target.PortName)
	}

	sessions := model.NewISCSISession()
	for _, session := range sessions {
		info.TargetPortals = append(info.TargetPortals, session.TargetPortal)
		info.TargetIqns = append(info.TargetIqns, session.TargetIqn)
	}

	return info, err
}
