package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAppConnect(t *testing.T) {
	goockApp := NewApp()
	err := goockApp.Run([]string{"goock", "connect", "192.168.1.8"})
	assert.IsType(t, &App{}, goockApp)
	assert.Error(t, err, "No path found")
}

func TestNewAppDisconnect(t *testing.T) {
	goockApp := NewApp()
	goockApp.Run([]string{"goock", "disconnect"})
	assert.IsType(t, &App{}, goockApp)
}

func TestNewAppInfo(t *testing.T) {
	goockApp := NewApp()
	goockApp.Run([]string{"goock", "info"})
	assert.IsType(t, &App{}, goockApp)
}
