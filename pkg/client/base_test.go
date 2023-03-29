package client

import (
	"testing"

	"github.com/peter-wangxu/goock/pkg/connector"
	"github.com/peter-wangxu/goock/pkg/exec"
	"github.com/peter-wangxu/goock/pkg/model"
	"github.com/peter-wangxu/goock/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitLogFalse(t *testing.T) {
	err := InitLog(false)
	assert.Nil(t, err)
	assert.Equal(t, logrus.InfoLevel, log.Level)
}
func TestInitLogTrue(t *testing.T) {
	err := InitLog(true)
	assert.Nil(t, err)
	assert.Equal(t, logrus.DebugLevel, log.Level)
}

func TestHandleExtendEmpty(t *testing.T) {
	err := HandleExtend()
	assert.Error(t, err)
}

func TestHandleExtendLocal(t *testing.T) {
	err := HandleExtend("/dev/sdm")
	assert.Error(t, err)
}

func TestHandleInfo(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	connector.SetExecutor(test.NewMockExecutor())
	err := HandleInfo()
	assert.Nil(t, err)
}
func TestHandleInfoFailed(t *testing.T) {
	model.SetExecutor(exec.New())
	connector.SetExecutor(exec.New())
	err := HandleInfo()
	assert.Error(t, err)
}
func TestValidateLunId_True(t *testing.T) {
	lunids, err := ValidateLunID([]string{"12", "113"})
	assert.Nil(t, err)
	assert.Len(t, lunids, 2)
}

func TestValidateLunId_False(t *testing.T) {
	lunids, err := ValidateLunID([]string{"a", "b"})
	assert.Error(t, err)
	assert.Len(t, lunids, 0)
}

func TestIsLunLike_False(t *testing.T) {
	r := IsLunLike("192.168.1.30")
	assert.False(t, r)
}

func TestIsLunLike_True(t *testing.T) {
	r := IsLunLike("192")
	assert.True(t, r)
}

func TestIsIpLike_True(t *testing.T) {
	r := IsIPLike("192.168.1.29")
	assert.True(t, r)
}

func TestIsIpLike_False(t *testing.T) {
	r := IsIPLike("192.168.1.")
	assert.False(t, r)

	r = IsIPLike("129")
	assert.False(t, r)
}

func TestIsFcLike_Wwpn(t *testing.T) {
	r := IsFcLike("5006016036e00e5a")
	assert.True(t, r)
}

func TestIsFcLike_wwn(t *testing.T) {
	r := IsFcLike("5006016089200925:5006016d09200925")
	assert.True(t, r)
}

func TestIsFcLike_False(t *testing.T) {
	r := IsFcLike("134134dda")
	assert.False(t, r)
}
