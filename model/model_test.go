package model

import (
	"github.com/peter-wangxu/goock/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHBA(t *testing.T) {
	old := executor
	executor = test.NewMockExecutor()
	defer func() {
		executor = old
	}()
	hbas := NewHBA()
	assert.Equal(t, "host7", hbas[0].Name)
	assert.Equal(t, "0x10000090fa534cd0", hbas[0].PortName)
	assert.Equal(t, 2, len(hbas))
}

func TestNewISCSISession(t *testing.T) {
	old := executor
	executor = test.NewMockExecutor()
	defer func() {
		executor = old
	}()
	sessions := NewISCSISession()
	assert.Equal(t, 2, len(sessions))
	assert.Contains(t, sessions[1].TargetIqn, "iqn.1992-04.com.emc:cx.fcnch097ae6ef3")
}

func TestNewMultipath(t *testing.T) {
	old := executor
	executor = test.NewMockExecutor()
	defer func() {
		executor = old
	}()
	multipaths := NewMultipath()
	assert.Equal(t, 5, len(multipaths))
	m := multipaths[0]
	assert.Equal(t, "36006016074e03a003dbe2a580510610a", m.Wwn)
	assert.Equal(t, "dm-17", m.DmDeviceName)
	assert.Equal(t, "DGC", m.Vendor)
	assert.Equal(t, "VRAID", m.Product)
	assert.Equal(t, 1.0, m.Size)
	assert.Equal(t, "2 queue_if_no_path retain_attached_hw_handler", m.Features)
	assert.Equal(t, "1 alua", m.HWHandler)
	assert.Equal(t, "rw", m.WritePermission)
	assert.Equal(t, 3, len(m.Paths))
	// check action
	m1 := multipaths[4]
	assert.Equal(t, "reload", m1.Action)
}


func TestDiscoverISCSISession(t *testing.T) {

	old := executor
	executor = test.NewMockExecutor()
	defer func() {
		executor = old
	}()
	discovered := DiscoverISCSISession([]string{"10.244.213.177", "10.244.213.179"})
	assert.Len(t, discovered, 4)
	assert.Equal(t, "iqn.1992-04.com.emc:cx.fnm00150600267.a0", discovered[0].TargetIqn)
	assert.Equal(t, "10.244.213.177:3260", discovered[0].TargetPortal)
	assert.Equal(t, "2", discovered[0].Tag)
}
func TestNewDeviceInfo(t *testing.T) {
	old := executor
	executor = test.NewMockExecutor()
	defer func() {
		executor = old
	}()
	devices := NewDeviceInfo("/dev/sdb")
	assert.Len(t, devices, 1)
	assert.Equal(t, devices[0].Device, "/dev/sdb")
	assert.Equal(t, devices[0].Host, "scsi0")
	assert.Equal(t, devices[0].Channel, 1)
	assert.Equal(t, devices[0].Target, 0)
	assert.Equal(t, devices[0].Lun, 0)
}
func TestRegSplit(t *testing.T) {
	var s = `|-+- policy='round-robin 0' prio=50 status=active
		 | -- 9:0:0:10   sdm  8:192   active ready  running
		  --+- policy='round-robin 0' prio=10 status=enabled
	  	 |- 9:0:2:10   sdap 66:144  active ready  running
	         -- 13:0:0:10  sdcd 69:16   active ready  running`
	ret := RegSplit(s, "\\|-\\+-")
	assert.Len(t, ret, 2)
}

func TestRegMatcher(t *testing.T) {
	var s = `36006016074e03a003dbe2a580510610a dm-17 DGC,VRAID
size=1.0G features='2 queue_if_no_path retain_attached_hw_handler' hwhandler='1 alua' wp=rw
|-+- policy='round-robin 0' prio=50 status=active
| |- 9:0:2:25   sdbc 67:96   active ready  running
| -- 13:0:0:25  sdcs 70:0    active ready  running
--+- policy='round-robin 0' prio=10 status=enabled
  - 9:0:0:25   sdaa 65:160  active ready  running
Mar 28 23:02:16 | sdef: alua not supported
Mar 28 23:02:16 | sdeg: alua not supported
3600601601290380036a00936cf13e711 dm-30 DGC,VRAID
size=2.0G features='1 retain_attached_hw_handler' hwhandler='1 alua' wp=rw
|-+- policy='round-robin 0' prio=0 status=active
| -- 11:0:0:151 sdef 128:112 failed faulty running
--+- policy='round-robin 0' prio=0 status=enabled
-- 12:0:0:151 sdeg 128:128 failed faulty running
36006016074e03a008dfd94ce623d4c0e dm-2 DGC,VRAID
size=2.0G features='2 queue_if_no_path retain_attached_hw_handler' hwhandler='1 alua' wp=rw
|-+- policy='round-robin 0' prio=50 status=active
| -- 9:0:0:10   sdm  8:192   active ready  running
--+- policy='round-robin 0' prio=10 status=enabled
  |- 9:0:2:10   sdap 66:144  active ready  running
  -- 13:0:0:10  sdcd 69:16   active ready  running`
	matched := RegMatcher(s, "\\w{33}")
	assert.Equal(t, 3, len(matched))
	assert.Contains(t, matched[0], "36006016074e03a003dbe2a580510610a")
}