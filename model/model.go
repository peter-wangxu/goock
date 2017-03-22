package model

import (
	"fmt"
	"reflect"
	"regexp"
	"github.com/peter-wangxu/goock/exec"
)

var executor = exec.New()

type Parser interface {
	Parse(output string, pat interface{}) []map[string]string
	filter(item map[string]string) bool
	splitter() string
}

type DefaultParser struct {
}

func (p *DefaultParser) Parse(output string, pat interface{}) []map[string]string {
	pat_str := pat.(string)
	pattern, _ := regexp.Compile(pat_str)
	// Split into lines by splitter
	var lines []string
	if p.splitter() != "" {
		lines = RegSplit(output, p.splitter())
	} else {
		lines = RegSplit(output, "\\n+")

	}
	data_map := make([]map[string]string, len(lines))
	actual := 0
	for _, line := range lines {
		re := pattern.FindStringSubmatch(line)
		if re == nil {
			continue
		}
		match_names := pattern.SubexpNames()
		data := make(map[string]string)
		for j := 1; j < len(re); j++ {
			n, v := match_names[j], re[j]
			data[n] = v
		}
		if p.filter(data) {
			data_map[actual] = data
			actual++
		}
	}

	return data_map[:actual]

}

// Returns true if a valid item parsed
func (p *DefaultParser) filter(item map[string]string) bool {
	return true
}

func (p *DefaultParser) splitter() string {
	return "\\n+"
}

type PairParser struct {
}

func (f *PairParser) Parse(output string, pat interface{}) []map[string]string {

	// Split into lines by splitter
	var lines []string
	if f.splitter() != "" {
		lines = RegSplit(output, f.splitter())
	} else {
		lines = RegSplit(output, "\\n+")

	}
	dataMap := make([]map[string]string, len(lines))
	actual := 0
	for _, line := range lines {
		pat_list := pat.([]string)
		data := make(map[string]string)
		for _, m_pat := range(pat_list){
			pattern, _ := regexp.Compile(m_pat)
			re := pattern.FindStringSubmatch(line)
			if re == nil {
				continue
			}
			match_names := pattern.SubexpNames()
			for j := 1; j < len(re); j++ {
				n, v := match_names[j], re[j]
				data[n] = v
			}
		}
		if f.filter(data) {
			dataMap[actual] = data
			actual++
		}
	}

	return dataMap[:actual]
}

func (f *PairParser) filter(item map[string]string) bool {
	return true
}

func (f *PairParser) splitter() string {
	return "\\n{2,}"
}

// Model declaration
type Model interface {
	GetPattern()  interface{}
	GetCommand() []string
	getOutput() string
	GetValue(key string) string
	setValue(key string, value string)
	setParser(parser Parser)
	Parse() []interface{}
}

// Implementation of ISCSISession
type ISCSISession struct {
	dataMap      map[string]string
	TargetIqn    string
	TargetPortal string
	TargetIp     string
	Tag          string
	parser       Parser
}

func (iscsi *ISCSISession) GetPattern() interface{} {
	return "\\s*(?P<TargetPortal>\\S+),(?P<Tag>\\d+)\\s+(?P<TargetIqn>\\S+)"
}

func (iscsi *ISCSISession) GetCommand() []string {
	return []string{"iscsiadm", "-m", "session"}
}

func (iscsi *ISCSISession) GetValue(key string) string {
	return iscsi.dataMap[key]
}

func (iscsi *ISCSISession) setValue(key string, value string) {
	ref := reflect.ValueOf(iscsi).Elem()
	field := ref.FieldByName(key)
	if field.IsValid() && field.CanSet(){
		ref.FieldByName(key).SetString(value)
	} else {
		fmt.Println("Invalid property name: ", key)
	}



}

func (iscsi *ISCSISession) Parse() []ISCSISession {
	// Why successful?
	parser := iscsi.parser
	data_list := parser.Parse(iscsi.getOutput(), iscsi.GetPattern())
	list := make([]ISCSISession, len(data_list))
	for i, each := range data_list {
		s := &ISCSISession{}
		for k, v := range each {
			s.setValue(k, v)
		}
		list[i] = *s
	}
	return list
}

func (iscsi *ISCSISession) getOutput() string {
	cmd := iscsi.GetCommand()
	out, err := executor.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if(nil != err){
		return ""
	}
	return string(out[:])
}
func NewISCSISession() *ISCSISession {
	return &ISCSISession{parser: &DefaultParser{}}
}

// end of implementation of ISCSISession


// (HBA) Subclass of Model
type HBA struct {
	dataMap         map[string]string
	Name            string
	Path            string
	FabricName      string
	NodeName        string
	PortName        string
	PortState       string
	Speed           string
	SupportedSpeeds string
	DevicePath 	string
	parser          Parser

}

func (s *HBA) GetPattern() interface{} {

	return []string{
		"^Class Device path\\s+=\\s+\"(?P<DevicePath>.*)\"",
		"^\\s+Device\\s+=\\s+\"(?P<Name>.*)\"",
		"fabric_name\\s+=\\s+\"(?P<FabricName>.*)\"",
		"node_name\\s+=\\s+\"(?P<NodeName>.*)\"",
		"port_name\\s+=\\s+\"(?P<PortName>.*)\"",
		"port_state\\s+=\\s+\"(?P<PortState>.*)\"",
		"speed\\s+=\\s+\"(?P<Speed>.*)\"",
		"supported_speeds\\s+=\\s+\"(?P<SupportedSpeeds>.*)\"",
	}
}

func (s *HBA) GetCommand() []string {
	return []string{"systool", "-c", "fc_host", "-v"}
}

func (s *HBA) GetValue(key string) string {
	return s.dataMap[key]
}

func (s *HBA) setValue(key string, value string) {
	ref := reflect.ValueOf(s).Elem()
	if ref.FieldByName(key).CanSet() {
		ref.FieldByName(key).SetString(value)
	} else {
		fmt.Println("Invalid property sepecified.")
	}
}

func (s *HBA) Parse() []HBA {
	// Why successful?
	parser := s.parser
	data_list := parser.Parse(s.getOutput(), s.GetPattern())
	list := make([]HBA, len(data_list))
	for i, each := range data_list {
		s := &HBA{}
		for k, v := range each {
			s.setValue(k, v)
		}
		list[i] = *s
	}
	return list
}

func (s *HBA) getOutput() string {
	cmd := s.GetCommand()
	out, err := executor.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if(nil != err){
		return ""
	}
	return string(out[:])
}

func NewHBA() *HBA {
	return &HBA{parser: &PairParser{}}
}

// Split by regexp specified by delimiter
func RegSplit(text string, delimiter string) []string {
	reg := regexp.MustCompile(delimiter)
	indexes := reg.FindAllStringIndex(text, -1)
	lastStart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[lastStart:element[0]]
		lastStart = element[1]
	}
	result[len(indexes)] = text[lastStart:len(text)]
	return result
}

