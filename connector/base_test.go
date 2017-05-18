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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnectionProperty_IsEmpty(t *testing.T) {
	prop := ConnectionProperty{}
	assert.Error(t, prop.IsEmpty())
}

func TestConnectionProperty_IsEmptyFalse(t *testing.T) {
	prop := ConnectionProperty{TargetPortals: []string{"192.168.1.2:3260"}, TargetLuns: []int{11}}
	assert.Nil(t, prop.IsEmpty())
}
