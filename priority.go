//

package atf

import (
    "strings"
)

var ValidPriorities = []string{"LOW", "NORMAL", "HIGH", "UNKNOWN"}

type Priority string

func (p Priority) String() string {
    return strings.ToUpper(string(p))
}

func IsValidPriority(prio Priority) bool {

    for _, p := range ValidPriorities {
        if p == strings.ToUpper(string(prio)) {
            return true
        }
    }
    return false
}


