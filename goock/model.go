package main

import (
	"fmt"
	"regexp"
)

type Parser interface {
	Parse(output string, raw string) []map[string]string
	filter(item map[string]string) bool
	splitter() string
}

type DefaultParser struct {
}

func (p *DefaultParser) Parse(output string, raw string) []map[string]string {
	pattern, _ := regexp.Compile(raw)
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

// Model decaration
type Model interface {
	GetPattern() string
	GetCommand() []string
	getOutput() string
	GetValue(key string) string
	setValue(key string, value string)
	setParser(parser Parser)
	Parse() []interface{}
}

// Subclass of Model
type ISCSISession struct {
	data_map      map[string]string
	target_iqn    string
	target_portal string
	target_ip     string
	tag           string
	parser        Parser
}

func (s *ISCSISession) GetPattern() string {
	return "\\s*(?P<target_portal>\\S+),(?P<tag>\\d+)\\s+(?P<target_iqn>\\S+)"
}

func (s *ISCSISession) GetCommand() []string {
	return []string{"iscsiadm", "-m", "session"}
}

func (s *ISCSISession) GetValue(key string) string {
	return s.data_map[key]
}

func (s *ISCSISession) setValue(key string, value string) {
	switch key {
	case "tag":
		s.tag = value
	case "target_iqn":
		s.target_iqn = value
	case "target_portal":
		s.target_portal = value
	default:
		fmt.Println("Invalid property sepecified.")
	}
}

func (s *ISCSISession) Parse() []ISCSISession {
	// Why successful?
	parser := s.parser.(*DefaultParser)
	data_list := parser.Parse(s.getOutput(), s.GetPattern())
	list := make([]ISCSISession, len(data_list))
	for i, each := range data_list {
		s := NewISCSISession()
		for k, v := range each {
			s.setValue(k, v)
		}
		list[i] = *s
	}
	return list
}

func (s *ISCSISession) getOutput() string {
	return `10.64.76.253:3260,1 iqn.1992-04.com.emc:cx.fcnch097ae5ef3.h1
 11.64.76.253:3260,1 iqn.1992-04.com.emc:cx.fcnch097ae6ef3.h2
asdfasd
asdfasdfasdf`
}
func NewISCSISession() *ISCSISession {
	return &ISCSISession{parser: &DefaultParser{}}
}

// Split by regexp specified by delimeter
func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}

func MyTest(m Model) {
	fmt.Printf("Peter: %s", m.GetValue("target_iqn"))
}

func main() {
	list := NewISCSISession().Parse()
	for i, each := range list {
		fmt.Printf("target_iqn[%d]: %s\t", i, each.target_iqn)
		fmt.Printf("target_portal[%d]: %s\n", i, each.target_portal)
		//MyTest(&each), TODO why cannot
	}

}
