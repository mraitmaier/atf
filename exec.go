/*
 * exec.go - a file implementing a simple script/program executor 
 *
 * The executor means that it is capable of executing different types of
 * scripts, including native programs and java jars.
 * Currently, it supports executing of Python, Perl, Tcl and Ruby scripts, as
 * well as executing the native (compiled) executables and executing of java
 * JARs (and only JARs!). This should suffice for some time...
 *
 * NOTE: there's one simple condition: interpreters MUST be in PATH; that
 * should not too difficult to fulfill since this a convenience.
 *
 * History:
 * 0.1  Apr10   MR  The first working version with limited testing 
 * 0.2  Mar12   MR  type ExecDisplayFnCback defined
 */
package atf

import (
	"os/exec"
	//"fmt"
	"path"
	"runtime"
)

/*
 * ExecDisplayFnCback -  This type is an alias for a closure that is used as a
 * parameter of Execute() method of the Executor interface (see below). It as
 * a callback function that mimics fmt.Printxy() functions to comment how the 
 * execution is proceeding; therefore more or less a log function.
 */
type ExecDisplayFnCback func(...string)

// Executor interface defining the Execute() method
type Executor interface {
	Execute(ExecDisplayFnCback) string
}

// String constants defining different script/program executors
const (
	pyExec     = "python"
	plExec     = "perl"
	tclExec    = "tclsh"
	expExec    = "expect"
	javaExec   = "java"
	rubyExec   = "ruby"
	groovyExec = "groovy"
)

// Define some executable types as enum
type ScriptType int
const (
	UnknownScript ScriptType = iota
	PythonScript
	PerlScript
	TclScript
	IxiaTclScript
	ExpectScript
	NativeExecutable
	JavaExecutable
	RubyScript
	GroovyScript
)

// Formats the output text from script/program.
func FmtOutput(o string) string {
	s := "Displaying output:\n################### OUTPUT ##################\n"
	s += o
	s += "################ OUTPUT END #################\n"
	return s
}

// Private function that actually executes the given script/program
// and returns the output and/or error code.
//
// Input:
//       exe - an interpreter for given script or program to be executed
//      args - arguments to the interpreter as slice of string; the script 
//          name is always included, of course. Any additional argument are to
//          be a part of this slice.
//
// eturns:
//      output - is the text output from the executed script/program
//         err - error code; if everything is OK, it should be nil
func execute(exe string, args []string) (output string, err error) {


	output = ""
    // simple error check
	if len(exe) < 1 {
		err = ATFError_Invalid_Value
        return
    }

	// prepare data for execution
	cmd := exec.Command(exe, args...)
	if cmd == nil {
		return
	}

    // run the command and wait for output text from STDIN and STDERR combined
	var out []byte
	out, err = cmd.CombinedOutput()
	output = string(out)
	return
}

// A private function that prepares arguments for executing the JARs.   
//
// Input:
//      jar  - a java JAR to be run 
//      args - additional arguments for the JAR as a slice of strings
//
// Returns:
//      out - is the text output from the executed script/program
//      err - error code; if everything is OK, it should be nil
func executeJava(jar string, args []string) (out string, err error) {
	realargs := make([]string, len(args)+3)
	realargs[0] = "-jar"
	realargs[1] = jar
	if len(args) > 0 {
		for ix, val := range args {
			realargs[ix+3] = val
		} // for
	} // if
	out, err = execute(javaExec, realargs)
	return out, err
}

// A private function that prepares arguments for executing the various scripts
// Script interpretter must be in PATH.
//
// Input:
//      exe - an executable that'll run the script (interpreter)
//      script  - a python script to be run 
//      args - additional arguments for the script as a slice of strings
// 
// Returns:
//      out - is the text output from the executed script/program
//      err - error code; if everything is OK, it should be nil
func executeScript(exe string, script string, args []string) (out string, err error) {
	// we need to insert an empty string before our args for python script to
	// run properly
	realargs := make([]string, len(args)+2)
	realargs[0] = script
	if len(args) > 0 {
		for ix, val := range args {
			realargs[ix+2] = val
		} // for
	} // if
	out, err = execute(exe, realargs)
	return out, err
}

// A private function that determines the type of script to be executed. This 
// is done by examining the file extension. If extension is not found (is empty
// string), the file is considered a native executable (true for POSIX OSes).
// 
// Input:
//      scr  - a file whose type is to be determined
//
// Returns:
//      a type of the script/program
func determineType(scr string) ScriptType {
	var t ScriptType
	e := path.Ext(scr)
	switch e {
	case "", ".exe", ".com", ".bat":
		t = NativeExecutable
	case ".py":
		t = PythonScript
	case ".pl":
		t = PerlScript
	case ".tcl":
		t = TclScript
	case ".exp":
		t = ExpectScript
	case ".rb":
		t = RubyScript
	case ".jar":
		t = JavaExecutable
	case ".groovy":
		t = GroovyScript
	default:
		t = UnknownScript
	}
	return t
}

// Executes the given script/program and returns the text output of the command
// (STDOUT & STDERR) and error code if something goes wrong.
// 
// Input:
//      script - a python script to be run 
//        args - additional arguments for the script as a slice of strings
// 
// Returns:
//      output - is the text output from the executed script/program
//         err - error code; if everything is OK, it should be nil
func Execute(script string, args []string) (output string, err error) {
	var scrtype ScriptType
	scrtype = determineType(script)
	switch scrtype {
	case PythonScript:
		output, err = executeScript(pyExec, script, args)
	case PerlScript:
		output, err = executeScript(plExec, script, args)
	case TclScript:
		output, err = executeScript(tclExec, script, args)
	case ExpectScript:
		// if we execute the script on WinXY, expect scripts are treated as
		// the TCL scripts; expect on Win is only a TCL extension, not the
		// separate interpreter
		if runtime.GOOS == "windows" {
			output, err = executeScript(tclExec, script, args)
		}
		output, err = executeScript(expExec, script, args)
	case NativeExecutable:
		output, err = execute(script, args)
	case JavaExecutable:
		output, err = executeJava(script, args)
	case RubyScript:
		output, err = executeScript(rubyExec, script, args)
	case GroovyScript:
		output, err = executeScript(groovyExec, script, args)
	default:
		output = "XXX: Invalid output"
		err = ATFError_Invalid_Value
	}
	return output, err
}
