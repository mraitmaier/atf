// 
//
// History:
// 1    Jun14   MR  Initial version, limited testing
//


package atf

import (
    "fmt"
    "strings"
    "encoding/xml"
    "encoding/json"
)

// A representation of the requirement.
type Requirement struct {

    //
    Id string

    // name of the requirement
    Name string                 `xml:"name,attr"`

    // short name of the requirement, usually a code of some sort or
    // abbreviation
    Short string

    // Longer description of the requirement
    Description string

    // Project
    *Project

    // current status
    Status RequirementStatus    `xml:"status:attr"`

    // priority (low, normal, high)
    Priority Priority           `xml:"priority,attr"`

    // A list of notes representing the changelog
    Notes []*Note               `xml:"Notes>Note"`
}

// Returns a string representation of the requirement.
func (r *Requirement) String() string {

    s:= fmt.Sprintf("Requirement: %s [%s]\n", r.Name, r.Short)
    s += fmt.Sprintf("Status: %s, Priority: %s\n", r.Status, r.Priority.String())
    s += fmt.Sprintf("Project: %s\n", r.Project.String())
    s += fmt.Sprintf("\n%s\n", r.Description)
    for _, n := range r.Notes {
        s += fmt.Sprintf("%s", n.String())
    }
    return s
}

// Returns an XML-encoded representation of the requirement.
func (r *Requirement) Xml() (string, error) {

    output, err := xml.MarshalIndent(r, "", "  ")
    if err != nil {
        return "", err
    }
    return string(output), nil

}

// Returns a JSON-encoded representation of the requirement.
func (r *Requirement) Json() (string, error) {

    output, err := json.Marshal(r)
    if err != nil {
        return "", err
    }
    return string(output), nil

}




// Custom requirement status type: it's basically a string but with limited set
// of values.
type RequirementStatus string

// The list of valid requirement statuses.
var ValidRequirementStatus = []string{"NEW", "ACKNOWLEDGED", "PENDING",
                                      "APPROVED", "REJECTED", "UNKNOWN" }

// Custom string representation for the RequirementStatus type.
func (r RequirementStatus) String() string {
    return strings.ToUpper(string(r))
}

// Check whether the given requirement status is valid.
func IsValidRequirementStatus(s RequirementStatus) bool {

    for _, status := range ValidRequirementStatus {
        if strings.ToUpper(string(s)) == status {
            return true
        }
    }
    return false
}
