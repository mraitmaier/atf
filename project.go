// project.go

package atf

import (
    "fmt"
//    "strings"
    "encoding/json"
    "encoding/xml"
)

type Project struct {

    //
    Name string

    //
    Short string

    //
    Description string

}

// Create a new instance of Project, name is given and short name (abbrviation) is needed.
func NewProject(name, short string) *Project {
    return &Project{ name, short, "" }
}

// Create a new instance of Project, all data is given.
func CreateProject(name, short, descr string) *Project {
    return &Project{ name, short, descr }
}

// Returns a string representation of the Project instance
func (p *Project) String() string {
    return fmt.Sprintf("%s (%s)", p.Name, p.Short)
}

// Returns an XML-encoded representation of the Project instance
func (p *Project) Xml() (string, error) {

    out, err := xml.MarshalIndent(p, "  ", "    ")
    if err != nil {
        return "", err
    }
    return string(out), err
}

// Returns a JSON-encoded representation of the Project instance
func (p *Project) Json() (string, error) {
    b, err := json.Marshal(p)
    if err != nil {
        return "", err
    }
    return string(b[:]), err
}
