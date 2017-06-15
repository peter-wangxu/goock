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
)

var fcConnector = connector.NewFibreChannelConnector()

func SetFcConnector(fc connector.FibreChannelInterface) {
	fcConnector = fc
}

func HandleFCConnect(args ...string) error {
	var err error
	if len(args) <= 0 {
		log.Error("[wwn] and/or [lun id] are required.")
		err = fmt.Errorf("[wwn] and/or [lun id] are required.")
	} else if len(args) == 1 {
		// User only supply the LUN ID, so did a wildcard scan for all connected targets
		fcConnector.GetHostInfo()
	} else {
		err = fmt.Errorf("Currently [wwn] with [lun id] is not supported.")
		log.WithError(err).Error("Unsupported parameters.")
	}

	return err
}

func HandleFCExtend(args ...string) error {
	return nil
}
