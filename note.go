package atf

//

import (
	"fmt"
	"time"
)

// Note is a type representing a single note: a string representing a note itself and a formatted timestamp
// (format: "2012-12-15 15:05:05")
type Note struct {

	/* a string representing a note */
	note string

	/* a string representing a formatted timestamp */
	created string
}

// String returns a human readable representation of the Note.
func (n *Note) String() string {
	return fmt.Sprintf("[%s] %s\n", n.created, n.note)
}

// AppendNote appends a new note to a list.
func AppendNote(notes []Note, s string) []Note {
	if notes != nil {
		t := time.Now()
		note := &Note{s, t.Format("2012-12-15 15:04:05")}
		notes = append(notes, *note)
	}
	return notes
}
