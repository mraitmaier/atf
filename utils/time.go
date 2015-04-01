package utils

/*
 * time.go -  misc utility functions for working with  date/time
 *
 * The collection of some easy but handy functions regarding time/date that I
 * need in GoATF.
 * History:
 *  1   Jul11   MR  The initial version
 */

import (
	"strings"
	"time"
)

// Now returns current timestamp as a string with the following format: "2006-01-02 15:04:05".
func Now() string {
	//	t := time.Now()
	//	return t.Format("2006-01-02 15:04:05")
	return time.Now().Format("2006-01-02 15:04:05")
}

// NowFile returns current timestamp as a string with the following format: "2006_01_02_15_04_05".
// Usually used as an extension for filenames so that existing files are not overwritten.
func NowFile() string {
	//t := time.Now()
	//return t.Format("2006_01_02_15_04_05")
	return time.Now().Format("2006_01_02_15_04_05")
}

// FileConv is a small string helper function that replaces " ", ":" and "-" with "_". Usually used for dynamically
// creating filenames.
func FileConv(o string) (n string) {
	n = strings.Replace(o, " ", "_", -1)
	n = strings.Replace(n, ":", "_", -1)
	n = strings.Replace(n, "-", "_", -1)
	return
}
