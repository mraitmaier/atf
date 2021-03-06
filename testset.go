package atf

/*
 * testset.go - implementation of the TestSet type
 *
 * The TestSet is an executable version of TestPlan: TestSet is executed, while
 * the TestPlan should serve only as a reference to a document that is stored
 * into database.
 *
 * NOTE: I'm not sure this is the right approach, but as a Go-learning step,
 * this is how it's currently done. It might even disappear in the future (or
 * the TestSet type).
 *
 * History:
 *  1   Apr10 MR Initial version, limited testing
 *  2   May14 MR Improved, simplified version: XML handling simplified,
 *               appending cases simplified, conversion to TestPlan added.
 */

import (
	//"github.com/mraitmaier/atf/utils"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// TestSet represents an executable set of test cases.
// Note that TestSet has a sibling type TestPlan. TestPlan is non-executable version of TestSet, otherwise holds the same
// information, with the exception of SUT data. Since a single TestSet can belong to more TestPlans, the TestSet holds the
// TestPlan property which is the name of the TestPlan it is associated with.
type TestSet struct {

	//  ID is a unique ID of the TestSet, used for DB access
	//ID string `bson:"_id, omitempty"`

	// Name is a test set name, of course; in XML, this is an attribute
	Name string `xml:"name,attr"`

	// Description is a arbitrary long text description of the test set
	Description string

	// TestPlan: test set is a subset of test plan; we remember its name
	//TestPlan string

	// Sut is a system under test description
	Sut *SysUnderTest `xml:"SystemUnderTest"`

	// Setup is a setup action
	Setup *Action `xml:"Setup"`

	// Cleanup is a cleanup action
	Cleanup *Action `xml:"Cleanup"`

	// Cases is a list of test cases; in XML, this is a list of <TestCase> tags
	Cases []*TestCase `xml:"Cases>TestCase"`
}

/*
// ToTestPlan converts a TestSet instance into TestPlan instance.
// Note that we force deep copy of the data. Also, SUT instance is not contained by TestPlan, so it must be omitted.
func (ts *TestSet) ToTestPlan() *TestPlan {

	tp := new(TestPlan)
	tp.Name = utils.CopyS(ts.TestPlan)
	tp.Description = utils.CopyS(ts.Description)
	*tp.Setup = *ts.Setup
	*tp.Cleanup = *ts.Cleanup
	for _, tcase := range ts.Cases {
		tp.Cases = append(tp.Cases, tcase)
	}
	return tp
}
*/

// Initialize initializes a new TestSet.
func (ts *TestSet) Initialize() {

	// Create empty actions for setup & cleanup, when empty
	if ts.Setup == nil {
		ts.Setup = CreateEmptyAction()
	}
	if ts.Cleanup == nil {
		ts.Cleanup = CreateEmptyAction()
	}

	for _, tcase := range ts.Cases {
		tcase.Initialize()
	}
}

// String returns a human-readable representation of the TestSet instance.
func (ts *TestSet) String() string {

	s := fmt.Sprintf("TestSet: %q", ts.Name)
	//s += fmt.Sprintf(" is owned by %q test plan.\n", ts.TestPlan)
	s += fmt.Sprintf("  Description:\n%q\n", ts.Description)
	if ts.Sut != nil {
		s += fmt.Sprintf("  SUT:\n%s\n\n", ts.Sut.String())
	}
	if ts.Setup != nil {
		s += fmt.Sprintf("  Setup: %s", ts.Setup.String())
	} else {
		s += fmt.Sprintln("  Setup: []")
	}
	if ts.Cleanup != nil {
		s += fmt.Sprintf("  Cleanup: %s", ts.Cleanup.String())
	} else {
		s += fmt.Sprintln(" Cleanup: []")
	}
	for _, v := range ts.Cases {
		s += fmt.Sprintf("\n%s", v.String())
	}
	return s
}

// XML returns an XML-encoded representation of the TestSet instance.
func (ts *TestSet) XML() (string, error) {

	output, err := xml.MarshalIndent(ts, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// JSON returns a JSON-encoded representation of the TestSet instance.
func (ts *TestSet) JSON() (string, error) {

	b, err := json.Marshal(ts)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

// HTML returns a HTML-encoded representation of the TestSet instance.
func (ts *TestSet) HTML() (string, error) {
	// TODO
	return "", nil
}

// Append one or more test cases to the list of cases.
func (ts *TestSet) Append(set ...*TestCase) {
	ts.Cases = append(ts.Cases, set...)
}

// CleanupAfterTsetSetupFail performs a clenaup of data when execution of the setup action fails.
func (ts *TestSet) CleanupAfterTsetSetupFail() string {

	o := "Setup has FAILED\n"
	o += "Stopping the complete test set execution.\n"
	// mark all tcs & cases as skipped
	for _, tc := range ts.Cases {
		for _, step := range tc.Steps {
			step.Status = "NotTested"
		}
	}
	o += fmt.Sprintf("<<< Leaving test set %q\n", ts.Name)
	return o
}

// Execute executes the entire TestSet.
func (ts *TestSet) Execute(display *ExecDisplayFnCback) {

	output := ""

	// define function from function pointer
	disp := *display

	// execute the cleanup action
	disp("notice", fmt.Sprintf(">>> Entering Test Set %q\n", ts.Name))
	if ts.Setup != nil && ts.Setup.Executable {
		disp("notice", fmt.Sprintf("Executing setup script: %q\n",
			ts.Setup.String()))
		output = ts.Setup.Execute()
		disp("info", FmtOutput(output))
		// if setup script has failed, there's no need to proceed...
		if ts.Setup.Result == "Fail" {
			disp("error", ts.CleanupAfterTsetSetupFail())
		}
	} else {
		disp("notice", fmt.Sprintln("Setup action is not defined."))
	}

	// execute test cases
	if ts.Cases != nil {
		for _, tc := range ts.Cases {
			tc.Execute(display)
		}
	}

	// execute the cleanup action
	if ts.Cleanup != nil && ts.Cleanup.Executable {
		disp("notice", fmt.Sprintf("Executing cleanup script: %q\n",
			ts.Cleanup.String()))
		disp("info", FmtOutput(ts.Cleanup.Execute()))
	} else {
		disp("notice", fmt.Sprintln("Cleanup action is not defined:"))
	}
	disp("notice", fmt.Sprintf("<<< Leaving test set %q\n", ts.Name))
}

// CreateTestSet creates a new instance of the TestSet type with given data.
func CreateTestSet(name, descr string, sut *SysUnderTest, setup, cleanup *Action) *TestSet {
	var tcs []*TestCase
	return &TestSet{name, descr, sut, setup, cleanup, tcs}
}
