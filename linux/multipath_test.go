package linux

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/peter-wangxu/goock/test"
)

func TestFlushPath(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	err := FlushPath("36006016074e03a00ee762958673eaf1b")
	assert.Nil(t, err)

	err = FlushPath("")
	assert.Nil(t, err)
}

func TestGetPaths(t *testing.T) {
	SetExecutor(test.NewMockExecutor())


}

func TestReconfigure(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	success := Reconfigure()
	assert.True(t, success)
}

func TestReload(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	err := Reload()
	assert.Nil(t, err)

}


func TestCheckDevice(t *testing.T) {

}


func TestResizeMpath(t *testing.T) {

}
