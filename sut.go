/*
 * sut.go - file defining SysUnderTest struct and its methods
 * 
 * SUT is just descriptive structure that keeps some information about the
 * TestSet currently executed (used in configuration and in reports), it 
 * doesn't have any influence on execution.
 */

package atf

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Represents a system under test: this either peice of SW or HW or a system
// built from both HW and SW. 
type SysUnderTest struct {

    // name of the SUT
	Name        string `xml:"name,attr"`

    // SUT System type: basically distinstion between HW and SW...
	Systype     string `xml:"Type"`

    // SUT version string (basically SUT HW or SW version)
	Version     string `xml:"Version"`

    // SUT description text
	Description string `xml:"Description"`

    // SUT IP address (if needed)
	IPaddr      string `xml:"IPAddress"`
}

// Create a new SUT instance.
func CreateSUT(name, systype, version, descr, ip string) *SysUnderTest {
	return &SysUnderTest{name, systype, version, descr, ip}
}

// Returns a plain text representation of the SUT instance.
func (s *SysUnderTest) String() string {
	txt := "SystemUnderTest:\n"
	txt += fmt.Sprintf("   Name: %s\n", s.Name)
	txt += fmt.Sprintf("   Type: %s\n", s.Systype)
	txt += fmt.Sprintf("   Version: %s\n", s.Version)
	txt += fmt.Sprintf("   IP address: %s\n", s.IPaddr)
	txt += fmt.Sprintf("   Description:\n%s", s.Description)
	return txt
}

// Returns a XML-encoded representation of the SUT instance.
func (s *SysUnderTest) Xml() (string, error) {

    output, err := xml.MarshalIndent(s, "  ", "    ")
    if err != nil {
        return "", err
    }
    return string(output), nil
}

// Returns an JSON-encoded representation of the SUT instance.
func (s *SysUnderTest) Json() (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}
