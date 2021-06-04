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

package connector

import (
	"github.com/peter-wangxu/goock/pkg/model"
	"github.com/peter-wangxu/goock/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectionProperty_IsEmpty(t *testing.T) {
	prop := ConnectionProperty{}
	assert.Error(t, prop.IsEmpty())
}

func TestConnectionProperty_IsEmptyFalse(t *testing.T) {
	prop := ConnectionProperty{
		StorageProtocol: IscsiProtocol,
		TargetPortals:   []string{"192.168.1.2:3260"},
		TargetLuns:      []int{11}}
	assert.Nil(t, prop.IsEmpty())
}

func TestConnectionProperty_FC_IsEmptyFalse(t *testing.T) {
	prop := ConnectionProperty{
		StorageProtocol: FcProtocol,
		TargetWwns:      []string{"fakeWwns"},
		TargetLuns:      []int{11}}
	assert.Nil(t, prop.IsEmpty())
}

func TestFormatLuns(t *testing.T) {
	formated := FormatLuns(10, 265)
	assert.Equal(t, "10", formated[0])
	assert.Equal(t, "0x0109000000000000", formated[1])
}

func TestGetHostInfo(t *testing.T) {
	SetExecutor(test.NewMockExecutor())
	model.SetExecutor(test.NewMockExecutor())
	info, _ := GetHostInfo()
	assert.NotEmpty(t, info.Hostname)
	assert.NotEmpty(t, info.OSType)
	assert.Equal(t, "iqn.1993-08.org.debian:01:b974ee37fea", info.Initiator)
	assert.Equal(t, []string{"20000090fa534cd0", "20000090fa534cd1"}, info.Wwnns)
	assert.Equal(t, []string{"10000090fa534cd0", "10000090fa534cd1"}, info.Wwpns)
	assert.Equal(t, []string{"5006016089200925", "50060160b6e00e5a",
		"5006016089200925", "50060160b6e00e5a"}, info.TargetWwnns)
	assert.Equal(t, []string{"5006016d09200925", "5006016036e00e5a",
		"5006016509200925", "5006016136e00e5a"}, info.TargetWwpns)
}
