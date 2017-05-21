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
package model

import (
	"github.com/peter-wangxu/goock/exec"
	"github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var log *logrus.Logger = logrus.New()

func SetLogger(l *logrus.Logger) {
	log = l
}

var executor = exec.New()

func SetExecutor(e exec.Interface) {
	executor = e
}

type Parser interface {
	Parse(output string, pat interface{}) []map[string]string
	filter(item map[string]string) bool
	Split(output string) []string
}

// LineParse parse each line as a single data map
// The data map could be used to initialize a new Object derived from Interface
type LineParser struct {
	Delimiter string
	Matcher   string
}

func (p *LineParser) Parse(output string, pat interface{}) []map[string]string {
	pat_str := pat.(string)
	pattern, _ := regexp.Compile(pat_str)
	// Split into lines by splitter
	//lines := p.Split(output)

	ret := pattern.FindAllStringSubmatch(output, -1)
	if ret == nil {
		return []map[string]string{}
	}
	dataMap := make([]map[string]string, len(ret))
	matchedNames := pattern.SubexpNames()
	for i, each := range ret {
		data := make(map[string]string)
		for j := 1; j < len(each); j++ {
			n, v := matchedNames[j], each[j]
			data[n] = v
		}
		if p.filter(data) {
			dataMap[i] = data
		}
	}
	return dataMap

}

// Returns true if a valid item parsed
func (p *LineParser) filter(item map[string]string) bool {
	return true
}

func (p *LineParser) Split(output string) []string {
	var lines []string
	if p.Delimiter != "" {
		lines = RegSplit(output, p.Delimiter)
	} else if p.Matcher != "" {
		lines = RegMatcher(output, p.Matcher)
	}
	return lines
}

type PairParser struct {
	Delimiter string
	Matcher   string
}

func (f *PairParser) Parse(output string, pat interface{}) []map[string]string {

	// Split into lines by splitter
	lines := f.Split(output)
	dataMap := make([]map[string]string, len(lines))
	actual := 0
	for _, line := range lines {
		pat_list := pat.([]string)
		data := make(map[string]string)
		for _, m_pat := range pat_list {
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

func (f *PairParser) Split(output string) []string {
	var lines []string
	if f.Delimiter != "" {
		lines = RegSplit(output, f.Delimiter)
	} else if f.Matcher != "" {
		lines = RegMatcher(output, f.Matcher)
	}
	return lines
}

// Interface declaration
type Interface interface {
	GetPattern() interface{}
	GetCommand() []string
	getOutput() string
	GetValue(key string) interface{}
	setValue(key string, value interface{})
	setParser(parser Parser)
	Parse() []interface{}
}

// Implementation of ISCSISession
type ISCSISession struct {
	dataMap      map[string]string
	params       []string
	TargetIqn    string
	TargetPortal string
	TargetIp     string
	Tag          string
	parser       Parser
}

func (iscsi *ISCSISession) GetPattern() interface{} {
	return "\\s*(?P<TargetPortal>\\S+:\\d*),(?P<Tag>\\d+)\\s+(?P<TargetIqn>\\S+)"
}

func (iscsi *ISCSISession) GetCommand() []string {
	if len(iscsi.params) == 0 {

		return []string{"iscsiadm", "-m", "session"}
	} else {
		return append([]string{"iscsiadm"}, iscsi.params...)
	}
}

func (iscsi *ISCSISession) GetValue(key string) interface{} {
	return iscsi.dataMap[key]
}

func (iscsi *ISCSISession) setValue(key string, value interface{}) {
	ref := reflect.ValueOf(iscsi).Elem()
	SetValue(ref.FieldByName(key), value)

}

func (iscsi *ISCSISession) Parse() []ISCSISession {
	parser := iscsi.parser
	dataList := parser.Parse(iscsi.getOutput(), iscsi.GetPattern())
	list := make([]ISCSISession, len(dataList))
	for i, each := range dataList {
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
	if nil != err {
		return ""
	}
	return string(out[:])
}
func NewISCSISession() []ISCSISession {
	return (&ISCSISession{parser: &LineParser{Delimiter: "\\n+"}}).Parse()
}

// Discover all the targets provided by targetPortals
// Use goroutine to accelerate the target discovery process.

func DiscoverISCSISession(targetPortals []string) []ISCSISession {
	var results []ISCSISession
	c := make(chan []ISCSISession, len(targetPortals))
	for _, portal := range targetPortals {
		discovery := []string{
			"-m", "discovery", "-t", "sendtargets", "-I", "default", "-p", portal,
		}
		go func() {
			session := ISCSISession{parser: &LineParser{Delimiter: "\\n+"}}
			session.params = discovery
			ret := session.Parse()
			// Aggregate the results
			c <- ret
		}()

	}
	// Wait for all discovery
	var discoveredTargets []string
	for i := 0; i < len(targetPortals); i++ {
		each := <-c
		results = append(results, each...)
		for _, d := range each {
			discoveredTargets = append(discoveredTargets, d.TargetPortal)
		}
		log.WithFields(
			logrus.Fields{
				"Target":     targetPortals,
				"Discovered": strings.Join(discoveredTargets, ", ")}).Debug(
			"Discover result")
	}
	return results
}

// end of implementation of ISCSISession

// (HBA) Subclass of Interface
type HBA struct {
	dataMap         map[string]string
	parser          Parser
	Name            string
	Path            string
	FabricName      string
	NodeName        string
	PortName        string
	PortState       string
	Speed           string
	SupportedSpeeds string
	DevicePath      string
}

func (s *HBA) GetPattern() interface{} {

	return []string{
		"Class Device\\s+=\\s+\"(?P<Name>.*)\"",
		"Class Device path\\s+=\\s+\"(?P<DevicePath>.*)\"",
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

func (s *HBA) GetValue(key string) interface{} {
	return s.dataMap[key]
}

func (s *HBA) setValue(key string, value interface{}) {
	ref := reflect.ValueOf(s).Elem()
	SetValue(ref.FieldByName(key), value)
}

func (s *HBA) Parse() []HBA {
	parser := s.parser
	dataList := parser.Parse(s.getOutput(), s.GetPattern())
	list := make([]HBA, len(dataList))
	for i, each := range dataList {
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
	if nil != err {
		return ""
	}
	return string(out[:])
}

func NewHBA() []HBA {
	return (&HBA{parser: &PairParser{Delimiter: "\\n{3,}"}}).Parse()
}

// (Multipath) Subclass of Interface
// Each Multipath contains one or more SinglePath
type Multipath struct {
	dataMap map[string]string
	parser  Parser
	params  []string
	// reload or reject
	Action          string
	Wwn             string
	DmDeviceName    string
	Vendor          string
	Product         string
	Size            float64
	Features        string
	HWHandler       string
	WritePermission string
	Paths           []SinglePath
}

func (s *Multipath) GetPattern() interface{} {

	return `((?P<Action>\w+):\s+)?(?P<Wwn>\w{33,})\s+(?P<DmDeviceName>\w+-?\d*)\s+(?P<Vendor>\w+)?,(?P<Product>\w+)?\r?\nsize=(?P<Size>[\d\.]+)G\s+features='(?P<Features>.*)'\s+hwhandler='(?P<HWHandler>.*)'\s+wp=(?P<WritePermission>\w+)(?P<Paths>.*)`
}

func (s *Multipath) GetCommand() []string {
	if len(s.params) == 0 {
		return []string{"multipath", "-ll"}
	} else {
		return append([]string{"multipath"}, s.params...)
	}
}

func (s *Multipath) GetValue(key string) interface{} {
	return s.dataMap[key]
}

func (s *Multipath) setValue(key string, value interface{}) {
	ref := reflect.ValueOf(s).Elem()
	SetValue(ref.FieldByName(key), value)
}

func (s *Multipath) Parse() []Multipath {
	parser := s.parser
	mOutput := s.getOutput()
	dataList := parser.Parse(mOutput, s.GetPattern())
	list := make([]Multipath, len(dataList))
	pathGroups := RegMatcher(mOutput, "\\w{33,}")
	for i, each := range dataList {
		s := &Multipath{}
		for k, v := range each {
			s.setValue(k, v)
		}
		// SinglePath
		s.Paths = NewSinglePath(pathGroups[i])
		list[i] = *s
	}
	return list
}

func (s *Multipath) getOutput() string {
	cmd := s.GetCommand()
	out, err := executor.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if nil != err {
		return ""
	}
	return string(out[:])
}

func (s *Multipath) SetParams(params []string) {
	s.params = params
}

func NewMultipath() []Multipath {
	return (&Multipath{parser: &LineParser{Matcher: "(\\w+:\\s+)?\\w{33,}"}}).Parse()
}

func FindMultipath(path string) []Multipath {
	m := &Multipath{parser: &LineParser{Matcher: "(\\w+:\\s+)?\\w{33,}"}}
	m.SetParams([]string{"-l", path})
	return m.Parse()
}

// (SinglePath) Subclass of Interface
type SinglePath struct {
	dataMap map[string]string
	parser  Parser
	output  string
	//Policy string
	//Priority string
	Host    int
	Channel int
	Id      int
	Lun     int
	DevNode string
	Major   int
	Minor   int
	// possible value: failed, active
	DmStatus string
	// possible value: ready, ghost, faulty, shaky
	PathStatus string
	// possible value: running, offline
	OnlineStatus string
}

func (single *SinglePath) GetPattern() interface{} {
	return `\-\s+(?P<Host>\d+):(?P<Channel>\d+):(?P<Id>\d+):(?P<Lun>\d+)\s+(?P<DevNode>\w+)\s+(?P<Major>\d+):(?P<Minor>\d+)\s+(?P<DmStatus>\w+)\s+(?P<PathStatus>\w+)\s+(?P<OnlineStatus>\w+)`

}

func (single *SinglePath) GetCommand() []string {
	return []string{}
}

func (single *SinglePath) getOutput() string {
	return single.output
}

func (single *SinglePath) SetOutput(output string) {
	single.output = output
}

func (single *SinglePath) GetValue(key string) interface{} {
	return single.dataMap[key]
}

func (single *SinglePath) setValue(key string, value interface{}) {
	ref := reflect.ValueOf(single).Elem()
	SetValue(ref.FieldByName(key), value)
}

func (single *SinglePath) setParser(parser Parser) {
	single.parser = parser
}

func (single *SinglePath) Parse() []SinglePath {
	parser := single.parser
	dataList := parser.Parse(single.getOutput(), single.GetPattern())
	list := make([]SinglePath, len(dataList))
	for i, each := range dataList {
		s := &SinglePath{}
		for k, v := range each {
			s.setValue(k, v)
		}
		list[i] = *s
	}
	return list
}

func NewSinglePath(output string) []SinglePath {
	rS := &SinglePath{parser: &LineParser{Delimiter: "\\n+"}}
	rS.SetOutput(output)
	return rS.Parse()
}

// DeviceInfo: subclass of Interface

type DeviceInfo struct {
	dataMap map[string]string
	parser  Parser
	params  []string
	Device  string
	Host    string
	// numbered host
	HostNumber int
	Channel    int
	Target     int
	Lun        int
}

func (d *DeviceInfo) GetPattern() interface{} {
	return `(?P<Device>.*):\s+(?P<Host>\w+)\s+channel=(?P<Channel>\d+)\s+id=(?P<Target>\d+)\s+lun=(?P<Lun>\d+)`
}

func (d *DeviceInfo) GetCommand() []string {
	if len(d.params) <= 0 {
		return []string{"sg_scan"}
	} else {
		return append([]string{"sg_scan"}, d.params...)
	}
}

func (d *DeviceInfo) getOutput() string {
	cmd := d.GetCommand()
	out, err := executor.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if nil != err {
		log.Debug("Failed to get device info: ", out)
	}
	return string(out[:])
}

func (d *DeviceInfo) GetValue(key string) interface{} {
	return d.dataMap[key]
}

func (d *DeviceInfo) setValue(key string, value interface{}) {
	ref := reflect.ValueOf(d).Elem()
	SetValue(ref.FieldByName(key), value)
}

func (d *DeviceInfo) setParser(parser Parser) {
	d.parser = parser
}

func (d *DeviceInfo) Parse() []DeviceInfo {
	parser := d.parser
	dataList := parser.Parse(d.getOutput(), d.GetPattern())
	list := make([]DeviceInfo, len(dataList))
	for i, each := range dataList {
		s := &DeviceInfo{}
		for k, v := range each {
			s.setValue(k, v)
		}
		list[i] = *s
	}
	return list
}

func NewDeviceInfo(path string) []DeviceInfo {
	rS := &DeviceInfo{parser: &LineParser{Delimiter: "\\n+"}, params: []string{path}}
	return rS.Parse()
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
	result[len(indexes)] = text[lastStart:]
	return result
}

// Return slices of the matched string
func RegMatcher(text string, matcher string) []string {
	reg := regexp.MustCompile(matcher)

	matchedIndex := reg.FindAllStringSubmatchIndex(text, -1)
	var result []string
	//matchedCount := len(matchedIndex)
	if matchedIndex == nil {
		return []string{}
	}
	if len(matchedIndex) == 1 {
		return []string{text}
	}

	for j := 0; j < len(matchedIndex)-1; j++ {
		m := text[matchedIndex[j][0]:matchedIndex[j+1][0]]
		result = append(result, m)
	}
	result = append(result, text[len(matchedIndex)-1:])
	return result
}

func SetValue(field reflect.Value, value interface{}) {
	if field.CanSet() {
		switch field.Kind() {
		case reflect.Int:
			i, err := strconv.ParseInt(value.(string), 10, 64)
			if err != nil {
				i = 0
			}
			field.SetInt(i)
		case reflect.Float64:
			f, err := strconv.ParseFloat(value.(string), 2)
			if err != nil {
				f = 0.0
			}
			field.SetFloat(f)
		case reflect.String:
			field.SetString(value.(string))

		case reflect.Bool:
			b, err := strconv.ParseBool(value.(string))
			if err != nil {
				b = false
			}
			field.SetBool(b)
		default:
			log.Debug("Unsupported data type ", value, "for field ", field.String())
		}
	} else {

		log.Debugf("Invalid value [%v] specified for property [%s]", value, field.String())
	}
}
