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
	"github.com/peter-wangxu/goock/pkg/model"
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

func TestFlushDeviceIO(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	err := FlushDeviceIO("/dev/sdm")
	assert.Nil(t, err)
}

func TestExtendDevice(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())

	newSize, err := ExtendDevice("/dev/sdg")

	assert.Nil(t, err)
	assert.Equal(t, 2147483648, newSize)

}

func TestExtendDeviceNoDeviceInfo(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	_, err := ExtendDevice("/dev/unknown")
	assert.Error(t, err)
}

func TestGetDeviceInfo(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	info, err := GetDeviceInfo("/dev/sdb")
	assert.Nil(t, err)
	assert.Equal(t, "/dev/sdb", info.Device)
	assert.Equal(t, "scsi0", info.Host)
	assert.Equal(t, 0, info.GetHostId())
	assert.Equal(t, "0:1:0:0", info.GetDeviceIdentifier())
}
func TestGetDeviceInfoNotFound(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	info, err := GetDeviceInfo("/dev/sdx")
	assert.Error(t, err)
	assert.Equal(t, "", info.Device)
}
