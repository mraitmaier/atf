package atf

/*
 * collect.go - implementation of the collector module
 *
 * Collector is a module that collects the configuration (from configuration
 * file) and builds the type hierarchy (that is: scripts) to be executed.
 * The configuration can be encoded as JSON or XML or plain text (that one is
 * not implemented yet and frankly I'm not sure that is actually needed; so it
 * might be omitted in the end...)
 *
 * History:
 *  1   Apr10   MR  The initial version
 *  1.1 Jul11   MR  JSON works, XML is out for good (too complex to handle)
 *  2   Mar12   MR  XML is back and it works (with xml.Unmarshal()!), too;
 *                  had to change XML schema and add an <Action> tag
 *                  into <TestStep>
 *  3   May14   MR  A refactoring and simplification of the collector code
 *  4   Sep14   MR  More simplification of the collector code
 */

import (
	"encoding/json"
	"encoding/xml"
	"github.com/mraitmaier/atf/utils"
	"io"
	"path"
)

// Collector defines the types that implement Collect() method.
type Collector interface {
	Collect(pth string, ts *TestSet) error
}

// JSONCollector defines the JSON collector type.
type JSONCollector string

// Collect implements the Collector interface.
func (c *JSONCollector) Collect(text string, ts *TestSet) error {
	return json.Unmarshal([]uint8(text), ts)
}

// XMLCollector defines the XML collector type.
type XMLCollector string

// Collect implements the Collector interface.
func (c *XMLCollector) Collect(text string, ts *TestSet) error { return xml.Unmarshal([]byte(text), ts) }

// TextCollector defines the plain text collector type.
type TextCollector string

// Collect implements the Collector interface.
func (c *TextCollector) Collect(pth string, ts *TestSet) error {
	// TODO: no implementation yet, returning empty pointer
	return nil
}

// Collect is a public factory function that resolves the right collector type and reads the config. The final result is the
// valid TestSet structure, ready to be executed.
func Collect(pth string) (ts *TestSet) {

	// let's create empty TestSet
	ts = new(TestSet)

	// we need one of the Collectors to get test set data
	var c Collector

	// determine the type of config file and unmarshal the data into TestSet
	switch path.Ext(pth) {
	case ".json":
		c = new(JSONCollector)
	case ".txt", ".cfg":
		c = new(TextCollector)
	case ".xml":
		c = new(XMLCollector)
	default:
		return nil
	}

	// read the text file
	text, err := utils.ReadTextFile(pth)
	if err != nil && err != io.EOF {
		return nil
	}

	// now collect the test set structure and update flags for actions
	c.Collect(text, ts)
	ts.Initialize()
	// silently drop error: if 'ts' is 'nil', it is an error already...

	return
}
