package model

import (
	"fmt"
	"reflect"
	"regexp"
	"os/exec"
)

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
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if(nil != err){
		return ""
	}
	return string(out[:])
	//return `10.64.76.253:3260,1 iqn.1992-04.com.emc:cx.fcnch097ae5ef3.h1
 //11.64.76.253:3260,1 iqn.1992-04.com.emc:cx.fcnch097ae6ef3.h2
//asdfasd
//asdfasdfasdf`
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
	return `Class = "fc_host"

  Class Device = "host7"
  Class Device path = "/sys/devices/pci0000:00/0000:00:03.0/0000:05:00.0/host7/fc_host/host7"
    active_fc4s         = "0x00 0x00 0x01 0x00 0x00 0x00 0x00 0x01 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 "
    dev_loss_tmo        = "30"
    fabric_name         = "0x100050eb1a033f59"
    issue_lip           = <store method only>
    max_npiv_vports     = "255"
    maxframe_size       = "2048 bytes"
    node_name           = "0x20000090fa534cd0"
    npiv_vports_inuse   = "0"
    port_id             = "0x010e00"
    port_name           = "0x10000090fa534cd0"
    port_state          = "Online"
    port_type           = "NPort (fabric via point-to-point)"
    speed               = "8 Gbit"
    supported_classes   = "Class 3"
    supported_fc4s      = "0x00 0x00 0x01 0x00 0x00 0x00 0x00 0x01 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 "
    supported_speeds    = "4 Gbit, 8 Gbit, 16 Gbit"
    symbolic_name       = "Emulex LPe16002B-E FV1.1.21.8 DV11.0.0.10. HN:(none) OS:Linux"
    tgtid_bind_type     = "wwpn (World Wide Port Name)"
    uevent              =
    vport_create        = <store method only>
    vport_delete        = <store method only>

    Device = "host7"
    Device path = "/sys/devices/pci0000:00/0000:00:03.0/0000:05:00.0/host7"
      uevent              = "DEVTYPE=scsi_host"


  Class Device = "host9"
  Class Device path = "/sys/devices/pci0000:00/0000:00:03.0/0000:05:00.1/host9/fc_host/host9"
    active_fc4s         = "0x00 0x00 0x01 0x00 0x00 0x00 0x00 0x01 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 "
    dev_loss_tmo        = "30"
    fabric_name         = "0x10000027f8c7928a"
    issue_lip           = <store method only>
    max_npiv_vports     = "255"
    maxframe_size       = "2048 bytes"
    node_name           = "0x20000090fa534cd1"
    npiv_vports_inuse   = "0"
    port_id             = "0x020d00"
    port_name           = "0x10000090fa534cd1"
    port_state          = "Online"
    port_type           = "NPort (fabric via point-to-point)"
    speed               = "16 Gbit"
    supported_classes   = "Class 3"
    supported_fc4s      = "0x00 0x00 0x01 0x00 0x00 0x00 0x00 0x01 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 "
    supported_speeds    = "4 Gbit, 8 Gbit, 16 Gbit"
    symbolic_name       = "Emulex LPe16002B-E FV1.1.21.8 DV11.0.0.10. HN:(none) OS:Linux"
    tgtid_bind_type     = "wwpn (World Wide Port Name)"
    uevent              =
    vport_create        = <store method only>
    vport_delete        = <store method only>

    Device = "host9"
    Device path = "/sys/devices/pci0000:00/0000:00:03.0/0000:05:00.1/host9"
      uevent              = "DEVTYPE=scsi_host"
`
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

func MyTest(m Model) {
	fmt.Printf("Peter: %s", m.GetValue("target_iqn"))
}

func main() {
	list := NewISCSISession().Parse()
	for i, each := range list {
		fmt.Printf("target_iqn[%d]: %s\t", i, each.TargetIqn)
		fmt.Printf("target_portal[%d]: %s\n", i, each.TargetPortal)
		//MyTest(&each), TODO why cannot
	}
	list2 := NewHBA().Parse()
	for _, each := range list2{
		fmt.Println(each.Name)
		fmt.Println(each.FabricName)
		fmt.Println(each.PortName)
	}

}
