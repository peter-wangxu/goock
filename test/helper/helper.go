package helper

import (
	"runtime"
	"testing"
)

func SkipIfWindows(t *testing.T) {
	osName := runtime.GOOS
	if osName == "windows" {
		t.Skip("Test case is skipped, as it's for non-windows platform.")
	}
}
