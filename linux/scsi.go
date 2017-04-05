package linux

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/peter-wangxu/goock/exec"
	"regexp"
	"strconv"
	"strings"
)

var executor = exec.New()

func SetExecutor(e exec.Interface) {
	executor = e
}

func GetWWN(path string) string {
	output, _ := executor.Command("/lib/udev/scsi_id", "--page", "0x83",
		"--whitelisted", path).CombinedOutput()
	return strings.Trim(string(output), " ")
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
	for _, result := range results {
		k, v := result[1], result[2]
		if k == path || k == wwn {
			if strings.Contains(v, "0") {
				return false
			}
		}
	}
	return true
}

// Get block device size
func GetDeviceSize(path string) int {
	output, err := executor.Command("blockdev", "--getsize64", path).CombinedOutput()
	if nil != err {
		logrus.WithError(err).Warn("Unable to get size of device %s", path)
	}
	trimmed := strings.Trim(string(output), " ")
	if trimmed == "" {
		return 0
	}
	i, _ := strconv.Atoi(trimmed)
	return i
}

// use echo "c t l" > to /sys/class/scsi_host/%s/scan
func ScanSCSIBus(path string, content string) error {
	if content == "" {
		// "hba_channel target_id target_lun"
		content = "- - -"
	}
	cmd := executor.Command("tee", "-a", path)
	cmd.SetStdin(strings.NewReader(content))
	_, err := cmd.CombinedOutput()
	if (err != nil) {
		logrus.WithError(err).Warn("Rescan Bus failed")

	}
	return err

}

// Use echo 1 > /sys/block/%s/device/delete to force delete the device
func RemoveSCSIDevice(path string) {
	if !strings.Contains(path, "/") {
		path = fmt.Sprintf("/sys/block/%s/device/delete", path)
	}

	cmd := executor.Command("tee", "-a", path)
	cmd.SetStdin(strings.NewReader("1"))
	cmd.CombinedOutput()
}

//TODO fix this after adding UT
func ExtendDevice(path string) bool {
	return true
}
