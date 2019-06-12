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
	WaitInterval int = 2
	MaxWait      int = 10
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
		time.Sleep(time.Second * time.Duration(WaitInterval))
		err := IsPathExists(path)
		if err == nil {
			return true
		}
	}
	log.Debugf("Path %s does not appear in %v seconds", path, maxWait*WaitInterval)
	return false
}

// Return immediately once any path found of hook is nil
// else run the hook and wait
func WaitForAnyPath(paths []string, hook func()) (string, error) {

	err := errors.New("No path found")
	maxWait := MaxWait
	if hook == nil {
		// Only run once for the paths.
		maxWait = 1
	}

	for x := 0; x < maxWait; x++ {
		time.Sleep(time.Second * time.Duration(WaitInterval))
		for _, path := range paths {
			err = IsPathExists(path)
			if err == nil {
				return path, err
			}
		}
		// Run the hook, such as, Rescan hosts
		if nil != hook {
			hook()
		}
	}
	return "", err
}

// FilterPath Filters out paths which are not existed.
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

// WaitForPathRemoval Returns the paths which are still existing
func WaitForPathRemoval(paths []string, maxWait int) []string {
	var left []string
	for x := 0; x < maxWait; x++ {
		time.Sleep(time.Second * time.Duration(WaitInterval))
		left, _ = FilterPath(paths)
		if len(left) == 0 {
			break
		}
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
