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
	}
	if len(args) == 1 {
		// User only supplies the local device name

	} else {
		// User specify TargetIP with LUN ID
	}
	if err != nil {
		return err
	}
	return HandleISCSIExtend(args...)
}
