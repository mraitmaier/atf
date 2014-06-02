
package utils

import (
//    "bytes"
)

// deep copy a string: this is simply a trick to force compiler to copy the 
// string it returns by slicing it.
// Source: https://groups.google.com/forum/#!topic/golang-nuts/naMCI9Jt6Qg
/*
func CopyS(a string) string {
	if len(a) == 0 {
		return ""
	}
	return a[0:1] + a[1:]
}
*/

// deep copy a string: this is simply a trick to force compiler to copy the 
// string by adding a space and returning a string without that additional
// space.
// Source: https://groups.google.com/forum/#!topic/golang-nuts/naMCI9Jt6Qg
/*
func CopyS(a string) string {
	return (a + " ")[:len(a)]
}
*/

// Deep copy a string by manipulating strings as byte slices and using built-in // copy function. As basic as it gets...
func CopyS(s string) string {
    a := []byte(s)
    b := []byte("")
    copy(b, a)
    return string(b)
}

