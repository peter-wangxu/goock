package linux

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/peter-wangxu/goock/test"
	"github.com/peter-wangxu/goock/model"
)

func TestGetDeviceInfo(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	info := GetDeviceInfo("/dev/sdb")
	assert.Equal(t, info.Device, "/dev/sdb")
}
func TestGetDeviceInfoNotFound(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	info := GetDeviceInfo("/dev/sdx")
	assert.Equal(t, info.Device, "")
}
