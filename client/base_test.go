package client

import (
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
