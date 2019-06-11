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
package test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/peter-wangxu/goock/exec"
	"github.com/stretchr/testify/mock"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
)

type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Command(cmd string, args ...string) exec.Cmd {
	return &MockCmd{Path: cmd, Args: args}
}

func (m *MockExecutor) LookPath(file string) (string, error) {
	return "", nil
}

func NewMockExecutor() exec.Interface {
	return &MockExecutor{}
}

type MockCmd struct {
	mock.Mock
	Path string
	Args []string
	Env  []string
	// Add more properties if mock needed
	Stdin []string
}

func (m *MockCmd) SetDir(dir string) {

}

func (m *MockCmd) SetStdin(in io.Reader) {
	var msg string
	if b, err := ioutil.ReadAll(in); err == nil {
		msg = string(b)
	}
	m.Stdin = strings.Split(msg, " ")
}

func (m *MockCmd) SetStdout(out io.Writer) {

}

func (m *MockCmd) CombinedOutput() ([]byte, error) {
	return m.mockOutput()
}

func (m *MockCmd) Output() ([]byte, error) {
	return m.mockOutput()
}

// This function returns the mocked output according
// to the joined commands and it's parameters.
func (m *MockCmd) mockOutput() ([]byte, error) {
	var cmds []string
	cmds = append(cmds, m.Path)
	cmds = append(cmds, m.Args...)
	// Append the Stdin as Args
	cmds = append(cmds, m.Stdin...)
	fileName := strings.Join(cmds, "_")
	// some commands contain "/" or "\" which may interfere the mock file
	// need to replace it _
	fileName = strings.Replace(fileName, "/", "_", -1)
	fileName = strings.Replace(fileName, "\\", "_", -1)
	fileName = strings.Replace(fileName, ":", "_", -1)
	fileName = path.Join(getMockDir(), fmt.Sprintf("%s.txt", fileName))

	// open a file
	if file, err := os.Open(fileName); err == nil {

		// make sure it gets closed
		defer func() {
			if file != nil {
				errClose := file.Close()
				if errClose != nil {
					fmt.Printf("failed to close file: %s: %s", file.Name(), errClose.Error())
				}
			}
		}()
		fmt.Printf("Reading mock file: %s\n", fileName)

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		var buffer bytes.Buffer
		// Default status code to -1 to prevent unexpected success.
		cmdStatus := "-1"
		isFirstLine := true
		for scanner.Scan() {
			if isFirstLine {
				cmdStatus = scanner.Text()
				isFirstLine = false
			} else {
				buffer.WriteString(scanner.Text() + "\n")
			}
		}
		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
		if cmdStatus == "0" {
			return []byte(buffer.String()), nil
		} else {
			cmdError := errors.New("Status code is " + cmdStatus)
			return []byte(buffer.String()), cmdError
		}

	} else {
		fmt.Printf("unable to read mock data from file %s(%q), default to empty string.\n", fileName, err.Error())
		return []byte(""), fmt.Errorf("unable to read mock file [%s]", fileName)
	}

}

var mockDataDir string

func init() {
	_, filename, _, _ := runtime.Caller(0)
	mockDataDir = path.Join(path.Dir(filename), "mock_data")
}

func getMockDir() string {
	return mockDataDir
}
