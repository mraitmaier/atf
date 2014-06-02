/*
 * testcase.go  - implementation of the TestCase type
 *
 * This type represents the test case and is the central data struct of the
 * complete application. TestCase is built from separate test steps (that are
 * self-evaluated: pass/fail) , including setup and cleanup actions, and the 
 * TestCase itself uses the evaluation algorithm to self evaluate (pass/fail) 
 * itself according to expected result.
 *
 * History:
 *  1   Apr10 MR Initial version, limited testing
 *  2   Mar12 MR heavy refactoring: changed the Execute() method to work with
 *                 registered closure; xml.Unmarshal() parsing definitions; 
 *  3   Mar12 MR case evaluation fixed
 *  4   May14 MR Improved and siplified version: XML handling simplified,
 *               appending steps simplified.
 */

package atf

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// Represents a single test case.
type TestCase struct {

	// a name of the test case; in XML, this is an attribute
	Name string         `xml:"name,attr"`

	// a test case setup action
	Setup *Action       `xml:"Setup"`

	// a test case cleanup action
	Cleanup *Action     `xml:"Cleanup"`

	// expected result for this test case: either pass or expected fail;
	// in XML, this is an attribute
	Expected TestResult `xml:"expected,attr"`

	// actual result for this test case after execution;
	// in XML, this is an attribute
	Status TestResult   `xml:"status,attr"`

	// a list of test steps; in XML, this is a sequence of <TestStep> tags
	Steps []*TestStep    `xml:"Steps>TestStep"`

	// a detailed description of the test case
	Description string
}

// Returns a plain text representation of the TestSet instance.
func (tc *TestCase) String() string {
	s := fmt.Sprintf("Test Case: %q\n\tstatus: %s \n", tc.Name, tc.Status)
	s += fmt.Sprintf("\tDescription: %q\n", tc.Description)
	s += fmt.Sprintf("\tExpected: %s \n", tc.Expected)
	if tc.Setup != nil {
		s += fmt.Sprintf("\tSetup: %s", tc.Setup.String())
	} else {
		s += fmt.Sprintf("\tSetup: none")
	}
	if tc.Cleanup != nil {
		s += fmt.Sprintf("\tCleanup: %s\n", tc.Cleanup.String())
	} else {
		s += fmt.Sprintf("\tCleanup: none\n")
	}
	if tc.Steps != nil {
		for _, step := range tc.Steps {
			s += fmt.Sprintf("%s\n", step.String())
		}
	} else {
		s += fmt.Sprintln("\tActions: empty\n")
	}
	return s
}

// Initialize the TestCase. This method is defined as a convenience.
// It is advisable to run it when TestCase instance is not defined using the
// "CreateTestCase()" method. For instance, when test cases are serialized
// (collected) from XML or JSON config file. 
func (tc *TestCase) Initialize() {

    // if setup and cleanup actions are empty....
    if tc.Setup == nil {
        tc.Setup = CreateEmptyAction()
    }
    if tc.Cleanup == nil {
        tc.Cleanup = CreateEmptyAction()
    }

    //
    for _, step := range tc.Steps {
        step.Initialize()
    }
}

// Returns an XML-encoded representation of the TestSet instance.
func (tc *TestCase) Xml() (string, error) {

    output, err := xml.MarshalIndent(tc, "  ", "    ")
    if err != nil {
        return "", err
    }
    return string(output), nil
}

