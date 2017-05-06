package linux

import (
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/test"
	testhelper "github.com/peter-wangxu/goock/test/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWWN(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	wwn := GetWWN("/dev/disk/by-path/ip-192.168.3.50:3260-iscsi-iqn.1992-04.com.emc:cx.apm00152904558.b12-lun-11")
	assert.Equal(t, "350060160b6e00e5a50060160b6e00e5a", wwn)
}

func TestCheckReadWrite(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	r := CheckReadWrite("sdb", "36006016003b03a00da41ad58e6ab1cc0")
	assert.True(t, r)
}

func TestCheckReadWritePartial(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	r := CheckReadWrite("sdg", "36006016015e03a00bea7c7588c91d581")
	assert.False(t, r)
}

func TestCheckReadWriteNonexistent(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	r := CheckReadWrite("sdx", "36006016015e03a00bea7c7588c91d581xxx")
	assert.False(t, r)
}

func TestRemoveSCSIDevice(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	RemoveSCSIDevice("sdb")
}

func TestRemoveSCSIDeviceWithPath(t *testing.T) {
	testhelper.SkipIfWindows(t)
	SetExecutor(test.NewMockExecutor())
	RemoveSCSIDevice("/dev/sdx")
}
func TestGetDeviceInfo(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	info := GetDeviceInfo("/dev/sdb")
	assert.Equal(t, "/dev/sdb", info.Device)
}
func TestGetDeviceInfoNotFound(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	info := GetDeviceInfo("/dev/sdx")
	assert.Equal(t, "", info.Device)
}
