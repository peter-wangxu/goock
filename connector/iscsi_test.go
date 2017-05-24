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
package connector

import (
	"fmt"
	"github.com/peter-wangxu/goock/linux"
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/test"
	goockutil "github.com/peter-wangxu/goock/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestISCSIConnector_GetHostInfo(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	iscsi := NewISCSIConnector()
	props, err := iscsi.GetHostInfo([]string{})
	assert.Nil(t, err)
	assert.Equal(t, "iqn.1993-08.org.debian:01:b974ee37fea", props.Initiator)
	assert.Contains(t, []string{"windows", "linux", "macOS"}, props.OSType)
	assert.NotEmpty(t, props.Hostname)
}

func TestISCSIConnector_LoginISCSIPortalAlreadyLoggedin(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	iscsi := NewISCSIConnector()
	err := iscsi.LoginPortal("10.64.76.253:3260", "iqn.1992-04.com.emc:cx.fcnch097ae5ef3.h1")
	assert.Nil(t, err)
}

func TestISCSIConnector_LoginISCSIPortal(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	iscsi := NewISCSIConnector()
	err := iscsi.LoginPortal("192.168.1.2:3260", "iqn.1992-04.com.emc:cx.fcnch097ae1234.h2")
	assert.Nil(t, err)
}

func TestISCSIConnector_SetNode2Auto(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	iscsi := NewISCSIConnector()
	err := iscsi.SetNode2Auto("110.244.213.177:3260", "iqn.1992-04.com.emc:cx.fnm00150600267.a0")
	assert.Nil(t, err)

}

func TestISCSIConnector_ConnectVolume_NoProp(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	goockutil.SetExecutor(test.NewMockExecutor())
	iscsi := NewISCSIConnector()
	fakeProperty := ConnectionProperty{}
	_, err := iscsi.ConnectVolume(fakeProperty)
	assert.EqualError(t, err, "No path found")
}

// Assert that all all are accessible
func TestISCSIConnector_ConnectVolume(t *testing.T) {
	// TODO add test here
}

func TestISCSIConnector_ConnectVolume_NotAll(t *testing.T) {
	goockutil.SetExecutor(test.NewMockExecutor())
	linux.SetExecutor(test.NewMockExecutor())
	iscsi := NewISCSIConnector()
	fakeProperty := ConnectionProperty{}
	fakeProperty.TargetIqns = []string{
		"iqn.1992-04.com.emc:cx.apm00152904558.b12",
		"iqn.1992-04.com.emc:cx.apm00152904558.a12",
	}
	fakeProperty.TargetPortals = []string{
		"192.168.3.50:3260",
		"192.168.3.49:3260",
	}
	fakeProperty.TargetLuns = []int{
		11,
		11,
	}
	device, err := iscsi.ConnectVolume(fakeProperty)
	assert.Nil(t, err)
	assert.Len(t, device.Paths, 1)
	assert.NotEmpty(t, device.Wwn)
	assert.Equal(t, fmt.Sprintf("/dev/disk/by-id/dm-uuid-mpath-%s", device.Wwn), device.Multipath)
	//assert.Equal(t, "/dev/disk/by-path/ip-192.168.3.50:3260-iscsi-iqn.1992-04.com.emc:cx.apm00152904558.b12-lun-11", device.Paths[0])
	assert.Equal(t, "/dev/disk/by-path/ip-192.168.3.49:3260-iscsi-iqn.1992-04.com.emc:cx.apm00152904558.a12-lun-11", device.Paths[0])
}

func TestISCSIConnector_ConnectVolumeNoMultipath(t *testing.T) {
	//TODO(peter) wait for test data feeding feature
}

func TestISCSIConnector_DisconnectVolume(t *testing.T) {
	goockutil.SetExecutor(test.NewMockExecutor())
	linux.SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	iscsi := NewISCSIConnector()
	fakeProperty := ConnectionProperty{}
	fakeProperty.TargetIqns = []string{
		"iqn.1992-04.com.emc:cx.apm00152904447.a17",
		"iqn.1992-04.com.emc:cx.apm00152904447.b17",
	}
	fakeProperty.TargetPortals = []string{
		"10.168.3.44:3260",
		"10.168.3.45:3260",
	}
	fakeProperty.TargetLuns = []int{
		11,
		11,
	}
	iscsi.DisconnectVolume(fakeProperty)
}

func TestISCSIConnector_DisconnectVolumeNoMultipath(t *testing.T) {
	//TODO(peter) wait for test data feeding feature

}

func TestISCSIConnector_ExtendVolume(t *testing.T) {
	linux.SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	iscsi := NewISCSIConnector()
	fakeProperty := ConnectionProperty{}
	fakeProperty.TargetIqns = []string{
		"iqn.1992-04.com.emc:cx.apm00141313414.a17",
		"iqn.1992-04.com.emc:cx.apm00141313414.b17",
	}
	fakeProperty.TargetPortals = []string{
		"10.168.7.14:3260",
		"10.168.7.15:3260",
	}
	fakeProperty.TargetLuns = []int{
		19,
		19,
	}
	err := iscsi.ExtendVolume(fakeProperty)
	assert.Nil(t, err)
}

func TestISCSIConnector_ExtendVolumeNoAnyPath(t *testing.T) {
	linux.SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	iscsi := NewISCSIConnector()
	fakeProperty := ConnectionProperty{}
	fakeProperty.TargetIqns = []string{
		"iqn.1992-04.com.emc:cx.apm55555555555.b17",
	}
	fakeProperty.TargetPortals = []string{
		"10.168.7.199:3260",
	}
	fakeProperty.TargetLuns = []int{
		19,
	}
	err := iscsi.ExtendVolume(fakeProperty)
	assert.Error(t, err)
}
