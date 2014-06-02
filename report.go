/*
 * report.go - implementation of the Reporter module
 *
 * This module is repsonsible for creating reports. According to input data,
 * different reports can be created: HTML, XML, JSON and plain text (the last
 * one has not been implemented yet and it might be omitted in the end, since
 * I'm not sure this is actually needed). These reports are written as files to
 * a specified path. By default, only HTML report is
 * created.
 *
 * History:
 *  1   Jul10   MR  The initial version
 */

package atf

import (
	"path"
	"path/filepath"
	//    "fmt"
	"bitbucket.org/miranr/goatf/atf/utils"
)

// Defined the interface for different report generators.
type Reporter interface {
	Create(tr *TestReport) (string, error)
}

// A report structure to rule them all...
// It wraps all types of reports that ATF is aware of and defines the
// operations on all of those reports.
type Report struct {
	reports map[string]string
}

// Create an empty report structure 
func CreateReport() *Report {
	var rpt = make(map[string]string)
	return &Report{rpt}
}

// Add a reference to HTML report 
func (r *Report) AddHtml() { r.reports["html"] = "" }

// Add a reference to XML report 
func (r *Report) AddXml() { r.reports["xml"] = "" }

// Add a reference to text report 
func (r *Report) AddJson() { r.reports["json"] = "" }

// Add a reference to JSON report 
func (r *Report) AddText() { r.reports["txt"] = "" }

// Private method that creates the report with given type.
func (r *Report) create(tr *TestReport, typ string) (rpt string, err error) {
	switch typ {
	case "html":
		rpt, err = tr.Html()
	case "xml":
		rpt, err = tr.Xml()
	case "txt": // TODO: TextReport not implemented yet
	case "json":
		rpt, err = tr.Json()
	default:
		rpt = "Unknown report type"
		err = ATFError_Unknown_Report_Type
	}
	return
}

// Create all the defined reports and write them
func (r *Report) Create(tr *TestReport, pth string) (err error) {
	// if path is empty, create the default path
	if pth == "" {
		pth = "."
	}
	// iterate through existing report (types), create them and write them as
	// "report.<type>" into given path
	for i, contents := range r.reports {
		contents, err = r.create(tr, i)
		if err != nil {
			return err
		}
		filename := filepath.ToSlash(path.Join(pth, "report."+i))
		err = utils.WriteTextFile(filename, contents)
		if err != nil {
			return err
		}
	}
	return
}
