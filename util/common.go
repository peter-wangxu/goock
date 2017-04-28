package util

import (
	"errors"
	"time"
	"github.com/Sirupsen/logrus"
	"github.com/peter-wangxu/goock/exec"
)

var executor = exec.New()

func SetExecutor(e exec.Interface) {
	executor = e
}

func WaitForPath(path string, maxWait int) bool {
	for x := 0; x < maxWait; x++ {
		err := IsPathExists(path)
		if (err == nil) {
			return true
		}
		time.Sleep(2)
	}
	logrus.Debug("Path ", path, " does not appear in ", maxWait * 2, " seconds.")
	return false
}

// Return immediately once any path found
func WaitForAnyPath(paths []string) (string, error) {

	err := errors.New("No path found")
	for _, path := range (paths) {
		err = IsPathExists(path)
		if (err == nil) {
			return path, err
		}
	}
	return "", err
}

func FilterPath(paths []string) ([]string, error) {
	var newPaths []string
	for _, path := range (paths) {
		err := IsPathExists(path)
		if (err == nil) {
			newPaths = append(newPaths, path)
		} else {
			logrus.WithError(err).Debugf("Unable to locate path: %s", path)
		}
	}
	return newPaths, nil
}

func IsPathExists(path string) error {
	_, err := executor.Command("ls", path).CombinedOutput()
	return err
}

func Contains(key string, all []string) bool {
	for _, item := range (all) {
		if (key == item) {
			return true
		}
	}
	return false
}