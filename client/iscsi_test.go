package client

import (
	"fmt"
	"github.com/peter-wangxu/goock/connector"
	"github.com/peter-wangxu/goock/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

type FakeISCSIConnector struct {
}

func (fake *FakeISCSIConnector) GetHostInfo(args []string) (connector.HostInfo, error) {
	return connector.HostInfo{}, nil
}

func (fake *FakeISCSIConnector) ConnectVolume(connectionProperty connector.ConnectionProperty) (connector.VolumeInfo, error) {
	if len(connectionProperty.TargetPortals) > 0 && connectionProperty.TargetPortals[0] == "10.244.244.244" {
		return connector.VolumeInfo{}, fmt.Errorf("Failed to connect volume.")
	}
	return connector.VolumeInfo{}, nil
}

func (fake *FakeISCSIConnector) DisconnectVolume(connectionProperty connector.ConnectionProperty) error {
	return nil
}
func (fake *FakeISCSIConnector) ExtendVolume(connectionProperty connector.ConnectionProperty) error {
	return nil
}

func (fake *FakeISCSIConnector) LoginPortal(targetPortal string, targetIqn string) error {
	return nil
}

func (fake *FakeISCSIConnector) DiscoverPortal(targetPortal ...string) []model.ISCSISession {
	if len(targetPortal) > 0 && targetPortal[0] == "10.244.244.244" {
		return []model.ISCSISession{
			{TargetPortal: "10.244.244.244", TargetIqn: "iqn.1992-05.com.redhat:sl7b92030000520000", Tag: "1"},
		}
	}
	return nil
}

// Testing

func TestSession2ConnectionProperty(t *testing.T) {

	session1 := model.ISCSISession{
		TargetIqn:    "iqn.1992-05.com.emc:sl7b92030000520000-2",
		TargetPortal: "192.168.0.10:3260",
		TargetIp:     "192.168.0.10",
		Tag:          "1",
	}

	conn := Session2ConnectionProperty([]model.ISCSISession{session1}, 4)
	assert.Equal(t, []int{4}, conn.TargetLuns)
	assert.Equal(t, []string{"iqn.1992-05.com.emc:sl7b92030000520000-2"}, conn.TargetIqns)
	assert.Equal(t, []string{"192.168.0.10:3260"}, conn.TargetPortals)
}

func TestHandleIscsiParam(t *testing.T) {
	fake := &FakeISCSIConnector{}
	SetISCSIConnector(fake)

	err := HandleISCSIConnect("192.168.1.17", "33")
	assert.Nil(t, err)
}
func TestHandleIscsiParamNoLun(t *testing.T) {
	fake := &FakeISCSIConnector{}
	SetISCSIConnector(fake)

	err := HandleISCSIConnect("192.168.1.19")
	assert.Nil(t, err)
}

func TestHandleIscsiParamNoLunFailed(t *testing.T) {
	fake := &FakeISCSIConnector{}
	SetISCSIConnector(fake)

	err := HandleISCSIConnect("10.244.244.244")
	assert.Error(t, err, "Failed to connect volume.")
}

func TestHandleIscsiParamInvalidLun(t *testing.T) {
	fake := &FakeISCSIConnector{}
	SetISCSIConnector(fake)

	err := HandleISCSIConnect("192.168.2.17", "invalid")
	assert.Error(t, err)
}

func TestHandleIscsiParamNoParam(t *testing.T) {
	fake := &FakeISCSIConnector{}
	SetISCSIConnector(fake)
	err := HandleISCSIConnect()
	assert.Error(t, err)
}

func TestBeautifyVolumeInfo(t *testing.T) {
	info := connector.VolumeInfo{Paths: []string{"/dev/disk/by-path/xxxxxxxxxxxxxxxx", "/dev/disk/by-path/yyyyyyyyyyyyyyyyyy"},
		MultipathId: "351160160b6e00e5a50060160b6e00e5a", Wwn: "351160160b6e00e5a50060160b6e00e5a",
		Multipath: "/dev/mapper/351160160b6e00e5a50060160b6e00e5a"}
	BeautifyVolumeInfo(info)
}
