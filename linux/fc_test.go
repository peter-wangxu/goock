/*
Copyright 2017 The Goock Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package linux

import (
	"github.com/peter-wangxu/goock/model"
	"github.com/peter-wangxu/goock/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsFCSupport(t *testing.T) {
	r := IsFCSupport()
	assert.Equal(t, false, r)
}

func TestGetFCHBA(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	hbas := GetFCHBA()
	assert.Len(t, hbas, 2)
}

func TestGetFCWWPN(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	wwpns := GetFcWwpn()
	assert.Len(t, wwpns, 2)
	assert.Equal(t, "10000090fa534cd0", wwpns[0])
}
func TestGetFCWWNN(t *testing.T) {
	model.SetExecutor(test.NewMockExecutor())
	wwnns := GetFcWwnn()
	assert.Len(t, wwnns, 2)
	assert.Equal(t, "20000090fa534cd0", wwnns[0])
}

func TestRescanHosts(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	RescanHosts([][]int{{9, 0, 1}, {7, 1, 0}}, 10)
}

func TestIsFCDeviceTrue(t *testing.T) {

}

func TestIsFCDeviceFalse(t *testing.T) {

}
