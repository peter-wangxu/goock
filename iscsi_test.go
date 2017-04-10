package goock

import (
	"testing"
	"github.com/peter-wangxu/goock/test"
	"github.com/stretchr/testify/assert"
	goockutil "github.com/peter-wangxu/goock/util"
	"github.com/peter-wangxu/goock/linux"
	"fmt"
)

func TestISCSIConnector_ConnectVolume_NoProp(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	goockutil.SetExecutor(test.NewMockExecutor())
	connector := New()
	fakeProperty := ConnectionProperty{}
	_, err := connector.ConnectVolume(fakeProperty)
	assert.EqualError(t, err, "No path found")
}

func TestISCSIConnector_ConnectVolume(t *testing.T) {
	goockutil.SetExecutor(test.NewMockExecutor())
	linux.SetExecutor(test.NewMockExecutor())
	connector := New()
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
	device, err := connector.ConnectVolume(fakeProperty)
	assert.Nil(t, err)
	assert.Len(t, device.paths, 2)
	assert.NotEmpty(t, device.Wwn)
	assert.Equal(t, device.Multipath, fmt.Sprintf("/dev/disk/by-id/dm-uuid-mpath-%s", device.Wwn))
	assert.Equal(t, device.paths[0], "/dev/disk/by-path/ip-192.168.3.50:3260-iscsi-iqn.1992-04.com.emc:cx.apm00152904558.b12-lun-11")
	assert.Equal(t, device.paths[1], "/dev/disk/by-path/ip-192.168.3.49:3260-iscsi-iqn.1992-04.com.emc:cx.apm00152904558.a12-lun-11")

}