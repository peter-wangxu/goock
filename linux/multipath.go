package linux

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/peter-wangxu/goock/model"
	goockutil "github.com/peter-wangxu/goock/util"
	"path/filepath"
)

func IsMultipathEnabled() bool {
	_, err := executor.Command("multipathd", "show", "status").CombinedOutput()
	if err != nil {
		return false
	}
	return true
}

// Flush device(s) via multipath -f <device>/-F
func FlushPath(path string) error {
	var err error
	if path != "" {
		_, err = executor.Command("multipath", "-f", path).CombinedOutput()
	} else {
		_, err = executor.Command("multipath", "-F").CombinedOutput()
	}
	return err
}

// Reconfigure multipath
func Reconfigure() error {
	output, err := executor.Command("multipathd", "reconfigure").CombinedOutput()
	if nil != err {
		logrus.WithError(err).Info(fmt.Sprintf("Failed to reconfigure the multipathd. %s", output))
	}
	return err
}

// Force multipath reloads devices via multipath -r
func Reload() error {
	output, err := executor.Command("multipath", "-r").Output()
	if nil != err {
		logrus.WithError(err).Debug(fmt.Sprintf("Reload multipath failed: %s", output))
	}
	return err
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

// Return the multipath by wwn
// 1) When multipath friendly names are ON:
// a device file will show up in
// /dev/disk/by-id/dm-uuid-mpath-<WWN>
// /dev/disk/by-id/dm-name-mpath<N>
// /dev/disk/by-id/scsi-mpath<N>
// /dev/mapper/mpath<N>
//
// 2) When multipath friendly names are OFF:
// /dev/disk/by-id/dm-uuid-mpath-<WWN>
// /dev/disk/by-id/scsi-<WWN>
// /dev/mapper/<WWN>
func FindMpathByWwn(wwn string) string {
	logrus.Info("Try to find multipath device for WWN: ", wwn)
	// Wait for its appearance under /dev/disk/by-id/dm-uuid-mpath
	potential1 := fmt.Sprintf("/dev/disk/by-id/dm-uuid-mpath-%s", wwn)
	existed := goockutil.WaitForPath(potential1, 10)
	if existed {
		return potential1
	}
	// Wait for its appearance under /dev/mapper/
	potential2 := fmt.Sprintf("/dev/mapper/%s", wwn)
	existed = goockutil.WaitForPath(potential2, 10)
	if existed {
		return potential2
	}
	return ""
}

// Use multipath -l <path> to discover multipath device
// Valid <path> could be WWN or /dev/sdb like path
func FindMpathByPath(path string) string {
	path, err := filepath.EvalSymlinks(path)
	logrus.WithError(err).Info("real path", path)
	logrus.Info("Try to find multipath device by multipath -l : ", path)
	models := model.FindMultipath(path)
	mPath := ""
	if len(models) > 0 {
		wwn := models[0].Wwn
		mPath = fmt.Sprintf("/dev/disk/by-id/dm-uuid-mpath-%s", wwn)
	}
	return mPath
}

func FindMultipathByWwn(wwn string) model.Multipath {
	models := model.FindMultipath(wwn)
	if len(models) >= 1 {
		return models[0]
	}
	return model.Multipath{}
}
