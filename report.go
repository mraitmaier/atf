package atf

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

import (
	"path"
	"path/filepath"
	//    "fmt"
	"github.com/mraitmaier/atf/utils"
)

// Reporter defines the interface for different report generators.
type Reporter interface {
	Create(tr *TestReport) (string, error)
}

// Report defines a report structure to rule them all...
// It wraps all types of reports that ATF is aware of and defines the operations on all of those reports.
type Report struct {
	reports map[string]string
}

// CreateReport creates an empty report structure
func CreateReport() *Report {
	var rpt = make(map[string]string)
	return &Report{rpt}
}

// AddHTML adds a reference to HTML report
func (r *Report) AddHTML() { r.reports["html"] = "" }

// AddXML adds a reference to XML report
func (r *Report) AddXML() { r.reports["xml"] = "" }

// AddJSON adds a reference to JSON report
func (r *Report) AddJSON() { r.reports["json"] = "" }

// AddText adds a reference to text report
func (r *Report) AddText() { r.reports["txt"] = "" }

// Private method that creates the report with given type.
func (r *Report) create(tr *TestReport, typ string) (rpt string, err error) {

	switch typ {
	case "html":
		rpt, err = tr.HTML()
	case "xml":
		rpt, err = tr.XML()
	case "txt": // TODO: TextReport not implemented yet
	case "json":
		rpt, err = tr.JSON()
	default:
		rpt = "Unknown report type"
		err = ErrorUnknownReportType
	}
	return
}

// Create all the defined reports and write them.
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
