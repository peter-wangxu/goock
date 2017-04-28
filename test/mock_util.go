package test

import (
	"fmt"
	"github.com/peter-wangxu/goock/exec"
	"github.com/stretchr/testify/mock"
	"io"
	"os"
	"strings"
	"bufio"
	"log"
	"bytes"
	"errors"
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
}

func (m *MockCmd) SetDir(dir string) {

}

func (m *MockCmd) SetStdin(in io.Reader) {

}

func (m *MockCmd) SetStdout(out io.Writer) {

}

func (m *MockCmd) CombinedOutput() ([]byte, error) {
	return m.mockOutput()
	//return []byte(s), nil
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
	fileName := strings.Join(cmds, "_")
	// some commands contain "/" or "\" which may interfere the mock file
	// need to replace it _
	fileName = strings.Replace(fileName, "/", "_", -1)
	fileName = strings.Replace(fileName, "\\", "_", -1)
	fileName = strings.Replace(fileName, ":", "_", -1)
	fileName = fmt.Sprintf("%s%s.txt", getMockDir(), fileName)

	// open a file
	if file, err := os.Open(fileName); err == nil {

		// make sure it gets closed
		defer file.Close()
		fmt.Printf("Reading mock file: %s\n", fileName)

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		var buffer bytes.Buffer
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
		fmt.Printf("Unable to read mock data from file %s, default to empty string.\n", fileName)
		return []byte(""), nil
	}

}

func getMockDir() string {
	goPath := os.Getenv("GOPATH")
	goPath = strings.Split(goPath, string(os.PathListSeparator))[0]
	var goProject string
	if (goPath == "") {
		goProject, _ = os.Getwd()
	} else {
		goProject = fmt.Sprintf("%s/src/github.com/peter-wangxu/goock", goPath)
	}

	mockDir := fmt.Sprintf("%s/%s/%s/", goProject, "test", "mock_data")
	return mockDir
}
