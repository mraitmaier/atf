package atf

/*
 * testrpt.go  - implementation of the TestReport type
 *
 * The TestReport represents the single TestSet run report. It is basically
 * extended TestSet, I'm only adding timestamps when the execution of the
 * TestSet was started and when it was finished. As such, this report is ready
 * to be saved directly into database (regardless of its form - HTML, XML...)
 *
 * History:
 *  1   jun11 MR Initial version, limited testing
 *  2   oct11 MR HTML report generation added
 *  3   may14 MR improved and cleaned version
 */

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// TestReport represents the test report (test set that has been executed).
type TestReport struct {

	// TestSet sctructure that will be executed
	TestSet *TestSet

	// Started is an execution start timestamp (as a string)
	Started string

	// Finished is an execution finish timestamp (as a string)
	Finished string
}

// String returns a human-readable representation of the TestReport
func (tr *TestReport) String() string {
	return fmt.Sprintf("TestReport: %s\nstarted: %s\nfinished: %s\n", tr.TestSet.String(), tr.Started, tr.Finished)
}

// Name returns the name of the TestReport (which is actually the name of the TestSet).
func (tr *TestReport) Name() string { return tr.TestSet.Name }

// XML creates an XML-encoded representation of the TestReport.
func (tr *TestReport) XML() (x string, err error) {

	b, err := xml.MarshalIndent(tr, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

// JSON creates a JSON representation of the TestReport.
func (tr *TestReport) JSON() (string, error) {

	if tr.TestSet != nil {
		b, err := json.Marshal(tr)
		if err != nil {
			return "", err
		}
		return string(b[:]), err
	}
	return "", nil
}

// HTML creates a HTML representation of the TestReport. Uses HTML5 standard.
func (tr *TestReport) HTML() (string, error) {

	var html = ""
	if tr.TestSet != nil {
		html += tr.addHeader2Html()
		for _, tc := range tr.TestSet.Cases {
			html += tr.addTestCase2Html(tc)
		}
	}
	return html, nil
}

// Add a <header> section to HTML report
func (tr *TestReport) addHeader2Html() string {

	html := fmt.Sprintln("<header>")
	html += fmt.Sprintf("<h1>Test Report: %s</h1>\n", tr.TestSet.Name)
	html += fmt.Sprintln("<table>")
	html += fmt.Sprintln("<tr><td><b>Execution Started</b></td>")
	html += fmt.Sprintf("<td>%s</td></tr>\n", tr.Started)
	html += fmt.Sprintln("<tr><td><b>Execution Finished</b></td>")
	html += fmt.Sprintf("<td>%s</td></tr>\n", tr.Finished)
	html += fmt.Sprintln("</table>")
	html += fmt.Sprintln("<p />")
	if tr.TestSet.Sut != nil {
		html += fmt.Sprintln(tr.addSut2Html(tr.TestSet.Sut))
	}
	html += fmt.Sprintln("<table>")
	if tr.TestSet.Setup != nil {
		html += fmt.Sprintf("<tr><td>Setup</td><td>%s</td>",
			tr.TestSet.Setup.String())
		html += fmt.Sprintf("<td class=%q>%s</td></tr>\n",
			resolveHTMLClass(tr.TestSet.Setup), tr.TestSet.Setup.Result)
	}
	if tr.TestSet.Cleanup != nil {
		html += fmt.Sprintf("<tr><td>Cleanup</td><td>%s</td>",
			tr.TestSet.Cleanup.String())
		html += fmt.Sprintf("<td class=%q>%s</td></tr>\n",
			resolveHTMLClass(tr.TestSet.Cleanup), tr.TestSet.Cleanup.Result)
	}
	html += fmt.Sprintln("</table>")
	html += fmt.Sprintln("</header>")
	return html
}

// Add a system under test data to HTML report.
func (tr *TestReport) addSut2Html(sut *SysUnderTest) string {

	html := fmt.Sprintln("<table>")
	html += fmt.Sprintf("<tr><th>System Under Test</th><th>%s</th></tr>\n",
		sut.Name)
	html += fmt.Sprintf("<tr><td>Type</td><td>%s</td></tr>", sut.Systype)
	html += fmt.Sprintf("<tr><td>Version</td><td>%s</td></tr>", sut.Version)
	html += fmt.Sprintf("<tr><td>IP Address</td><td>%s</td></tr>", sut.IPaddr)
	html += fmt.Sprintf("<tr><td>Description</td><td>%s</td></tr>",
		sut.Description)
	html += fmt.Sprintln("</table>")
	html += fmt.Sprintln("<p />")
	return html
}

// Add a test case data to HTML report.
func (tr *TestReport) addTestCase2Html(tc *TestCase) string {

	html := "<article>\n"
	html += fmt.Sprintf("<h3>Test Case: %s</h3>", tc.Name)
	html += "<table>\n"
	html += fmt.Sprintf("<tr><th class=%q>Name</th><th>Action</th>", "name")
	html += fmt.Sprintf("<th class=%q>Expected Status</th>", "status")
	html += fmt.Sprintf("<th class=%q>Status</th></tr>\n", "status")
	if tc.Setup != nil {
		html += fmt.Sprintf("<tr><td>Setup</td><td>%s</td><td>Pass</td>",
			tc.Setup.String())
		html += fmt.Sprintf("<td class=%q>%s</td></tr>\n",
			resolveHTMLClass(tc.Setup), tc.Setup.Result)
	}
	for _, step := range tc.Steps {
		html += tr.addStep2Html(step)
	}
	if tc.Cleanup != nil {
		html += fmt.Sprintf("<tr><td>Cleanup</td><td>%s</td><td>Pass</td>",
			tc.Cleanup.String())
		html += fmt.Sprintf("<td class=%q>%s</td></tr>\n",
			resolveHTMLClass(tc.Cleanup), tc.Cleanup.Result)
	}
	html += fmt.Sprintln("</table><p />")
	html += "</article>\n"
	return html
}

// Add a test step data to HTML report.
func (tr *TestReport) addStep2Html(step *TestStep) string {

	// let's see if step has passed and set the HTML class accordingly
	//fmt.Printf("DEBUG step: %s\n", step.String()) // DEBUG
	class := resolveHTMLClass(step)
	html := fmt.Sprintf("<tr><td>%s</td>", step.Name)
	html += fmt.Sprintf("<td>%s</td><td>%s</td>",
		step.Action.String(), step.Expected)
	html += fmt.Sprintf("<td class=%q>%s</td></tr>\n", class, step.Status)
	return html
}

// Takes a structure and determines which CSS class should be used in HTML
// report. Only 'Action' (for setup and cleanup actions) and 'TestStep' types
// are evaluated. The CSS classes are used to define background color according
// to status of the Action/TestStep: red, green etc.
func resolveHTMLClass(structure interface{}) (cls string) {

	cls = ""
	switch t := structure.(type) {

	case *Action:
		switch t.Result {
		case "Pass":
			cls = "passed"
		case "Fail":
			cls = "failed"
		case "NotTested":
			cls = "nottested"
		}

	case *TestStep:
		switch t.Status {
		case "Pass":
			cls = "passed"
		case "Fail":
			cls = "failed"
		case "NotTested":
			cls = "nottested"
		}
	}
	return cls
}

// CreateTestReport creates a new TestReport instance with given TestSet.
func CreateTestReport(ts *TestSet) *TestReport { return &TestReport{ts, "", ""} }
