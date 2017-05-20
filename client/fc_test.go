package client

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHandleFCConnect(t *testing.T) {
	err := HandleFCConnect()
	assert.Nil(t, err)
}

func TestHandleFCExtend(t *testing.T) {
	err := HandleFCExtend()
	assert.Nil(t, err)
}
