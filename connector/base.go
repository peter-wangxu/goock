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
	"github.com/peter-wangxu/goock/model"
	"github.com/sirupsen/logrus"
)

type StringEnum string

const (
	READWRITE StringEnum = "rw"
	READONLY  StringEnum = "ro"
)

const (
	ISCSI_PROTOCOL StringEnum = "iscsi"
	FC_PROTOCOL    StringEnum = "fc"
)

type ConnectionProperty struct {
	TargetIqns      []string
	TargetPortals   []string
	TargetLuns      []int
	StorageProtocol string
	AccessMode      StringEnum
}

// Validate whether the ConnectionProperty is empty
func (prop ConnectionProperty) IsEmpty() error {
	if len(prop.TargetPortals) == 0 || len(prop.TargetLuns) == 0 {
		return fmt.Errorf("An empty ConnectionProperty is specified, forget target IP or LUN id?")
	}
	return nil
}

type HostInfo struct {
	Initiator string
	Ip        string
	Hostname  string
	OSType    string
}

type VolumeInfo struct {
	MultipathId string
	Paths       []string
	Wwn         string
	Multipath   string
}

type Interface interface {
	GetHostInfo(args []string) (HostInfo, error)
	ConnectVolume(connectionProperty ConnectionProperty) (VolumeInfo, error)
	DisconnectVolume(connectionProperty ConnectionProperty) error
	ExtendVolume(connectionProperty ConnectionProperty) error
	LoginPortal(targetPortal string, targetIqn string) error
	DiscoverPortal(targetPortal ...string) []model.ISCSISession
}

var log *logrus.Logger = logrus.New()

func SetLogger(l *logrus.Logger) {
	log = l
}
