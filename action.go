package atf

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

// actioner interface
type Actioner interface {
	IsExecutable() bool /* is action executable? */
	IsManual() bool     /* is action Manual? */
}

// Represents a single action.
// There are 3 types of actions:
// - Automated action: this is the executable one (either script or executable
//   program).
// - Manual action: it contains only description.
// - Empty action: empty action that does nothing, defined for convenience.
// The different types are defined by means of two boolean flags: 'manual' and
// 'executable' flag. If 'executable' flag is set, this automated action, if
// 'manual' is set, this is manual action. If both are reset, we deal with
// empty action. Note that both flags must NOT be set. 
type Action struct {

	// script to be executed
	Script string

	// arguments to script (if needed)
	Args string

	// script execution success
	Result TestResult `xml:"result,attr"`

	// script execution output text
	Output string

	// description text, used mainly for manual actions
	Description string

	// is this action executable?
	executable bool `xml:"executable,attr"`

	// is this action manual?
	manual bool `xml:"manual,attr"`
}

// Return a string represenation of the Action instance.
func (a *Action) String() string {
	if a.IsManual() {
		return fmt.Sprintf("Manual Action:\n%s", a.Description)
	} else {
		if a.IsExecutable() {
			s := fmt.Sprintf("%s %s\n", a.Script, a.Args)
			return s
		} // if isexecutable
	} // if ismanual
	return fmt.Sprint(a.Script, " ", a.Args)
}

// Initialize Action: check the manual and executable flags and set them
// properly.
// This method is defined for convenience: it is advisable to run it when the
// action has NOT been defined using the 'Create*' methods. This is the case
// when actions are defined by marshaling from XML or JSON config file.
func (a *Action) Init() {

    // default result is always set to "not tested".
    a.Result = "NotTested"

	// initialy, action is neither executable not manual
	a.executable = false
	a.manual = false

	// if the action script is defined, action is executable
	// we like executable actions, so we gave them precedence
	if a.Script != "" {
		a.executable = true
		a.manual = false
	} else {
		// otherwise, if only Description is defined, we have a manual action
		if a.Description != "" {
			a.executable = false
			a.manual = true
		}
	}
}

// Is this action an executable one?
func (a *Action) IsExecutable() bool { return a.executable }

// Is this action a manual one?
func (a *Action) IsManual() bool { return a.manual }

// Returns an XML-encoded representation of the Action.
func (a *Action) Xml() (string, error) {

    output, err := xml.MarshalIndent(a, "  ", "    ")
    if err != nil {
        return "", err
    }
    return string(output), nil
}

// Returns a JSON-encoded representation of the Action.
func (a *Action) Json() (string, error) {
	b, err := json.Marshal(a) // marshal returns a []byte, not string!
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

// Execute the action.
// The action is executed only if 'executed' flag is set: consequently 
// this means that a particular action is an executable script or a program. 
// If 'manual' flag is set, the action is considered manual. If both arguments 
// are reset, that action is considered an empty (do-nothing) action.
// If we deal with non-executable action, 'description' is simply copied to
// 'output' field. Also, 'success' has a meaning only if action is executed;
// if not, 'Result' is always set to "not tested".
func (a *Action) Execute() string {

	a.Result = "NotTested" // we assume neutral status

	// We execute the action only if it's marked executable
	if a.IsExecutable() {

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

// Create a new Automated (executable) action.
// The 'script' fields is mandatory, the 'args' field can be empty string. 
// Also, the 'executed' flag must be set and the 'manual' flag reset. 
// The 'Result' flag is set to 'NotTested' by default. The 'description' field 
// has no special meaning with automated action.
func CreateAction(script string, args string) *Action {
	return &Action{script, args, "NotTested", "", "", true, false}
}

// Create a manual action.
// This is creation function for a manual action. The 'script' and 'args'
// fields are left empty, only 'description' is needed.
// The 'manual' flag is set and 'executable' flag is reset.
// Since this action is not executable, the success is set to "not tested".
func CreateManualAction(descr string) *Action {
	return &Action{"", "", "NotTested", "", descr, false, true}
}

// Create empty (do-nothing) action.
// This is creation function for empty (do-nothing) action. All fields are set
// apropriately: only flags are actually needed. The 'manual' and 'executable'
// flags are reset, 'success' flag is set to "not tested".
func CreateEmptyAction() *Action {
	return &Action{ "No action", "", "NotTested", "", "", false, false }
}
