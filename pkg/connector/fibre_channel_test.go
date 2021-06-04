package connector

import (
	"github.com/peter-wangxu/goock/pkg/linux"
	"github.com/peter-wangxu/goock/pkg/model"
	"github.com/peter-wangxu/goock/test"
	goockutil "github.com/peter-wangxu/goock/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFibreChannelConnector_ConnectVolume(t *testing.T) {
	goockutil.SetExecutor(test.NewMockExecutor())
	SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	linux.SetExecutor(test.NewMockExecutor())
	fc := NewFibreChannelConnector()
	fakeProperty := ConnectionProperty{}
	fakeProperty.TargetWwns = []string{
		"5006016d09200925",
		"5006016136e00e5a",
	}
	fakeProperty.TargetLun = 11
	info, err := fc.ConnectVolume(fakeProperty)
	assert.NotNil(t, info)
	assert.Nil(t, err)
	assert.Equal(t, "/dev/disk/by-id/dm-uuid-mpath-350060160b6e00e5a50060160b6e11317", info.Multipath)
	assert.Equal(t, "350060160b6e00e5a50060160b6e11317", info.MultipathId)
	assert.Equal(t, "350060160b6e00e5a50060160b6e11317", info.Wwn)
	assert.Equal(t,
		[]string{"/dev/disk/by-path/pci-0000:05:00.0-fc-0x5006016d09200925-lun-11",
			"/dev/disk/by-path/pci-0000:05:00.1-fc-0x5006016136e00e5a-lun-11"},
		info.Paths)
}
