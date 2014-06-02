// 

package atf

import (
    "strings"
)

type Requirement struct {

    //
    Name string

    //
    Short string

    //
    Description string

    //
    Status RequirementStatus

    // 
    Priority Priority

    // A list of notes representing the changelog
    Note []Note
}

type RequirementStatus string

var ValidRequirementStatus = []string{"NEW", "ACKNOWLEDGED", "PENDING",
                                      "APPROVED", "REJECTED", "UNKNOWN" }

func (r RequirementStatus) String() string {
    return strings.ToUpper(string(r))
}

func IsValidRequirementStatus(s RequirementStatus) bool {

    for _, status := range ValidRequirementStatus {
        if strings.ToUpper(string(s)) == status {
            return true
        }
    }
    return false
}
