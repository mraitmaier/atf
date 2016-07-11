package atf

/*
 * atferror.go
 */

// Error is an enum defining custom error values
type Error int

const (
	// zero value should be empty
	_ Error = iota
	// ErrorUnknown represents the unknown error
	ErrorUnknown
	// ErrorInvalidValue is a substitute for EINVAL
	ErrorInvalidValue
	// ErrorUnknownReportType is FIXME
	ErrorUnknownReportType
	// ErrorInvalidTestResult is FIXME
	ErrorInvalidTestResult
)

// Error implements the 'error' interface
func (e Error) Error() string {
	msg := "Unknown error"
	switch e {
	case ErrorUnknown:
		msg = "Unknown Error"
	case ErrorInvalidValue:
		msg = "Invalid value"
	case ErrorUnknownReportType:
		msg = "Unknown report type"
	case ErrorInvalidTestResult:
		msg = "Invalid test result value"
	}
	return msg
}
