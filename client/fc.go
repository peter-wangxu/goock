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
	"github.com/peter-wangxu/goock/connector"
	"strconv"
)

var fcConnector = connector.NewFibreChannelConnector()

// SetFcConnector sets the connector for FC connection
func SetFcConnector(fc connector.FibreChannelInterface) {
	fcConnector = fc
}

// Convert2ConnectionProperty converts wwn and lunid pair into ConnnectionProperty
func Convert2ConnectionProperty(wwns []string, lunID string) connector.ConnectionProperty {
	var property connector.ConnectionProperty
	property.TargetWwns = wwns
	property.TargetLun, _ = strconv.Atoi(lunID)
	property.StorageProtocol = connector.FcProtocol

	return property
}

// HandleFCConnect handles the connection of FC target
func HandleFCConnect(args ...string) error {
	var err error
	if len(args) == 1 {
		// User only supply the LUN ID, so did a wildcard scan for all connected targets
		err = fmt.Errorf("currently [LUN ID] is not supported")
		log.WithError(err).Error("Unsupported parameters.")
	} else {
		targets := args[:len(args)-1]

		conn := Convert2ConnectionProperty(targets, args[len(args)-1])

		var info connector.VolumeInfo

		if info, err = fcConnector.ConnectVolume(conn); err == nil {
			BeautifyVolumeInfo(info)
		}
	}

	return err
}

// HandleFCExtend handle the request to extend the FC devices.
func HandleFCExtend(args ...string) error {
	return nil
}
