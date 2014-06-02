/*
 * teststep.go - implementation of the TestStep type 
 *
 * This data structure represents the single test case (executable) step (or
 * action). It is always expected for a step to pass, so the self-evaluation is
 * as simple as possible.
 *
 * History:
 *  1   Apr10 MR Initial version, limited testing
 *  2   May14 MR Improved version, action and status handling is now accurate.
 */

package atf

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Represents a single test step (action with additional data).
type TestStep struct {

	/* name of the test step; in XML, this is an attribute */
	Name string         `xml:"name,attr"`

	/* expected status of the step; in XML, this is an attribute */
	Expected TestResult `xml:"expected,attr"`

	/* status of the step; in XML, this is an attribute */
	Status TestResult   `xml:"status,attr"`

	/* every test step needs an action: either manual or executable */
	Action *Action      `xml:"Action"`
}

// Returns a string representation of the TestStep instance.
func (ts *TestStep) String() string {

	var act string
	// let's check the action first...
	if ts.Action != nil {
		act = ts.Action.String()
	} else {
		act = "none"
	}

	return fmt.Sprintf(
		"TestStep: %q expected: %q status: %q action: %q\n",
		ts.Name, ts.Expected, ts.Status, act)
}

// Displays a TestStep. Meant mainly for testing & debugging purposes.
func (ts *TestStep) Display() string {
	txt := fmt.Sprintf("TestStep: %q\n", ts.Name)
	txt += fmt.Sprintf("Expected status: %q\n", ts.Expected)
	txt += fmt.Sprintf("Status: %q\n", ts.Status)
	if ts.Action != nil {
		txt += fmt.Sprintf("Action: %q\n", ts.Action.String())
	} else {
		txt += "Action: N/A\n"
	}
	return txt
}

// Returns an XML-encoded represenation of the TestStep instance.
func (ts *TestStep) Xml() (string, error) {

    output, err := xml.MarshalIndent(ts, "", "  ")
    if err != nil {
        return "", nil
    }

    return string(output), nil
}

// Returns a JSON-encoded represenation of the TestStep instance.
func (ts *TestStep) Json() (string, error) {
	b, err := json.Marshal(ts)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

// Returns a HTML-encoded represenation of the TestStep instance.
func (ts *TestStep) Html() (string, error) {
	// TODO
	return "", nil
}

// Initialize the test step.
// Note that when step's action is empty, the method will panic (this is
// unacceptable condition!).
func (ts *TestStep) Initialize() {

    // if action is empty, just panic, this is not acceptable...
    if ts.Action == nil {
        panic("Test step action is empty!")
    }
    ts.Action.Init()

    // default step status is "not tested"
    ts.Status = "NotTested"

    // if expected status is empty for executable action, force "Pass"
    if ts.Action.IsExecutable() && ts.Expected == "" {
        ts.Expected = "Pass"
    }
}

// Execute the TestStep.
func (ts *TestStep) Execute(display *ExecDisplayFnCback) {

	// we turn the function ptr back to function 
	disp := *display

	// and start the execution
	disp("info", fmt.Sprintf(">>> Entering test step %q\n", ts.Name))

	// we execute the action when it's not empty
	if ts.Action != nil && ts.Action.IsExecutable() {
	    disp("notice", fmt.Sprintf("Executing test step action: %q\n",
                ts.Action.String()))
		disp("info", FmtOutput(ts.Action.Execute()))
	} else {
		disp("error", fmt.Sprintln("Action is EMPTY?????"))
	}

	// let's evaluate expectations and final status of the step
	switch ts.Expected {
	case "Pass":
		if ts.Action.Result == "Pass" {
			ts.Status = "Pass"
		} else {
			ts.Status = "Fail"
		}
	case "XFail":
		if ts.Action.Result == "Pass" {
			ts.Status = "Fail"
		} else {
			ts.Status = "Pass"
		}
	default:
		//only Pass & XFail are allowed as expected status 
		ts.Status = "NotTested"
	}
	disp("notice", fmt.Sprintf("Test step evaluated to %q\n", ts.Status))
    disp("info", fmt.Sprintf("<<< Leaving test step %q\n", ts.Name))
}

// Create a new TestStep instance.
func CreateTestStep(name string, descr string, expected TestResult,
	status TestResult, act *Action) *TestStep {
	return &TestStep{name, expected, status, act}
}
