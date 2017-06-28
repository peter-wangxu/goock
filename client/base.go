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
	"github.com/peter-wangxu/goock/exec"
	"github.com/peter-wangxu/goock/linux"
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/util"
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger = logrus.New()

var VolumeFormat = `Volume Information:
Multipath    : %s
Single paths :
%s
Multipath ID : %s
WWN          : %s
`

var HostInfoFormat = `Volume Information:
iSCSI Qualified Name(IQN)      :
%s
Host Bus Adapter               :
%s
Connected Fibre Channel Target :
%s
Connected iSCSI sessions       :
%s
`

// Enable the console log for client
func InitLog(debug bool) error {
	if debug {
		log = logrus.New()
		//logrus.SetFormatter(&logrus.JSONFormatter{})

		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		log.Out = os.Stdout

		// Only log the warning severity or above.
		log.Level = logrus.DebugLevel
		log.Formatter = &logrus.TextFormatter{}

	} else {
		log = logrus.New()
	}

	// Set logger for all modules
	//cmd.SetLogger(log)
	connector.SetLogger(log)
	exec.SetLogger(log)
	linux.SetLogger(log)
	model.SetLogger(log)
	util.SetLogger(log)
	return nil
}

// Handle the Extend request based the device type
func HandleExtend(args ...string) error {
	var err error
	if len(args) <= 0 {
		err = fmt.Errorf("Need device name or Target IP with LUN ID.")
	} else if len(args) == 1 {
		// User only supplies the local device name
		err = fmt.Errorf("Currently device name is not supported.")
	} else {
		// User specify TargetIP with LUN ID
		err = HandleISCSIExtend(args...)
	}
	return err

}

func HandleInfo(args ...string) error {
	hostInfo, err := connector.GetHostInfo()
	if err == nil {
		BeautifyHostInfo(hostInfo)
	} else {
		log.WithError(err).Warn("Unable to get host information, not privileged or tools not installed?")
	}
	return err
}

// BeautifyHostInfo prints the output to console
func BeautifyHostInfo(info connector.HostInfo) {
	var wwns []string
	var targetWwns []string

	for i, wwnns := range info.Wwnns {
		wwns = append(wwns, wwnns+":"+info.Wwpns[i])
	}

	for j, targetWwnns := range info.TargetWwnns {
		targetWwns = append(targetWwns, targetWwnns+":"+info.TargetWwpns[j])
	}
	fmt.Printf(HostInfoFormat, info.Initiator, wwns, targetWwns, "fakeIP")
}
