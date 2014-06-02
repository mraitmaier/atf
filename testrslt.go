/*
 * testrslt.go - implementation of the TestResult type
 *
 * This type defines the valid test results (pass/fail/xfail...) and valid 
 * operations on them.
 */

package atf

import (
    "encoding/xml"
)

// A slice of valid test result (string) values
var ValidTestResults = []string{"UnknownResult", "Pass", "Fail",
	"XFail", "NotTested"}

// Checks the validity of the test result value.
func IsValidTestResult(val string) bool {
	status := false
	for _, v := range ValidTestResults {
		if v == val {
			status = true
			break
		}
	}
	return status
}

// Custom type for handling test results.
type TestResult string

func (tr *TestResult) Xml() (x string, err error) {

	x = ""
	b, err := xml.MarshalIndent(tr, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

