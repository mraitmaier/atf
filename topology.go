package atf

/*
 * topology.go - file defining Topology struct and its methods
 *
 */

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Topology represents a list of SystemUnderTest instances.
type Topology []*SysUnderTest

// NewTopology returns new empty instance of Topology.
func NewTopology() []*SysUnderTest { return make([]*SysUnderTest, 0) }

// String returns a human-readable representation of the SUT instance.
func (t Topology) String() string {

	txt := "TOPOLOGY\n"
    for _, s := range t {
        txt += fmt.Sprintf("%s\n%s", s.String())
    }
	return txt
}

// XML returns a XML-encoded representation of the Topology instance.
func (t Topology) XML() (string, error) {

	output, err := xml.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// JSON returns an JSON-encoded representation of the SUT instance.
func (t Topology) JSON() (string, error) {

	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}
