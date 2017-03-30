package model

import (
	"fmt"
	"github.com/peter-wangxu/goock/exec"
	"reflect"
	"regexp"
	"github.com/Sirupsen/logrus"
	"strconv"
)

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
// The data map could be used to initialize a new Object derived from Model
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
	if (p.Delimiter != "") {
		lines = RegSplit(output, p.Delimiter)
	} else if (p.Matcher != "") {
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
	if (f.Delimiter != "") {
		lines = RegSplit(output, f.Delimiter)
	} else if (f.Matcher != "") {
		lines = RegMatcher(output, f.Matcher)
	}
	return lines
}


// Model declaration
type Model interface {
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

func (iscsi *ISCSISession) GetValue(key string) interface{} {
	return iscsi.dataMap[key]
}

func (iscsi *ISCSISession) setValue(key string, value interface{}) {
	ref := reflect.ValueOf(iscsi).Elem()
	field := ref.FieldByName(key)
	if field.IsValid() && field.CanSet() {
		//ref.FieldByName(key).Set(value)
		SetValue(ref.FieldByName(key), value)
	} else {
		fmt.Println("Invalid property name: ", key)
	}

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
func NewISCSISession() *ISCSISession {
	return &ISCSISession{parser: &LineParser{Delimiter:"\\n+"}}
}

// end of implementation of ISCSISession

// (HBA) Subclass of Model
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
	if ref.FieldByName(key).CanSet() {
		//ref.FieldByName(key).Set(value)
		SetValue(ref.FieldByName(key), value)
	} else {
		fmt.Println("Invalid property sepecified.")
	}
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

func NewHBA() *HBA {
	return &HBA{parser: &PairParser{Delimiter:"\\n{3,}"}}
}

// (Multipath) Subclass of Model
// Each Multipath contains one or more SinglePath
type Multipath struct {
	dataMap         map[string]string
	parser          Parser
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

	return `((?P<Action>\w+):\s+)?(?P<Wwn>\w{33})\s+(?P<DmDeviceName>\w+-?\d*)\s+(?P<Vendor>\w+),(?P<Product>\w+)\r?\nsize=(?P<Size>[\d\.]+)G\s+features='(?P<Features>.*)'\s+hwhandler='(?P<HWHandler>.*)'\s+wp=(?P<WritePermission>\w+)(?P<Paths>.*)`
}

func (s *Multipath) GetCommand() []string {
	return []string{"multipath", "-ll"}
}

func (s *Multipath) GetValue(key string) interface{} {
	return s.dataMap[key]
}

func (s *Multipath) setValue(key string, value interface{}) {
	ref := reflect.ValueOf(s).Elem()
	if ref.FieldByName(key).CanSet() {
		SetValue(ref.FieldByName(key), value)
	} else {
		fmt.Println("Invalid property sepecified.")
	}
}

func (s *Multipath) Parse() []Multipath {
	parser := s.parser
	mOutput := s.getOutput()
	dataList := parser.Parse(mOutput, s.GetPattern())
	list := make([]Multipath, len(dataList))
	pathGroups := RegMatcher(mOutput, "\\w{33}")
	for i, each := range dataList {
		s := &Multipath{}
		for k, v := range each {
			s.setValue(k, v)
		}
		// SinglePath
		s.Paths = NewSinglePath(pathGroups[i]).Parse()
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

func NewMultipath() *Multipath {
	return &Multipath{parser: &LineParser{Matcher:"(\\w+:\\s+)?\\w{33}"}}
}

// (SinglePath) Subclass of Model
type SinglePath struct {
	dataMap      map[string]string
	parser       Parser
	output       string
	//Policy string
	//Priority string
	Host         int
	Channel      int
	Id           int
	Lun          int
	DevNode      string
	Major        int
	Minor        int
	// possible value: failed, active
	DmStatus     string
	// possible value: ready, ghost, faulty, shaky
	PathStatus   string
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
	if ref.FieldByName(key).CanSet() {
		SetValue(ref.FieldByName(key), value)
	} else {
		fmt.Println("Invalid property sepecified.")
	}
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

func NewSinglePath(output string) *SinglePath {
	rS := &SinglePath{parser: &LineParser{Delimiter:"\\n+"}}
	rS.SetOutput(output)
	return rS
}

// Split by regexp specified by delimiter
func RegSplit(text string, delimiter string) []string {
	reg := regexp.MustCompile(delimiter)
	indexes := reg.FindAllStringIndex(text, -1)
	lastStart := 0
	result := make([]string, len(indexes) + 1)
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
	if (matchedIndex == nil) {
		return []string{}
	}
	if (len(matchedIndex) == 1) {
		return []string{text}
	}

	for j := 0; j < len(matchedIndex) - 1; j++ {
		m := text[matchedIndex[j][0]:matchedIndex[j + 1][0]]
		result = append(result, m)
	}
	result = append(result, text[len(matchedIndex) - 1:])
	return result
}

func SetValue(field reflect.Value, value interface{}) {
	switch field.Kind() {
	case reflect.Int:
		i, err := strconv.ParseInt(value.(string), 10, 64)
		if (err != nil) {
			i = 0
		}
		field.SetInt(i)
	case reflect.Float64:
		f, err := strconv.ParseFloat(value.(string), 2)
		if (err != nil) {
			f = 0.0
		}
		field.SetFloat(f)
	case reflect.String:
		field.SetString(value.(string))

	case reflect.Bool:
		b, err := strconv.ParseBool(value.(string))
		if (err != nil) {
			b = false
		}
		field.SetBool(b)
	default:
		logrus.Debug("Unsupported data type ", value, "for field ", field.String())
	}
}
