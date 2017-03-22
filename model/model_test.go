package model

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"io"
	"github.com/peter-wangxu/goock/exec"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"fmt"
)


type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Command(cmd string, args ...string) exec.Cmd {
	return &MockCmd{Path:cmd, Args: args}
}

func (m *MockExecutor) LookPath(file string) (string ,error){
	return "", nil
}
type MockCmd struct {
	mock.Mock
	Path string
	Args []string
	Env []string
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
	return []byte{'a', 'b'}, nil
}

// This function returns the mocked output according
// to the joined commands and it's parameters.
func (m *MockCmd) mockOutput() ([]byte, error){
	var cmds []string
	cmds = append(cmds, m.Path)
	cmds = append(cmds, m.Args...)
	fileName := strings.Join(cmds, "_")
	file, err := os.Open("../testing/mock_data/" + fileName + ".txt")
	if(nil != err){
		fmt.Printf("Unable to read mock data from file %s, default to empty string.", fileName)
		return []byte(""), nil
	}
	fstate, _ := file.Stat()
	fsize := fstate.Size()
	mock_data := make([]byte, fsize)
	file.Read(mock_data)
	return mock_data, nil

}


func TestNewHBA(t *testing.T) {
	old := executor
	executor = new(MockExecutor)
	defer func (){
		executor = old
	}()
	hbas := NewHBA().Parse()
	assert.Equal(t, "host7", hbas[1].Path)
	assert.Equal(t, 2, len(hbas))
}

func TestNewISCSISession(t *testing.T) {
	old := executor
	executor = new(MockExecutor)
	defer func (){
		executor = old
	}()
	sessions := NewISCSISession().Parse()
	assert.Equal(t, 2, len(sessions))
	assert.Contains(t, sessions[1].TargetIqn, "iqn.1992-04.com.emc:cx.fcnch097ae6ef3")
}