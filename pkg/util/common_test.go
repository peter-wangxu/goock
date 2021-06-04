package util

import (
	"fmt"
	"github.com/peter-wangxu/goock/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func hook() {
	fmt.Println("Called!")
}

func TestWaitForAnyPathWithHook(t *testing.T) {

	SetExecutor(test.NewMockExecutor())
	r, err := WaitForAnyPath([]string{"/real/path"}, hook)
	assert.Equal(t, "/real/path", r)
	assert.Nil(t, err)
}

func TestWaitForPathWithHook_NotFound(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	r, err := WaitForAnyPath([]string{"/real/not_created_yet"}, hook)
	assert.Empty(t, r)
	assert.Error(t, err)
}

func TestWaitForAnyPath(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	r, err := WaitForAnyPath([]string{"/fake/path", "/fake/path2"}, nil)
	assert.Error(t, err)
	assert.Empty(t, r)
}
