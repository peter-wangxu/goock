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
	"regexp"
	"strconv"
	"strings"
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

// HandleConnect dispatches the cli to iscsi/fc respectively.
func HandleConnect(args ...string) error {
	var err error
	if len(args) <= 0 {
		log.Error("Target IP or wwn is required.")
		err = fmt.Errorf("Target IP or wwn is required.")
	} else if len(args) == 1 {
		// User only supply the LUN ID, so did a wildcard scan for all connected targets
		err = fmt.Errorf("Currently [lun id] is not supported.")
		log.WithError(err).Error("Unsupported parameters.")
		log.Error("%s", args)
	} else {
		target := args[0]
		// Make sure the last param is LUN ID.
		if _, err = ValidateLunId(args[len(args)-1:]); err == nil {
			if IsIpLike(target) {
				return HandleISCSIConnect(args...)
			}
			if IsFcLike(target) {
				return HandleFCConnect(args...)
			}
		}
	}
	return err
}

// HandleDisconnect dispatches the cli to iscsi/fc respectively.
func HanddleDisconnect(args ...string) error {
	return nil
}

// HandleExtend handles the Extend request based the device type
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
		log.WithError(err).Warn("Unable to get host information, permission denied or tools not installed?")
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

func ValidateLunId(lunIDs []string) ([]int, error) {
	var err error
	re, _ := regexp.Compile("\\d+")
	var ret []int
	for _, lun := range lunIDs {
		if re.MatchString(lun) == false {
			err = fmt.Errorf("%s does not look like a LUN ID.", lun)
			break
		}
		i, _ := strconv.Atoi(lun)
		ret = append(ret, i)
	}
	if len(ret) <= 0 {
		log.Warnf("No lun ID specified, correct and retry.")
	}
	return ret, err
}

// IsLunLike tests if *data* is a lun id.
func IsLunLike(data string) bool {

	if _, err := strconv.Atoi(data); err != nil {
		return false
	}
	return true
}

// IsIpLike tests if *data* is a ipv4 address.
func IsIpLike(data string) bool {
	// IPv4 match
	if m, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$", data); !m {
		return false
	}
	return true
}

// IsFcLike tests if *data* is a fc wwn.
func IsFcLike(data string) bool {
	// Replace the colons if presents
	data = strings.Replace(data, ":", "", -1)
	// Matches the wwpn
	if m, _ := regexp.MatchString("^\\w{16}$", data); m == true {
		return true
	}
	// Matches the wwnn + wwpn
	if m, _ := regexp.MatchString("^\\w{32}$", data); m == true {
		return true
	}
	return false
}
