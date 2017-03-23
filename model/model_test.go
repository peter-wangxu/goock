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
	hbas := NewHBA().Parse()
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
	sessions := NewISCSISession().Parse()
	assert.Equal(t, 2, len(sessions))
	assert.Contains(t, sessions[1].TargetIqn, "iqn.1992-04.com.emc:cx.fcnch097ae6ef3")
}
