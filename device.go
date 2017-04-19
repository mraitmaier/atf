package atf
/*
 * device.go - file defining Device struct and its methods
 *
 */

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// DeviceType is an enum defining different device types.
type DeviceType int
const (
    DevUnknown DeviceType = iota
    DevServer
    DevClient
    DevEthernet
    DevSwitch
    DevRouter
    DevPowerSwitch
    DevPhySwitch
    DevTrafficGenerator
    DevTrafficSniffer
    DevAttenuator
    DevODU
    DevCTR8300
    DevCTR8540
    DevCtr8560
    DevWTM3300
    DevWTM4100
    DevWTM4200
)

// Device is a generic interface for all types of devices
type Device interface {
    // TODO
}


// GenericDevice represents a system...
type GenericDevice struct {
	// Name of the SUT
	Name string `xml:"Name"`
	// SysType is a SUT System type: basically distinction between HW and SW...
	Dtype DeviceType `xml:"Type"`
	// Description is a SUT description text
	Description string `xml:"Description"`
    //
    Family string
    //
    Model string
    //
    Management []string
    //
    Location string
    // Is this device a DUT (Device under test)?
    IsDUT bool
}

// NewgenericDevice creates a new GenericDevice instance.
func NewGenericDevice(name string, dtype DeviceType) *GenericDevice {
	return &GenericDevice{
        Name: name,
        Dtype: dtype,
        Description: "",
        Family: "",
        Model: "",
        Management: nil,
        Location: make([]string, 0),
        IsDut: false,
        }
}

//
type EthernetDevice struct {
    //
    GenericDevice
    // Ports is a list of ports
    Ports []Port
}

// NewEthernetDevice creates a new EthernetDevice instance.
func NewEthernetDevice(name string) *EthernetDevice {
	return &EthernetDevice{
        Name: name,
        Dtype: DeviceType.DevEthernet,
        Description: "",
        Family: "",
        Model: "",
        Management: nil,
        Location: make([]string, 0),
        IsDut: false,
        Ports: make([]Port, 0),
        }
}

//
type Server struct {
    //
    GenericDevice
    // Ports is a list of ports
    URI string
}

// NewServer creates a new Server instance.
func NewServer(name string) *Server{
	return &Server{
        Name: name,
        Dtype: DeviceType.DevEthernet,
        Description: "",
        Family: "",
        Model: "",
        Management: nil,
        Location: make([]string, 0),
        IsDut: false,
        URI: '',
        }
}

// PortType is an enum defining a device port type.
type PortType int
const (
    PortTypeUnknown PortType = iota << 1
    PortCopper
    PortFiber
    PortHDX
    PortFDX
    Port1M
    Port10M
    Port100M
    Port1G
    Port10G
    Port40G
    Port100G
)

//
func (p PortType) IsCopper() bool { return p & p.PortCopper }
//
func (p PortType) IsFiber() bool { return p & p.PortFiber }
//
func (p PortType) IsHalfDuplex() bool { return p & p.PortHDX }
//
func (p PortType) IsFullDuplex() bool { return p & p.PortFDX }
//
func (p PortType) IsMegabit() bool { return p & p.Port1M }
//
func (p PortType) Is10Megabit() bool { return p & p.Port10M }
//
func (p PortType) Is100Megabit() bool { return p & p.Port100M }
//
func (p PortType) IsGigabit() bool { return p & p.Port1G }
//
func (p PortType) Is10Gigabit() bool { return p & p.Port10G }
//
func (p PortType) Is40Gigabit() bool { return p & p.Port40G }
//
func (p PortType) Is100Gigabit() bool { return p & p.Port100G }

// Port is ...
type Port struct {
    //
    Name string
    //
    Description string
    //
    PortType
}

// NewPort creates a new instance of Port.
func NewPort() *Port {
    return &Port{
        Name: '',
        Description: '',
        PortType: PortTypeUnknown
    }
}

// CreatePort creates a new instance of Port from known parameters.
func CreatePort(name, desc string, ptype PortType) *Port {
    return &Port{
        Name: name,
        Description: desc,
        PortType: ptype
    }
}

/*
// String returns a human-readable representation of the SUT instance.
func (s *SysUnderTest) String() string {

	txt := "SystemUnderTest:\n"
	txt += fmt.Sprintf("          Name: %s\n", s.Name)
	txt += fmt.Sprintf("          Type: %s\n", s.Systype)
	txt += fmt.Sprintf("       Version: %s\n", s.Version)
	txt += fmt.Sprintf("    IP address: %s\n", s.IPaddr)
	txt += fmt.Sprintf("   Description: \n%s\n", s.Description)
	txt += fmt.Sprintf("         is Up? %s\n", s.IsUp)
	return txt
}

// XML returns a XML-encoded representation of the SUT instance.
func (s *SysUnderTest) XML() (string, error) {

	output, err := xml.MarshalIndent(s, "  ", "    ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// JSON returns an JSON-encoded representation of the SUT instance.
func (s *SysUnderTest) JSON() (string, error) {

	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}
*/
