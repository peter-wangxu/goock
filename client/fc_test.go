package client

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleFCConnectNoParam(t *testing.T) {
	err := HandleFCConnect()
	assert.Error(t, err)
}

func TestHandleFCConnect(t *testing.T) {
	err := HandleFCConnect("12")
	assert.Nil(t, err)
}

func TestHandleFCConnect2Param(t *testing.T) {
	err := HandleFCConnect("fc", "11")
	assert.Error(t, err)
}

func TestHandleFCExtend(t *testing.T) {
	err := HandleFCExtend()
	assert.Nil(t, err)
}
