package client

import (
	"github.com/peter-wangxu/goock/connector"
	"github.com/peter-wangxu/goock/linux"
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/test"
	"github.com/peter-wangxu/goock/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleFCConnectOnlyLun(t *testing.T) {
	err := HandleFCConnect("12")
	assert.Error(t, err)
}

func TestHandleFCConnect2Param(t *testing.T) {
	err := HandleFCConnect("fc", "11")
	assert.Error(t, err)
}

func TestHandleFCConnect(t *testing.T) {
	connector.SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	linux.SetExecutor(test.NewMockExecutor())
	util.SetExecutor(test.NewMockExecutor())
	err := HandleFCConnect("5006016d09200925", "11")
	assert.Nil(t, err)
}

func TestHandleFCExtend(t *testing.T) {
	err := HandleFCExtend()
	assert.Nil(t, err)
}
