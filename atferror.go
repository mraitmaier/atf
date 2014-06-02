/*
 * atferror.go
 */

package atf

// An enum defining custom error values
type AtfError int

const (
	_ AtfError = iota // zero value should be empty
	ATFError_Unknown
	ATFError_Invalid_Value // substitute for EINVAL
	ATFError_Unknown_Report_Type
	ATFError_Invalid_Test_Result
)

// implementing the 'error' interface
func (e AtfError) Error() string {
	msg := "Unknown error"
	switch e {
	case ATFError_Unknown:
		msg = "Unknown Error"
	case ATFError_Invalid_Value:
		msg = "Invalid value"
	case ATFError_Unknown_Report_Type:
		msg = "Unknown report type"
	case ATFError_Invalid_Test_Result:
		msg = "Invalid test result value"
	}
	return msg
}
