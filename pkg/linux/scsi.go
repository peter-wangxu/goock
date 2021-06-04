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
package linux

import (
	"fmt"
	"github.com/peter-wangxu/goock/pkg/model"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GetWWN(path string) string {
	output, _ := executor.Command("/lib/udev/scsi_id", "--page", "0x83",
		"--whitelisted", path).CombinedOutput()
	return strings.Trim(string(output), "\n")
}

// Check if path is already RW/RO
// For a multipath, pass wwn here to validate

// NAME                                RO
// sdb                                  0
// └─36006016003b03a00da41ad58e6ab1cc0  0
// sdd                                  0
// └─36006016015e03a00bea7c7588c91d581  0
// sde                                  0
// sdf                                  0
// └─36006016003b03a00da41ad58e6ab1cc0  0
// sdg                                  0
// └─36006016015e03a00bea7c7588c91d581  0
// sr0                                  0
// vda                                  0
// └─vda1                               0
func CheckReadWrite(path string, wwn string) bool {
	output, _ := executor.Command("lsblk", "-o", "NAME,RO", "-l", "-n").CombinedOutput()
	pattern, _ := regexp.Compile("(\\w+)\\s+([01])\\s?")
	results := pattern.FindAllStringSubmatch(string(output), -1)
	readWrite := false
	for _, result := range results {
		k, v := result[1], result[2]
		if k == path || k == wwn {
			if strings.Contains(v, "1") {
				readWrite = false
			} else {
				readWrite = true
			}
		}
	}
	return readWrite
}

// Get block device size
func GetDeviceSize(path string) int {
	output, err := executor.Command("blockdev", "--getsize64", path).CombinedOutput()
	if nil != err {
		log.WithError(err).Warnf("Unable to get size of device %s", path)
	}
	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return 0
	}
	i, _ := strconv.Atoi(trimmed)
	return i
}

// use echo "c t l" > to /sys/class/scsi_host/%s/scan
func ScanSCSIBus(path string, content string) error {
	cmd := executor.Command("tee", "-a", path)
	cmd.SetStdin(strings.NewReader(content))
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.WithError(err).Warn("Rescan Bus failed")
	}
	return err

}

// path = "/dev/sdb" or "sdb"
// Use echo 1 > /sys/block/%s/device/delete to force delete the device
func RemoveSCSIDevice(path string) {
	if strings.Contains(path, string(filepath.Separator)) {
		// Before remove the device from host, flush buffers to disk
		FlushDeviceIO(path)
		// Get the file name from the full path, ex : /dev/sdb -> sdb
		_, path = filepath.Split(path)
	} else {
		FlushDeviceIO(fmt.Sprintf("/dev/%s", path))
	}

	path = fmt.Sprintf("/sys/block/%s/device/delete", path)
	ScanSCSIBus(path, "1")
	log.Debugf("Removed device [%s].", path)
}

// path = "/dev/sdb" or "
// "/dev/disk/by-path/ip-10.244.213.177:3260-iscsi-iqn.1992-04.com.emc:cx.fnm00150600267.a0-lun-10"
func FlushDeviceIO(path string) error {
	cmd := executor.Command("blockdev", "-v", "--flushbufs", path)
	_, err := cmd.CombinedOutput()
	return err
}

// Commands example:
// echo 1 > /sys/bus/scsi/drivers/sd/9:0:0:6/rescan
func ExtendDevice(path string) (int, error) {

	info, err := GetDeviceInfo(path)
	if err != nil {
		return 0, fmt.Errorf("Unable to extend device %s, device info not found", path)
	}
	deviceId := info.GetDeviceIdentifier()
	rescanPath := fmt.Sprintf("/sys/bus/scsi/drivers/sd/%s/rescan", deviceId)
	deviceSize := GetDeviceSize(path)
	log.WithFields(logrus.Fields{
		"path":     path,
		"device":   deviceId,
		"original": deviceSize,
	}).Debug("Begin to extend the device.")

	ScanSCSIBus(rescanPath, "1")
	newSize := GetDeviceSize(path)
	log.WithFields(logrus.Fields{
		"path":    path,
		"newSize": newSize,
	}).Info("Extend device finished.")
	return newSize, err
}

// output:
// sudo sg_scan /dev/disk/by-path/pci-0000:05:00.1-fc-0x5006016d09200925-lun-0
// /dev/disk/by-path/pci-0000:05:00.1-fc-0x5006016d09200925-lun-0: scsi9 channel=0 id=0 lun=0 [em]
func GetDeviceInfo(path string) (model.DeviceInfo, error) {
	devices := model.NewDeviceInfo(path)
	if len(devices) <= 0 {
		log.Warn("Unable to get device info for device ", path)
		return model.DeviceInfo{}, fmt.Errorf("Unable to get device info.")
	}
	return devices[0], nil
}
