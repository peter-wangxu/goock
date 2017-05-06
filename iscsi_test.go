package goock

import (
	"fmt"
	"github.com/peter-wangxu/goock/connector"
	"github.com/peter-wangxu/goock/linux"
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/test"
	goockutil "github.com/peter-wangxu/goock/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestISCSIConnector_GetHostInfo(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	iscsi := New()
	props, err := iscsi.GetHostInfo([]string{})
	assert.Nil(t, err)
	assert.Equal(t, "iqn.1993-08.org.debian:01:b974ee37fea", props.Initiator)
	assert.Contains(t, []string{"windows", "linux", "macOS"}, props.OSType)
	assert.NotEmpty(t, props.Hostname)
}

func TestISCSIConnector_LoginISCSIPortalAlreadyLoggedin(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	iscsi := New()
	err := iscsi.LoginISCSIPortal("10.64.76.253:3260", "iqn.1992-04.com.emc:cx.fcnch097ae5ef3.h1")
	assert.Nil(t, err)
}

func TestISCSIConnector_LoginISCSIPortal(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	iscsi := New()
	err := iscsi.LoginISCSIPortal("192.168.1.2:3260", "iqn.1992-04.com.emc:cx.fcnch097ae1234.h2")
	assert.Nil(t, err)
}

func TestISCSIConnector_ConnectVolume_NoProp(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	goockutil.SetExecutor(test.NewMockExecutor())
	iscsi := New()
	fakeProperty := connector.ConnectionProperty{}
	_, err := iscsi.ConnectVolume(fakeProperty)
	assert.EqualError(t, err, "No path found")
}

func TestISCSIConnector_ConnectVolume(t *testing.T) {
	goockutil.SetExecutor(test.NewMockExecutor())
	linux.SetExecutor(test.NewMockExecutor())
	iscsi := New()
	fakeProperty := connector.ConnectionProperty{}
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
	assert.Len(t, device.Paths, 2)
	assert.NotEmpty(t, device.Wwn)
	assert.Equal(t, device.Multipath, fmt.Sprintf("/dev/disk/by-id/dm-uuid-mpath-%s", device.Wwn))
	assert.Equal(t, device.Paths[0], "/dev/disk/by-path/ip-192.168.3.50:3260-iscsi-iqn.1992-04.com.emc:cx.apm00152904558.b12-lun-11")
	assert.Equal(t, device.Paths[1], "/dev/disk/by-path/ip-192.168.3.49:3260-iscsi-iqn.1992-04.com.emc:cx.apm00152904558.a12-lun-11")

}

func TestISCSIConnector_ConnectVolumeNoMultipath(t *testing.T) {
	//TODO(peter) wait for test data feeding feature
}

func TestISCSIConnector_DisconnectVolume(t *testing.T) {
	goockutil.SetExecutor(test.NewMockExecutor())
	linux.SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	iscsi := New()
	fakeProperty := connector.ConnectionProperty{}
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
