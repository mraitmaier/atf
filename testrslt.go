package atf

/*
 * testrslt.go - implementation of the TestResult type
 *
 * This type defines the valid test results (pass/fail/xfail...) and valid
 * operations on them.
 */

import (
	"encoding/xml"
)

// ValidTestResults is a slice of valid test result (string) values
var ValidTestResults = []string{"UnknownResult", "Pass", "Fail", "XFail", "NotTested"}

// IsValidTestResult checks for the validity of the given test result value.
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

// TestResult is a custom type for handling test results.
type TestResult string

// XML returns an XML-encoded representation of the TestResult
func (tr *TestResult) XML() (x string, err error) {

	x = ""
	b, err := xml.MarshalIndent(tr, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}
