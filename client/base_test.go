package client

import (
	"github.com/peter-wangxu/goock/connector"
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/test"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
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
	HandleInfo()
}
func TestHandleInfoFailed(t *testing.T) {
	HandleInfo()
}
