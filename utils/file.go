/*
 * file.go -  misc utility functions for working with files 
 *
 * History:
 *  1   Jul11   MR  The initial version
 */

package utils

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

/*
 * LoadFile - read a file with 'filename' and return the contents as a string
 */
func LoadFile(path string) (text string, err error) {
	text = ""
	// open the file as read-only
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close() // always close the file

	// read the file line by line
	read := bufio.NewReader(file)
	str, err := read.ReadString('\n')
	text += str
	for err != io.EOF {
		str, err = read.ReadString('\n')
		text += str
	}
	return
}

/*
 * ReadTextFile - read a text file and return the contents as a string 
 *
 * If an error occurs during file read, we return an empty string (and 
 * an os.Error, of course).
 */
func ReadTextFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), err
}

/*
 * ReadLines - read a text file and return a list of lines 
 *
 * If an error occurs during file read, we return only a list with single empty
 * string (and an os.Error, of course).
 */
func ReadLines(filename string) (lines []string, err error) {
	// we read a file
	data, err := ioutil.ReadFile(filename)
	// if there's an error reading a file, we return a list with single empty
	// string and error
	if err != nil {
		return []string{""}, err
	}
	// now we convert the text into an array of lines
	lines = strings.Split(string(data), "\n")
	return
}

/*
 * WriteTextFile - write a text file with given path
 *
 */
func WriteTextFile(path string, contents string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(contents))
	if err != nil {
		return
	}
	return
}

/*
 * CopyFile - copy a file from source 'src' to destination (dst)
 *
 */
func CopyFile(dst, src string) (int64, error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()

	return io.Copy(df, sf)
}

/* Reports whether the named file or directory exists. */
func FileExists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

// Reports whether d is a directory.
func IsDir(d string) (status bool) {
    if fi, err := os.Stat(d); err == nil {
        if fi.IsDir() {
            status = true
        }
    }
    return
}
