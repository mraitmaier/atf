package atf

//
//
// History:
// 1    Jun14   MR  Initial version, limited testing
//

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

// Requirement represents a single requirement.
type Requirement struct {

	// Name represents the name of the requirement
	Name string `xml:"name,attr"`

	// Short represents a short name of the requirement, usually a code of some sort or abbreviation
	Short string

	// Description is a detailed description of the requirement
	Description string

	// Project represents a project that is related to the requirement
	Project

	// Status represents the current status
	Status ReqStatus `xml:"status:attr"`

	// Priority represents the priority (low, normal, high) of the requirement
	Priority Priority `xml:"priority,attr"`
}

// String returns a human-readable representation of the requirement.
func (r *Requirement) String() string {

	s := fmt.Sprintf("Requirement: %s [%s]\n", r.Name, r.Short)
	s += fmt.Sprintf("Status: %s, Priority: %s\n", r.Status, r.Priority.String())
	s += fmt.Sprintf("Project: %s\n", r.Project.String())
	s += fmt.Sprintf("\n%s\n", r.Description)
	return s
}

// XML returns an XML-encoded representation of the requirement.
func (r *Requirement) XML() (string, error) {

	output, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil

}

// JSON returns a JSON-encoded representation of the requirement.
func (r *Requirement) JSON() (string, error) {

	output, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(output), nil

}

// ReqStatus defines the custom requirement status type: it's basically a string but with limited set of values.
type ReqStatus string

// ValidRequirementStatus is the list of valid requirement statuses.
var ValidReqStatus = []string{"NEW", "ACKNOWLEDGED", "PENDING", "APPROVED", "REJECTED", "UNKNOWN"}

// String returns a human-readable representation for the ReqStatus type.
func (r ReqStatus) String() string { return strings.ToUpper(string(r)) }

// IsValidRequirementStatus checks whether the given requirement status is valid or not.
func IsValidReqStatus(s ReqStatus) bool {

	for _, status := range ValidReqStatus {
		if strings.ToUpper(string(s)) == status {
			return true
		}
	}
	return false
}
