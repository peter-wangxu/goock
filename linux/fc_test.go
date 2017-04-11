package linux

import (
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsFCSupport(t *testing.T) {
	r := IsFCSupport()
	assert.Equal(t, false, r)
}

func TestGetFCHBA(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	hbas := GetFCHBA()
	assert.Len(t, hbas, 2)
}

func TestGetFCWWPN(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	wwpns := GetFCWWPN()
	assert.Len(t, wwpns, 2)
	assert.Equal(t, "10000090fa534cd0", wwpns[0])
}
func TestGetFCWWNN(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	wwnns := GetFCWWNN()
	assert.Len(t, wwnns, 2)
	assert.Equal(t, "20000090fa534cd0", wwnns[0])
}

func TestRescanHosts(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	SetExecutor(test.NewMockExecutor())
	RescanHosts()
}