// Returns a JSON-encoded representation of the TestSet instance.
func (tc *TestCase) Json() (string, error) {
	b, err := json.Marshal(tc)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

// Returns an HTML-encoded representation of the TestSet instance.
func (tc *TestCase) Html() (string, error) {
	// TODO
	return "", nil
}

// Append one or more test steps to a list of steps.
func (tc *TestCase) Append(steps ...*TestStep) {
    tc.Steps = append(tc.Steps, steps...)
}

// Cleanup data when execution of the setup action fails.
func (tc *TestCase) cleanupAfterCaseSetupFail() string {
	output := "Setup action has FAILED.\n"
	output += "Skipping the rest of the case...\n"
	output += fmt.Sprintf("<<< Leaving TestCase %q\n", tc.Name)
	tc.Status = "Fail"
	// set all steps' status to NotTested
	for _, step := range tc.Steps {
		step.Status = "NotTested"
	}
	return output
}


// Execute the entire TestCase.
func (tc *TestCase) Execute(display *ExecDisplayFnCback) {

	// we turn function ptr back to function
	disp := *display

	// and start with execution...
	disp("notice", fmt.Sprintf(">>> Entering TestCase %q\n", tc.Name))

	// let's execute setup action (if not empty)
	if tc.Setup != nil && tc.Setup.IsExecutable() {
		disp("notice", fmt.Sprintf("Executing case setup action: %q\n",
                tc.Setup.String()))
		disp("info", FmtOutput(tc.Setup.Execute()))
		// if setup action has failed, skip the rest of the case
		if tc.Setup.Result == "Fail" {
			disp("error", tc.cleanupAfterCaseSetupFail())
		}
	} else {
		disp("notice", fmt.Sprintln("Setup action is not defined.\n"))
	}

	// now we execute the steps...
	if tc.Steps != nil {
		for _, step := range tc.Steps {
			step.Execute(display)
		}
	}

	// let's execute cleanup action (if not empty)
	if tc.Cleanup != nil && tc.Cleanup.IsExecutable() {
		disp("notice", fmt.Sprintf("Executing case cleanup action: %q\n",
                tc.Cleanup.String()))
        if tc.Setup != nil {
		    disp("info", FmtOutput(tc.Setup.Execute()))
        }
	} else {
		disp("notice", fmt.Sprintln("Cleanup action is not defined.\n"))
	}
	// now we evaluate the complete test case
	tc.evaluate()
	disp("notice", fmt.Sprintf("Test case evaluated to %q\n", tc.Status))
	disp("notice", fmt.Sprintf("<<< Leaving TestCase %q\n", tc.Name))
}

// Evaluate results after the case was executed.
// There is a simple algorithm how expected status and actual statuses are
// treated. Expected status can be either Pass or XFail (expected fail).
// According to expected status, test case is evaluated as follows:
// - if setup action fails, the whole test case fails (steps are even not
//   executed...). 
// - if cleanup action fails, the whole test case fails.
// - if expected status is Pass and any of the steps fails, the whole test case
//   fails. Test case passes only if all actions pass (including setup and
//   cleanup).
// - if expected status is XFail and any of the steps passes, the whole test
//   case is evaluated to Fail. Test case passes only if all actions fail.
// - The NotTested status is treated neutral.
func (tc *TestCase) evaluate() {

	tc.Status = "Pass" // initial values is NotTested

	// otherwise compare steps' expected and final results
	switch tc.Expected {

	case "Pass": tc.evaluateExpectedPass()

	case "XFail": tc.evaluateExpectedFail()

	default:
		// by definition, only PASS & XFAIL are allowed as expected results 
		tc.Status = "NotTested"

	} // switch 
}

// Evaluate the test case status when expected status is XFail.
func (tc *TestCase) evaluateExpectedFail() {

    // evaluate setup and cleanup actions  
	// if setup or cleanup have passed, the complete test case fails 
    if tc.Setup != nil && tc.Setup.Result == "Pass" {
	    tc.Status = "Fail"
	    return
    }
    if tc.Cleanup != nil && tc.Cleanup.Result == "Pass" {
	    tc.Status = "Fail"
	    return
    }

    // If any of the steps passes, the whole test case fails.
    not_tested := 0 // we count the NotTested occurences
	for _, step := range tc.Steps {
		switch step.Status {

        case "Pass":
			tc.Status = "Fail"
			return

        case "NotTested":
            not_tested += 1
		}
	}

    // If all steps' statuses are NotTested, the whole case is obviously 
    // evaluated to NotTested.
    if not_tested == len(tc.Steps) {
        tc.Status = "NotTested"
    }
}

// Evaluate the test case status when expected status is Pass.
func (tc *TestCase) evaluateExpectedPass() {

    // evaluate setup and cleanup actions  
    if tc.Setup != nil && tc.Setup.Result == "Fail" {
	    tc.Status = "Fail"
	    return
    }
    if tc.Cleanup != nil && tc.Cleanup.Result == "Fail" {
	    tc.Status = "Fail"
	    return
    }

    // If any of the steps fails, the whole test case fails.
    not_tested := 0 // we count NotTested occurences
	for _, step := range tc.Steps {
		switch step.Status {
        case "Fail":
            tc.Status = "Fail"
			return
        case "NotTested":
            not_tested += 1
        }
	}

    // If all steps' statuses are NotTested, the whole case is obviously 
    // evaluated to NotTested.
    if not_tested == len(tc.Steps) {
        tc.Status = "NotTested"
    }
}

// Create a new instance of TestCase.
func CreateTestCase(name, descr string, setup, cleanup *Action,
	                expected, status TestResult) *TestCase {
	steps := make([]*TestStep, 0)
	return &TestCase{name, setup, cleanup, expected, status, steps, descr}
}
