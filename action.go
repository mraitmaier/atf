package atf

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

// Action represents a single action.
// There are 3 types of actions:
// - Automated action: this is the executable one (either script or executable program).
// - Manual action: it contains only description.
// - Empty action: empty action that does nothing, defined for convenience.
// The different types are defined by means of two boolean flags: Manual and Executable flag. If Executable flag is set,
// this automated action, if Manual is set, this is manual action. If both are reset, we deal with empty action. Note that
// both flags must NOT be set.
type Action struct {

	// Script to be executed
	Script string

	// Args represents arguments to script (if needed)
	Args string

	// Result is script execution success
	Result TestResult `xml:"result,attr"`

	// Output is script execution output text
	Output string

	// Description text, used mainly for manual actions
	Description string

	// Executable: is this action executable?
	Executable bool `xml:"executable,attr"`

	// Manual: is this action manual?
	Manual bool `xml:"manual,attr"`
}

// String returns a human-readable represenation of the Action instance.
func (a *Action) String() string {

	if a.Manual {
		return fmt.Sprintf("Manual Action:\n%s", a.Description)
	} else if a.Executable {
		s := fmt.Sprintf("%s %s\n", a.Script, a.Args)
		return s
	} // if isexecutable
	return fmt.Sprint(a.Script, " ", a.Args)
}

// Init initializes the action: check the manual and executable flags and set them properly.
// This method is defined for convenience: it is advisable to run it when the action has NOT been defined using the 'Create*'
// methods. This is the case when actions are defined by marshaling from XML or JSON config file.
func (a *Action) Init() {

	// default result is always set to "not tested".
	a.Result = "NotTested"

	// initialy, action is neither executable not manual
	a.Executable = false
	a.Manual = false

	// if the action script is defined, action is executable
	// we like executable actions, so we gave them precedence
	if a.Script != "" {
		a.Executable = true
		a.Manual = false
	} else {
		// otherwise, if only Description is defined, we have a manual action
		if a.Description != "" {
			a.Executable = false
			a.Manual = true
		}
	}
}

// XML returns an XML-encoded representation of the Action.
func (a *Action) XML() (string, error) {

	output, err := xml.MarshalIndent(a, "  ", "    ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// JSON returns a JSON-encoded representation of the Action.
func (a *Action) JSON() (string, error) {

	b, err := json.Marshal(a) // marshal returns a []byte, not string!
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

// Execute executes the action.
// The action is executed only if 'executed' flag is set: consequently this means that a particular action is an executable
// script or a program. If 'manual' flag is set, the action is considered manual. If both arguments are reset, that action is
// considered an empty (do-nothing) action. If we deal with non-executable action, 'description' is simply copied to
// 'output' field. Also, 'success' has a meaning only if action is executed; if not, 'Result' is always set to "not tested".
func (a *Action) Execute() string {

	a.Result = "NotTested" // we assume neutral status

	// We execute the action only if it's marked executable
	if a.Executable {

		var err error
		a.Output, err = Execute(a.Script, strings.Split(a.Args, " "))

		// if error has accured, script has failed; otherwise, it's OK
		if err != nil {
			a.Result = "Fail"
		} else {
			a.Result = "Pass"
		}
	} else {
		// otherwise we just put description into output, success is already set
		a.Output = a.Description
	}
	return a.Output
}

// CreateAction creates a new Automated (executable) action.
// The 'script' fields is mandatory, the 'args' field can be empty string. Also, the 'executed' flag must be set and the
// 'manual' flag reset. The 'Result' flag is set to 'NotTested' by default. The 'description' field has no special meaning
// with automated action.
func CreateAction(script string, args string) *Action {
	return &Action{script, args, "NotTested", "", "", true, false}
}

// CreateManualAction creates new a manual action.
// This is creation function for a manual action. The 'script' and 'args' fields are left empty, only 'description' is needed.
// The 'manual' flag is set and 'executable' flag is reset. Since this action is not executable, the success is set to
// "not tested".
func CreateManualAction(descr string) *Action {
	return &Action{"", "", "NotTested", "", descr, false, true}
}

// CreateEmptyAction creates a new empty (do-nothing) action.
// This is creation function for empty (do-nothing) action. All fields are set apropriately: only flags are actually needed. The i
// 'manual' and 'executable' flags are reset, 'success' flag is set to "not tested".
func CreateEmptyAction() *Action { return &Action{"No action", "", "NotTested", "", "", false, false} }
