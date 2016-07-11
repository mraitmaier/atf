package atf

// priority.go -

import (
	"strings"
)

// ValidPriorities is a list of valid priority values.
var ValidPriorities = []string{"LOW", "NORMAL", "HIGH", "UNKNOWN"}

// Priority defines the priority.
type Priority string

// String returns a human-readable representation of the Priority.
func (p Priority) String() string { return strings.ToUpper(string(p)) }

// IsValidPriority returns indication if given priority value is valid or not.
func IsValidPriority(prio Priority) bool {

	for _, p := range ValidPriorities {
		if p == strings.ToUpper(string(prio)) {
			return true
		}
	}
	return false
}
