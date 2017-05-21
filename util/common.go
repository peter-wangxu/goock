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
package util

import (
	"errors"
	"github.com/peter-wangxu/goock/exec"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	WAIT_INTERVAL int = 2
)

var log *logrus.Logger = logrus.New()

func SetLogger(l *logrus.Logger) {
	log = l
}

var executor = exec.New()

func SetExecutor(e exec.Interface) {
	executor = e
}

func WaitForPath(path string, maxWait int) bool {
	for x := 0; x < maxWait; x++ {
		err := IsPathExists(path)
		if err == nil {
			return true
		}
		time.Sleep(time.Second * time.Duration(WAIT_INTERVAL))
	}
	log.Debugf("Path %s does not appear in %s seconds", path, maxWait*WAIT_INTERVAL)
	return false
}

// Return immediately once any path found
func WaitForAnyPath(paths []string) (string, error) {

	err := errors.New("No path found")
	for _, path := range paths {
		err = IsPathExists(path)
		if err == nil {
			return path, err
		}
	}
	return "", err
}

func FilterPath(paths []string) ([]string, error) {
	var newPaths []string
	for _, path := range paths {
		err := IsPathExists(path)
		if err == nil {
			newPaths = append(newPaths, path)
		} else {
			log.WithError(err).Debugf("Unable to locate path: %s", path)
		}
	}
	return newPaths, nil
}

// Returns the paths which are still existing
func WaitForPathRemoval(paths []string, maxWait int) []string {
	var left []string
	for x := 0; x < maxWait; x++ {
		left, _ = FilterPath(paths)
		if len(left) == 0 {
			break
		}
		time.Sleep(time.Second * time.Duration(WAIT_INTERVAL))
	}
	return left
}

func IsPathExists(path string) error {
	_, err := executor.Command("ls", path).CombinedOutput()
	return err
}

func Contains(key string, all []string) bool {
	for _, item := range all {
		if key == item {
			return true
		}
	}
	return false
}
