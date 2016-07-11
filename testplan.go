package atf

/*
 * testplan.go
 *
 * History:
 *  1   May11 MR Initial version, limited testing
 *  2   May14 MR Updated version: XML handling simplified, added conversion to
 *               TestSet, appending test cases simplified
 */

import (
	"github.com/mraitmaier/atf/utils"
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// TestPlan defines a single test plan.
// Note that the TestPlan has a sibling in TestSet type: TestSet is an executable version of TestPlan.
type TestPlan struct {
	// Name of the test plan
	Name string `xml:"name,attr"`
	// Description is a detailed description of the test plan
	Description string `xml:"Description"`
	// Setup is a setup action
	Setup *Action `xml:"Setup"`
	// Cleanup is a cleanup action
	Cleanup *Action `xml:"Cleanup"`
	// Cases is a list of test cases
	Cases []*TestCase `xml:"Cases>TestCase"`
}

// String returns a human-readable representation of the TestPlan instance.
func (tp *TestPlan) String() string { return fmt.Sprintf("TestPlan: %q\n", tp.Name) }

// XML returns a XML-encoded representation of the TestPlan instance.
func (tp *TestPlan) XML() (string, error) {

	output, err := xml.MarshalIndent(tp, "  ", "    ")
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// JSON returns a JSON-encoded representation of the TestPlan instance.
func (tp *TestPlan) JSON() (string, error) {

	b, err := json.Marshal(tp)
	if err != nil {
		return "", err
	}
	return string(b[:]), err
}

// Append appends one or more test cases to the list of test cases.
func (tp *TestPlan) Append(cases ...*TestCase) { tp.Cases = append(tp.Cases, cases...) }

// ToTestSet converts a TestPlan into a TestSet instance. Note that we force deep copy of data.
func (tp *TestPlan) ToTestSet() *TestSet {

	ts := new(TestSet)
	ts.Name = utils.CopyS(tp.Name) // TestSet name can (and should) be changed
	ts.Description = utils.CopyS(tp.Description)
	//ts.TestPlan = utils.CopyS(tp.Name)
	*ts.Setup = *tp.Setup
	*ts.Cleanup = *tp.Cleanup
	ts.Sut = new(SysUnderTest) // return empty instance
	//copy(ts.Cases, tp.Cases)
	for _, tcase := range tp.Cases {
		ts.Cases = append(ts.Cases, tcase)
	}

	return ts
}

// CreateTestPlan creates a new TestPlan instance.
func CreateTestPlan(name, descr string, setup, cleanup *Action) *TestPlan {
	var tcs []*TestCase
	return &TestPlan{name, descr, setup, cleanup, tcs}
}
