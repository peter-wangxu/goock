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
	"github.com/peter-wangxu/goock/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsMultipathEnabled(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	enabled := IsMultipathEnabled()
	assert.True(t, enabled)
}

func TestFlushPath(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	err := FlushPath("36006016074e03a00ee762958673eaf1b")
	assert.Nil(t, err)

}

func TestFlushPathAll(t *testing.T) {

	SetExecutor(test.NewMockExecutor())
	err := FlushPath("")
	assert.Nil(t, err)
}

func TestReconfigure(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	success := Reconfigure()
	assert.Nil(t, success)
}

func TestReconfigureError(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	// TODO need to support specifying test data
	//err := Reconfigure()
	//assert.Error(t, err, "failed to reconfigure.")
}

func TestReload(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	err := Reload()
	assert.Nil(t, err)

}

func TestCheckDevice(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	var ret = CheckDevice("/dev/sdx")
	assert.Equal(t, true, ret, "The return of CheckDevice is not true.")

}

func TestCheckDeviceNotFound(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	var ret = CheckDevice("/dev/invalid/path")
	assert.Equal(t, false, ret)
}

func TestResizeMpath(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	var ret = ResizeMpath("36006016074e03a003dbe2a580510610b")
	assert.Nil(t, ret)
}

func TestFindMpathByPath(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	FindMpathByPath("/dev/disk/by-path/ip-192.168.3.50:3260-iscsi-iqn.1992-04.com.emc:cx.apm00152904558.b12-lun-11")
	//TODO how to mock "path, err := filepath.EvalSymlinks(path)"
	//assert.NotEmpty(t, ret)
}
