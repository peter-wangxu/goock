package connector

type StringEnum string

const (
	READWRITE StringEnum = "rw"
	READONLY  StringEnum = "ro"
)

const (
	ISCSI_PROTOCOL StringEnum = "iscsi"
	FC_PROTOCOL    StringEnum = "fc"
)

type ConnectionProperty struct {
	TargetIqns      []string
	TargetPortals   []string
	TargetLuns      []int
	StorageProtocol string
	AccessMode      StringEnum
}

type HostInfo struct {
	Initiator string
	Ip        string
	Hostname  string
	OSType    string
}

type DeviceInfo struct {
	MultipathId string
	Paths       []string
	Wwn         string
	Multipath   string
}

type Interface interface {
	ConnectVolume(connectionProperty ConnectionProperty) (DeviceInfo, error)
	DisconnectVolume(connectionProperty ConnectionProperty) error
	ExtendVolume(connectionProperty ConnectionProperty) error
}
