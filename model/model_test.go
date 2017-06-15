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
	hostId, err := hbas[0].GetHostId()
	assert.Nil(t, err)
	assert.Equal(t, 7, hostId)
	assert.Equal(t, "", hbas[0].Path)
	assert.Equal(t, "100050eb1a033f59", hbas[0].FabricName)
	assert.Equal(t, "20000090fa534cd0", hbas[0].NodeName)
	assert.Equal(t, "10000090fa534cd0", hbas[0].PortName)
	assert.Equal(t, "Online", hbas[0].PortState)
	assert.Equal(t, "8 Gbit", hbas[0].Speed)
	assert.Equal(t, "4 Gbit, 8 Gbit, 16 Gbit", hbas[0].SupportedSpeeds)
	assert.Equal(t, "/sys/devices/pci0000:00/0000:00:03.0/0000:05:00.0/host7", hbas[0].DevicePath)
	assert.Equal(t, 2, len(hbas))
}

func TestNewFibreChannelTarget(t *testing.T) {
	executor = test.NewMockExecutor()

	targets := NewFibreChannelTarget()
	assert.Len(t, targets, 4)
	assert.Equal(t, "0:0", targets[0].ClassDevice)
	assert.Equal(t, "/sys/devices/pci0000:00/0000:00:03.0/0000:05:00.1/host9/rport-9:0-2/target9:0:0/fc_transport/target9:0:0", targets[0].ClassDevicePath)
	assert.Equal(t, "0x020500", targets[0].PortId)
	assert.Equal(t, "5006016089200925", targets[0].NodeName)
	assert.Equal(t, "5006016d09200925", targets[0].PortName)
	assert.Equal(t, "target9:0:0", targets[0].Device)
	assert.Equal(t, "/sys/devices/pci0000:00/0000:00:03.0/0000:05:00.1/host9/rport-9:0-2/target9:0:0", targets[0].DevicePath)
	hcl, err := targets[0].GetHostChannelTarget()
	assert.Nil(t, err)
	assert.EqualValues(t, []int{9, 0, 0}, hcl)
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
	for _, path := range m.Paths {
		assert.Regexp(t, "^\\w+$", path.DevNode)
	}
	// check action
	m1 := multipaths[4]
	assert.Equal(t, "reload", m1.Action)
}

// Test that multipath still works when WWN is longer thant 33 chars.
// Added for iscsitarget package
func TestFindMultipath(t *testing.T) {
	old := executor
	executor = test.NewMockExecutor()
	defer func() {
		executor = old
	}()
	multipaths := FindMultipath("149455400000000003592265eae69d00a6e8560cd2833744e")
	assert.Len(t, multipaths, 1)
	m := multipaths[0]
	assert.Equal(t, "149455400000000003592265eae69d00a6e8560cd2833744e", m.Wwn)
	assert.Equal(t, "dm-3", m.DmDeviceName)
	assert.Equal(t, 1.0, m.Size)
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
	assert.Equal(t, "scsi0", devices[0].Host)
	assert.Equal(t, 1, devices[0].Channel)
	assert.Equal(t, 0, devices[0].Target)
	assert.Equal(t, 0, devices[0].Lun)
	assert.Equal(t, 0, devices[0].GetHostId())
	assert.Equal(t, "0:1:0:0", devices[0].GetDeviceIdentifier())
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
