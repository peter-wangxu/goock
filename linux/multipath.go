package linux

import (
	"fmt"
	"github.com/Sirupsen/logrus"
)

// Flush device(s) via multipath -f <device>/-F
func FlushPath(path string) error {
	var err error
	if path != "" {
		_, err = executor.Command("multipath", "-f", path).CombinedOutput()
	} else {
		_, err = executor.Command("multipaht", "-F").CombinedOutput()
	}
	return err
}

// Get paths by multipath -ll
func GetPaths(path string) ([]string, error) {
	// TODO wait for the multipath parser
	output, err := executor.Command("multipath", "-l", path).CombinedOutput()
	if nil != err {
		logrus.WithError(err).Warn("Got error when multipath -l.")
	} else {
		// TODO parse and return
	}
	return []string{string(output)}, nil
}

// Reconfigure multipath
func Reconfigure() bool {
	output, err := executor.Command("multipathd", "reconfigure").CombinedOutput()
	if nil != err {
		logrus.WithError(err).Info(fmt.Sprintf("Failed to reconfigure the multipathd. %s", output))
		return false
	}
	return true
}

// Force multipath reloads devices via multipath -r
func Reload() error {
	output, err := executor.Command("multipath", "-r").Output()
	if nil != err {
		logrus.WithError(err).Debug(fmt.Sprintf("Reload multipath failed: %s", output))
	}
	return nil
}

// Check if the path is a multipath device
func CheckDevice(path string) bool {
	output, err := executor.Command("multipath", "-c", path).CombinedOutput()
	if nil != err {
		logrus.WithError(err).Debug(fmt.Sprintf("The specified path doesn't exist: %s", output))
		return false
	}
	return true
}

func ResizeMpath(mpathId string) error {
	output, err := executor.Command("multipathd", "resize", "map", mpathId).CombinedOutput()
	if nil != err {
		logrus.WithError(err).Debug(fmt.Sprintf("Resize %s failed due to [%s]", mpathId, output))
	}
	return err
}
